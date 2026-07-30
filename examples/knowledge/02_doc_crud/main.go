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

func addDocV1ByURL(kc *knowledge.CollectionClient, docID, docName, docType, url string, tagList []kmodel.MetaItem) (*kmodel.AddDocResponse, error) {
	req := kmodel.AddDocRequest{
		AddType:    "url",
		DocID:      &docID,
		DocName:    &docName,
		DocType:    &docType,
		URL:        &url,
		Meta:       tagList,
		PathPrefix: []string{"google", "2025"},
	}
	return kc.AddDoc(context.Background(), req)
}

func addDocV2ByURL(kc *knowledge.CollectionClient, docID, docName, docType, url string, tagList []kmodel.MetaItem) (*kmodel.AddDocV2Response, error) {
	req := kmodel.AddDocV2Request{
		DocID:        &docID,
		DocName:      &docName,
		DocType:      &docType,
		URI:          &url,
		TagList:      tagList,
		PathSegments: []string{"google", "2026"},
	}
	return kc.AddDocV2(context.Background(), req)
}

func runDocCrud() error {
	client, err := initClient()
	if err != nil {
		return err
	}
	kc := initCollection(client)

	// Add doc v1
	docID := "google-report-2025-q1"
	docName := "Google 2025 Q1 Financial Report"
	docType := "pdf"
	url := "https://pdf.dfcfw.com/pdf/H3_AP202504281663850212_1.pdf"
	meta := []kmodel.MetaItem{
		{FieldName: strPtr("category"), FieldType: strPtr("string"), FieldValue: "financial_report"},
		{FieldName: strPtr("quarter"), FieldType: strPtr("string"), FieldValue: "Q1"},
		{FieldName: strPtr("year"), FieldType: strPtr("int64"), FieldValue: 2025},
	}

	addV1Res, err := addDocV1ByURL(kc, docID, docName, docType, url, meta)
	if err != nil {
		return fmt.Errorf("add doc v1 failed: %w", err)
	}
	fmt.Println("add_doc_v1:", addV1Res)

	// Add doc v2
	docID = "google-report-2026-q1"
	docName = "Google 2026 Q1 Financial Report"
	docType = "pdf"
	url = "https://s206.q4cdn.com/479360582/files/doc_financials/2026/q2/2026q2-alphabet-earnings-release.pdf"
	meta = []kmodel.MetaItem{
		{FieldName: strPtr("category"), FieldType: strPtr("string"), FieldValue: "financial_report"},
		{FieldName: strPtr("quarter"), FieldType: strPtr("string"), FieldValue: "Q1"},
		{FieldName: strPtr("year"), FieldType: strPtr("int64"), FieldValue: 2025},
	}
	addV2Res, err := addDocV2ByURL(kc, docID, docName, docType, url, meta)
	if err != nil {
		return fmt.Errorf("add doc v2 failed: %w", err)
	}
	fmt.Println("add_doc_v2:", addV2Res)

	info, err := kc.GetDoc(context.Background(), docID, true)
	if err != nil {
		return fmt.Errorf("get doc failed: %w", err)
	}
	fmt.Println("get_doc:", info)

	meta = append(meta, kmodel.MetaItem{FieldName: strPtr("updated_at"), FieldType: strPtr("int64"), FieldValue: 1714560000})
	updateMetaResp, err := kc.UpdateDocMeta(context.Background(), docID, meta)
	if err != nil {
		return fmt.Errorf("update doc meta failed: %w", err)
	}
	fmt.Println("update_doc_meta:", updateMetaResp)

	time.Sleep(30 * time.Second)

	newName := docName + "-updated"
	updateDocResp, err := kc.UpdateDoc(context.Background(), docID, newName)
	if err != nil {
		return fmt.Errorf("update doc failed: %w", err)
	}
	fmt.Println("update_doc:", updateDocResp)

	listReq := kmodel.ListDocsRequest{
		Offset: 0,
		Limit:  10,
		Filter: &kmodel.ListDocsFilter{
			DocIDList: []string{docID},
		},
		ReturnTokenUsage: boolPtr(true),
	}
	listRes, err := kc.ListDocs(context.Background(), listReq)
	if err != nil {
		return fmt.Errorf("list docs failed: %w", err)
	}
	fmt.Println("list_docs:", listRes)

	limit := 2
	listV2Res, err := kc.ListDocsV2(context.Background(), kmodel.ListDocsV2Request{
		Limit: &limit,
	})
	if err != nil {
		return fmt.Errorf("list docs v2 failed: %w", err)
	}
	fmt.Println("list_docs_v2:", listV2Res)

	filterLimit := 10
	searchByFilterRes, err := kc.SearchDocsByFilter(context.Background(), kmodel.SearchDocsByFilterRequest{
		Filter: map[string]interface{}{
			"op":    "must",
			"field": "category",
			"conds": []string{"financial_report"},
		},
		Limit: &filterLimit,
	})
	if err != nil {
		return fmt.Errorf("search docs by filter failed: %w", err)
	}
	fmt.Println("search_docs_by_filter:", searchByFilterRes)

	// Optionally delete:
	// if resp, err := kc.DeleteDoc(context.Background(), docID); err != nil {
	// 	return err
	// }
	// fmt.Println("delete_doc:", resp)

	return nil
}

func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }

func main() {
	if err := runDocCrud(); err != nil {
		panic(err)
	}
}
