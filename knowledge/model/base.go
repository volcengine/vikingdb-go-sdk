// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package model

// CommonResponse represents the shared response envelope returned by Knowledge APIs.
// Note: code field varies in type across endpoints; exclude it here to avoid decoding issues.
type CommonResponse struct {
	API       string `json:"api,omitempty"`
	Message   string `json:"message,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// CollectionMeta describes knowledge collection locator.
type CollectionMeta struct {
	CollectionName string `json:"collection_name,omitempty"`
	ProjectName    string `json:"project,omitempty"`
	ResourceID     string `json:"resource_id,omitempty"`
}
