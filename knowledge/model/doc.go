// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package model

type ListDocsRequest struct {
	Offset           int                    `json:"offset"`
	Limit            int                    `json:"limit"`
	DocType          *string                `json:"doc_type,omitempty"`
	Filter           map[string]interface{} `json:"filter,omitempty"`
	ReturnTokenUsage *bool                  `json:"return_token_usage,omitempty"`
}

type AddDocRequest struct {
	AddType     string        `json:"add_type"`
	DocID       *string       `json:"doc_id,omitempty"`
	DocName     *string       `json:"doc_name,omitempty"`
	DocType     *string       `json:"doc_type,omitempty"`
	Description *string       `json:"description,omitempty"`
	TOSPath     *string       `json:"tos_path,omitempty"`
	URL         *string       `json:"url,omitempty"`
	LarkFile    *string       `json:"lark_file,omitempty"`
	Meta        []MetaItem    `json:"meta,omitempty"`
	Dedup       *DedupOptions `json:"dedup,omitempty"`
}

type DedupOptions struct {
	ContentDedup *bool `json:"content_dedup,omitempty"`
	DocNameDedup *bool `json:"doc_name_dedup,omitempty"`
	AutoSkip     *bool `json:"auto_skip,omitempty"`
}

type MetaItem struct {
	FieldName  *string     `json:"field_name,omitempty"`
	FieldType  *string     `json:"field_type,omitempty"`
	FieldValue interface{} `json:"field_value,omitempty"`
}

type AddDocV2Request struct {
	DocID        *string    `json:"doc_id,omitempty"`
	DocName      *string    `json:"doc_name,omitempty"`
	DocType      *string    `json:"doc_type,omitempty"`
	Description  *string    `json:"description,omitempty"`
	TagList      []MetaItem `json:"tag_list,omitempty"`
	URI          *string    `json:"uri,omitempty"`
	PathSegments []string   `json:"path_segments,omitempty"`
}

type DedupInfo struct {
	Skip       *bool    `json:"skip,omitempty"`
	SameDocIDs []string `json:"same_doc_ids,omitempty"`
}

type DocStatus struct {
	ProcessStatus *int    `json:"process_status,omitempty"`
	FailedCode    *int    `json:"failed_code,omitempty"`
	FailedMsg     *string `json:"failed_msg,omitempty"`
}

type DocPremiumStatus struct {
	DocSummaryStatusCode *int `json:"doc_summary_status_code,omitempty"`
}

type DocInfo struct {
	CollectionName   *string                `json:"collection_name,omitempty"`
	DocName          *string                `json:"doc_name,omitempty"`
	DocID            *string                `json:"doc_id,omitempty"`
	DocHash          *string                `json:"doc_hash,omitempty"`
	AddType          *string                `json:"add_type,omitempty"`
	DocType          *string                `json:"doc_type,omitempty"`
	DocSummary       *string                `json:"doc_summary,omitempty"`
	BriefSummary     *string                `json:"brief_summary,omitempty"`
	DocSize          *int                   `json:"doc_size,omitempty"`
	Description      *string                `json:"description,omitempty"`
	CreateTime       *int64                 `json:"create_time,omitempty"`
	AddedBy          *string                `json:"added_by,omitempty"`
	UpdateTime       *int64                 `json:"update_time,omitempty"`
	URL              *string                `json:"url,omitempty"`
	TOSPath          *string                `json:"tos_path,omitempty"`
	PointNum         *int                   `json:"point_num,omitempty"`
	Status           *DocStatus             `json:"status,omitempty"`
	Title            *string                `json:"title,omitempty"`
	Source           *string                `json:"source,omitempty"`
	TotalTokens      interface{}            `json:"total_tokens,omitempty"`
	DocSummaryTokens *int                   `json:"doc_summary_tokens,omitempty"`
	DocPremiumStatus *DocPremiumStatus      `json:"doc_premium_status,omitempty"`
	Meta             *string                `json:"meta,omitempty"`
	Labels           map[string]string      `json:"labels,omitempty"`
	VideoOutline     map[string]interface{} `json:"video_outline,omitempty"`
	AudioOutline     map[string]interface{} `json:"audio_outline,omitempty"`
	Statistics       map[string]interface{} `json:"statistics,omitempty"`
	Project          *string                `json:"project,omitempty"`
	ResourceID       *string                `json:"resource_id,omitempty"`
}

type ListDocsResult struct {
	DocList  []DocInfo `json:"doc_list"`
	Count    *int      `json:"count,omitempty"`
	TotalNum *int      `json:"total_num,omitempty"`
}

type ListDocsResponse struct {
	Code      int             `json:"code,omitempty"`
	Message   string          `json:"message,omitempty"`
	RequestID string          `json:"request_id,omitempty"`
	Data      *ListDocsResult `json:"data,omitempty"`
}

type AddDocResponseData struct {
	CollectionName *string    `json:"collection_name,omitempty"`
	ResourceID     *string    `json:"resource_id,omitempty"`
	Project        *string    `json:"project,omitempty"`
	DocID          *string    `json:"doc_id,omitempty"`
	TaskID         *int64     `json:"task_id,omitempty"`
	DedupInfo      *DedupInfo `json:"dedup_info,omitempty"`
	MoreInfo       *string    `json:"more_info,omitempty"`
}

type AddDocResponse struct {
	Code      int                 `json:"code,omitempty"`
	Message   string              `json:"message,omitempty"`
	RequestID string              `json:"request_id,omitempty"`
	Data      *AddDocResponseData `json:"data,omitempty"`
}
