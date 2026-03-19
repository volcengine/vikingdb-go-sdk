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

	nowTs := time.Now().UnixNano() / int64(time.Millisecond)
	messages := []map[string]interface{}{
		{
			"role":    "user",
			"content": "Hello, how is the weather today?",
		},
		{
			"role":    "assistant",
			"content": "Today is sunny with a temperature of 22 degrees, perfect for going out.",
		},
	}
	metadata := map[string]interface{}{
		"default_user_id":        "user_3",
		"default_user_name":      "Li",
		"default_assistant_id":   "111",
		"default_assistant_name": "Smart Assistant",
		"time":                   nowTs,
	}

	resp, err := collection.AddSession(context.Background(), mmodel.AddSessionRequest{
		SessionID: "session_001",
		Messages:  messages,
		Metadata:  metadata,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Add session result:", resp)
}
