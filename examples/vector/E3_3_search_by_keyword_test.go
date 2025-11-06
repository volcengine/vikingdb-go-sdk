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

// TestSnippetSearchKeywords illustrates a straightforward keyword query over scoped documents.
func TestSnippetSearchKeywords(t *testing.T) {
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
	documents := []model.MapStr{
		{
			"title":     "Signal playbook",
			"paragraph": baseParagraph,
			"score":     86.0,
			"text":      "Signal insights tailored for keyword demonstrations.",
		},
		{
			"title":     "Session recap",
			"paragraph": baseParagraph + 1,
			"score":     78.0,
			"text":      "Recap without keywords to show contrast in results.",
		},
	}

	for _, doc := range documents {
		upsertReq := model.UpsertDataRequest{
			WriteDataBase: model.WriteDataBase{
				Data: []model.MapStr{doc},
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
		"lt":    baseParagraph + 2,
	}
	keywordsReq := model.SearchByKeywordsRequest{
		Keywords: []string{"playbook"},
		SearchBase: model.SearchBase{
			RecallBase: model.RecallBase{
				Filter: filter,
			},
			Limit:        intPtr(2),
			OutputFields: []string{"title", "score"},
		},
	}

	searchResp, err := indexClient.SearchByKeywords(ctx, keywordsReq)
	if err != nil {
		t.Fatal(err)
	}
	if searchResp == nil || searchResp.Result == nil || len(searchResp.Result.Data) == 0 {
		t.Fatalf("SearchByKeywords returned no hits")
	}

	for _, item := range searchResp.Result.Data {
		t.Logf("SearchByKeywords hit id=%v title=%v score=%v", item.ID, item.Fields["title"], item.Fields["score"])
	}
}

// Scenario 4 – Search Extensions & Analytics
//
// Building on the Atlas journey, this guide strings together advanced discovery APIs:
//   - SearchByMultiModal keeps session prompts scoped with filters.
//   - SearchByKeywords surfaces content by explicit tags.
//   - SearchByRandom reminds us what's already in rotation.
//   - Aggregate captures score analytics over the current session.
//   - Sort reorders specific chapters using a query vector.
func TestScenarioSearchKeywords(t *testing.T) {
	env := requireEnv(t)
	sdk := mustNewClient(t, env)
	collectionClient := sdk.Collection(collectionBase(env))
	indexClient := sdk.Index(indexBase(env))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sessionTag := newSessionTag("index-extensions")
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

	time.Sleep(5 * time.Second)

	limit := 3
	sessionFilter := sessionParagraphBounds(baseParagraph, len(chapters))

	keywordsReq := model.SearchByKeywordsRequest{
		Keywords: []string{"Signal"},
		SearchBase: model.SearchBase{
			RecallBase: model.RecallBase{
				Filter: sessionFilter,
			},
			Limit:        &limit,
			OutputFields: []string{"title", "score", "text"},
		},
	}
	keyResp, keyErr := indexClient.SearchByKeywords(ctx, keywordsReq)
	require.NoError(t, keyErr, "SearchByKeywords failed")
	require.NotNil(t, keyResp)
	require.NotNil(t, keyResp.Result)
	require.NotEmpty(t, keyResp.Result.Data, "keywords search should surface tagged chapters")
	for _, item := range keyResp.Result.Data {
		title := requireStringField(t, item.Fields, "title")
		marshal, _ := json.Marshal(item)
		t.Logf(string(marshal))
		// search_by_keyword only take effect on rank, possible returned record don't contain keywords.
		// require.Contains(t, title, "Signal", "keywords search should return Signal chapters")
		require.Contains(t, title, sessionTag, "keywords search should remain scoped to this session")
	}
	t.Logf("SearchByKeywords request_id=%s hits=%d", keyResp.RequestID, len(keyResp.Result.Data))
}
