// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/volcengine/vikingdb-go-sdk/vector"
	"github.com/volcengine/vikingdb-go-sdk/vector/model"
)

// TestSnippetEmbeddingMultiModalPipeline keeps the multimodal embedding call inline for quick copy/paste.
func TestSnippetEmbeddingMultiModalPipeline(t *testing.T) {
	client, err := vector.New(
		vector.AuthIAM(os.Getenv("VIKINGDB_AK"), os.Getenv("VIKINGDB_SK")),
		vector.WithEndpoint("https://"+os.Getenv("VIKINGDB_HOST")),
		vector.WithRegion(os.Getenv("VIKINGDB_REGION")),
	)
	if err != nil {
		t.Fatal(err)
	}

	embeddingClient := client.Embedding()

	ctx := context.Background()

	text := "Short multimodal prompt for the reference embedding pipeline."
	modelName := "doubao-embedding-vision"
	modelVersion := "250615"

	request := model.EmbeddingRequest{
		DenseModel: &model.EmbeddingModelOpt{
			ModelName:    stringPtr(modelName),
			ModelVersion: stringPtr(modelVersion),
		},
		Data: []*model.EmbeddingData{
			{
				FullModalSeq: []model.FullModalData{
					{Text: stringPtr(text)},
				},
			},
		},
	}

	resp, err := embeddingClient.Embedding(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Result == nil || len(resp.Result.Data) == 0 {
		t.Fatalf("embedding response missing data")
	}
	if len(resp.Result.Data[0].DenseVectors) == 0 {
		t.Fatalf("embedding response missing dense vector")
	}
	t.Logf("Embedding request_id=%s dense_dims=%d", resp.RequestID, len(resp.Result.Data[0].DenseVectors))
}

// TestSnippetEmbeddingDSPipeline shows dense+sparse embedding retrieval inline.
func TestSnippetEmbeddingDSPipeline(t *testing.T) {
	client, err := vector.New(
		vector.AuthIAM(os.Getenv("VIKINGDB_AK"), os.Getenv("VIKINGDB_SK")),
		vector.WithEndpoint("https://"+os.Getenv("VIKINGDB_HOST")),
		vector.WithRegion(os.Getenv("VIKINGDB_REGION")),
	)
	if err != nil {
		t.Fatal(err)
	}

	embeddingClient := client.Embedding()

	ctx := context.Background()

	text := "Reference dense and sparse embedding request."
	modelName := "bge-m3"

	request := model.EmbeddingRequest{
		DenseModel: &model.EmbeddingModelOpt{
			ModelName: stringPtr(modelName),
		},
		SparseModel: &model.EmbeddingModelOpt{
			ModelName: stringPtr(modelName),
		},
		Data: []*model.EmbeddingData{
			{Text: stringPtr(text)},
		},
	}

	resp, err := embeddingClient.Embedding(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Result == nil || len(resp.Result.Data) == 0 {
		t.Fatalf("embedding response missing data")
	}
	if len(resp.Result.Data[0].DenseVectors) == 0 {
		t.Fatalf("embedding response missing dense vector")
	}
	t.Logf("Embedding request_id=%s dense_dims=%d", resp.RequestID, len(resp.Result.Data[0].DenseVectors))
}

// Scenario 5 – Embedding Pipeline
//
// This guide demonstrates how to obtain embeddings from VikingDB:
//   - Configure dense (and optional sparse) models.
//   - Provide text/image/video content, including multimodal sequences.
//   - Inspect the response for dense/sparse vectors and token usage metadata.
//
// As with other guides, ensure the models referenced are available in your region/account.
func TestScenarioEmbeddingMultiModalPipeline(t *testing.T) {
	env := requireEnv(t)
	client := mustNewClient(t, env).Embedding()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	text := "generate embeddings with VikingDB"
	modelName := "doubao-embedding-vision"
	modelVersion := "250615"

	bge := "bge-m3"
	req := model.EmbeddingRequest{
		DenseModel: &model.EmbeddingModelOpt{
			ModelName:    &modelName,
			ModelVersion: &modelVersion,
		},
		SparseModel: &model.EmbeddingModelOpt{
			ModelName: &bge,
		},
		Data: []*model.EmbeddingData{
			{
				//Text: &text,
				FullModalSeq: []model.FullModalData{
					{Text: &text},
				},
			},
		},
	}
	d, _ := json.Marshal(req)
	t.Logf(string(d))

	resp, err := client.Embedding(ctx, req)
	require.NoError(t, err, "embedding request failed")
	require.NotNil(t, resp.Result)
	require.NotEmpty(t, resp.Result.Data)
	require.NotEmpty(t, resp.Result.Data[0].DenseVectors)
	//require.NotEmpty(t, resp.Result.Data[0].SparseVectors)
	t.Logf("Embedding request_id=%s dense_dims=%d token_usage=%v", resp.RequestID, len(resp.Result.Data[0].DenseVectors), resp.Result.TokenUsage)
	t.Logf("Dense[:5]=%v, Sparse=%v", resp.Result.Data[0].DenseVectors[:5], resp.Result.Data[0].SparseVectors)
}

func TestScenarioEmbeddingDSPipeline(t *testing.T) {
	env := requireEnv(t)
	client := mustNewClient(t, env).Embedding()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	text := "generate text dense&sparse embeddings with VikingDB"
	modelName := "bge-m3"

	req := model.EmbeddingRequest{
		DenseModel: &model.EmbeddingModelOpt{
			ModelName: &modelName,
			//ModelVersion: &modelVersion, // use default
		},
		SparseModel: &model.EmbeddingModelOpt{
			ModelName: &modelName,
		},
		Data: []*model.EmbeddingData{
			{
				Text: &text,
			},
		},
	}

	resp, err := client.Embedding(ctx, req)
	require.NoError(t, err, "embedding request failed")
	require.NotNil(t, resp.Result)
	require.NotEmpty(t, resp.Result.Data)
	require.NotEmpty(t, resp.Result.Data[0].DenseVectors)
	//require.NotEmpty(t, resp.Result.Data[0].SparseVectors)
	t.Logf("Embedding request_id=%s dense_dims=%d token_usage=%v", resp.RequestID, len(resp.Result.Data[0].DenseVectors), resp.Result.TokenUsage)
	t.Logf("Dense[:5]=%v, Sparse=%v", resp.Result.Data[0].DenseVectors[:5], resp.Result.Data[0].SparseVectors)
}
