// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package vector

import (
	"context"

	"github.com/volcengine/vikingdb-go-sdk/vector/model"
)

// CollectionClient provides collection-scoped data operations.
type CollectionClient interface {
	Upsert(ctx context.Context, request model.UpsertDataRequest, opts ...RequestOption) (*model.UpsertDataResponse, error)
	Update(ctx context.Context, request model.UpdateDataRequest, opts ...RequestOption) (*model.UpdateDataResponse, error)
	Delete(ctx context.Context, request model.DeleteDataRequest, opts ...RequestOption) (*model.DeleteDataResponse, error)
	Fetch(ctx context.Context, request model.FetchDataInCollectionRequest, opts ...RequestOption) (*model.FetchDataInCollectionResponse, error)

	CollectionName() string
	ResourceID() string
	ProjectName() string
}

// IndexClient provides index-level search and metadata operations.
type IndexClient interface {
	Fetch(ctx context.Context, request model.FetchDataInIndexRequest, opts ...RequestOption) (*model.FetchDataInIndexResponse, error)
	SearchByVector(ctx context.Context, request model.SearchByVectorRequest, opts ...RequestOption) (*model.SearchResponse, error)
	SearchByMultiModal(ctx context.Context, request model.SearchByMultiModalRequest, opts ...RequestOption) (*model.SearchResponse, error)
	SearchByID(ctx context.Context, request model.SearchByIDRequest, opts ...RequestOption) (*model.SearchResponse, error)
	SearchByScalar(ctx context.Context, request model.SearchByScalarRequest, opts ...RequestOption) (*model.SearchResponse, error)
	SearchByKeywords(ctx context.Context, request model.SearchByKeywordsRequest, opts ...RequestOption) (*model.SearchResponse, error)
	SearchByRandom(ctx context.Context, request model.SearchByRandomRequest, opts ...RequestOption) (*model.SearchResponse, error)
	Aggregate(ctx context.Context, request model.AggRequest, opts ...RequestOption) (*model.AggResponse, error)

	CollectionName() string
	IndexName() string
	ResourceID() string
	ProjectName() string
}

// EmbeddingClient provides embedding operations.
type EmbeddingClient interface {
	Embedding(ctx context.Context, request model.EmbeddingRequest, opts ...RequestOption) (*model.EmbeddingResponse, error)
}
