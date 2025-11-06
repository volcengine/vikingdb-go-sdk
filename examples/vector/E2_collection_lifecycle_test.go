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

// TestSnippetCollectionLifecycle keeps the lifecycle flow inline so readers can see each request payload.
func TestSnippetCollectionLifecycle(t *testing.T) {
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

	ctx := context.Background()

	paragraph := time.Now().UnixNano() % 1_000_000
	chapter := model.MapStr{
		"title":     "Lifecycle quickstart",
		"paragraph": paragraph,
		"score":     42.5,
		"text":      "Simple lifecycle payload written inline for the reference flow.",
	}

	upsertReq := model.UpsertDataRequest{
		WriteDataBase: model.WriteDataBase{
			Data: []model.MapStr{chapter},
		},
	}
	upsertResp, err := collectionClient.Upsert(ctx, upsertReq)
	if err != nil {
		t.Fatal(err)
	}
	if upsertResp != nil {
		t.Logf("Upsert request_id=%s", upsertResp.RequestID)
	}

	time.Sleep(2 * time.Second)

	searchReq := model.SearchByMultiModalRequest{
		Text:            stringPtr("Need the lifecycle quickstart chapter overview"),
		NeedInstruction: boolPtr(false),
		SearchBase: model.SearchBase{
			Limit:        intPtr(1),
			OutputFields: []string{"title", "score"},
		},
	}
	searchResp, err := indexClient.SearchByMultiModal(ctx, searchReq)
	if err != nil {
		t.Fatal(err)
	}
	if searchResp == nil || searchResp.Result == nil || len(searchResp.Result.Data) == 0 {
		t.Fatalf("SearchByMultiModal returned no hits")
	}
	chapterID := searchResp.Result.Data[0].ID
	if chapterID == nil {
		t.Fatalf("SearchByMultiModal response missing chapter id")
	}
	t.Logf("SearchByMultiModal request_id=%s id=%v", searchResp.RequestID, chapterID)

	newScore := 47.0
	updateReq := model.UpdateDataRequest{
		WriteDataBase: model.WriteDataBase{
			Data: []model.MapStr{
				{
					"__AUTO_ID__": chapterID,
					"text":        "Updated lifecycle payload for reference clarity.",
					"score":       newScore,
				},
			},
		},
	}
	updateResp, err := collectionClient.Update(ctx, updateReq)
	if err != nil {
		t.Fatal(err)
	}
	if updateResp != nil {
		t.Logf("Update request_id=%s", updateResp.RequestID)
	}

	fetchReq := model.FetchDataInCollectionRequest{
		IDs: []interface{}{chapterID},
	}
	fetchResp, err := collectionClient.Fetch(ctx, fetchReq)
	if err != nil {
		t.Fatal(err)
	}
	if fetchResp != nil && fetchResp.Result != nil && len(fetchResp.Result.Items) > 0 {
		score := fetchResp.Result.Items[0].Fields["score"]
		if v, ok := score.(json.Number); ok {
			score, _ = v.Float64()
		}
		t.Logf("Fetch request_id=%s score=%v", fetchResp.RequestID, score)
	}

	deleteReq := model.DeleteDataRequest{
		IDs: []interface{}{chapterID},
	}
	deleteResp, err := collectionClient.Delete(ctx, deleteReq)
	if err != nil {
		t.Fatal(err)
	}
	if deleteResp != nil {
		t.Logf("Delete request_id=%s removed_id=%v", deleteResp.RequestID, chapterID)
	}
}

// Scenario 2 – Collection Lifecycle
//
// Demonstrates how to manage data within a VikingDB collection:
//   - Upsert documents with scalar fields and vectors.
//   - Update existing documents to adjust scalar values or TTL.
//   - Fetch documents by primary key and inspect missing IDs.
//   - Delete documents when no longer needed.
func TestScenarioCollectionLifecycle(t *testing.T) {
	env := requireEnv(t)
	sdk := mustNewClient(t, env)
	collectionClient := sdk.Collection(collectionBase(env))
	indexClient := sdk.Index(indexBase(env))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sessionTag := newSessionTag("collection-lifecycle")
	baseParagraph := time.Now().UnixNano() % 1_000_000
	chapters := buildStoryChapters(sessionTag, baseParagraph)

	// Note: Collection with vectorize does not support upsert data more than one.
	for _, c := range chaptersToUpsert(chapters) {
		upsertReq := model.UpsertDataRequest{
			WriteDataBase: model.WriteDataBase{
				Data: []model.MapStr{c},
			},
		}
		upsertResp, err := collectionClient.Upsert(ctx, upsertReq)
		require.NoError(t, err, "upsert failed")
		require.NotNil(t, upsertResp.Result)
		t.Logf("Upsert request_id=%s token_usage=%v chapters_written=%d", upsertResp.RequestID, upsertResp.Result.TokenUsage, len(chapters))
	}

	assignChapterIDsViaSearch(ctx, t, indexClient, chapters, []string{"title", "paragraph", "score"})

	targetChapter := findChapter(t, chapters, "retrieval-lab")
	require.NotNil(t, targetChapter.ID, "retrieval lab chapter must have an id assigned")
	t.Logf("Managing lifecycle for chapter_id=%v title=%q", targetChapter.ID, targetChapter.Title)

	newScore := targetChapter.Score + 4.25
	updateReq := model.UpdateDataRequest{
		WriteDataBase: model.WriteDataBase{
			Data: []model.MapStr{
				{
					"__AUTO_ID__": targetChapter.ID,
					"text":        "updated content",
					"score":       newScore,
				},
			},
		},
	}
	updateResp, err := collectionClient.Update(ctx, updateReq)
	require.NoError(t, err, "update failed")
	require.NotNil(t, updateResp.Result)
	t.Logf("Update request_id=%s new_score=%.2f", updateResp.RequestID, newScore)

	fetchReq := model.FetchDataInCollectionRequest{
		IDs: []interface{}{targetChapter.ID},
	}
	fetchResp, err := collectionClient.Fetch(ctx, fetchReq)
	require.NoError(t, err, "fetch failed")
	require.NotNil(t, fetchResp.Result)
	require.Len(t, fetchResp.Result.Items, 1, "expected one document back")
	fetched := fetchResp.Result.Items[0]
	cFetched, _ := fetched.Fields["score"].(json.Number).Float64()
	require.EqualValues(t, newScore, cFetched, "score should reflect update")
	t.Logf("Fetch request_id=%s record=%+v missing=%v", fetchResp.RequestID, fetched, fetchResp.Result.NotFoundIDs)

	deleteReq := model.DeleteDataRequest{
		IDs: []interface{}{targetChapter.ID},
	}
	deleteResp, err := collectionClient.Delete(ctx, deleteReq)
	require.NoError(t, err, "delete failed")
	require.NotNil(t, deleteResp)
	t.Logf("Delete request_id=%s removed_id=%v", deleteResp.RequestID, targetChapter.ID)
}
