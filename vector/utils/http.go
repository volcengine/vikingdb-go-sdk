// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"io"
	"net/http"

	"github.com/volcengine/vikingdb-go-sdk/vector/model"
)

// DoHTTPRequest executes the HTTP request and wraps transport errors in an SDK error.
func DoHTTPRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, model.NewErrorWithCause(model.ErrCodeHTTPRequestFailed, "failed to execute http request", err, http.StatusServiceUnavailable)
	}
	return resp, nil
}

// ParseResponse reads the HTTP response body, decoding JSON into result when provided.
func ParseResponse(resp *http.Response, result interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.NewErrorWithCause(model.ErrCodeUnknown, "failed to read response body", err, http.StatusInternalServerError)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		var errResp struct {
			Code      string `json:"code"`
			Message   string `json:"message"`
			RequestID string `json:"request_id"`
		}
		if parseErr := ParseJSONUseNumber(body, &errResp); parseErr == nil && (errResp.Code != "" || errResp.Message != "") {
			return model.NewErrorWithRequestID(model.ErrorCode(errResp.Code), errResp.Message, errResp.RequestID, resp.StatusCode)
		}
		return model.NewErrorWithCause(model.ErrCodeUnknown, fmt.Sprintf("unexpected %d response: %s", resp.StatusCode, string(body)), err, resp.StatusCode)
	}

	if result == nil || len(body) == 0 {
		return nil
	}

	if err := ParseJSONUseNumber(body, result); err != nil {
		return model.NewErrorWithCause(model.ErrCodeUnknown, "failed to unmarshal response body", err, resp.StatusCode)
	}

	return nil
}
