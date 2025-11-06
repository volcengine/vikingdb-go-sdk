// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package model

// CommonResponse represents the shared response envelope returned by VikingDB APIs.
type CommonResponse struct {
	API       string `json:"api,omitempty"`
	Message   string `json:"message,omitempty"`
	Code      string `json:"code,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// CollectionLocator carries general collection level identifiers.
type CollectionLocator struct {
	CollectionName string `json:"collection_name"`
	ProjectName    string `json:"project_name,omitempty"`
	ResourceID     string `json:"resource_id"`
}

// IndexLocator extends collection metadata with the index name.
type IndexLocator struct {
	CollectionLocator
	IndexName string `json:"index_name"`
}

type Refer struct {
	AccountID  string `json:"account_id"`
	InstanceNO string `json:"instance_no"`
	ResourceID string `json:"resource_id"`
}

type MapStr map[string]interface{}

// PaginationRequest represents pagination inputs.
type PaginationRequest struct {
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}

// PaginationResponse represents pagination metadata.
type PaginationResponse struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}
