// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package memory

type RequestOptions struct {
	MaxRetries int
	Headers    map[string]string
	Query      map[string]string
	RequestID  string
}

type RequestOption func(*RequestOptions)

func defaultRequestOptions() *RequestOptions {
	return &RequestOptions{
		Headers: make(map[string]string),
		Query:   make(map[string]string),
	}
}

func WithRequestMaxRetries(maxRetries int) RequestOption {
	return func(o *RequestOptions) { o.MaxRetries = maxRetries }
}

func WithRequestHeader(key, value string) RequestOption {
	return func(o *RequestOptions) { o.Headers[key] = value }
}

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

func WithRequestQueryParam(key, value string) RequestOption {
	return func(o *RequestOptions) { o.Query[key] = value }
}

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

func WithRequestID(requestID string) RequestOption {
	return func(o *RequestOptions) { o.RequestID = requestID }
}
