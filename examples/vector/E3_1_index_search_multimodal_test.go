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

// TestSnippetIndexSearchMultiModal shows a minimal multi-modal search with inline payloads.
func TestSnippetIndexSearchMultiModal(t *testing.T) {
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

	baseParagraph := time.Now().UnixNano() % 1_000_000
	chapters := []model.MapStr{
		{
			"title":     "Session kickoff",
			"paragraph": baseParagraph,
			"score":     73.0,
			"text":      "This kickoff chapter keeps the reference search grounded.",
		},
		{
			"title":     "Session filters",
			"paragraph": baseParagraph + 1,
			"score":     91.0,
			"text":      "Filter walkthrough for multi-modal search.",
		},
	}

	for _, chapter := range chapters {
		upsertReq := model.UpsertDataRequest{
			WriteDataBase: model.WriteDataBase{
				Data: []model.MapStr{chapter},
			},
		}
		resp, upsertErr := collectionClient.Upsert(ctx, upsertReq)
		if upsertErr != nil {
			t.Fatal(upsertErr)
		}
		if resp != nil {
			t.Logf("Upsert request_id=%s", resp.RequestID)
		}
	}

	time.Sleep(3 * time.Second)

	filter := model.MapStr{
		"op":    "range",
		"field": "paragraph",
		"gte":   baseParagraph,
		"lte":   baseParagraph + 1,
	}

	searchReq := model.SearchByMultiModalRequest{
		Text:            stringPtr("Which chapter explains the session filters?"),
		NeedInstruction: boolPtr(false),
		SearchBase: model.SearchBase{
			RecallBase: model.RecallBase{
				Filter: filter,
			},
			Limit:        intPtr(2),
			OutputFields: []string{"title", "score", "paragraph"},
		},
	}

	searchResp, err := indexClient.SearchByMultiModal(ctx, searchReq)
	if err != nil {
		t.Fatal(err)
	}
	if searchResp == nil || searchResp.Result == nil || len(searchResp.Result.Data) == 0 {
		t.Fatalf("SearchByMultiModal returned no hits")
	}
	for _, hit := range searchResp.Result.Data {
		t.Logf("SearchByMultiModal hit id=%v title=%v score=%v paragraph=%v", hit.ID, hit.Fields["title"], hit.Fields["score"], hit.Fields["paragraph"])
	}
}

// Scenario 3.1 – Multi-Modal Retrieval With Filters
//
// This chapter shows how to blend multi-modal search with scalar filters to focus on the
// most relevant sessions from our Atlas journey:
//  1. Upsert several themed chapters into the collection.
//  2. Use SearchByMultiModal with a narrative prompt to surface related chapters.
//  3. Apply score/paragraph filters so only the current session's highlights appear.
func TestScenarioIndexSearchMultiModal(t *testing.T) {
	env := requireEnv(t)
	sdk := mustNewClient(t, env)
	collectionClient := sdk.Collection(collectionBase(env))
	indexClient := sdk.Index(indexBase(env))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sessionTag := newSessionTag("index-multimodal")
	baseParagraph := time.Now().UnixNano() % 1_000_000
	chapters := buildStoryChapters(sessionTag, baseParagraph)

	// only one record at a time.
	for _, chapter := range chaptersToUpsert(chapters) {
		upsertReq := model.UpsertDataRequest{
			WriteDataBase: model.WriteDataBase{
				Data: []model.MapStr{chapter},
			},
		}
		upsertResp, err := collectionClient.Upsert(ctx, upsertReq)
		require.NoError(t, err, "upsert failed")
		require.NotNil(t, upsertResp.Result)
		t.Logf("Upsert request_id=%s token_usage=%v chapters_written=%d", upsertResp.RequestID, upsertResp.Result.TokenUsage, len(chapters))
	}

	time.Sleep(3 * time.Second)

	filter := andFilters(
		sessionParagraphBounds(baseParagraph, len(chapters)),
		scoreAtLeastFilter(85.0),
	)

	query := findChapter(t, chapters, "retrieval-lab").Text
	limit := 3
	searchReq := model.SearchByMultiModalRequest{
		Text: &query,
		SearchBase: model.SearchBase{
			RecallBase: model.RecallBase{
				Filter: filter,
			},
			Limit:        &limit,
			OutputFields: []string{"title", "score", "paragraph"},
		},
	}

	d, _ := json.Marshal(searchReq)
	t.Logf("Search query: %s", string(d))

	var searchResp *model.SearchResponse
	var lastErr error
	resp, err := indexClient.SearchByMultiModal(ctx, searchReq)
	if err != nil {
		lastErr = err
		t.Logf("SearchByMultiModal failed: %v", err)
		return
	} else if resp != nil && resp.Result != nil && len(resp.Result.Data) > 0 {
		searchResp = resp
	} else {
		lastErr = nil
		t.Errorf("SearchByMultiModal returned no matches")
		return
	}

	require.NotNilf(t, searchResp, "expected multi-modal search results (last error: %v)", lastErr)
	require.NotNil(t, searchResp.Result)

	for idx, item := range searchResp.Result.Data {
		t.Logf("SearchByMultiModal request_id=%s hit[%d]=id:%v title:%v score:%v paragraph:%v",
			searchResp.RequestID, idx, item.ID, item.Fields["title"], item.Fields["score"], item.Fields["paragraph"])
	}

	require.Equal(t, 2, len(searchResp.Result.Data), "score filter should surface the two advanced chapters")
}
