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

// RerankMultiModal keeps the multimodal embedding call inline for quick copy/paste.
func RerankMultiModal() {
	client, err := vector.New(
		vector.AuthIAM(os.Getenv("VIKINGDB_AK"), os.Getenv("VIKINGDB_SK")),
		vector.WithEndpoint("https://"+os.Getenv("VIKINGDB_HOST")),
		vector.WithRegion(os.Getenv("VIKINGDB_REGION")),
	)
	if err != nil {
		panic(err)
	}

	rerankClient := client.Rerank()

	ctx := context.Background()

	txt := "Here is a mountain and it is very beautiful."
	image := "https://ark-project.tos-cn-beijing.volces.com/images/view.jpeg"
	video := "https://ark-project.tos-cn-beijing.volces.com/doc_video/ark_vlm_video_input.mp4"
	fps := 0.2
	data := make([][]model.FullModalData, 0)
	data = append(data, []model.FullModalData{{Text: &txt}})
	data = append(data, []model.FullModalData{{Image: &image}})
	data = append(data, []model.FullModalData{{Video: map[string]interface{}{"value": video, "fps": fps}}})
	queryContent := "This is iceberg."
	instruction := "Whether the Document answers the Query or matches the content retrieval intent"

	request := model.RerankRequest{
		ModelName:    "doubao-seed-rerank",
		ModelVersion: "251028",
		Data:         data,
		Query:        []model.FullModalData{{Text: &queryContent}},
		Instruction:  &instruction,
	}

	resp, err := rerankClient.Rerank(ctx, request)
	if err != nil {
		panic(err)
	}
	if resp == nil || resp.Result == nil || len(resp.Result.Data) == 0 {
		panic("rerank response missing data")
	}
	for i, rerankItem := range resp.Result.Data {
		log.Printf("rerank item %d:  id=%d  score=%v", i, rerankItem.ID, rerankItem.Score)
	}
}
