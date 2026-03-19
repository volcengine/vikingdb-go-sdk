// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package memory

import (
	"context"
	"fmt"
	"net/http"

	mmodel "github.com/volcengine/vikingdb-go-sdk/memory/model"
	vectorutils "github.com/volcengine/vikingdb-go-sdk/vector/utils"
)

type CollectionClient struct {
	transport *transport
	meta      mmodel.CollectionMeta
}

func (c *CollectionClient) metaPayload() map[string]interface{} {
	payload := map[string]interface{}{}
	if c.meta.CollectionName != "" {
		payload["collection_name"] = c.meta.CollectionName
	}
	if c.meta.ProjectName != "" {
		payload["project_name"] = c.meta.ProjectName
	}
	if c.meta.ResourceID != "" {
		payload["resource_id"] = c.meta.ResourceID
	}
	return payload
}

func (c *CollectionClient) mergePayload(request interface{}) (map[string]interface{}, error) {
	payload := c.metaPayload()
	if request == nil {
		return payload, nil
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
	return payload, nil
}

func (c *CollectionClient) AddEvent(ctx context.Context, request mmodel.AddEventRequest, opts ...RequestOption) (*mmodel.Response, error) {
	if request.EventType == "" {
		return nil, fmt.Errorf("event_type is required")
	}
	payload, err := c.mergePayload(request)
	if err != nil {
		return nil, err
	}
	var resp mmodel.Response
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/memory/event/add", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *CollectionClient) UpdateEvent(ctx context.Context, request mmodel.UpdateEventRequest, opts ...RequestOption) (*mmodel.Response, error) {
	if request.EventID == "" {
		return nil, fmt.Errorf("event_id is required")
	}
	payload, err := c.mergePayload(request)
	if err != nil {
		return nil, err
	}
	var resp mmodel.Response
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/memory/event/update", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *CollectionClient) DeleteEvent(ctx context.Context, request mmodel.DeleteEventRequest, opts ...RequestOption) (*mmodel.Response, error) {
	if request.EventID == "" {
		return nil, fmt.Errorf("event_id is required")
	}
	payload, err := c.mergePayload(request)
	if err != nil {
		return nil, err
	}
	var resp mmodel.Response
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/memory/event/delete", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *CollectionClient) BatchDeleteEvent(ctx context.Context, request mmodel.BatchDeleteEventRequest, opts ...RequestOption) (*mmodel.Response, error) {
	payload, err := c.mergePayload(request)
	if err != nil {
		return nil, err
	}
	var resp mmodel.Response
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/memory/event/batch_delete", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *CollectionClient) AddProfile(ctx context.Context, request mmodel.AddProfileRequest, opts ...RequestOption) (*mmodel.Response, error) {
	if request.ProfileType == "" {
		return nil, fmt.Errorf("profile_type is required")
	}
	payload, err := c.mergePayload(request)
	if err != nil {
		return nil, err
	}
	var resp mmodel.Response
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/memory/profile/add", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *CollectionClient) UpdateProfile(ctx context.Context, request mmodel.UpdateProfileRequest, opts ...RequestOption) (*mmodel.Response, error) {
	if request.ProfileID == "" {
		return nil, fmt.Errorf("profile_id is required")
	}
	payload, err := c.mergePayload(request)
	if err != nil {
		return nil, err
	}
	var resp mmodel.Response
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/memory/profile/update", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *CollectionClient) DeleteProfile(ctx context.Context, request mmodel.DeleteProfileRequest, opts ...RequestOption) (*mmodel.Response, error) {
	if request.ProfileID == "" {
		return nil, fmt.Errorf("profile_id is required")
	}
	payload, err := c.mergePayload(request)
	if err != nil {
		return nil, err
	}
	var resp mmodel.Response
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/memory/profile/delete", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *CollectionClient) BatchDeleteProfile(ctx context.Context, request mmodel.BatchDeleteProfileRequest, opts ...RequestOption) (*mmodel.Response, error) {
	payload, err := c.mergePayload(request)
	if err != nil {
		return nil, err
	}
	var resp mmodel.Response
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/memory/profile/batch_delete", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *CollectionClient) TriggerUpdateProfile(ctx context.Context, request mmodel.TriggerUpdateProfileRequest, opts ...RequestOption) (*mmodel.Response, error) {
	payload, err := c.mergePayload(request)
	if err != nil {
		return nil, err
	}
	var resp mmodel.Response
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/memory/profile/trigger_update", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *CollectionClient) AddSession(ctx context.Context, request mmodel.AddSessionRequest, opts ...RequestOption) (*mmodel.Response, error) {
	if request.SessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}
	payload, err := c.mergePayload(request)
	if err != nil {
		return nil, err
	}
	var resp mmodel.Response
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/memory/session/add", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *CollectionClient) GetSessionInfo(ctx context.Context, request mmodel.GetSessionInfoRequest, opts ...RequestOption) (*mmodel.Response, error) {
	if request.SessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}
	payload, err := c.mergePayload(request)
	if err != nil {
		return nil, err
	}
	var resp mmodel.Response
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/memory/session/info", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *CollectionClient) SearchMemory(ctx context.Context, request mmodel.SearchMemoryRequest, opts ...RequestOption) (*mmodel.Response, error) {
	payload, err := c.mergePayload(request)
	if err != nil {
		return nil, err
	}
	var resp mmodel.Response
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/memory/search", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *CollectionClient) SearchEventMemory(ctx context.Context, request mmodel.SearchEventMemoryRequest, opts ...RequestOption) (*mmodel.Response, error) {
	payload, err := c.mergePayload(request)
	if err != nil {
		return nil, err
	}
	var resp mmodel.Response
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/memory/event/search", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *CollectionClient) SearchProfileMemory(ctx context.Context, request mmodel.SearchProfileMemoryRequest, opts ...RequestOption) (*mmodel.Response, error) {
	payload, err := c.mergePayload(request)
	if err != nil {
		return nil, err
	}
	var resp mmodel.Response
	if err := c.transport.doRequest(ctx, http.MethodPost, "/api/memory/profile/search", payload, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}
