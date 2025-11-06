// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"github.com/volcengine/vikingdb-go-sdk/vector"
	"github.com/volcengine/vikingdb-go-sdk/vector/model"
)

// TestSnippetConnectivity keeps setup inline and skips helper functions so the flow reads sequentially for reference.
// This file does not load `.env` values programmatically; export them in the shell before running:
// > env $(grep -v '^#' ./examples/vector/.env | xargs) go test -v ./examples/vector -run TestSnippetConnectivity
// ---
// TestSnippetConnectivity 省略辅助函数，便于按顺序参考阅读案例代码。
// 本文件不会在 Go 代码中自动加载 `.env`；运行前请在 shell 中先导出环境变量（命令同上）。
func TestSnippetConnectivity(t *testing.T) {
	client, err := vector.New(
		vector.AuthIAM(os.Getenv("VIKINGDB_AK"), os.Getenv("VIKINGDB_SK")),
		vector.WithEndpoint("https://"+os.Getenv("VIKINGDB_HOST")),
		vector.WithRegion(os.Getenv("VIKINGDB_REGION")),
	)
	if err != nil {
		t.Fatal(err)
	}

	index := client.Index(model.IndexLocator{
		CollectionLocator: model.CollectionLocator{
			CollectionName: os.Getenv("VIKINGDB_COLLECTION"),
		},
		IndexName: os.Getenv("VIKINGDB_INDEX"),
	})

	resp, err := index.SearchByRandom(context.Background(), model.SearchByRandomRequest{
		SearchBase: model.SearchBase{Limit: intPtr(1)},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("request_id=%s hits=%d", resp.RequestID, len(resp.Result.Data))
}

// Scenario 1 – Connecting to VikingDB
//
// This guide shows how to bootstrap the Go SDK with credentials sourced from the environment
// and validate connectivity via a lightweight SearchByRandom request. Treat the test as executable documentation:
//  1. Collect AK/SK, host, and region from the environment (set VIKINGDB_* before running).
//  2. Build a shared client configuration with endpoint, region, and retry tuning.
//  3. Instantiate collection/index/embedding clients through helper constructors.
//  4. Execute a small SearchByRandom call to confirm authentication and network reachability.
//  5. Reuse the helpers in other scenarios to keep configuration consistent.
func TestScenarioConnectivity(t *testing.T) {
	env := requireEnv(t)

	client := mustNewClient(t, env)
	collectionClient := client.Collection(collectionBase(env))
	indexClient := client.Index(indexBase(env))
	embeddingClient := client.Embedding()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Logf("Checking VikingDB connectivity host=%s region=%s collection=%s index=%s", env.Host, env.Region, env.Collection, env.Index)

	limit := 1
	randomReq := model.SearchByRandomRequest{
		SearchBase: model.SearchBase{
			Limit: &limit,
		},
	}
	randomResp, err := indexClient.SearchByRandom(ctx, randomReq)
	require.NoError(t, err, "SearchByRandom health check failed")
	if randomResp != nil && randomResp.Result != nil {
		t.Logf("SearchByRandom request_id=%s hits=%d", randomResp.RequestID, len(randomResp.Result.Data))
	}

	require.NotNil(t, collectionClient)
	require.NotNil(t, embeddingClient)
}

type guideEnv struct {
	AccessKey  string
	SecretKey  string
	Host       string
	Region     string
	Collection string
	Index      string
}

func requireEnv(t *testing.T) guideEnv {
	t.Helper()
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}
	if os.Getenv("VIKINGDB_AK") == "" {
		_ = godotenv.Load(".env.local")
	}

	env := guideEnv{
		AccessKey:  os.Getenv("VIKINGDB_AK"),
		SecretKey:  os.Getenv("VIKINGDB_SK"),
		Host:       os.Getenv("VIKINGDB_HOST"),
		Region:     os.Getenv("VIKINGDB_REGION"),
		Collection: os.Getenv("VIKINGDB_COLLECTION"),
		Index:      os.Getenv("VIKINGDB_INDEX"),
	}

	missing := []string{}
	if env.AccessKey == "" {
		missing = append(missing, "VIKINGDB_AK")
	}
	if env.SecretKey == "" {
		missing = append(missing, "VIKINGDB_SK")
	}
	if env.Host == "" {
		missing = append(missing, "VIKINGDB_HOST")
	}
	if env.Region == "" {
		missing = append(missing, "VIKINGDB_REGION")
	}
	if env.Collection == "" {
		missing = append(missing, "VIKINGDB_COLLECTION")
	}
	if env.Index == "" {
		missing = append(missing, "VIKINGDB_INDEX")
	}

	if len(missing) > 0 {
		t.Skipf("missing required environment variables: %v", missing)
	}

	return env
}

func mustNewClient(t *testing.T, env guideEnv) *vector.Client {
	t.Helper()

	client, err := vector.New(
		vector.AuthIAM(env.AccessKey, env.SecretKey),
		sharedClientOptions(env)...,
	)
	require.NoError(t, err)
	return client
}

func sharedClientOptions(env guideEnv) []vector.ClientOption {
	// all has default values
	return []vector.ClientOption{
		vector.WithEndpoint(fmt.Sprintf("https://%s", env.Host)),
		vector.WithRegion(env.Region),
		vector.WithMaxRetries(3),
		vector.WithTimeout(30 * time.Second),
		vector.WithUserAgent("vikingdb-go-sdk-guide"),
	}
}

func collectionBase(env guideEnv) model.CollectionLocator {
	return model.CollectionLocator{
		CollectionName: env.Collection,
		ProjectName:    "",
	}
}

func indexBase(env guideEnv) model.IndexLocator {
	return model.IndexLocator{
		CollectionLocator: model.CollectionLocator{
			CollectionName: env.Collection,
		},
		IndexName: env.Index,
	}
}
