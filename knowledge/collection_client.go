// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package knowledge

import (
	"context"
	"net/http"

	vectorutils "github.com/volcengine/vikingdb-go-sdk/vector/utils"

	kmodel "github.com/volcengine/vikingdb-go-sdk/knowledge/model"
)

// CollectionClient provides collection-scoped operations for KnowledgeBase.
type CollectionClient struct {
	transport *transport
	meta      kmodel.CollectionMeta
}

func (c *CollectionClient) metaPayload() map[string]interface{} {
	payload := map[string]interface{}{}
	if c.meta.CollectionName != "" {
		payload["collection_name"] = c.meta.CollectionName
	}
	if c.meta.ProjectName != "" {
		payload["project"] = c.meta.ProjectName
	}
	if c.meta.ResourceID != "" {
		payload["resource_id"] = c.meta.ResourceID
	}
	return payload
}

// AddDoc adds a document (deprecated, use AddDocV2).
func (c *CollectionClient) AddDoc(ctx context.Context, request kmodel.AddDocRequest, opts ...RequestOption) (*kmodel.AddDocResponse, error) {
	meta := c.metaPayload()
	payload := map[string]interface{}{}
	for k, v := range meta {
		payload[k] = v
	}
	// merge request struct into payload
	reqBytes, err := vectorutils.SerializeToJSON(request)
	if err != nil {
		return nil, err
	}
	var reqMap map[string]interface{}
	if err := vectorutils.ParseJSONUseNumber(reqBytes, &reqMap); err != nil {
		return nil, err
	}
	for k, v := range reqMap {
		payload[k] = v
	}
	var response kmodel.AddDocResponse
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/doc/add", payload, &response, opts...); err != nil {
		return nil, err
	}
	return &response, nil
}

// AddDocV2 adds a document using v2 API.
func (c *CollectionClient) AddDocV2(ctx context.Context, request kmodel.AddDocV2Request, opts ...RequestOption) (*kmodel.AddDocResponse, error) {
	meta := c.metaPayload()
	payload := map[string]interface{}{}
	for k, v := range meta {
		payload[k] = v
	}
	reqBytes, err := vectorutils.SerializeToJSON(request)
	if err != nil {
		return nil, err
	}
	var reqMap map[string]interface{}
	if err := vectorutils.ParseJSONUseNumber(reqBytes, &reqMap); err != nil {
		return nil, err
	}
	for k, v := range reqMap {
		payload[k] = v
	}
	var response kmodel.AddDocResponse
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/doc/v2/add", payload, &response, opts...); err != nil {
		return nil, err
	}
	return &response, nil
}

// DeleteDoc deletes a document by id.
func (c *CollectionClient) DeleteDoc(ctx context.Context, docID string, opts ...RequestOption) (*kmodel.CommonResponse, error) {
	payload := c.metaPayload()
	payload["doc_id"] = docID
	var resp kmodel.CommonResponse
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/doc/delete", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDoc retrieves document info.
func (c *CollectionClient) GetDoc(ctx context.Context, docID string, returnTokenUsage bool, opts ...RequestOption) (*kmodel.DocInfo, error) {
	payload := c.metaPayload()
	payload["doc_id"] = docID
	if returnTokenUsage {
		payload["return_token_usage"] = true
	}
	var raw struct {
		Code      int                    `json:"code,omitempty"`
		Message   string                 `json:"message,omitempty"`
		RequestID string                 `json:"request_id,omitempty"`
		Data      map[string]interface{} `json:"data,omitempty"`
	}
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/doc/info", payload, &raw, opts...); err != nil {
		return nil, err
	}
	docBytes, err := vectorutils.SerializeToJSON(raw.Data)
	if err != nil {
		return nil, err
	}
	var doc kmodel.DocInfo
	if err := vectorutils.ParseJSONUseNumber(docBytes, &doc); err != nil {
		return nil, err
	}
	if doc.Project == nil && c.meta.ProjectName != "" {
		p := c.meta.ProjectName
		doc.Project = &p
	}
	if doc.ResourceID == nil && c.meta.ResourceID != "" {
		r := c.meta.ResourceID
		doc.ResourceID = &r
	}
	return &doc, nil
}

// ListDocs lists documents in collection.
func (c *CollectionClient) ListDocs(ctx context.Context, request kmodel.ListDocsRequest, opts ...RequestOption) (*kmodel.ListDocsResponse, error) {
	meta := c.metaPayload()
	payload := map[string]interface{}{}
	for k, v := range meta {
		payload[k] = v
	}
	reqBytes, err := vectorutils.SerializeToJSON(request)
	if err != nil {
		return nil, err
	}
	var reqMap map[string]interface{}
	if err := vectorutils.ParseJSONUseNumber(reqBytes, &reqMap); err != nil {
		return nil, err
	}
	for k, v := range reqMap {
		payload[k] = v
	}
	var response kmodel.ListDocsResponse
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/doc/list", payload, &response, opts...); err != nil {
		return nil, err
	}
	return &response, nil
}

// ListDocsV2 lists documents in collection with cursor-based pagination.
// This SDK surface intentionally omits pipeline_name and only queries the main version.
func (c *CollectionClient) ListDocsV2(ctx context.Context, request kmodel.ListDocsV2Request, opts ...RequestOption) (*kmodel.ListDocsV2Response, error) {
	meta := c.metaPayload()
	payload := map[string]interface{}{}
	for k, v := range meta {
		payload[k] = v
	}
	reqBytes, err := vectorutils.SerializeToJSON(request)
	if err != nil {
		return nil, err
	}
	var reqMap map[string]interface{}
	if err := vectorutils.ParseJSONUseNumber(reqBytes, &reqMap); err != nil {
		return nil, err
	}
	for k, v := range reqMap {
		payload[k] = v
	}
	var response kmodel.ListDocsV2Response
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/doc/v2/list", payload, &response, opts...); err != nil {
		return nil, err
	}
	return &response, nil
}

// SearchDocsByFilter searches documents by tag filter on the main version only.
func (c *CollectionClient) SearchDocsByFilter(ctx context.Context, request kmodel.SearchDocsByFilterRequest, opts ...RequestOption) (*kmodel.SearchDocsByFilterResponse, error) {
	meta := c.metaPayload()
	payload := map[string]interface{}{}
	for k, v := range meta {
		payload[k] = v
	}
	reqBytes, err := vectorutils.SerializeToJSON(request)
	if err != nil {
		return nil, err
	}
	var reqMap map[string]interface{}
	if err := vectorutils.ParseJSONUseNumber(reqBytes, &reqMap); err != nil {
		return nil, err
	}
	for k, v := range reqMap {
		payload[k] = v
	}
	var response kmodel.SearchDocsByFilterResponse
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/doc/v2/search_by_filter", payload, &response, opts...); err != nil {
		return nil, err
	}
	return &response, nil
}

// UpdateDocMeta updates document metadata.
func (c *CollectionClient) UpdateDocMeta(ctx context.Context, docID string, meta []kmodel.MetaItem, opts ...RequestOption) (*kmodel.CommonResponse, error) {
	payload := c.metaPayload()
	payload["doc_id"] = docID
	payload["meta"] = meta
	var resp kmodel.CommonResponse
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/doc/update_meta", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateDoc updates document name.
func (c *CollectionClient) UpdateDoc(ctx context.Context, docID, docName string, opts ...RequestOption) (*kmodel.CommonResponse, error) {
	payload := c.metaPayload()
	payload["doc_id"] = docID
	payload["doc_name"] = docName
	var resp kmodel.CommonResponse
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/doc/update", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPoint retrieves point info.
func (c *CollectionClient) GetPoint(ctx context.Context, pointID string, getAttachmentLink bool, opts ...RequestOption) (*kmodel.PointInfo, error) {
	payload := c.metaPayload()
	payload["point_id"] = pointID
	if getAttachmentLink {
		payload["get_attachment_link"] = true
	}
	var raw struct {
		Code      int                    `json:"code,omitempty"`
		Message   string                 `json:"message,omitempty"`
		RequestID string                 `json:"request_id,omitempty"`
		Data      map[string]interface{} `json:"data,omitempty"`
	}
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/point/info", payload, &raw, opts...); err != nil {
		return nil, err
	}
	pointBytes, err := vectorutils.SerializeToJSON(raw.Data)
	if err != nil {
		return nil, err
	}
	var point kmodel.PointInfo
	if err := vectorutils.ParseJSONUseNumber(pointBytes, &point); err != nil {
		return nil, err
	}
	if point.Project == nil && c.meta.ProjectName != "" {
		p := c.meta.ProjectName
		point.Project = &p
	}
	if point.ResourceID == nil && c.meta.ResourceID != "" {
		r := c.meta.ResourceID
		point.ResourceID = &r
	}
	return &point, nil
}

// ListPoints lists points for the collection.
func (c *CollectionClient) ListPoints(ctx context.Context, request kmodel.ListPointsRequest, opts ...RequestOption) (*kmodel.ListPointsResponse, error) {
	meta := c.metaPayload()
	payload := map[string]interface{}{}
	for k, v := range meta {
		payload[k] = v
	}
	reqBytes, err := vectorutils.SerializeToJSON(request)
	if err != nil {
		return nil, err
	}
	var reqMap map[string]interface{}
	if err := vectorutils.ParseJSONUseNumber(reqBytes, &reqMap); err != nil {
		return nil, err
	}
	for k, v := range reqMap {
		payload[k] = v
	}
	var response kmodel.ListPointsResponse
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/point/list", payload, &response, opts...); err != nil {
		return nil, err
	}
	return &response, nil
}

// AddPoint adds a point to a document.
func (c *CollectionClient) AddPoint(ctx context.Context, request kmodel.AddPointRequest, opts ...RequestOption) (*kmodel.PointAddResponse, error) {
	meta := c.metaPayload()
	payload := map[string]interface{}{}
	for k, v := range meta {
		payload[k] = v
	}
	reqBytes, err := vectorutils.SerializeToJSON(request)
	if err != nil {
		return nil, err
	}
	var reqMap map[string]interface{}
	if err := vectorutils.ParseJSONUseNumber(reqBytes, &reqMap); err != nil {
		return nil, err
	}
	for k, v := range reqMap {
		payload[k] = v
	}
	var response kmodel.PointAddResponse
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/point/add", payload, &response, opts...); err != nil {
		return nil, err
	}
	return &response, nil
}

// UpdatePoint updates an existing point.
func (c *CollectionClient) UpdatePoint(ctx context.Context, pointID string, update kmodel.UpdatePointRequest, opts ...RequestOption) (*kmodel.CommonResponse, error) {
	payload := c.metaPayload()
	payload["point_id"] = pointID
	payload["chunk_title"] = update.ChunkTitle
	payload["content"] = update.Content
	payload["question"] = update.Question
	if len(update.Fields) > 0 {
		payload["fields"] = update.Fields
	}
	var resp kmodel.CommonResponse
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/point/update", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeletePoint deletes a point.
func (c *CollectionClient) DeletePoint(ctx context.Context, request kmodel.DeletePointRequest, opts ...RequestOption) (*kmodel.CommonResponse, error) {
	meta := c.metaPayload()
	payload := map[string]interface{}{}
	for k, v := range meta {
		payload[k] = v
	}
	payload["point_id"] = request.PointID
	var resp kmodel.CommonResponse
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/point/delete", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SearchCollection performs collection search.
func (c *CollectionClient) SearchCollection(ctx context.Context, request kmodel.SearchCollectionRequest, opts ...RequestOption) (*kmodel.SearchResponse, error) {
	meta := c.metaPayload()
	payload := map[string]interface{}{}
	for k, v := range meta {
		payload[k] = v
	}
	payload["name"] = c.meta.CollectionName
	reqBytes, err := vectorutils.SerializeToJSON(request)
	if err != nil {
		return nil, err
	}
	var reqMap map[string]interface{}
	if err := vectorutils.ParseJSONUseNumber(reqBytes, &reqMap); err != nil {
		return nil, err
	}
	for k, v := range reqMap {
		payload[k] = v
	}
	var response kmodel.SearchResponse
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/collection/search", payload, &response, opts...); err != nil {
		return nil, err
	}
	return &response, nil
}

// SearchKnowledge performs knowledge search.
func (c *CollectionClient) SearchKnowledge(ctx context.Context, request kmodel.SearchKnowledgeRequest, opts ...RequestOption) (*kmodel.SearchKnowledgeResponse, error) {
	meta := c.metaPayload()
	payload := map[string]interface{}{}
	for k, v := range meta {
		payload[k] = v
	}
	payload["name"] = c.meta.CollectionName
	reqBytes, err := vectorutils.SerializeToJSON(request)
	if err != nil {
		return nil, err
	}
	var reqMap map[string]interface{}
	if err := vectorutils.ParseJSONUseNumber(reqBytes, &reqMap); err != nil {
		return nil, err
	}
	for k, v := range reqMap {
		payload[k] = v
	}
	var response kmodel.SearchKnowledgeResponse
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/knowledge/collection/search_knowledge", payload, &response, opts...); err != nil {
		return nil, err
	}
	return &response, nil
}
