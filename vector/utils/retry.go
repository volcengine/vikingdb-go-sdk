// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"math/rand"
	"time"

	"github.com/volcengine/vikingdb-go-sdk/vector/model"
)

const (
	defaultInitialBackoff = 100 * time.Millisecond
	defaultMaxBackoff     = 10 * time.Second
	backoffMultiplier     = 2.0
)

// Retry executes fn with exponential backoff. Retries stop when fn returns nil, the max retry count is reached,
// or shouldRetry returns false for the latest error.
func Retry(maxRetries int, fn func() error, shouldRetry func(error) bool) error {
	if maxRetries < 0 {
		maxRetries = 0
	}
	var lastErr error
	delay := defaultInitialBackoff

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			jitter := time.Duration(rand.Int63n(int64(delay)))
			sleepFor := delay + jitter
			if sleepFor > defaultMaxBackoff {
				sleepFor = defaultMaxBackoff
			}
			time.Sleep(sleepFor)
			next := time.Duration(float64(delay) * backoffMultiplier)
			if next > defaultMaxBackoff {
				next = defaultMaxBackoff
			}
			delay = next
		}

		if err := fn(); err != nil {
			lastErr = err
			if shouldRetry != nil && !shouldRetry(err) {
				return err
			}
			if attempt == maxRetries {
				return err
			}
			continue
		}
		return nil
	}

	return lastErr
}

// IsRetryableError delegates to model.IsRetryableError for backward compatibility.
func IsRetryableError(err error) bool {
	return model.IsRetryableError(err)
}
