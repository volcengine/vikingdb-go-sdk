package main

import (
	"context"
	"fmt"
	"os"
	"time"

	memory "github.com/volcengine/vikingdb-go-sdk/memory"
	mmodel "github.com/volcengine/vikingdb-go-sdk/memory/model"
)

var eventTemplates = []string{"event_customer_support"}
var profileTemplates = []string{"profile_customer_support"}

func initClient() (*memory.Client, error) {
	ak := os.Getenv("VIKINGDB_AK")
	sk := os.Getenv("VIKINGDB_SK")
	return memory.New(
		memory.AuthIAM(ak, sk),
		memory.WithEndpoint("http://api-knowledgebase.mlp.cn-beijing.volces.com"),
		memory.WithRegion("cn-beijing"),
	)
}

func addSupportSession(collection *memory.CollectionClient) (*mmodel.Response, error) {
	nowTs := time.Now().UnixNano() / int64(time.Millisecond)
	messages := []map[string]interface{}{
		{"role": "system", "content": "你是某电商平台的AI客服助手。你的目标是快速理解用户问题，先澄清关键信息，再给出可执行的解决方案。必要时引导用户提供订单号/手机号后四位等信息（注意隐私保护：不要索要完整身份证号/银行卡号）。语气友好、专业、简洁。"},
		{"role": "user", "content": "我昨天买的耳机到货了，但只有左边有声音，右边没声，怎么处理？"},
		{"role": "assistant", "content": "抱歉影响你使用了。我先帮你快速排查一下，通常是连接/音源或单只故障导致。\n\n请你确认3点：\n1）你用的是蓝牙还是有线？\n2）换一台手机/电脑连接，右耳仍然没声音吗？\n3）把左右耳互换佩戴/重置配对后有变化吗？\n\n另外，为了继续帮你发起售后，请发我：订单号（或截图打码）+ 收货手机号后四位。"},
		{"role": "user", "content": "蓝牙的。我换了手机还是右边没声。重置配对也不行。订单号是 1234567890，手机号后四位 8899。"},
		{"role": "assistant", "content": "收到。根据你的排查结果（换设备仍无声、重置无效），基本可判断为单只硬件异常。\n\n我这边给你两种处理方式，你选一个就行：\n- 方案A：换货（推荐）——寄回后我们给你换一副新的。\n- 方案B：退货退款——寄回验收后原路退款。\n\n补充确认两点以便我帮你提交：\n1）你更倾向换货还是退款？\n2）商品外观是否完好、配件（充电盒/线材/说明书）是否齐全？"},
		{"role": "user", "content": "换货吧。外观完好，配件都在。"},
		{"role": "assistant", "content": "好的，我帮你发起【换货】。\n\n流程说明：\n1）我将为订单 1234567890 创建换货申请；\n2）系统会生成上门取件/寄回地址（通常1小时内）；\n3）你把耳机和全部配件装回原包装（没有原箱也可用结实纸箱），外箱贴面单；\n4）仓库验收通过后，我们会安排补发新件。\n\n注意事项：\n- 请勿寄回个人赠品或非本订单物品；\n- 寄回前建议拍一张装箱照片，便于后续核对。\n\n你希望取件时间是：今天还是明天？大概哪个时间段（上午/下午/晚上）更方便？"},
		{"role": "user", "content": "明天下午吧。"},
		{"role": "assistant", "content": "可以的。我已为你提交换货申请，并备注【明天下午上门取件】。稍后你会在订单页看到换货单号和取件信息。\n\n如果明天取件前你想再自查一次，也可以试试：把耳机放回充电盒充满电后再重置一次（有时低电量会导致单耳失联）。不过按你目前情况，走换货会更省心。\n\n还有其他需要我一起处理的吗？"},
	}
	metadata := map[string]interface{}{
		"default_user_id":        "user_support_001",
		"default_user_name":      "SupportUser",
		"default_assistant_id":   "assistant_support_001",
		"default_assistant_name": "Customer Support",
		"time":                   nowTs,
	}
	return collection.AddSession(context.Background(), mmodel.AddSessionRequest{
		SessionID: "support_session_001",
		Messages:  messages,
		Metadata:  metadata,
	})
}

func searchSupportMemories(collection *memory.CollectionClient) (map[string]*mmodel.Response, error) {
	result := map[string]*mmodel.Response{}
	if len(eventTemplates) > 0 {
		eventFilter := map[string]interface{}{
			"user_id":      "user_support_001",
			"assistant_id": "assistant_support_001",
			"memory_type":  eventTemplates,
		}
		eventResp, err := collection.SearchEventMemory(context.Background(), mmodel.SearchEventMemoryRequest{
			Query:  "耳机 换货",
			Filter: eventFilter,
			Limit:  5,
		})
		if err != nil {
			return nil, err
		}
		result["event_memories"] = eventResp
	}
	if len(profileTemplates) > 0 {
		profileFilter := map[string]interface{}{
			"user_id":      "user_support_001",
			"assistant_id": "assistant_support_001",
			"memory_type":  profileTemplates,
		}
		profileResp, err := collection.SearchProfileMemory(context.Background(), mmodel.SearchProfileMemoryRequest{
			Filter: profileFilter,
			Limit:  5,
		})
		if err != nil {
			return nil, err
		}
		result["profile_memories"] = profileResp
	}
	return result, nil
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

	fmt.Println("=== Customer Support Best Practice ===")
	fmt.Println("Step 1: Adding support session...")
	addResult, err := addSupportSession(collection)
	if err != nil {
		panic(err)
	}
	fmt.Println("Add session result:", addResult)

	fmt.Println("Step 2: Waiting 30 seconds for data processing...")
	time.Sleep(30 * time.Second)
	fmt.Println("Wait completed")

	fmt.Println("Step 3: Searching support memories...")
	searchResult, err := searchSupportMemories(collection)
	if err != nil {
		panic(err)
	}
	fmt.Println("Search result:", searchResult)
	fmt.Println("Customer Support workflow completed")
}
