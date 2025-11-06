// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package vector

import (
	"context"
	"net/http"

	"github.com/volcengine/vikingdb-go-sdk/vector/model"
)

type indexClient struct {
	transport *transport
	indexBase model.IndexLocator
}

func (i *indexClient) Fetch(ctx context.Context, request model.FetchDataInIndexRequest, opts ...RequestOption) (*model.FetchDataInIndexResponse, error) {
	response := &model.FetchDataInIndexResponse{}
	req := struct {
		model.IndexLocator
		model.FetchDataInIndexRequest
	}{
		IndexLocator:            i.indexBase,
		FetchDataInIndexRequest: request,
	}
	err := i.transport.doRequest(ctx, http.MethodPost, "/api/vikingdb/data/fetch_in_index", req, response, opts...)
	return response, err
}

func (i *indexClient) SearchByVector(ctx context.Context, request model.SearchByVectorRequest, opts ...RequestOption) (*model.SearchResponse, error) {
	response := &model.SearchResponse{}
	req := struct {
		model.IndexLocator
		model.SearchByVectorRequest
	}{
		IndexLocator:          i.indexBase,
		SearchByVectorRequest: request,
	}
	err := i.transport.doRequest(ctx, http.MethodPost, "/api/vikingdb/data/search/vector", req, response, opts...)
	return response, err
}

func (i *indexClient) SearchByMultiModal(ctx context.Context, request model.SearchByMultiModalRequest, opts ...RequestOption) (*model.SearchResponse, error) {
	response := &model.SearchResponse{}
	req := struct {
		model.IndexLocator
		model.SearchByMultiModalRequest
	}{
		IndexLocator:              i.indexBase,
		SearchByMultiModalRequest: request,
	}
	err := i.transport.doRequest(ctx, http.MethodPost, "/api/vikingdb/data/search/multi_modal", req, response, opts...)
	return response, err
}

func (i *indexClient) SearchByID(ctx context.Context, request model.SearchByIDRequest, opts ...RequestOption) (*model.SearchResponse, error) {
	response := &model.SearchResponse{}
	req := struct {
		model.IndexLocator
		model.SearchByIDRequest
	}{
		IndexLocator:      i.indexBase,
		SearchByIDRequest: request,
	}
	err := i.transport.doRequest(ctx, http.MethodPost, "/api/vikingdb/data/search/id", req, response, opts...)
	return response, err
}

func (i *indexClient) SearchByScalar(ctx context.Context, request model.SearchByScalarRequest, opts ...RequestOption) (*model.SearchResponse, error) {
	response := &model.SearchResponse{}
	req := struct {
		model.IndexLocator
		model.SearchByScalarRequest
	}{
		IndexLocator:          i.indexBase,
		SearchByScalarRequest: request,
	}
	err := i.transport.doRequest(ctx, http.MethodPost, "/api/vikingdb/data/search/scalar", req, response, opts...)
	return response, err
}

func (i *indexClient) SearchByKeywords(ctx context.Context, request model.SearchByKeywordsRequest, opts ...RequestOption) (*model.SearchResponse, error) {
	response := &model.SearchResponse{}
	req := struct {
		model.IndexLocator
		model.SearchByKeywordsRequest
	}{
		IndexLocator:            i.indexBase,
		SearchByKeywordsRequest: request,
	}
	err := i.transport.doRequest(ctx, http.MethodPost, "/api/vikingdb/data/search/keywords", req, response, opts...)
	return response, err
}

func (i *indexClient) SearchByRandom(ctx context.Context, request model.SearchByRandomRequest, opts ...RequestOption) (*model.SearchResponse, error) {
	response := &model.SearchResponse{}
	req := struct {
		model.IndexLocator
		model.SearchByRandomRequest
	}{
		IndexLocator:          i.indexBase,
		SearchByRandomRequest: request,
	}
	err := i.transport.doRequest(ctx, http.MethodPost, "/api/vikingdb/data/search/random", req, response, opts...)
	return response, err
}

func (i *indexClient) Aggregate(ctx context.Context, request model.AggRequest, opts ...RequestOption) (*model.AggResponse, error) {
	response := &model.AggResponse{}
	req := struct {
		model.IndexLocator
		model.AggRequest
	}{
		IndexLocator: i.indexBase,
		AggRequest:   request,
	}
	err := i.transport.doRequest(ctx, http.MethodPost, "/api/vikingdb/data/agg", req, response, opts...)
	return response, err
}

func (i *indexClient) CollectionName() string {
	return i.indexBase.CollectionName
}

func (i *indexClient) IndexName() string {
	return i.indexBase.IndexName
}

func (i *indexClient) ResourceID() string {
	return i.indexBase.ResourceID
}

func (i *indexClient) ProjectName() string {
	return i.indexBase.ProjectName
}
