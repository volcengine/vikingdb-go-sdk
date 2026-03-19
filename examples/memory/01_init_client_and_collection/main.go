package main

import (
	"context"
	"fmt"
	"os"

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
		resourceId     = os.Getenv("VIKING_COLLECTION_RID")
	)
	client, err := initClient()
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println("Client initialized successfully")

	collection1, err := client.GetCollection(collectionName, projectName)
	if err != nil {
		panic(err)
	}
	fmt.Println("collection1:", collection1 != nil)

	collection2, err := client.Collection(mmodel.CollectionMeta{
		ResourceID: resourceId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("collection2:", collection2 != nil)
}
