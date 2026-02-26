package main

import (
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

func main() {
	client, err := initClient()
	if err != nil {
		panic(err)
	}
	collection := initCollection(client)
	fmt.Println("client:", "Client")
	if collection != nil {
		fmt.Println("collection:", "CollectionClient")
	}
}
