package main

import (
	"context"
	"fmt"
	"os"
	"time"

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

func addDocByURL(kc *knowledge.CollectionClient, docID, docName, docType, url string, tagList []kmodel.MetaItem) (*kmodel.AddDocResponse, error) {
	req := kmodel.AddDocV2Request{
		DocID:   &docID,
		DocName: &docName,
		DocType: &docType,
		URI:     &url,
		TagList: tagList,
	}
	return kc.AddDocV2(context.Background(), req)
}

func runDocCrud() error {
	client, err := initClient()
	if err != nil {
		return err
	}
	kc := initCollection(client)

	docID := "google-report-2025-q1"
	docName := "Google 2025 Q1 Financial Report"
	docType := "pdf"
	url := "https://pdf.dfcfw.com/pdf/H3_AP202504281663850212_1.pdf"
	meta := []kmodel.MetaItem{
		{FieldName: strPtr("category"), FieldType: strPtr("string"), FieldValue: "financial_report"},
		{FieldName: strPtr("quarter"), FieldType: strPtr("string"), FieldValue: "Q1"},
		{FieldName: strPtr("year"), FieldType: strPtr("int64"), FieldValue: 2025},
	}

	addRes, err := addDocByURL(kc, docID, docName, docType, url, meta)
	if err != nil {
		return err
	}
	fmt.Println("add_doc:", addRes)

	info, err := kc.GetDoc(context.Background(), docID, true)
	if err != nil {
		return err
	}
	fmt.Println("get_doc:", info)

	meta = append(meta, kmodel.MetaItem{FieldName: strPtr("updated_at"), FieldType: strPtr("int64"), FieldValue: 1714560000})
	if err := kc.UpdateDocMeta(context.Background(), docID, meta); err != nil {
		return err
	}
	fmt.Println("update_doc_meta: ok")

	time.Sleep(30 * time.Second)

	newName := docName + "-updated"
	if err := kc.UpdateDoc(context.Background(), docID, newName); err != nil {
		return err
	}
	fmt.Println("update_doc: ok")

	listReq := kmodel.ListDocsRequest{Offset: 0, Limit: 10, ReturnTokenUsage: boolPtr(true)}
	listRes, err := kc.ListDocs(context.Background(), listReq)
	if err != nil {
		return err
	}
	fmt.Println("list_docs:", listRes)

	// Optionally delete:
	// if err := kc.DeleteDoc(context.Background(), docID); err != nil {
	// 	return err
	// }
	// fmt.Println("delete_doc: ok")

	return nil
}

func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }

func main() {
	if err := runDocCrud(); err != nil {
		panic(err)
	}
}
