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

func initClient() (*knowledge.Client, error) {
	ak := os.Getenv("VOLC_AK")
	sk := os.Getenv("VOLC_SK")
	return knowledge.New(knowledge.AuthIAM(ak, sk))
}

func initClientByAPIKey() (*knowledge.Client, error) {
	apiKey := os.Getenv("VIKING_SERVICE_API_KEY")
	return knowledge.New(knowledge.AuthAPIKey(apiKey))
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

func runSearchCollection() error {
	client, err := initClient()
	if err != nil {
		return err
	}
	kc := initCollection(client)
	req := kmodel.SearchCollectionRequest{
		Query:           "2025 Q1 revenue growth",
		Limit:           10,
		DenseWeight:     floatPtr(0.5),
		RerankSwitch:    boolPtr(false),
		RetrieveCount:   intPtr(25),
		EndpointID:      nil,
		RerankModel:     strPtr("Doubao-pro-4k-rerank"),
		RerankOnlyChunk: boolPtr(false),
	}
	res, err := kc.SearchCollection(context.Background(), req)
	if err != nil {
		fmt.Println("search_collection_error:", err)
		return nil
	}
	fmt.Println("search_collection:", res)
	return nil
}

func runSearchKnowledge() (*kmodel.SearchKnowledgeResponse, error) {
	client, err := initClient()
	if err != nil {
		return nil, err
	}
	kc := initCollection(client)
	req := kmodel.SearchKnowledgeRequest{
		Query:          "2025 Q1 revenue growth",
		ImageQuery:     nil,
		PreProcessing:  nil,
		PostProcessing: nil,
		QueryParam: &kmodel.QueryParam{
			DocFilter: map[string]interface{}{
				"op":    "must",
				"field": "quarter",
				"conds": []string{"Q1"},
			},
			IncludePathList: &[]string{"/google/2025/"},
		},
		Limit:       intPtr(10),
		DenseWeight: floatPtr(0.5),
	}
	res, err := kc.SearchKnowledge(context.Background(), req)
	if err != nil {
		fmt.Println("search_knowledge_error:", err)
		return nil, nil
	}
	fmt.Println("search_knowledge:", res)

	for _, item := range res.Data.ResultList {
		b, _ := json.Marshal(item)
		fmt.Println(string(b))
	}
	return res, nil
}

func makeMessages(sk *kmodel.SearchKnowledgeResponse, query string) []kmodel.ChatMessage {
	top := []string{}
	if sk != nil && sk.Data != nil {
		for i, item := range sk.Data.ResultList {
			if i >= 5 {
				break
			}
			title := ""
			if item.ChunkTitle != nil {
				title = *item.ChunkTitle
			}
			content := ""
			if item.Content != nil {
				content = *item.Content
			}
			top = append(top, fmt.Sprintf("【%s】\n%s", title, content))
		}
	}
	contextText := "（检索结果为空或不可用）"
	if len(top) > 0 {
		contextText = ""
		for idx, s := range top {
			if idx > 0 {
				contextText += "\n\n"
			}
			contextText += s
		}
	}
	systemPrompt := "你是一位专业的财报分析师，你需要根据「参考资料」来回答接下来的「用户问题」，这些信息在 <context></context> XML 标签之内。回答必须在参考资料范围内，尽可能简洁，无法回答时请礼貌说明并引导提供更多信息。\n\n<context>\n" + contextText + "\n</context>"
	return []kmodel.ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: query},
	}
}

func runSearchKnowledgeAndChatCompletion() error {
	skRes, err := runSearchKnowledge()
	if err != nil {
		return err
	}
	msgs := makeMessages(skRes, "总结下 2025 Q1 财报数据")
	client, err := initClient()
	if err != nil {
		return err
	}
	req := kmodel.ChatCompletionRequest{
		Model:            "Doubao-1-5-pro-32k",
		Messages:         msgs,
		MaxTokens:        intPtr(4096),
		Temperature:      floatPtr(0.1),
		ReturnTokenUsage: boolPtr(true),
		APIKey:           strPtr(os.Getenv("VIKING_CHAT_API_KEY")),
		Stream:           boolPtr(false),
	}
	res, err := client.ChatCompletion(context.Background(), req)
	if err != nil {
		fmt.Println("chat_completion_error:", err)
		return nil
	}
	fmt.Println("chat_completion:", res)
	return nil
}

func runSearchKnowledgeAndChatCompletionStream() error {
	skRes, err := runSearchKnowledge()
	if err != nil {
		return err
	}
	msgs := makeMessages(skRes, "总结下 2025 Q1 财报数据")
	client, err := initClient()
	if err != nil {
		return err
	}
	req := kmodel.ChatCompletionRequest{
		Model:            "Doubao-1-5-pro-32k",
		Messages:         msgs,
		MaxTokens:        intPtr(4096),
		Temperature:      floatPtr(0.1),
		ReturnTokenUsage: boolPtr(true),
		APIKey:           strPtr(os.Getenv("VIKING_CHAT_API_KEY")),
		Stream:           boolPtr(true),
	}
	ch, err := client.ChatCompletionStream(context.Background(), req)
	if err != nil {
		fmt.Println("chat_completion_stream_error:", err)
		return nil
	}
	fmt.Println("chat_completion_stream:")
	for item := range ch {
		if item.Data != nil && item.Data.GeneratedAnswer != nil {
			fmt.Print(*item.Data.GeneratedAnswer)
		}
	}
	fmt.Print("\n")
	return nil
}

func runServiceChat() error {
	client, err := initClientByAPIKey()
	if err != nil {
		return err
	}
	serviceRID := os.Getenv("VIKING_SERVICE_RID")
	req := kmodel.ServiceChatRequest{
		ServiceResourceID: serviceRID,
		Messages:          []kmodel.ChatMessage{{Role: "user", Content: "列举 2025 Q1 财报里的三项亮点"}},
		Stream:            boolPtr(false),
	}
	res, err := client.ServiceChat(context.Background(), req)
	if err != nil {
		fmt.Println("service_chat_error:", err)
		return nil
	}
	fmt.Println("service_chat:", res)
	return nil
}

func runServiceChatStream() error {
	client, err := initClientByAPIKey()
	if err != nil {
		return err
	}
	serviceRID := os.Getenv("VIKING_SERVICE_RID")
	req := kmodel.ServiceChatRequest{
		ServiceResourceID: serviceRID,
		Messages:          []kmodel.ChatMessage{{Role: "user", Content: "列举 2025 Q1 财报里的三项亮点"}},
		Stream:            boolPtr(true),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	ch, err := client.ServiceChatStream(ctx, req)
	if err != nil {
		fmt.Println("service_chat_stream_error:", err)
		return nil
	}
	fmt.Println("service_chat_stream:")
	for item := range ch {
		if item.Data != nil && item.Data.GeneratedAnswer != nil {
			fmt.Print(*item.Data.GeneratedAnswer)
		}
	}
	fmt.Print("\n")
	return nil
}

func runRerankOps() error {
	client, err := initClient()
	if err != nil {
		return err
	}
	query := "2025 Q1 revenue growth"
	datas := []kmodel.RerankDataItem{
		{Query: query, Content: strPtr("Revenue grew 12% YoY to $3.4B."), Title: strPtr("Revenue")},
		{Query: query, Content: strPtr("Operating margin improved by 1.5pp to 17%."), Title: strPtr("Margin")},
	}
	request := kmodel.RerankRequest{
		Datas:       datas,
		RerankModel: strPtr("Doubao-pro-4k-rerank"),
	}

	resp, err := client.Rerank(context.Background(), request)
	if err != nil {
		fmt.Println("rerank_error:", err)
		return nil
	}
	fmt.Println("rerank:", resp)
	return nil
}

func intPtr(i int) *int           { return &i }
func floatPtr(f float64) *float64 { return &f }
func boolPtr(b bool) *bool        { return &b }
func strPtr(s string) *string     { return &s }

func main() {
	_ = runSearchCollection()
	_, _ = runSearchKnowledge()
	_ = runServiceChat()
	_ = runServiceChatStream()
	_ = runRerankOps()
	_ = runSearchKnowledgeAndChatCompletion()
	_ = runSearchKnowledgeAndChatCompletionStream()
}
