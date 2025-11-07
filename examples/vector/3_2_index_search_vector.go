// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/volcengine/vikingdb-go-sdk/vector"
	"github.com/volcengine/vikingdb-go-sdk/vector/model"
)

// IndexSearchVector demonstrates SearchByVector end-to-end without helper functions.
// Note: the extra V in TestVSnippet is because this test require different index and collection,
// for test conveniently I changed prefix.
func IndexSearchVector(collection, index string) {
	client, err := vector.New(
		vector.AuthIAM(os.Getenv("VIKINGDB_AK"), os.Getenv("VIKINGDB_SK")),
		vector.WithEndpoint("https://"+os.Getenv("VIKINGDB_HOST")),
		vector.WithRegion(os.Getenv("VIKINGDB_REGION")),
	)
	if err != nil {
		panic(err)
	}

	collectionClient := client.Collection(model.CollectionLocator{CollectionName: collection})
	indexClient := client.Index(model.IndexLocator{
		CollectionLocator: model.CollectionLocator{CollectionName: collection},
		IndexName:         index,
	})
	embeddingClient := client.Embedding()

	ctx := context.Background()

	baseParagraph := time.Now().UnixNano() % 1_000_000
	chapters := []struct {
		title string
		text  string
	}{
		{title: "Vector intro", text: "Inline vector search example for the reference suite."},
		{title: "Vector deep dive", text: "Demonstrates embedding reuse for query vectors."},
	}

	modelName := "bge-m3"
	modelVersion := "default"

	embedReq := model.EmbeddingRequest{
		DenseModel: &model.EmbeddingModelOpt{
			ModelName:    stringPtr(modelName),
			ModelVersion: stringPtr(modelVersion),
		},
		Data: make([]*model.EmbeddingData, 0, len(chapters)),
	}
	for _, chapter := range chapters {
		embedReq.Data = append(embedReq.Data, &model.EmbeddingData{Text: stringPtr(chapter.text)})
	}

	embedResp, err := embeddingClient.Embedding(ctx, embedReq)
	if err != nil {
		panic(err)
	}
	if embedResp == nil || embedResp.Result == nil || len(embedResp.Result.Data) != len(chapters) {
		panic("unexpected embedding response")
	}

	toFloat64 := func(src []float32) []float64 {
		out := make([]float64, len(src))
		for i, v := range src {
			out[i] = float64(v)
		}
		return out
	}

	upsertPayload := make([]model.MapStr, 0, len(chapters))
	for idx, chapter := range chapters {
		upsertPayload = append(upsertPayload, model.MapStr{
			"title":     chapter.title,
			"paragraph": baseParagraph + int64(idx),
			"score":     80.0 + float64(idx),
			"text":      chapter.text,
			"vector":    toFloat64(embedResp.Result.Data[idx].DenseVectors),
		})
	}

	upsertReq := model.UpsertDataRequest{
		WriteDataBase: model.WriteDataBase{
			Data: upsertPayload,
		},
	}
	upsertResp, err := collectionClient.Upsert(ctx, upsertReq)
	if err != nil {
		panic(err)
	}
	if upsertResp != nil {
		log.Printf("Upsert request_id=%s", upsertResp.RequestID)
	}

	time.Sleep(3 * time.Second)

	queryText := "Show me the chapter that demonstrates embedding reuse for query vectors."
	queryReq := model.EmbeddingRequest{
		DenseModel: &model.EmbeddingModelOpt{
			ModelName:    stringPtr(modelName),
			ModelVersion: stringPtr(modelVersion),
		},
		Data: []*model.EmbeddingData{{Text: stringPtr(queryText)}},
	}
	queryResp, err := embeddingClient.Embedding(ctx, queryReq)
	if err != nil {
		panic(err)
	}
	if queryResp == nil || queryResp.Result == nil || len(queryResp.Result.Data) == 0 {
		panic("query embedding returned no vectors")
	}

	filter := model.MapStr{
		"op":    "range",
		"field": "paragraph",
		"gte":   baseParagraph,
		"lt":    baseParagraph + int64(len(chapters)),
	}
	searchReq := model.SearchByVectorRequest{
		SearchBase: model.SearchBase{
			RecallBase: model.RecallBase{
				Filter: filter,
			},
			Limit:        intPtr(3),
			OutputFields: []string{"title", "score", "paragraph"},
		},
		DenseVector: toFloat64(queryResp.Result.Data[0].DenseVectors),
	}

	searchResp, err := indexClient.SearchByVector(ctx, searchReq)
	if err != nil {
		panic(err)
	}
	if searchResp == nil || searchResp.Result == nil || len(searchResp.Result.Data) == 0 {
		panic("SearchByVector returned no hits")
	}
	for _, hit := range searchResp.Result.Data {
		log.Printf("SearchByVector hit id=%v title=%v score=%v", hit.ID, hit.Fields["title"], hit.Fields["score"])
	}
}
