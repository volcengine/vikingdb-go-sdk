// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package model

// EmbeddingModelOpt describes the model configuration for dense or sparse embeddings.
type EmbeddingModelOpt struct {
	ModelName    *string `json:"name"`
	ModelVersion *string `json:"version,omitempty"`
	Dim          *int    `json:"dim,omitempty"`
}

// FullModalData represents a single multimodal element that can be embedded.
type FullModalData struct {
	Text  *string     `json:"text,omitempty"`
	Image *string     `json:"image,omitempty"`
	Video interface{} `json:"video,omitempty"`
}

// EmbeddingData captures the payload to embed, supporting multimodal sequences.
type EmbeddingData struct {
	Text         *string         `json:"text,omitempty"`
	Image        interface{}     `json:"image,omitempty"`
	Video        interface{}     `json:"video,omitempty"`
	FullModalSeq []FullModalData `json:"full_modal_seq,omitempty"`
}

// EmbeddingRequest mirrors the Java SDK request payload.
type EmbeddingRequest struct {
	ProjectName *string            `json:"project_name,omitempty"`
	DenseModel  *EmbeddingModelOpt `json:"dense_model,omitempty"`
	SparseModel *EmbeddingModelOpt `json:"sparse_model,omitempty"`
	Data        []*EmbeddingData   `json:"data"`
}

type EmbeddingResponse struct {
	CommonResponse
	Result *EmbeddingResult `json:"result,omitempty"`
}

type EmbeddingResult struct {
	Data       []*Embedding `json:"data"`
	TokenUsage interface{}  `json:"token_usage,omitempty"`
}

// Embedding contains the generated dense and sparse vectors.
type Embedding struct {
	DenseVectors  []float32          `json:"dense,omitempty"`
	SparseVectors map[string]float32 `json:"sparse,omitempty"`
}
