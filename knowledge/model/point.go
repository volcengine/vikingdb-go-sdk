// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package model

type AddPointRequest struct {
	DocID      string                   `json:"doc_id"`
	ChunkType  string                   `json:"chunk_type"`
	ChunkTitle *string                  `json:"chunk_title,omitempty"`
	Content    *string                  `json:"content,omitempty"`
	Question   *string                  `json:"question,omitempty"`
	Fields     []map[string]interface{} `json:"fields,omitempty"`
}

type PointAddResult struct {
	CollectionName *string `json:"collection_name,omitempty"`
	Project        *string `json:"project,omitempty"`
	ResourceID     *string `json:"resource_id,omitempty"`
	DocID          *string `json:"doc_id,omitempty"`
	ChunkID        *int64  `json:"chunk_id,omitempty"`
	PointID        *string `json:"point_id,omitempty"`
}

type PointAddResponse struct {
	Code      int             `json:"code,omitempty"`
	Message   string          `json:"message,omitempty"`
	RequestID string          `json:"request_id,omitempty"`
	Data      *PointAddResult `json:"data,omitempty"`
}

type UpdatePointRequest struct {
	ChunkTitle *string                  `json:"chunk_title,omitempty"`
	Content    *string                  `json:"content,omitempty"`
	Question   *string                  `json:"question,omitempty"`
	Fields     []map[string]interface{} `json:"fields,omitempty"`
}

type ListPointsRequest struct {
	Offset            int                    `json:"offset"`
	Limit             int                    `json:"limit"`
	DocIDs            map[string]interface{} `json:"doc_ids,omitempty"`
	PointIDs          map[string]interface{} `json:"point_ids,omitempty"`
	GetAttachmentLink *bool                  `json:"get_attachment_link,omitempty"`
}

type DeletePointRequest struct {
	PointID string `json:"point_id"`
}

type PointDocInfo struct {
	DocID      *string    `json:"doc_id,omitempty"`
	DocName    *string    `json:"doc_name,omitempty"`
	CreateTime *int64     `json:"create_time,omitempty"`
	DocType    *string    `json:"doc_type,omitempty"`
	DocMeta    *string    `json:"doc_meta,omitempty"`
	Source     *string    `json:"source,omitempty"`
	Title      *string    `json:"title,omitempty"`
	Status     *DocStatus `json:"status,omitempty"`
}

type ChunkAttachment struct {
	UUID       *string `json:"uuid,omitempty"`
	Caption    *string `json:"caption,omitempty"`
	Type       *string `json:"type,omitempty"`
	Link       *string `json:"link,omitempty"`
	InfoLink   *string `json:"info_link,omitempty"`
	ColumnName *string `json:"column_name,omitempty"`
}

type PointTableChunkField struct {
	FieldName  *string     `json:"field_name,omitempty"`
	FieldValue interface{} `json:"field_value,omitempty"`
}

type PointInfo struct {
	CollectionName   *string                `json:"collection_name,omitempty"`
	PointID          *string                `json:"point_id,omitempty"`
	ProcessTime      *int64                 `json:"process_time,omitempty"`
	OriginText       *string                `json:"origin_text,omitempty"`
	MDContent        *string                `json:"md_content,omitempty"`
	HTMLContent      *string                `json:"html_content,omitempty"`
	ChunkTitle       *string                `json:"chunk_title,omitempty"`
	ChunkType        *string                `json:"chunk_type,omitempty"`
	Content          *string                `json:"content,omitempty"`
	ChunkID          *int64                 `json:"chunk_id,omitempty"`
	OriginalQuestion *string                `json:"original_question,omitempty"`
	DocInfo          *PointDocInfo          `json:"doc_info,omitempty"`
	RerankScore      *float64               `json:"rerank_score,omitempty"`
	Score            *float64               `json:"score,omitempty"`
	ChunkSource      *string                `json:"chunk_source,omitempty"`
	ChunkAttachment  []ChunkAttachment      `json:"chunk_attachment,omitempty"`
	TableChunkFields []PointTableChunkField `json:"table_chunk_fields,omitempty"`
	UpdateTime       *int64                 `json:"update_time,omitempty"`
	Description      *string                `json:"description,omitempty"`
	ChunkStatus      *string                `json:"chunk_status,omitempty"`
	VideoFrame       *string                `json:"video_frame,omitempty"`
	VideoURL         *string                `json:"video_url,omitempty"`
	VideoStartTime   *int64                 `json:"video_start_time,omitempty"`
	VideoEndTime     *int64                 `json:"video_end_time,omitempty"`
	VideoOutline     map[string]interface{} `json:"video_outline,omitempty"`
	AudioStartTime   *int64                 `json:"audio_start_time,omitempty"`
	AudioEndTime     *int64                 `json:"audio_end_time,omitempty"`
	AudioOutline     map[string]interface{} `json:"audio_outline,omitempty"`
	SheetName        *string                `json:"sheet_name,omitempty"`
	Project          *string                `json:"project,omitempty"`
	ResourceID       *string                `json:"resource_id,omitempty"`
}

type ListPointsResult struct {
	PointList []PointInfo `json:"point_list"`
}

type ListPointsResponse struct {
	Code      int               `json:"code,omitempty"`
	Message   string            `json:"message,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
	Data      *ListPointsResult `json:"data,omitempty"`
	Count     *int              `json:"count,omitempty"`
	TotalNum  *int              `json:"total_num,omitempty"`
}
