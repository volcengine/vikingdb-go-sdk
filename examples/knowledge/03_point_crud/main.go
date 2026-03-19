package main

import (
	"context"
	"fmt"
	"os"

	knowledge "github.com/volcengine/vikingdb-go-sdk/knowledge"
	kmodel "github.com/volcengine/vikingdb-go-sdk/knowledge/model"
)

func initClient() (*knowledge.Client, error) {
	ak := os.Getenv("VOLC_AK")
	sk := os.Getenv("VOLC_SK")
	return knowledge.New(knowledge.AuthIAM(ak, sk))
}

func initCollection(c *knowledge.Client) *knowledge.CollectionClient {
	resourceID := os.Getenv("VIKING_COLLECTION_RID")
	collectionName := os.Getenv("VIKING_COLLECTION_NAME")
	if collectionName == "" {
		collectionName = "financial_reports"
	}
	projectName := os.Getenv("VIKING_PROJECT")
	if projectName == "" {
		projectName = "default"
	}
	return c.Collection(kmodel.CollectionMeta{
		ResourceID:     resourceID,
		CollectionName: collectionName,
		ProjectName:    projectName,
	})
}

func runPointCrud() error {
	client, err := initClient()
	if err != nil {
		return err
	}
	kc := initCollection(client)

	docID := "google-report-2025-q1"
	chunkType := "text"
	chunkTitle := "Revenue Highlights"
	content := "Revenue grew 12% YoY to $3.4B."
	fields := []map[string]interface{}{
		{"field_name": "topic", "field_type": "string", "field_value": "revenue"},
		{"field_name": "year", "field_type": "int64", "field_value": 2025},
		{"field_name": "quarter", "field_type": "string", "field_value": "Q1"},
	}

	addReq := kmodel.AddPointRequest{
		DocID:      docID,
		ChunkType:  chunkType,
		ChunkTitle: &chunkTitle,
		Content:    &content,
		Fields:     fields,
	}
	addRes, err := kc.AddPoint(context.Background(), addReq)
	if err != nil {
		return err
	}
	fmt.Println("add_point:", addRes)
	pointID := ""
	if addRes != nil && addRes.Data != nil && addRes.Data.PointID != nil {
		pointID = *addRes.Data.PointID
	}

	info, err := kc.GetPoint(context.Background(), pointID, true)
	if err != nil {
		return err
	}
	fmt.Println("get_point:", info)

	updatedContent := content + " Updated."
	updatedTitle := chunkTitle + " (Updated)"
	updReq := kmodel.UpdatePointRequest{
		Content:    &updatedContent,
		ChunkTitle: &updatedTitle,
	}
	updatePointResp, err := kc.UpdatePoint(context.Background(), pointID, updReq)
	if err != nil {
		return err
	}
	fmt.Println("update_point_content:", updatePointResp)

	updFields := []map[string]interface{}{
		{"field_name": "topic", "field_type": "string", "field_value": "revenue"},
		{"field_name": "revised", "field_type": "bool", "field_value": true},
	}
	updateFieldsResp, err := kc.UpdatePoint(context.Background(), pointID, kmodel.UpdatePointRequest{Fields: updFields})
	if err != nil {
		return err
	}
	fmt.Println("update_point_fields:", updateFieldsResp)

	listReq := kmodel.ListPointsRequest{Offset: 0, Limit: 10}
	getLink := true
	listReq.GetAttachmentLink = &getLink
	listRes, err := kc.ListPoints(context.Background(), listReq)
	if err != nil {
		return err
	}
	fmt.Println("list_points:", listRes)

	deleteResp, err := kc.DeletePoint(context.Background(), kmodel.DeletePointRequest{PointID: pointID})
	if err != nil {
		return err
	}
	fmt.Println("delete_point:", deleteResp)

	return nil
}

func main() {
	if err := runPointCrud(); err != nil {
		panic(err)
	}
}
