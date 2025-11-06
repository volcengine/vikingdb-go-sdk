// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package vector

// RequestOptions captures per-request overrides for retries, headers, and query params.
type RequestOptions struct {
	MaxRetries int
	Headers    map[string]string
	Query      map[string]string
	RequestID  string
}

// RequestOption mutates RequestOptions when constructing a request.
type RequestOption func(*RequestOptions)

// defaultRequestOptions builds a RequestOptions instance with empty header and query maps.
func defaultRequestOptions() *RequestOptions {
	return &RequestOptions{
		Headers: make(map[string]string),
		Query:   make(map[string]string),
	}
}

// WithRequestMaxRetries limits the retry count for the current request.
func WithRequestMaxRetries(maxRetries int) RequestOption {
	return func(o *RequestOptions) {
		o.MaxRetries = maxRetries
	}
}

// WithRequestHeader sets a single header value for the request.
func WithRequestHeader(key, value string) RequestOption {
	return func(o *RequestOptions) {
		o.Headers[key] = value
	}
}

// WithRequestHeaders merges the provided headers into the request.
func WithRequestHeaders(headers map[string]string) RequestOption {
	return func(o *RequestOptions) {
		if len(headers) == 0 {
			return
		}
		for k, v := range headers {
			o.Headers[k] = v
		}
	}
}

// WithRequestQueryParam adds a single query parameter to the request.
func WithRequestQueryParam(key, value string) RequestOption {
	return func(o *RequestOptions) {
		o.Query[key] = value
	}
}

// WithRequestQueryParams merges the provided query parameters into the request.
func WithRequestQueryParams(params map[string]string) RequestOption {
	return func(o *RequestOptions) {
		if len(params) == 0 {
			return
		}
		for k, v := range params {
			o.Query[k] = v
		}
	}
}

// WithRequestID sets the request id which will be propagated as X-Tt-Logid.
func WithRequestID(requestID string) RequestOption {
	return func(o *RequestOptions) {
		o.RequestID = requestID
	}
}
