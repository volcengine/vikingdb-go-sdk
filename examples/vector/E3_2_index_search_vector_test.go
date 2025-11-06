// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/volcengine/vikingdb-go-sdk/vector"
	"github.com/volcengine/vikingdb-go-sdk/vector/model"
)

// TestVSnippetIndexSearchVector demonstrates SearchByVector end-to-end without helper functions.
// Note: the extra V in TestVSnippet is because this test require different index and collection,
// for test conveniently I changed prefix.
func TestVSnippetIndexSearchVector(t *testing.T) {
	client, err := vector.New(
		vector.AuthIAM(os.Getenv("VIKINGDB_AK"), os.Getenv("VIKINGDB_SK")),
		vector.WithEndpoint("https://"+os.Getenv("VIKINGDB_HOST")),
		vector.WithRegion(os.Getenv("VIKINGDB_REGION")),
	)
	if err != nil {
		t.Fatal(err)
	}

	collectionClient := client.Collection(model.CollectionLocator{CollectionName: os.Getenv("VIKINGDB_COLLECTION")})
	indexClient := client.Index(model.IndexLocator{
		CollectionLocator: model.CollectionLocator{CollectionName: os.Getenv("VIKINGDB_COLLECTION")},
		IndexName:         os.Getenv("VIKINGDB_INDEX"),
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
		t.Fatal(err)
	}
	if embedResp == nil || embedResp.Result == nil || len(embedResp.Result.Data) != len(chapters) {
		t.Fatalf("unexpected embedding response")
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
		t.Fatal(err)
	}
	if upsertResp != nil {
		t.Logf("Upsert request_id=%s", upsertResp.RequestID)
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
		t.Fatal(err)
	}
	if queryResp == nil || queryResp.Result == nil || len(queryResp.Result.Data) == 0 {
		t.Fatalf("query embedding returned no vectors")
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
		t.Fatal(err)
	}
	if searchResp == nil || searchResp.Result == nil || len(searchResp.Result.Data) == 0 {
		t.Fatalf("SearchByVector returned no hits")
	}
	for _, hit := range searchResp.Result.Data {
		t.Logf("SearchByVector hit id=%v title=%v score=%v", hit.ID, hit.Fields["title"], hit.Fields["score"])
	}
}

// Scenario 3.2 – Vector Retrieval With Embeddings
//
// This guide mirrors the Python quickstart: write chapters that include a dense vector field,
// generate a query embedding, and retrieve matching chapters with SearchByVector.
func TestScenarioIndexSearchVector(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	env := requireEnv(t)
	env.Collection = "vector"
	env.Index = "vector_index"

	client := mustNewClient(t, env)
	collectionClient := client.Collection(collectionBase(env))
	indexClient := client.Index(indexBase(env))
	embeddingClient := client.Embedding()

	sessionTag := newSessionTag("search-vector")
	baseParagraph := currentParagraphSeed()
	chapters := buildStoryChapters(sessionTag, baseParagraph)

	// 1. Embed the chapter narratives so we can upsert them with a pre-computed dense vector.
	modelName := "bge-m3"
	modelVersion := "default"
	chapterEmbeddings := batchEmbedTexts(t, ctx, embeddingClient, chapters, modelName, modelVersion)

	upsertPayload := make([]model.MapStr, 0, len(chapters))
	for idx, chapter := range chapters {
		upsertPayload = append(upsertPayload, model.MapStr{
			"title":     chapter.Title,
			"paragraph": chapter.Paragraph,
			"score":     chapter.Score,
			"text":      chapter.Text,
			"vector":    chapterEmbeddings[idx],
		})
	}

	upsertReq := model.UpsertDataRequest{
		WriteDataBase: model.WriteDataBase{
			Data: upsertPayload,
		},
	}
	upsertResp, err := collectionClient.Upsert(ctx, upsertReq)
	require.NoError(t, err, "vector upsert failed")
	t.Logf("Upsert request_id=%s", upsertResp.RequestID)

	// 2. Allow the index a moment to surface the newly written vectors.
	time.Sleep(3 * time.Second)

	// 3. Embed the retrieval chapter again to build our query vector.
	targetChapter := findChapter(t, chapters, "retrieval-lab")
	queryVector := embedSingleText(t, ctx, embeddingClient, targetChapter.Text, modelName, modelVersion)

	// 4. Run SearchByVector over the chapter window we just created.
	filter := sessionParagraphBounds(baseParagraph, len(chapters))
	limit := 5
	searchReq := model.SearchByVectorRequest{
		SearchBase: model.SearchBase{
			RecallBase: model.RecallBase{
				Filter: filter,
			},
			Limit:        &limit,
			OutputFields: []string{"title", "score", "paragraph"},
		},
		DenseVector: queryVector,
	}

	searchResp, err := indexClient.SearchByVector(ctx, searchReq)
	require.NoError(t, err, "vector search failed")
	require.NotNil(t, searchResp.Result, "vector search should return results")
	require.Greater(t, len(searchResp.Result.Data), 0, "vector search should produce hits")

	hits := searchResp.Result.Data
	titles := make([]string, 0, len(hits))
	for _, hit := range hits {
		if title, ok := hit.Fields["title"].(string); ok {
			titles = append(titles, title)
		}
	}
	require.Containsf(t, titles, targetChapter.Title, "expected %q to appear in vector search results", targetChapter.Title)
}

func batchEmbedTexts(t *testing.T, ctx context.Context, embeddingClient vector.EmbeddingClient, chapters []*storyChapter, modelName, modelVersion string) [][]float64 {
	t.Helper()

	req := model.EmbeddingRequest{
		DenseModel: &model.EmbeddingModelOpt{
			ModelName:    &modelName,
			ModelVersion: &modelVersion,
		},
		Data: make([]*model.EmbeddingData, 0, len(chapters)),
	}
	for _, chapter := range chapters {
		text := chapter.Text
		req.Data = append(req.Data, &model.EmbeddingData{Text: &text})
	}

	resp, err := embeddingClient.Embedding(ctx, req)
	require.NoError(t, err, "embedding batch request failed")
	require.NotNil(t, resp.Result, "embedding batch should return a result")
	require.Len(t, resp.Result.Data, len(chapters), "embedding batch must mirror chapter count")

	out := make([][]float64, len(resp.Result.Data))
	for idx, item := range resp.Result.Data {
		require.NotEmptyf(t, item.DenseVectors, "missing dense vector for chapter %s", chapters[idx].Key)
		out[idx] = float32SliceToFloat64(item.DenseVectors)
	}
	return out
}

func embedSingleText(t *testing.T, ctx context.Context, embeddingClient vector.EmbeddingClient, text, modelName, modelVersion string) []float64 {
	t.Helper()

	req := model.EmbeddingRequest{
		DenseModel: &model.EmbeddingModelOpt{
			ModelName:    &modelName,
			ModelVersion: &modelVersion,
		},
		Data: []*model.EmbeddingData{{Text: &text}},
	}

	resp, err := embeddingClient.Embedding(ctx, req)
	require.NoError(t, err, "query embedding request failed")
	require.NotNil(t, resp.Result, "embedding response should include a result")
	require.NotEmpty(t, resp.Result.Data, "embedding response should include data")
	require.NotEmpty(t, resp.Result.Data[0].DenseVectors, "embedding response should include dense vectors")

	return float32SliceToFloat64(resp.Result.Data[0].DenseVectors)
}

func currentParagraphSeed() int64 {
	return time.Now().Unix() % 1_000_000
}
