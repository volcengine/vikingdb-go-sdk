// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package vector

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/volcengine/vikingdb-go-sdk/vector/model"
	"github.com/volcengine/vikingdb-go-sdk/vector/utils"
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
		return nil, model.NewInvalidParameterError("access key and secret key cannot be empty")
	}
	return utils.SignRequestWithRegion(req, a.ak, a.sk, a.region), nil
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
		return nil, model.NewInvalidParameterError("endpoint cannot be empty")
	}

	defaults := DefaultConfig()

	baseURL, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return nil, model.NewErrorWithCause(model.ErrCodeInvalidParameter, "invalid endpoint", err, http.StatusBadRequest)
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
		userAgent = fmt.Sprintf("vikingdb-go-sdk/%s", Version)
	}

	var auth authenticator = noAuth{}
	switch authConfig.kind {
	case authKindIAM:
		if authConfig.accessKey == "" || authConfig.secretKey == "" {
			return nil, model.NewInvalidParameterError("access key and secret key cannot be empty")
		}
		auth = iamAuth{ak: authConfig.accessKey, sk: authConfig.secretKey, region: cfg.Region}
	case authKindAPIKey:
		if authConfig.apiKey == "" {
			return nil, model.NewInvalidParameterError("api key cannot be empty")
		}
		auth = apiKeyAuth{token: authConfig.apiKey}
	default:
		return nil, model.NewInvalidParameterError("no auth")
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

// Client represents the entry point for interacting with VikingDB services.
type Client struct {
	transport *transport
}

// New constructs a Client with the provided authentication settings and options.
func New(auth Auth, opts ...ClientOption) (*Client, error) {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	transport, err := newTransport(cfg, auth)
	if err != nil {
		return nil, err
	}

	return &Client{transport: transport}, nil
}

// Collection scopes the client to collection operations using the supplied locator metadata.
func (c *Client) Collection(base model.CollectionLocator) CollectionClient {
	if c == nil || c.transport == nil {
		return nil
	}
	return &collectionClient{
		client:         c.transport,
		collectionBase: base,
	}
}

// Index scopes the client to index operations using the supplied locator metadata.
func (c *Client) Index(base model.IndexLocator) IndexClient {
	if c == nil || c.transport == nil {
		return nil
	}
	return &indexClient{
		transport: c.transport,
		indexBase: base,
	}
}

// Embedding exposes embedding operations.
func (c *Client) Embedding() EmbeddingClient {
	if c == nil || c.transport == nil {
		return nil
	}
	return &embeddingClient{client: c.transport}
}

// Rerank exposes rerank operations.
func (c *Client) Rerank() RerankClient {
	if c == nil || c.transport == nil {
		return nil
	}
	return &rerankClient{client: c.transport}
}

func (c *transport) doRequest(ctx context.Context, method, path string, request, response interface{}, opts ...RequestOption) error {
	if ctx == nil {
		ctx = context.Background()
	}

	requestOpts := defaultRequestOptions()
	for _, opt := range opts {
		opt(requestOpts)
	}

	retries := requestOpts.MaxRetries
	if retries <= 0 {
		retries = c.config.MaxRetries
	}
	if retries < 0 {
		retries = 0
	}

	var body []byte
	if request != nil {
		serialized, err := utils.SerializeToJSON(request)
		if err != nil {
			return model.NewErrorWithCause(model.ErrCodeInvalidParameter, "failed to marshal request", err, http.StatusBadRequest)
		}
		body = serialized
	}

	return utils.Retry(retries, func() error {
		req, err := c.buildRequest(ctx, method, path, body, requestOpts)
		if err != nil {
			return err
		}

		resp, err := utils.DoHTTPRequest(c.httpClient, req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		return utils.ParseResponse(resp, response)
	}, utils.IsRetryableError)
}

func (c *transport) buildRequest(ctx context.Context, method, path string, body []byte, opts *RequestOptions) (*http.Request, error) {
	targetURL := c.baseURL.ResolveReference(&url.URL{Path: path})
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
		return nil, model.NewErrorWithCause(model.ErrCodeUnknown, "failed to create request", err, http.StatusBadRequest)
	}

	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}
	if opts.RequestID != "" {
		req.Header.Set(requestIDHeader, opts.RequestID)
	}

	signedReq, err := c.auth.apply(req)
	if err != nil {
		return nil, err
	}
	return signedReq, nil
}
