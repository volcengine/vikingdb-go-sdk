// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package model

// CommonResponse represents the shared response envelope returned by Knowledge APIs.
type CommonResponse struct {
	Code      int32       `json:"code,omitempty"`
	Message   string      `json:"message,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

// CollectionMeta describes knowledge collection locator.
type CollectionMeta struct {
	CollectionName string `json:"collection_name,omitempty"`
	ProjectName    string `json:"project,omitempty"`
	ResourceID     string `json:"resource_id,omitempty"`
}
