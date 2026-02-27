// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package knowledge

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	vectorutils "github.com/volcengine/vikingdb-go-sdk/vector/utils"
	"github.com/volcengine/volc-sdk-golang/base"

	kmodel "github.com/volcengine/vikingdb-go-sdk/knowledge/model"
)

const requestIDHeader = "X-Tt-Logid"

type authKind int

const (
	authKindNone authKind = iota
	authKindIAM
	authKindAPIKey
)

// Auth describes how the SDK should sign outgoing requests.
type Auth struct {
	kind      authKind
	accessKey string
	secretKey string
	apiKey    string
}

// AuthNone disables request signing.
func AuthNone() Auth {
	return Auth{kind: authKindNone}
}

// AuthIAM configures AK/SK signing.
func AuthIAM(accessKey, secretKey string) Auth {
	return Auth{
		kind:      authKindIAM,
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

// AuthAPIKey configures API key authentication.
func AuthAPIKey(apiKey string) Auth {
	return Auth{
		kind:   authKindAPIKey,
		apiKey: apiKey,
	}
}

type authenticator interface {
	apply(req *http.Request) (*http.Request, error)
}

type noAuth struct{}

func (noAuth) apply(req *http.Request) (*http.Request, error) {
	return req, nil
}

type apiKeyAuth struct {
	token string
}

func (a apiKeyAuth) apply(req *http.Request) (*http.Request, error) {
	if a.token == "" {
		return req, nil
	}
	req.Header.Set("Authorization", "Bearer "+a.token)
	return req, nil
}

type iamAuth struct {
	ak     string
	sk     string
	region string
}

func (a iamAuth) apply(req *http.Request) (*http.Request, error) {
	if a.ak == "" || a.sk == "" {
		return nil, fmt.Errorf("access key and secret key cannot be empty")
	}
	cred := base.Credentials{
		AccessKeyID:     a.ak,
		SecretAccessKey: a.sk,
		Service:         "air",
		Region:          a.region,
	}
	return cred.Sign(req), nil
}

type transport struct {
	config     Config
	httpClient *http.Client
	baseURL    *url.URL
	auth       authenticator
	userAgent  string
}

func newTransport(cfg Config, authConfig Auth) (*transport, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("endpoint cannot be empty")
	}

	defaults := DefaultConfig()
	baseURL, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint: %w", err)
	}
	if baseURL.Scheme == "" {
		baseURL.Scheme = "https"
	}
	if cfg.Region == "" {
		cfg.Region = defaults.Region
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = defaults.Timeout
	}
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: cfg.Timeout}
	}
	userAgent := cfg.UserAgent
	if userAgent == "" {
		userAgent = defaults.UserAgent
	}
	var auth authenticator = noAuth{}
	switch authConfig.kind {
	case authKindIAM:
		auth = iamAuth{ak: authConfig.accessKey, sk: authConfig.secretKey, region: cfg.Region}
	case authKindAPIKey:
		auth = apiKeyAuth{token: authConfig.apiKey}
	default:
		auth = noAuth{}
	}
	if cfg.MaxRetries < 0 {
		cfg.MaxRetries = 0
	}
	return &transport{
		config:     cfg,
		httpClient: httpClient,
		baseURL:    baseURL,
		auth:       auth,
		userAgent:  userAgent,
	}, nil
}

// Client represents the entry point for interacting with Knowledge APIs.
type Client struct {
	transport *transport
}

// New constructs a Client with the provided authentication settings and options.
func New(auth Auth, opts ...ClientOption) (*Client, error) {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	tr, err := newTransport(cfg, auth)
	if err != nil {
		return nil, err
	}
	return &Client{transport: tr}, nil
}

func (t *transport) doRequest(ctx context.Context, method, path string, request interface{}, response interface{}, opts ...RequestOption) error {
	if ctx == nil {
		ctx = context.Background()
	}
	requestOpts := defaultRequestOptions()
	for _, opt := range opts {
		opt(requestOpts)
	}
	retries := requestOpts.MaxRetries
	if retries <= 0 {
		retries = t.config.MaxRetries
	}
	if retries < 0 {
		retries = 0
	}
	var body []byte
	if request != nil {
		serialized, err := vectorutils.SerializeToJSON(request)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		body = serialized
	}
	return vectorutils.Retry(retries, func() error {
		req, err := t.buildRequest(ctx, method, path, body, requestOpts)
		if err != nil {
			return err
		}
		resp, err := vectorutils.DoHTTPRequest(t.httpClient, req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return vectorutils.ParseResponse(resp, response)
	}, vectorutils.IsRetryableError)
}

func (t *transport) buildRequest(ctx context.Context, method, path string, body []byte, opts *RequestOptions) (*http.Request, error) {
	targetURL := t.baseURL.ResolveReference(&url.URL{Path: path})
	if len(opts.Query) > 0 {
		query := targetURL.Query()
		for k, v := range opts.Query {
			query.Set(k, v)
		}
		targetURL.RawQuery = query.Encode()
	}
	var buf *bytes.Reader
	if len(body) > 0 {
		buf = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, targetURL.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if t.userAgent != "" {
		req.Header.Set("User-Agent", t.userAgent)
	}
	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}
	if opts.RequestID != "" {
		req.Header.Set(requestIDHeader, opts.RequestID)
	}
	signedReq, err := t.auth.apply(req)
	if err != nil {
		return nil, err
	}
	return signedReq, nil
}

// Collection scopes the client to knowledge collection operations.
func (c *Client) Collection(meta kmodel.CollectionMeta) *CollectionClient {
	if c == nil || c.transport == nil {
		return nil
	}
	return &CollectionClient{
		transport: c.transport,
		meta:      meta,
	}
}

// Rerank performs rerank service call.
func (c *Client) Rerank(ctx context.Context, request kmodel.RerankRequest, opts ...RequestOption) (*kmodel.RerankResponse, error) {
	if len(request.Datas) == 0 {
		return nil, fmt.Errorf("datas list is empty")
	}
	if len(request.Datas) > 50 {
		return nil, fmt.Errorf("datas list too large")
	}
	var response kmodel.RerankResponse
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/service/rerank", request, &response, opts...); err != nil {
		return nil, err
	}
	return &response, nil
}

// ChatCompletion calls chat completions API (non-stream).
func (c *Client) ChatCompletion(ctx context.Context, request kmodel.ChatCompletionRequest, opts ...RequestOption) (*kmodel.ChatCompletionResponse, error) {
	if request.Stream != nil && *request.Stream {
		return nil, fmt.Errorf("stream=true: use ChatCompletionStream for streaming responses")
	}
	var response kmodel.ChatCompletionResponse
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/chat/completions", request, &response, opts...); err != nil {
		return nil, err
	}
	return &response, nil
}

// ServiceChat calls service chat API (non-stream).
func (c *Client) ServiceChat(ctx context.Context, request kmodel.ServiceChatRequest, opts ...RequestOption) (*kmodel.ServiceChatResponse, error) {
	if request.Stream != nil && *request.Stream {
		return nil, fmt.Errorf("stream=true: use ServiceChatStream for streaming responses")
	}
	var response kmodel.ServiceChatResponse
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/service/chat", request, &response, opts...); err != nil {
		return nil, err
	}
	return &response, nil
}

// ChatCompletionStream streams chat completion responses when request.Stream is true.
func (c *Client) ChatCompletionStream(ctx context.Context, request kmodel.ChatCompletionRequest, opts ...RequestOption) (<-chan kmodel.ChatCompletionResponse, error) {
	enabled := request.Stream != nil && *request.Stream
	if !enabled {
		flag := true
		request.Stream = &flag
	}
	requestOpts := defaultRequestOptions()
	for _, opt := range opts {
		opt(requestOpts)
	}
	body, err := vectorutils.SerializeToJSON(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	req, err := c.transport.buildRequest(ctx, http.MethodPost, "/api/knowledge/chat/completions", body, requestOpts)
	if err != nil {
		return nil, err
	}
	resp, err := vectorutils.DoHTTPRequest(c.transport.httpClient, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		defer resp.Body.Close()
		raw, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("unexpected %d response and failed to read body: %w", resp.StatusCode, readErr)
		}
		var errResp struct {
			Code      interface{} `json:"code"`
			Message   string      `json:"message"`
			RequestID string      `json:"request_id"`
		}
		_ = vectorutils.ParseJSONUseNumber(raw, &errResp)
		return nil, fmt.Errorf("chat completion stream error: code=%v, message=%s, request_id=%s, status=%d", errResp.Code, errResp.Message, errResp.RequestID, resp.StatusCode)
	}
	out := make(chan kmodel.ChatCompletionResponse)
	go func() {
		defer resp.Body.Close()
		defer close(out)

		scanner := bufio.NewScanner(resp.Body)
		const maxCapacity = 4 * 1024 * 1024
		buf := make([]byte, maxCapacity)
		scanner.Buffer(buf, maxCapacity)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			if !strings.HasPrefix(line, "data:") {
				continue
			}
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if data == "" || data == "[DONE]" {
				return
			}

			var item kmodel.ChatCompletionResponse
			if err := vectorutils.ParseJSONUseNumber([]byte(data), &item); err != nil {
				continue
			}

			select {
			case <-ctx.Done():
				return
			case out <- item:
			}
		}
	}()
	return out, nil
}

// ServiceChatStream streams service chat responses when request.Stream is true.
func (c *Client) ServiceChatStream(ctx context.Context, request kmodel.ServiceChatRequest, opts ...RequestOption) (<-chan kmodel.ServiceChatResponse, error) {
	enabled := request.Stream != nil && *request.Stream
	if !enabled {
		flag := true
		request.Stream = &flag
	}
	requestOpts := defaultRequestOptions()
	for _, opt := range opts {
		opt(requestOpts)
	}
	body, err := vectorutils.SerializeToJSON(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	req, err := c.transport.buildRequest(ctx, http.MethodPost, "/api/knowledge/service/chat", body, requestOpts)
	if err != nil {
		return nil, err
	}
	resp, err := vectorutils.DoHTTPRequest(c.transport.httpClient, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		defer resp.Body.Close()
		raw, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("unexpected %d response and failed to read body: %w", resp.StatusCode, readErr)
		}
		var errResp struct {
			Code      interface{} `json:"code"`
			Message   string      `json:"message"`
			RequestID string      `json:"request_id"`
		}
		_ = vectorutils.ParseJSONUseNumber(raw, &errResp)
		return nil, fmt.Errorf("service chat stream error: code=%v, message=%s, request_id=%s, status=%d", errResp.Code, errResp.Message, errResp.RequestID, resp.StatusCode)
	}
	out := make(chan kmodel.ServiceChatResponse)
	go func() {
		defer resp.Body.Close()
		defer close(out)

		scanner := bufio.NewScanner(resp.Body)
		const maxCapacity = 4 * 1024 * 1024
		buf := make([]byte, maxCapacity)
		scanner.Buffer(buf, maxCapacity)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			if !strings.HasPrefix(line, "data:") {
				continue
			}
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if data == "" || data == "[DONE]" {
				return
			}

			var item kmodel.ServiceChatResponse
			if err := vectorutils.ParseJSONUseNumber([]byte(data), &item); err != nil {
				continue
			}

			select {
			case <-ctx.Done():
				return
			case out <- item:
			}

			if item.Data != nil && item.Data.End != nil && *item.Data.End {
				return
			}
		}
	}()
	return out, nil
}
