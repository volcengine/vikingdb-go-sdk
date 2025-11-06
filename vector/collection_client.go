// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package vector

import (
	"context"
	"net/http"

	"github.com/volcengine/vikingdb-go-sdk/vector/model"
)

type collectionClient struct {
	client         *transport
	collectionBase model.CollectionLocator
}

func (c *collectionClient) Upsert(ctx context.Context, request model.UpsertDataRequest, opts ...RequestOption) (*model.UpsertDataResponse, error) {
	response := &model.UpsertDataResponse{}
	req := struct {
		model.CollectionLocator
		model.UpsertDataRequest
	}{
		CollectionLocator: c.collectionBase,
		UpsertDataRequest: request,
	}
	err := c.client.doRequest(ctx, http.MethodPost, "/api/vikingdb/data/upsert", req, response, opts...)
	return response, err
}

func (c *collectionClient) Update(ctx context.Context, request model.UpdateDataRequest, opts ...RequestOption) (*model.UpdateDataResponse, error) {
	response := &model.UpdateDataResponse{}
	req := struct {
		model.CollectionLocator
		model.UpdateDataRequest
	}{
		CollectionLocator: c.collectionBase,
		UpdateDataRequest: request,
	}
	err := c.client.doRequest(ctx, http.MethodPost, "/api/vikingdb/data/update", req, response, opts...)
	return response, err
}

func (c *collectionClient) Delete(ctx context.Context, request model.DeleteDataRequest, opts ...RequestOption) (*model.DeleteDataResponse, error) {
	response := &model.DeleteDataResponse{}
	req := struct {
		model.CollectionLocator
		model.DeleteDataRequest
	}{
		CollectionLocator: c.collectionBase,
		DeleteDataRequest: request,
	}
	err := c.client.doRequest(ctx, http.MethodPost, "/api/vikingdb/data/delete", req, response, opts...)
	return response, err
}

func (c *collectionClient) Fetch(ctx context.Context, request model.FetchDataInCollectionRequest, opts ...RequestOption) (*model.FetchDataInCollectionResponse, error) {
	response := &model.FetchDataInCollectionResponse{}
	req := struct {
		model.CollectionLocator
		model.FetchDataInCollectionRequest
	}{
		CollectionLocator:            c.collectionBase,
		FetchDataInCollectionRequest: request,
	}
	err := c.client.doRequest(ctx, http.MethodPost, "/api/vikingdb/data/fetch_in_collection", req, response, opts...)
	return response, err
}

func (c *collectionClient) CollectionName() string {
	return c.collectionBase.CollectionName
}

func (c *collectionClient) ProjectName() string {
	return c.collectionBase.ProjectName
}

func (c *collectionClient) ResourceID() string {
	return c.collectionBase.ResourceID
}
