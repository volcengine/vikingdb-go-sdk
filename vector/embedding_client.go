// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package vector

import (
	"context"
	"net/http"

	"github.com/volcengine/vikingdb-go-sdk/vector/model"
)

type embeddingClient struct {
	client *transport
}

func (e *embeddingClient) Embedding(ctx context.Context, request model.EmbeddingRequest, opts ...RequestOption) (*model.EmbeddingResponse, error) {
	response := &model.EmbeddingResponse{}
	err := e.client.doRequest(ctx, http.MethodPost, "/api/vikingdb/embedding", request, response, opts...)
	return response, err
}
