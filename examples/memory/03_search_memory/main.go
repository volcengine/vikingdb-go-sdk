package main

import (
	"context"
	"fmt"
	"os"
	"time"

	memory "github.com/volcengine/vikingdb-go-sdk/memory"
	mmodel "github.com/volcengine/vikingdb-go-sdk/memory/model"
)

func initClient() (*memory.Client, error) {
	ak := os.Getenv("VIKINGDB_AK")
	sk := os.Getenv("VIKINGDB_SK")
	return memory.New(
		memory.AuthIAM(ak, sk),
		memory.WithEndpoint("http://api-knowledgebase.mlp.cn-beijing.volces.com"),
		memory.WithRegion("cn-beijing"),
	)
}

func searchUserProfile(collection *memory.CollectionClient) (*mmodel.Response, error) {
	filter := map[string]interface{}{
		"user_id":      "user_3",
		"assistant_id": "111",
		"memory_type":  []string{"profile_v1"},
	}
	return collection.SearchMemory(context.Background(), mmodel.SearchMemoryRequest{
		Filter: filter,
		Limit:  1,
	})
}

func searchEventsByQuery(collection *memory.CollectionClient, query string) (*mmodel.Response, error) {
	filter := map[string]interface{}{
		"user_id":      "user_3",
		"assistant_id": "111",
		"memory_type":  []string{"event_v1"},
	}
	return collection.SearchMemory(context.Background(), mmodel.SearchMemoryRequest{
		Query:  query,
		Filter: filter,
		Limit:  10,
	})
}

func searchRecentEvents(collection *memory.CollectionClient) (*mmodel.Response, error) {
	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	oneDayAgo := currentTime - 24*60*60*1000
	filter := map[string]interface{}{
		"user_id":      "user_3",
		"assistant_id": "111",
		"memory_type":  []string{"event_v1"},
		"start_time":   oneDayAgo,
		"end_time":     currentTime,
	}
	return collection.SearchMemory(context.Background(), mmodel.SearchMemoryRequest{
		Filter: filter,
		Limit:  10,
	})
}

func main() {
	var (
		collectionName = os.Getenv("VIKING_COLLECTION_NAME")
		projectName    = os.Getenv("VIKING_PROJECT")
	)
	client, err := initClient()
	if err != nil {
		panic(err)
	}
	collection, err := client.GetCollection(collectionName, projectName)
	if err != nil {
		panic(err)
	}

	fmt.Println("=== Viking Memory Search Examples ===")

	res1, err := searchUserProfile(collection)
	if err != nil {
		panic(err)
	}
	fmt.Println("1. User profile search result:", res1)

	res2, err := searchEventsByQuery(collection, "how is the weather today")
	if err != nil {
		panic(err)
	}
	fmt.Println("2. Event search by query result:", res2)

	res3, err := searchRecentEvents(collection)
	if err != nil {
		panic(err)
	}
	fmt.Println("3. Recent events search result:", res3)
}
