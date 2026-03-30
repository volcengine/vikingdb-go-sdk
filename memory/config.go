// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package memory

import (
	"github.com/volcengine/vikingdb-go-sdk/common"
	"net/http"
	"time"
)

type Config struct {
	Endpoint   string
	Region     string
	Timeout    time.Duration
	MaxRetries int
	HTTPClient *http.Client
	UserAgent  string
}

func DefaultConfig() Config {
	return Config{
		Endpoint:   "https://api-knowledgebase.mlp.cn-beijing.volces.com",
		Region:     "cn-beijing",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		UserAgent:  "vikingdb-go-sdk-memory/" + common.Version,
	}
}

type ClientOption func(*Config)

func WithEndpoint(endpoint string) ClientOption {
	return func(c *Config) { c.Endpoint = endpoint }
}

func WithRegion(region string) ClientOption {
	return func(c *Config) { c.Region = region }
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Config) { c.Timeout = timeout }
}

func WithMaxRetries(maxRetries int) ClientOption {
	return func(c *Config) { c.MaxRetries = maxRetries }
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Config) { c.HTTPClient = httpClient }
}

func WithUserAgent(userAgent string) ClientOption {
	return func(c *Config) { c.UserAgent = userAgent }
}
