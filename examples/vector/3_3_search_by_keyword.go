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

// IndexSearchKeywords illustrates a straightforward keyword query over scoped documents.
func IndexSearchKeywords() {
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
		panic(err)
	}
	if searchResp == nil || searchResp.Result == nil || len(searchResp.Result.Data) == 0 {
		panic("SearchByKeywords returned no hits")
	}

	for _, item := range searchResp.Result.Data {
		log.Printf("SearchByKeywords hit id=%v title=%v score=%v", item.ID, item.Fields["title"], item.Fields["score"])
	}
}
