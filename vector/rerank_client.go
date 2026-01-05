// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package vector

import (
	"context"
	"net/http"

	"github.com/volcengine/vikingdb-go-sdk/vector/model"
)

type rerankClient struct {
	client *transport
}

func (r *rerankClient) Rerank(ctx context.Context, request model.RerankRequest, opts ...RequestOption) (*model.RerankResponse, error) {
	response := &model.RerankResponse{}
	err := r.client.doRequest(ctx, http.MethodPost, "/api/vikingdb/rerank", request, response, opts...)
	return response, err
}
