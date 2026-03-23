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
	DenseDim     int                `json:"dense_dim,omitempty"`
	DenseVector  []float32          `json:"dense_vector,omitempty"`
	SparseVector map[string]float32 `json:"sparse_vector,omitempty"`
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

	ReturnSchema         *bool `json:"return_schema,omitempty"`
	ReturnDownloadURL    *bool `json:"return_download_url,omitempty"`
	ReturnAnalyzedResult *bool `json:"return_analyzed_result,omitempty"`
	ReturnDetailInfo     *bool `json:"return_detail_info,omitempty"`
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

	Schema              MapStr  `json:"schema,omitempty"`
	RerankError         *string `json:"rerank_error,omitempty"`
	InstForDense        *string `json:"instruction_for_dense,omitempty"`
	InstForSparse       *string `json:"instruction_for_sparse,omitempty"`
	InstForTensor       *string `json:"instruction_for_tensor,omitempty"`
	EmbeddingTimeCostMS *int    `json:"embedding_time_cost_ms,omitempty"`
	RecallTimeCostMS    *int    `json:"recall_time_cost_ms,omitempty"`
	RerankTimeCostMS    *int    `json:"rerank_time_cost_ms,omitempty"`
}

// SearchItemResult represents a single hit within a search response.
type SearchItemResult struct {
	ID       interface{} `json:"id"`
	Fields   MapStr      `json:"fields,omitempty"`
	ANNScore float32     `json:"ann_score,omitempty"`
	Score    float32     `json:"score,omitempty"`

	OriginScore   *float32 `json:"origin_score,omitempty"`
	AdditionScore *float32 `json:"addition_score,omitempty"`
}

// TensorRerank defines tensor-based rerank settings.
type TensorRerank struct {
	Tensor            *[][]float32 `json:"tensor,omitempty"`
	InputLimit        *int         `json:"input_limit,omitempty"`
	MaxSimilarityAlgo *string      `json:"max_similarity_algo,omitempty"`
}

// ModelRerank defines model-based rerank settings.
type ModelRerank struct {
	ModelName      *string  `json:"model_name,omitempty"`
	ModelVersion   *string  `json:"model_version,omitempty"`
	Instruction    *string  `json:"instruction,omitempty"`
	InputLimit     *int     `json:"input_limit,omitempty"`
	ScoreThreshold *float64 `json:"score_threshold,omitempty"`
	FailStrategy   *string  `json:"fail_strategy,omitempty"`
	TimeoutMs      *int     `json:"timeout_ms,omitempty"`
}

// Instruction configures auto-fill instruction generation.
type Instruction struct {
	AutoFill *bool `json:"auto_fill,omitempty"`
}

// SearchByVectorRequest performs vector similarity search.
type SearchByVectorRequest struct {
	SearchBase
	DenseVector  []float64          `json:"dense_vector"`
	SparseVector map[string]float64 `json:"sparse_vector,omitempty"`
	TensorRerank *TensorRerank      `json:"tensor_rerank,omitempty"`
}

// SearchByMultiModalRequest performs multimodal search.
type SearchByMultiModalRequest struct {
	SearchBase
	Text            *string       `json:"text,omitempty"`
	Image           interface{}   `json:"image,omitempty"`
	Video           interface{}   `json:"video,omitempty"`
	NeedInstruction *bool         `json:"need_instruction,omitempty"`
	Instruction     *Instruction  `json:"instruction,omitempty"`
	TensorRerank    *TensorRerank `json:"tensor_rerank,omitempty"`
	ModelRerank     *ModelRerank  `json:"rerank,omitempty"`
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
	Mode          string   `json:"mode,omitempty"`
	Keywords      []string `json:"keywords,omitempty"`
	Query         string   `json:"query,omitempty"`
	CaseSensitive bool     `json:"case_sensitive,omitempty"`
	Fields        []string `json:"fields,omitempty"`
	BM25K1        *float64 `json:"bm25_k1,omitempty"`
	BM25B         *float64 `json:"bm25_b,omitempty"`
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
