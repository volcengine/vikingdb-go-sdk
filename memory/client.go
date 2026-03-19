// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package memory

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	vectorutils "github.com/volcengine/vikingdb-go-sdk/vector/utils"
	"github.com/volcengine/volc-sdk-golang/base"

	mmodel "github.com/volcengine/vikingdb-go-sdk/memory/model"
)

const requestIDHeader = "X-Tt-Logid"

type authKind int

const (
	authKindNone authKind = iota
	authKindIAM
	authKindAPIKey
)

type Auth struct {
	kind      authKind
	accessKey string
	secretKey string
	apiKey    string
}

func AuthNone() Auth {
	return Auth{kind: authKindNone}
}

func AuthIAM(accessKey, secretKey string) Auth {
	return Auth{
		kind:      authKindIAM,
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

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

type Client struct {
	transport *transport
}

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

func (c *Client) Ping(ctx context.Context, opts ...RequestOption) error {
	if c == nil || c.transport == nil {
		return fmt.Errorf("client is nil")
	}
	return c.transport.doRequest(ctx, http.MethodGet, "/api/memory/ping", nil, nil, opts...)
}

func (c *Client) Collection(meta mmodel.CollectionMeta) (*CollectionClient, error) {
	if c == nil || c.transport == nil {
		return nil, fmt.Errorf("client is nil")
	}
	if meta.ResourceID == "" && meta.CollectionName == "" {
		return nil, fmt.Errorf("either CollectionName or ResourceID must be provided")
	}
	if meta.ProjectName == "" && meta.CollectionName != "" && meta.ResourceID == "" {
		meta.ProjectName = "default"
	}
	return &CollectionClient{
		transport: c.transport,
		meta:      meta,
	}, nil
}

func (c *Client) GetCollection(collectionName, projectName string) (*CollectionClient, error) {
	meta := mmodel.CollectionMeta{
		CollectionName: collectionName,
		ProjectName:    projectName,
	}
	return c.Collection(meta)
}

func (c *Client) GetCollectionByResourceID(resourceID string) (*CollectionClient, error) {
	meta := mmodel.CollectionMeta{
		ResourceID: resourceID,
	}
	return c.Collection(meta)
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
			return &mmodel.VikingMemException{
				Code:       1000028,
				Message:    "failed to marshal request",
				StatusCode: http.StatusBadRequest,
				Err:        err,
			}
		}
		body = serialized
	}
	return vectorutils.Retry(retries, func() error {
		req, err := t.buildRequest(ctx, method, path, body, requestOpts)
		if err != nil {
			return err
		}
		resp, err := t.httpClient.Do(req)
		if err != nil {
			return mmodel.PromoteException(&mmodel.VikingMemException{
				Code:       1000028,
				Message:    "failed to execute http request",
				StatusCode: http.StatusServiceUnavailable,
				Err:        err,
			})
		}
		defer resp.Body.Close()
		return parseResponse(resp, response)
	}, isRetryableError)
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
	var buf io.Reader
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

func parseResponse(resp *http.Response, result interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mmodel.PromoteException(&mmodel.VikingMemException{
			Code:       1000028,
			Message:    "failed to read response body",
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		})
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return mmodel.PromoteException(parseErrorEnvelope(body, resp.StatusCode, err))
	}

	if result == nil || len(body) == 0 {
		return nil
	}

	if err := vectorutils.ParseJSONUseNumber(body, result); err != nil {
		return mmodel.PromoteException(&mmodel.VikingMemException{
			Code:       1000028,
			Message:    "failed to unmarshal response body",
			StatusCode: resp.StatusCode,
			Err:        err,
		})
	}
	return nil
}

func parseErrorEnvelope(body []byte, statusCode int, readErr error) *mmodel.VikingMemException {
	var env struct {
		Code      interface{} `json:"code"`
		Message   string      `json:"message"`
		RequestID string      `json:"request_id"`
	}
	if err := vectorutils.ParseJSONUseNumber(body, &env); err != nil {
		return &mmodel.VikingMemException{
			Code:       1000028,
			Message:    fmt.Sprintf("unexpected %d response: %s", statusCode, string(body)),
			StatusCode: statusCode,
			Err:        err,
		}
	}
	code := parseErrorCode(env.Code)
	message := env.Message
	if message == "" {
		message = string(body)
	}
	if code == 0 {
		code = 1000028
	}
	return &mmodel.VikingMemException{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		RequestID:  env.RequestID,
		Err:        readErr,
	}
}

func parseErrorCode(raw interface{}) int {
	switch v := raw.(type) {
	case nil:
		return 0
	case json.Number:
		i, err := v.Int64()
		if err != nil {
			return 0
		}
		return int(i)
	case float64:
		return int(v)
	case float32:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	case string:
		if v == "" {
			return 0
		}
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		return i
	default:
		return 0
	}
}

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	var memErr *mmodel.VikingMemException
	if errors.As(err, &memErr) {
		switch memErr.StatusCode {
		case http.StatusTooManyRequests,
			http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			return true
		}
		if memErr.Code == 1000029 {
			return true
		}
		return false
	}
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return true
	}
	return false
}
