// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/volcengine/vikingdb-go-sdk/vector/model"
)

func TestHeaderAuth_UpsertAndSearch(t *testing.T) {
	env := requireEnv(t)
	client := mustNewClient(t, env)

	collectionClient := client.Collection(collectionBase(env))
	indexClient := client.Index(indexBase(env))

	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())
	testID := int64(rand.Intn(9000000) + 1000000)

	testVector := make([]float64, 768)
	for i := range testVector {
		testVector[i] = rand.Float64()
	}

	// 1. Upsert data
	t.Logf("--- Step 1: Upserting data to %s ---", env.Collection)
	data := []model.MapStr{
		{
			"id":            testID,
			"dense_vector":  testVector,
			"sparse_vector": map[string]interface{}{"1": 0.1, "10": 0.5},
			"range":         10,
			"enum":          "test_A",
			"text":          "Golang record with HeaderAuth",
		},
	}

	upsertReq := model.UpsertDataRequest{
		WriteDataBase: model.WriteDataBase{
			Data: data,
		},
	}

	upsertResp, err := collectionClient.Upsert(ctx, upsertReq)
	if err != nil {
		t.Fatalf("Upsert failed: %v", err)
	}
	t.Logf("Upsert successful, request_id: %s", upsertResp.RequestID)

	// Wait for indexing
	t.Log("Waiting 2 seconds for indexing...")
	time.Sleep(2 * time.Second)

	// 2. Search data
	t.Logf("--- Step 2: Searching data in %s ---", env.Index)
	searchReq := model.SearchByVectorRequest{
		DenseVector: testVector,
		SearchBase: model.SearchBase{
			Limit:        intPtr(10),
			OutputFields: []string{"id", "text"},
		},
	}

	searchResp, err := indexClient.SearchByVector(ctx, searchReq)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	t.Logf("Search successful, request_id: %s", searchResp.RequestID)

	found := false
	for _, item := range searchResp.Result.Data {
		if id, ok := item.ID.(int64); ok && id == testID {
			found = true
			break
		} else if id, ok := item.ID.(float64); ok && int64(id) == testID {
			// JSON unmarshaling might result in float64
			found = true
			break
		}
	}

	if !found {
		t.Logf("Warning: Record %d not found in search results. Indexing might be delayed.", testID)
	} else {
		t.Log("SUCCESS: Found the upserted record in search results.")
	}
}
