// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package model

type RerankDataItem struct {
	Query   interface{} `json:"query"`
	Content *string     `json:"content"`
	Title   *string     `json:"title,omitempty"`
	Image   interface{} `json:"image,omitempty"` // string or []string
}

type RerankRequest struct {
	Datas       []RerankDataItem `json:"datas,required"`
	EndpointId  *string          `json:"endpoint_id,omitempty"`
	RerankModel *string          `json:"rerank_model,omitempty"`
	// 重排指令（seed rerank 使用）
	RerankInstruction *string `json:"rerank_instruction,omitempty"`
}

type RerankResult struct {
	Scores     []float64 `json:"scores"`
	TokenUsage *int      `json:"token_usage,omitempty"`
}

type RerankResponse struct {
	Code      int           `json:"code,omitempty"`
	Message   string        `json:"message,omitempty"`
	RequestID string        `json:"request_id,omitempty"`
	Data      *RerankResult `json:"data,omitempty"`
}
