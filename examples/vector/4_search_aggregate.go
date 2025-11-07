// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/volcengine/vikingdb-go-sdk/vector"
	"github.com/volcengine/vikingdb-go-sdk/vector/model"
)

// IndexSearchAggregate keeps aggregation parameters inlined for quick scanning.
func IndexSearchAggregate() {
	client, err := vector.New(
		vector.AuthIAM(os.Getenv("VIKINGDB_AK"), os.Getenv("VIKINGDB_SK")),
		vector.WithEndpoint("https://"+os.Getenv("VIKINGDB_HOST")),
		vector.WithRegion(os.Getenv("VIKINGDB_REGION")),
	)
	if err != nil {
		panic(err)
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
			panic(upsertErr)
		}
		if resp != nil {
			log.Printf("Upsert request_id=%s", resp.RequestID)
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
		panic(err)
	}
	if aggResp == nil || aggResp.Result == nil {
		panic("aggregate returned empty response")
	}

	aggJSON, _ := json.Marshal(aggResp.Result.Agg)
	log.Printf("Aggregate request_id=%s agg=%s", aggResp.RequestID, string(aggJSON))
}
