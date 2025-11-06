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

// TestSnippetSearchAggregate keeps aggregation parameters inlined for quick scanning.
func TestSnippetSearchAggregate(t *testing.T) {
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
			"title":     "Aggregate intro",
			"paragraph": baseParagraph,
			"score":     70.0,
			"text":      "Reference aggregation payload number one.",
		},
		{
			"title":     "Aggregate follow-up",
			"paragraph": baseParagraph + 1,
			"score":     82.0,
			"text":      "Reference aggregation payload number two.",
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

	time.Sleep(2 * time.Second)

	aggReq := model.AggRequest{
		Field: stringPtr("paragraph"),
		Op:    "count",
		Cond:  model.MapStr{"gte": baseParagraph},
	}

	aggResp, err := indexClient.Aggregate(ctx, aggReq)
	if err != nil {
		t.Fatal(err)
	}
	if aggResp == nil || aggResp.Result == nil {
		t.Fatalf("aggregate returned empty response")
	}

	aggJSON, _ := json.Marshal(aggResp.Result.Agg)
	t.Logf("Aggregate request_id=%s agg=%s", aggResp.RequestID, string(aggJSON))
}

// Scenario 4 – Search Aggregations
//
// Building on the Atlas journey, discover counts, group by paragraphs:
//   - Aggregate captures score analytics over the current session.
func TestScenarioSearchExtensionsAndAnalytics(t *testing.T) {
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

	time.Sleep(3 * time.Second)

	field := "paragraph"
	op := "count"
	aggReq := model.AggRequest{
		Field: &field,
		Op:    op,
		Cond:  model.MapStr{"gt": 1},
	}
	aggResp, aggErr := indexClient.Aggregate(ctx, aggReq)
	require.NoError(t, aggErr, "Aggregate failed")
	require.NotNil(t, aggResp)
	require.NotNil(t, aggResp.Result)
	require.NotNil(t, aggResp.Result.Agg)
	d, _ := json.Marshal(aggResp.Result.Agg)

	require.NotZero(t, len(aggResp.Result.Agg))
	t.Logf("Aggregate request_id=%s agg=%s", aggResp.RequestID, string(d))
}
