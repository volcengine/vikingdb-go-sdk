// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package model

import (
	"fmt"
)

type CollectionMeta struct {
	CollectionName string `json:"collection_name,omitempty"`
	ProjectName    string `json:"project_name,omitempty"`
	ResourceID     string `json:"resource_id,omitempty"`
}

type Response struct {
	Code      int         `json:"code,omitempty"`
	Message   string      `json:"message,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type VikingMemException struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code,omitempty"`
	RequestID  string `json:"request_id,omitempty"`
	Err        error  `json:"-"`
}

func (e *VikingMemException) Error() string {
	if e == nil {
		return "vikingdb memory error: <nil>"
	}
	if e.RequestID != "" {
		return fmt.Sprintf("vikingdb memory error: code=%d, message=%s, status_code=%d, err=%v, request_id=%s", e.Code, e.Message, e.StatusCode, e.Err, e.RequestID)
	}
	return fmt.Sprintf("vikingdb memory error: code=%d, message=%s, status_code=%d, err=%v", e.Code, e.Message, e.StatusCode, e.Err)
}

func (e *VikingMemException) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

type UnauthorizedException struct{ *VikingMemException }
type InvalidRequestException struct{ *VikingMemException }
type CollectionExistException struct{ *VikingMemException }
type CollectionNotExistException struct{ *VikingMemException }
type IndexExistException struct{ *VikingMemException }
type IndexNotExistException struct{ *VikingMemException }
type DataNotFoundException struct{ *VikingMemException }
type DelOpFailedException struct{ *VikingMemException }
type UpsertOpFailedException struct{ *VikingMemException }
type InvalidVectorException struct{ *VikingMemException }
type InvalidPrimaryKeyException struct{ *VikingMemException }
type InvalidFilterException struct{ *VikingMemException }
type IndexSearchException struct{ *VikingMemException }
type IndexFetchException struct{ *VikingMemException }
type IndexInitializingException struct{ *VikingMemException }
type EmbeddingException struct{ *VikingMemException }
type InternalServerException struct{ *VikingMemException }
type QuotaLimiterException struct{ *VikingMemException }

func (e *UnauthorizedException) Unwrap() error           { return e.VikingMemException }
func (e *InvalidRequestException) Unwrap() error         { return e.VikingMemException }
func (e *CollectionExistException) Unwrap() error        { return e.VikingMemException }
func (e *CollectionNotExistException) Unwrap() error     { return e.VikingMemException }
func (e *IndexExistException) Unwrap() error             { return e.VikingMemException }
func (e *IndexNotExistException) Unwrap() error          { return e.VikingMemException }
func (e *DataNotFoundException) Unwrap() error           { return e.VikingMemException }
func (e *DelOpFailedException) Unwrap() error            { return e.VikingMemException }
func (e *UpsertOpFailedException) Unwrap() error         { return e.VikingMemException }
func (e *InvalidVectorException) Unwrap() error          { return e.VikingMemException }
func (e *InvalidPrimaryKeyException) Unwrap() error      { return e.VikingMemException }
func (e *InvalidFilterException) Unwrap() error          { return e.VikingMemException }
func (e *IndexSearchException) Unwrap() error            { return e.VikingMemException }
func (e *IndexFetchException) Unwrap() error             { return e.VikingMemException }
func (e *IndexInitializingException) Unwrap() error      { return e.VikingMemException }
func (e *EmbeddingException) Unwrap() error              { return e.VikingMemException }
func (e *InternalServerException) Unwrap() error         { return e.VikingMemException }
func (e *QuotaLimiterException) Unwrap() error           { return e.VikingMemException }

var ExceptionMap = map[int]func(*VikingMemException) error{
	1000001: func(e *VikingMemException) error { return &UnauthorizedException{e} },
	1000003: func(e *VikingMemException) error { return &InvalidRequestException{e} },
	1000004: func(e *VikingMemException) error { return &CollectionExistException{e} },
	1000005: func(e *VikingMemException) error { return &CollectionNotExistException{e} },
	1000007: func(e *VikingMemException) error { return &IndexExistException{e} },
	1000008: func(e *VikingMemException) error { return &IndexNotExistException{e} },
	1000011: func(e *VikingMemException) error { return &DataNotFoundException{e} },
	1000013: func(e *VikingMemException) error { return &DelOpFailedException{e} },
	1000014: func(e *VikingMemException) error { return &UpsertOpFailedException{e} },
	1000016: func(e *VikingMemException) error { return &InvalidVectorException{e} },
	1000017: func(e *VikingMemException) error { return &InvalidPrimaryKeyException{e} },
	1000019: func(e *VikingMemException) error { return &InvalidFilterException{e} },
	1000021: func(e *VikingMemException) error { return &IndexSearchException{e} },
	1000022: func(e *VikingMemException) error { return &IndexFetchException{e} },
	1000023: func(e *VikingMemException) error { return &IndexInitializingException{e} },
	1000025: func(e *VikingMemException) error { return &EmbeddingException{e} },
	1000028: func(e *VikingMemException) error { return &InternalServerException{e} },
	1000029: func(e *VikingMemException) error { return &QuotaLimiterException{e} },
}

func PromoteException(err *VikingMemException) error {
	if err == nil {
		return nil
	}
	if ctor, ok := ExceptionMap[err.Code]; ok {
		return ctor(err)
	}
	return err
}

type AddEventRequest struct {
	EventType      string      `json:"event_type"`
	MemoryInfo     interface{} `json:"memory_info"`
	UserID         string      `json:"user_id,omitempty"`
	AssistantID    string      `json:"assistant_id,omitempty"`
	GroupID        string      `json:"group_id,omitempty"`
	UpdateProfiles interface{} `json:"update_profiles,omitempty"`
}

type UpdateEventRequest struct {
	EventID      string      `json:"event_id"`
	MemoryInfo   interface{} `json:"memory_info"`
	UserID       string      `json:"user_id,omitempty"`
	AssistantID  string      `json:"assistant_id,omitempty"`
}

type DeleteEventRequest struct {
	EventID string `json:"event_id"`
}

type BatchDeleteEventRequest struct {
	Filter     interface{} `json:"filter,omitempty"`
	DeleteType string      `json:"delete_type,omitempty"`
}

type AddProfileRequest struct {
	ProfileType  string      `json:"profile_type"`
	MemoryInfo   interface{} `json:"memory_info"`
	IsUpsert     bool        `json:"is_upsert"`
	UserID       string      `json:"user_id,omitempty"`
	AssistantID  string      `json:"assistant_id,omitempty"`
	GroupID      string      `json:"group_id,omitempty"`
}

type UpdateProfileRequest struct {
	ProfileID   string      `json:"profile_id"`
	MemoryInfo  interface{} `json:"memory_info"`
}

type DeleteProfileRequest struct {
	ProfileID string `json:"profile_id"`
}

type BatchDeleteProfileRequest struct {
	Filter interface{} `json:"filter,omitempty"`
}

type TriggerUpdateProfileRequest struct {
	UpdateProfileType interface{} `json:"update_profile_type,omitempty"`
	Filters           interface{} `json:"filters,omitempty"`
}

type AddSessionRequest struct {
	SessionID  string      `json:"session_id"`
	Messages   interface{} `json:"messages"`
	Metadata   interface{} `json:"metadata,omitempty"`
	Profiles   interface{} `json:"profiles,omitempty"`
	StoreFile  interface{} `json:"store_file,omitempty"`
}

type GetSessionInfoRequest struct {
	SessionID string `json:"session_id"`
}

type SearchMemoryRequest struct {
	Query  string      `json:"query,omitempty"`
	Filter interface{} `json:"filter,omitempty"`
	Limit  int         `json:"limit,omitempty"`
}

type SearchEventMemoryRequest struct {
	Query           string      `json:"query,omitempty"`
	Filter          interface{} `json:"filter,omitempty"`
	TimeDecayConfig interface{} `json:"time_decay_config,omitempty"`
	Limit           int         `json:"limit,omitempty"`
}

type SearchProfileMemoryRequest struct {
	Query  string      `json:"query,omitempty"`
	Filter interface{} `json:"filter,omitempty"`
	Limit  int         `json:"limit,omitempty"`
}
