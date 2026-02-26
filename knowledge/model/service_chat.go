// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package model

type ServiceChatRequest struct {
	ServiceResourceID string                 `json:"service_resource_id"`
	Messages          []ChatMessage          `json:"messages"`
	QueryParam        map[string]interface{} `json:"query_param,omitempty"`
	Stream            *bool                  `json:"stream,omitempty"`
}

type ServiceChatRetrieveItemDocInfo struct {
	DocID      *string `json:"doc_id,omitempty"`
	DocName    *string `json:"doc_name,omitempty"`
	CreateTime *int64  `json:"create_time,omitempty"`
	DocType    *string `json:"doc_type,omitempty"`
	DocMeta    *string `json:"doc_meta,omitempty"`
	Source     *string `json:"source,omitempty"`
	Title      *string `json:"title,omitempty"`
}

type ServiceChatRetrieveItem struct {
	ID                 *string                         `json:"id,omitempty"`
	Content            *string                         `json:"content,omitempty"`
	MDContent          *string                         `json:"md_content,omitempty"`
	Score              *float64                        `json:"score,omitempty"`
	PointID            *string                         `json:"point_id,omitempty"`
	OriginText         *string                         `json:"origin_text,omitempty"`
	OriginalQuestion   *string                         `json:"original_question,omitempty"`
	ChunkTitle         *string                         `json:"chunk_title,omitempty"`
	ChunkID            *int64                          `json:"chunk_id,omitempty"`
	ProcessTime        *int64                          `json:"process_time,omitempty"`
	RerankScore        *float64                        `json:"rerank_score,omitempty"`
	DocInfo            *ServiceChatRetrieveItemDocInfo `json:"doc_info,omitempty"`
	RecallPosition     *int                            `json:"recall_position,omitempty"`
	RerankPosition     *int                            `json:"rerank_position,omitempty"`
	ChunkType          *string                         `json:"chunk_type,omitempty"`
	ChunkSource        *string                         `json:"chunk_source,omitempty"`
	UpdateTime         *int64                          `json:"update_time,omitempty"`
	ChunkAttachment    []ChunkAttachment               `json:"chunk_attachment,omitempty"`
	TableChunkFields   []PointTableChunkField          `json:"table_chunk_fields,omitempty"`
	OriginalCoordinate map[string]interface{}          `json:"original_coordinate,omitempty"`
}

type ServiceChatData struct {
	Count            *int                      `json:"count,omitempty"`
	RewriteQuery     *string                   `json:"rewrite_query,omitempty"`
	TokenUsage       interface{}               `json:"token_usage,omitempty"`
	ResultList       []ServiceChatRetrieveItem `json:"result_list,omitempty"`
	GeneratedAnswer  *string                   `json:"generated_answer,omitempty"`
	ReasoningContent *string                   `json:"reasoning_content,omitempty"`
	Prompt           *string                   `json:"prompt,omitempty"`
	End              *bool                     `json:"end,omitempty"`
}

type ServiceChatResponse struct {
	Code      int              `json:"code,omitempty"`
	Message   string           `json:"message,omitempty"`
	RequestID string           `json:"request_id,omitempty"`
	Data      *ServiceChatData `json:"data,omitempty"`
}
