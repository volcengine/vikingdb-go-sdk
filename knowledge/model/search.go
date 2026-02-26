// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package model

type SearchCollectionRequest struct {
	Query           interface{}            `json:"query"`
	Limit           int                    `json:"limit"`
	DenseWeight     float64                `json:"dense_weight"`
	RerankSwitch    bool                   `json:"rerank_switch"`
	QueryParam      map[string]interface{} `json:"query_param,omitempty"`
	RetrieveCount   *int                   `json:"retrieve_count,omitempty"`
	EndpointID      *string                `json:"endpoint_id,omitempty"`
	RerankModel     *string                `json:"rerank_model,omitempty"`
	RerankOnlyChunk bool                   `json:"rerank_only_chunk"`
}

type SearchResult struct {
	ResultList []PointInfo `json:"result_list"`
}

type SearchResponse struct {
	Code      int           `json:"code,omitempty"`
	Message   string        `json:"message,omitempty"`
	RequestID string        `json:"request_id,omitempty"`
	Data      *SearchResult `json:"data,omitempty"`
}

type SearchKnowledgeRequest struct {
	Query          interface{}            `json:"query"`
	ImageQuery     interface{}            `json:"image_query,omitempty"`
	PreProcessing  map[string]interface{} `json:"pre_processing,omitempty"`
	PostProcessing map[string]interface{} `json:"post_processing,omitempty"`
	QueryParam     map[string]interface{} `json:"query_param,omitempty"`
	Limit          int                    `json:"limit"`
	DenseWeight    float64                `json:"dense_weight"`
}

type SearchKnowledgeResult struct {
	Count        *int                   `json:"count,omitempty"`
	RewriteQuery *string                `json:"rewrite_query,omitempty"`
	TokenUsage   map[string]interface{} `json:"token_usage,omitempty"`
	ResultList   []PointInfo            `json:"result_list"`
}

type SearchKnowledgeResponse struct {
	Code      int                    `json:"code,omitempty"`
	Message   string                 `json:"message,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	Data      *SearchKnowledgeResult `json:"data,omitempty"`
}
