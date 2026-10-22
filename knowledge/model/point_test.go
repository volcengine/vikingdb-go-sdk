package model

import (
	"encoding/json"
	"testing"
)

func TestListPointsResponseUnmarshalNestedServiceFields(t *testing.T) {
	payload := []byte(`{
		"code": 0,
		"message": "ok",
		"request_id": "req-1",
		"data": {
			"collection_name": "knowledge-a",
			"count": 1,
			"total_num": 2,
			"point_list": [
				{
					"point_id": "point-1",
					"content": "hello"
				}
			]
		}
	}`)

	var response ListPointsResponse
	if err := json.Unmarshal(payload, &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if response.Data == nil {
		t.Fatal("expected data to be present")
	}
	if response.Data.CollectionName == nil || *response.Data.CollectionName != "knowledge-a" {
		t.Fatalf("collection_name: got %v", response.Data.CollectionName)
	}
	if response.Data.Count == nil || *response.Data.Count != 1 {
		t.Fatalf("count: got %v", response.Data.Count)
	}
	if response.Data.TotalNum == nil || *response.Data.TotalNum != 2 {
		t.Fatalf("total_num: got %v", response.Data.TotalNum)
	}
	if len(response.Data.PointList) != 1 {
		t.Fatalf("point_list length: got %d", len(response.Data.PointList))
	}
	if response.Data.PointList[0].PointID == nil || *response.Data.PointList[0].PointID != "point-1" {
		t.Fatalf("point_id: got %v", response.Data.PointList[0].PointID)
	}
}
