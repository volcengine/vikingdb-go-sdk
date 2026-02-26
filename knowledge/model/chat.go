// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package model

type ChatMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type ChatCompletionRequest struct {
	Model            string        `json:"model"`
	Messages         []ChatMessage `json:"messages"`
	ModelVersion     *string       `json:"model_version,omitempty"`
	Thinking         interface{}   `json:"thinking,omitempty"`
	MaxTokens        *int          `json:"max_tokens,omitempty"`
	Temperature      *float64      `json:"temperature,omitempty"`
	ReturnTokenUsage *bool         `json:"return_token_usage,omitempty"`
	APIKey           *string       `json:"api_key,omitempty"`
	Stream           *bool         `json:"stream,omitempty"`
}

type ChatCompletionResult struct {
	ReasoningContent *string     `json:"reasoning_content,omitempty"`
	GeneratedAnswer  *string     `json:"generated_answer,omitempty"`
	Usage            interface{} `json:"usage,omitempty"`
	Prompt           *string     `json:"prompt,omitempty"`
	Model            *string     `json:"model,omitempty"`
	FinishReason     *string     `json:"finish_reason,omitempty"`
	TotalTokens      interface{} `json:"total_tokens,omitempty"`
}

type ChatCompletionResponse struct {
	Code      int                   `json:"code,omitempty"`
	Message   string                `json:"message,omitempty"`
	RequestID string                `json:"request_id,omitempty"`
	Data      *ChatCompletionResult `json:"data,omitempty"`
}
