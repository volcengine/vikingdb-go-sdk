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

// Connectivity keeps setup inline and skips helper functions so the flow reads sequentially for reference.
// This file does not load `.env` values programmatically; export them in the shell before running:
// > env $(grep -v '^#' ./examples/vector/.env | xargs) go test -v ./examples/vector -run Connectivity
// ---
// Connectivity 省略辅助函数，便于按顺序参考阅读案例代码。
// 本文件不会在 Go 代码中自动加载 `.env`；运行前请在 shell 中先导出环境变量（命令同上）。
func Connectivity() {
	client, err := vector.New(
		vector.AuthIAM(os.Getenv("VIKINGDB_AK"), os.Getenv("VIKINGDB_SK")),
		vector.WithEndpoint("https://"+os.Getenv("VIKINGDB_HOST")),
		vector.WithRegion(os.Getenv("VIKINGDB_REGION")),
	)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	log.Printf("request_id=%s hits=%d", resp.RequestID, len(resp.Result.Data))
}
