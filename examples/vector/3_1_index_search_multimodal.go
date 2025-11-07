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

// IndexSearchMultiModal shows a minimal multi-modal search with inline payloads.
func IndexSearchMultiModal() {
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
			panic(upsertErr)
		}
		if resp != nil {
			log.Printf("Upsert request_id=%s", resp.RequestID)
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
		panic(err)
	}
	if searchResp == nil || searchResp.Result == nil || len(searchResp.Result.Data) == 0 {
		panic("SearchByMultiModal returned no hits")
	}
	for _, hit := range searchResp.Result.Data {
		log.Printf("SearchByMultiModal hit id=%v title=%v score=%v paragraph=%v", hit.ID, hit.Fields["title"], hit.Fields["score"], hit.Fields["paragraph"])
	}
}
