// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"context"
	"log"
	"os"

	"github.com/volcengine/vikingdb-go-sdk/vector"
	"github.com/volcengine/vikingdb-go-sdk/vector/model"
)

// EmbeddingMultiModal keeps the multimodal embedding call inline for quick copy/paste.
func EmbeddingMultiModal() {
	client, err := vector.New(
		vector.AuthIAM(os.Getenv("VIKINGDB_AK"), os.Getenv("VIKINGDB_SK")),
		vector.WithEndpoint("https://"+os.Getenv("VIKINGDB_HOST")),
		vector.WithRegion(os.Getenv("VIKINGDB_REGION")),
	)
	if err != nil {
		panic(err)
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
		panic(err)
	}
	if resp == nil || resp.Result == nil || len(resp.Result.Data) == 0 {
		panic("embedding response missing data")
	}
	if len(resp.Result.Data[0].DenseVectors) == 0 {
		panic("embedding response missing dense vector")
	}
	log.Printf("Embedding request_id=%s dense_dims=%d", resp.RequestID, len(resp.Result.Data[0].DenseVectors))
}

// EmbeddingDenseSparse shows dense+sparse embedding retrieval inline.
func EmbeddingDenseSparse() {
	client, err := vector.New(
		vector.AuthIAM(os.Getenv("VIKINGDB_AK"), os.Getenv("VIKINGDB_SK")),
		vector.WithEndpoint("https://"+os.Getenv("VIKINGDB_HOST")),
		vector.WithRegion(os.Getenv("VIKINGDB_REGION")),
	)
	if err != nil {
		panic(err)
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
		panic(err)
	}
	if resp == nil || resp.Result == nil || len(resp.Result.Data) == 0 {
		panic("embedding response missing data")
	}
	if len(resp.Result.Data[0].DenseVectors) == 0 {
		panic("embedding response missing dense vector")
	}
	log.Printf("Embedding request_id=%s dense_dims=%d", resp.RequestID, len(resp.Result.Data[0].DenseVectors))
}
