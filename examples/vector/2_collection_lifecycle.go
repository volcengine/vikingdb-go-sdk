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

// CollectionLifecycle keeps the lifecycle flow inline so readers can see each request payload.
func CollectionLifecycle() {
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
		panic(err)
	}
	if upsertResp != nil {
		log.Printf("Upsert request_id=%s", upsertResp.RequestID)
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
		panic(err)
	}
	if searchResp == nil || searchResp.Result == nil || len(searchResp.Result.Data) == 0 {
		panic("SearchByMultiModal returned no hits")
	}
	chapterID := searchResp.Result.Data[0].ID
	if chapterID == nil {
		panic("SearchByMultiModal response missing chapter id")
	}
	log.Printf("SearchByMultiModal request_id=%s id=%v", searchResp.RequestID, chapterID)

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
		panic(err)
	}
	if updateResp != nil {
		log.Printf("Update request_id=%s", updateResp.RequestID)
	}

	fetchReq := model.FetchDataInCollectionRequest{
		IDs: []interface{}{chapterID},
	}
	fetchResp, err := collectionClient.Fetch(ctx, fetchReq)
	if err != nil {
		panic(err)
	}
	if fetchResp != nil && fetchResp.Result != nil && len(fetchResp.Result.Items) > 0 {
		score := fetchResp.Result.Items[0].Fields["score"]
		if v, ok := score.(json.Number); ok {
			score, _ = v.Float64()
		}
		log.Printf("Fetch request_id=%s score=%v", fetchResp.RequestID, score)
	}

	deleteReq := model.DeleteDataRequest{
		IDs: []interface{}{chapterID},
	}
	deleteResp, err := collectionClient.Delete(ctx, deleteReq)
	if err != nil {
		panic(err)
	}
	if deleteResp != nil {
		log.Printf("Delete request_id=%s removed_id=%v", deleteResp.RequestID, chapterID)
	}
}
