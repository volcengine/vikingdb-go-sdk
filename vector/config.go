// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package vector

import (
	"net/http"
	"time"
)

// Version denotes the SDK version reported via the User-Agent header.
const Version = "0.1.0"

// Config carries shared settings for all clients.
type Config struct {
	Endpoint   string
	Region     string
	Timeout    time.Duration
	MaxRetries int
	HTTPClient *http.Client
	UserAgent  string
}

// DefaultConfig returns the baseline configuration.
func DefaultConfig() Config {
	return Config{
		Endpoint:   "https://api.vector.bytedance.com",
		Region:     "cn-beijing",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
	}
}

// ClientOption mutates the SDK configuration before client creation.
type ClientOption func(*Config)

func WithEndpoint(endpoint string) ClientOption {
	return func(c *Config) {
		c.Endpoint = endpoint
	}
}

func WithRegion(region string) ClientOption {
	return func(c *Config) {
		c.Region = region
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

func WithMaxRetries(maxRetries int) ClientOption {
	return func(c *Config) {
		c.MaxRetries = maxRetries
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Config) {
		c.HTTPClient = httpClient
	}
}

func WithUserAgent(userAgent string) ClientOption {
	return func(c *Config) {
		c.UserAgent = userAgent
	}
}
