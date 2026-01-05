// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package model

type RerankRequest struct {
	ModelName        string            `json:"model_name"`
	ModelVersion     string            `json:"model_version"`
	Data             [][]FullModalData `json:"data"`
	Query            []FullModalData   `json:"query"`
	Instruction      *string           `json:"instruction,omitempty"`
	ReturnOriginData *bool             `json:"return_origin_data,omitempty"`
}

type RerankResponse struct {
	CommonResponse
	Result *RerankResult `json:"result,omitempty"`
}

type RerankResult struct {
	Data       []RerankItem `json:"data"`
	TokenUsage interface{}  `json:"token_usage,omitempty"`
}

// RerankItem contains the id, score and origin data.
type RerankItem struct {
	ID         int64           `json:"id"`
	Score      float32         `json:"score"`
	OriginData []FullModalData `json:"origin_data,omitempty"`
}
