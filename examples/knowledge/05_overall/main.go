package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	knowledge "github.com/volcengine/vikingdb-go-sdk/knowledge"
	kmodel "github.com/volcengine/vikingdb-go-sdk/knowledge/model"
)

func requireEnv(name string) (string, error) {
	v := os.Getenv(name)
	if v == "" {
		return "", fmt.Errorf("missing required env: %s", name)
	}
	return v, nil
}

func initClient() (*knowledge.Client, error) {
	apiKey, err := requireEnv("VIKING_API_KEY")
	if err != nil {
		return nil, err
	}
	endpoint := "https://api-knowledgebase.mlp.cn-beijing.volces.com"
	region := "cn-beijing"

	return knowledge.New(
		knowledge.AuthAPIKey(apiKey),
		knowledge.WithEndpoint(endpoint),
		knowledge.WithRegion(region),
	)
}

func initCollection(c *knowledge.Client) (*knowledge.CollectionClient, error) {
	collectionName, err := requireEnv("VIKING_COLLECTION_NAME")
	if err != nil {
		return nil, err
	}
	projectName := os.Getenv("VIKING_PROJECT")
	if projectName == "" {
		projectName = "default"
	}
	resourceID := os.Getenv("VIKING_COLLECTION_RID")
	return c.Collection(kmodel.CollectionMeta{
		ResourceID:     resourceID,
		CollectionName: collectionName,
		ProjectName:    projectName,
	}), nil
}

func addDocV2(ctx context.Context, kc *knowledge.CollectionClient, docID, docName, uri string) (*kmodel.AddDocResponse, error) {
	req := kmodel.AddDocV2Request{
		DocID:   &docID,
		DocName: nil,
		URI:     &uri,
	}
	if docName != "" {
		req.DocName = &docName
	}
	resp, err := kc.AddDocV2(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func runOverall() error {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	client, err := initClient()
	if err != nil {
		return err
	}
	kc, err := initCollection(client)
	if err != nil {
		return err
	}
	if kc == nil {
		return fmt.Errorf("collection client is nil")
	}

	tosURI := "tos://your-bucket/your-path/your-file.pdf"
	tosDocID := "your-doc-id"
	tosResp, err := addDocV2(ctx, kc, tosDocID, "", tosURI)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else if tosResp.Data != nil && tosResp.Data.DocID != nil {
		fmt.Println("add tos doc id:", *tosResp.Data.DocID)
	}

	urlURI := "https://your-url.pdf"
	urlDocID := "your-doc-id"
	urlDocName := "your-file-name.pdf"

	urlResp, err := addDocV2(ctx, kc, urlDocID, urlDocName, urlURI)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else if urlResp.Data != nil && urlResp.Data.DocID != nil {
		fmt.Println("add url doc id:", *urlResp.Data.DocID)
	}

	query := "your query"
	serviceRID, err := requireEnv("VIKING_SERVICE_RID")
	if err != nil {
		return err
	}

	stream := false
	serviceChatReq := kmodel.ServiceChatRequest{
		ServiceResourceID: serviceRID,
		Messages:          []kmodel.ChatMessage{{Role: "user", Content: query}},
		Stream:            &stream,
	}
	serviceChatResp, err := client.ServiceChat(ctx, serviceChatReq)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(serviceChatResp, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println("service_chat:", string(b))

	return nil
}

func main() {
	if err := runOverall(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
