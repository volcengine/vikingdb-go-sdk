// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package model

// FetchDataInIndexRequest fetches documents (and optional vectors) from an index.
type FetchDataInIndexRequest struct {
	IDs          []interface{} `json:"ids"`
	Partition    string        `json:"partition,omitempty"` // advanced feature, support string&int partition
	OutputFields []string      `json:"output_fields,omitempty"`
}

type IndexDataItem struct {
	DataItem
	DenseDim    int       `json:"dense_dim,omitempty"`
	DenseVector []float32 `json:"dense_vector,omitempty"`
}

// FetchDataInIndexResponse mirrors DataApiResponse<FetchDataInIndexResult>.
type FetchDataInIndexResponse struct {
	CommonResponse
	Result *FetchDataInIndexResult `json:"result,omitempty"`
}

type FetchDataInIndexResult struct {
	Items       []IndexDataItem `json:"fetch,omitempty"`
	NotFoundIDs []interface{}   `json:"ids_not_exist,omitempty"`
}

// RecallBase carries shared search filters.
type RecallBase struct {
	Filter    MapStr `json:"filter,omitempty"`
	Partition string `json:"partition,omitempty"` // advanced feature, support string&int partition
}

// SearchBase enriches recall filters with pagination/output hints.
type SearchBase struct {
	RecallBase

	OutputFields []string       `json:"output_fields,omitempty"`
	Limit        *int           `json:"limit,omitempty"`
	Offset       *int           `json:"offset,omitempty"`
	Advance      *SearchAdvance `json:"advance,omitempty"`
}

// SearchAdvance maps to Java's SearchAdvance DTO.
type SearchAdvance struct {
	DenseWeight           *float64      `json:"dense_weight,omitempty"`
	IDsIn                 []interface{} `json:"ids_in,omitempty"`
	IDsNotIn              []interface{} `json:"ids_not_in,omitempty"`
	PostProcessOps        []MapStr      `json:"post_process_ops,omitempty"`
	PostProcessInputLimit *int          `json:"post_process_input_limit,omitempty"`
	ScaleK                *float64      `json:"scale_k,omitempty"`
	FilterPreAnnLimit     *int          `json:"filter_pre_ann_limit,omitempty"`
	FilterPreAnnRatio     *float64      `json:"filter_pre_ann_ratio,omitempty"`
}

type SearchResponse struct {
	CommonResponse
	Result *SearchResult `json:"result,omitempty"`
}

type SearchResult struct {
	Data               []SearchItemResult `json:"data,omitempty"`
	FilterMatchedCount int                `json:"filter_matched_count,omitempty"`
	TotalReturnCount   int                `json:"total_return_count,omitempty"`
	RealTextQuery      string             `json:"real_text_query,omitempty"`
	TokenUsage         MapStr             `json:"token_usage,omitempty"`
}

// SearchItemResult represents a single hit within a search response.
type SearchItemResult struct {
	ID       interface{} `json:"id"`
	Fields   MapStr      `json:"fields,omitempty"`
	ANNScore float32     `json:"ann_score,omitempty"`
	Score    float32     `json:"score,omitempty"`
}

// SearchByVectorRequest performs vector similarity search.
type SearchByVectorRequest struct {
	SearchBase
	DenseVector  []float64          `json:"dense_vector"`
	SparseVector map[string]float64 `json:"sparse_vector,omitempty"`
}

// SearchByMultiModalRequest performs multimodal search.
type SearchByMultiModalRequest struct {
	SearchBase
	Text            *string     `json:"text,omitempty"`
	Image           interface{} `json:"image,omitempty"`
	Video           interface{} `json:"video,omitempty"`
	NeedInstruction *bool       `json:"need_instruction,omitempty"`
}

// SearchByIDRequest looks up a document by primary key.
type SearchByIDRequest struct {
	SearchBase
	ID interface{} `json:"id"`
}

// ScalarOrder represents the sort direction for scalar search.
type ScalarOrder string

const (
	ScalarOrderAsc  ScalarOrder = "asc"
	ScalarOrderDesc ScalarOrder = "desc"
)

// SearchByScalarRequest orders results by scalar field.
type SearchByScalarRequest struct {
	SearchBase
	Field *string     `json:"field,omitempty"`
	Order ScalarOrder `json:"order,omitempty"`
}

// SearchByKeywordsRequest matches documents by keywords.
type SearchByKeywordsRequest struct {
	SearchBase
	Keywords      []string `json:"keywords,omitempty"`
	Query         string   `json:"query,omitempty"`
	CaseSensitive bool     `json:"case_sensitive,omitempty"`
}

// SearchByRandomRequest randomly samples documents.
type SearchByRandomRequest struct {
	SearchBase
}

// AggRequest performs aggregations on search results.
type AggRequest struct {
	RecallBase
	Op    string      `json:"op"`
	Field *string     `json:"field,omitempty"`
	Cond  MapStr      `json:"cond,omitempty"`
	Order ScalarOrder `json:"order,omitempty"`
}

type AggResponse struct {
	CommonResponse
	Result *AggResult `json:"result,omitempty"`
}

type AggResult struct {
	Agg   MapStr `json:"agg,omitempty"`
	Op    string `json:"op,omitempty"`
	Field string `json:"field,omitempty"`
}
