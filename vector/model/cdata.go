// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package model

// DataItem represents a document stored in the collection.
type DataItem struct {
	ID     interface{} `json:"id"`
	Fields MapStr      `json:"fields"`
}

// WriteDataBase holds common fields for data writes.
type WriteDataBase struct {
	Data                []MapStr `json:"data"`
	TTL                 *int32   `json:"ttl,omitempty"`
	IgnoreUnknownFields bool     `json:"ignore_unknown_fields,omitempty"`
}

// UpsertDataRequest creates or updates documents within a collection.
type UpsertDataRequest struct {
	WriteDataBase
	Async bool `json:"async,omitempty"`
}

// UpsertDataResponse mirrors the Java DataApiResponse<UpsertDataResult>.
type UpsertDataResponse struct {
	CommonResponse
	Result *UpsertDataResult `json:"result,omitempty"`
}

type UpsertDataResult struct {
	TokenUsage interface{} `json:"token_usage,omitempty"`
}

// UpdateDataRequest updates existing documents.
type UpdateDataRequest struct {
	WriteDataBase
}

type UpdateDataResponse struct {
	CommonResponse
	Result *UpdateDataResult `json:"result,omitempty"`
}

type UpdateDataResult struct {
	TokenUsage interface{} `json:"token_usage,omitempty"`
}

// DeleteDataRequest removes documents by primary key.
type DeleteDataRequest struct {
	IDs    []interface{} `json:"ids"`
	DelAll bool          `json:"del_all,omitempty"`
}

type DeleteDataResponse struct {
	CommonResponse
}

// FetchDataInCollectionRequest fetches documents by primary key from a collection.
type FetchDataInCollectionRequest struct {
	IDs []interface{} `json:"ids"`
}

// FetchDataInCollectionResponse returns fetched documents and missing IDs.
type FetchDataInCollectionResponse struct {
	CommonResponse
	Result *FetchDataInCollectionResult `json:"result,omitempty"`
}

type FetchDataInCollectionResult struct {
	Items       []DataItem    `json:"fetch,omitempty"`
	NotFoundIDs []interface{} `json:"ids_not_exist,omitempty"`
}
