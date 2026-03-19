package main

import (
	"context"
	"fmt"
	"os"
	"time"

	memory "github.com/volcengine/vikingdb-go-sdk/memory"
	mmodel "github.com/volcengine/vikingdb-go-sdk/memory/model"
)

var eventTemplates = []string{"event_education"}
var profileTemplates = []string{"profile_education", "knowledge_mastery"}

func initClient() (*memory.Client, error) {
	ak := os.Getenv("VIKINGDB_AK")
	sk := os.Getenv("VIKINGDB_SK")
	return memory.New(
		memory.AuthIAM(ak, sk),
		memory.WithEndpoint("http://api-knowledgebase.mlp.cn-beijing.volces.com"),
		memory.WithRegion("cn-beijing"),
	)
}

func addEducationSession(collection *memory.CollectionClient) (*mmodel.Response, error) {
	nowTs := time.Now().UnixNano() / int64(time.Millisecond)
	messages := []map[string]interface{}{
		{"role": "system", "content": "你是教育助手（Education Assistant）。你的目标是：记录学生偏好与学习进度，持续个性化辅导。教学风格：循序渐进、先结论后例题、鼓励学生用自己的话复述；遇到不确定要先澄清。每次讲解包含：核心概念→步骤/方法→例题→小练习→检查理解。"},
		{"role": "user", "content": "我总是不会把二次函数化成顶点式，尤其是配方法。能教我吗？"},
		{"role": "assistant", "content": "可以。我们就盯住一个知识点：**用配方法把二次函数化成顶点式**。\n\n**1）你要记住的目标形式**\n二次函数顶点式：\n> y = a(x - h)^2 + k\n这样顶点就是 (h, k)，且 a 决定开口方向与胖瘦。\n\n**2）配方法的固定套路（适用于 a=1 或先把 a 提出去）**\n对 y = ax^2 + bx + c：\n- 第一步：如果 a≠1，先把 a 从 x^2 和 x 项里提出：\n  y = a(x^2 + (b/a)x) + c\n- 第二步：在括号内补一个平方：\n  x^2 + px = (x + p/2)^2 - (p/2)^2\n  这里 p = b/a\n- 第三步：把补进去的减回来，并整理成 a(x-h)^2 + k。\n\n**3）例题：把 y = x^2 + 6x + 5 化成顶点式**\n(1) 先看 a=1，直接配：\n  x^2 + 6x + 5\n(2) 把 x^2+6x 配成完全平方：\n  x^2 + 6x = (x+3)^2 - 9\n(3) 代回去：\n  y = (x+3)^2 - 9 + 5 = (x+3)^2 - 4\n所以顶点式：\n> y = (x+3)^2 - 4\n顶点：(-3, -4)\n\n**4）你常见的坑（提前避开）**\n- 坑A：只加了 (p/2)^2，却忘了减回来。\n- 坑B：p 是 x 的系数，不是 b b本身；要先确保 x^2 系数为 1（或先提 a）。\n\n**5）小练习（你来做，我来帮你改）**\n把下面这个化成顶点式：\n> y = 2x^2 + 8x + 1\n提示：先把 2 提出去。\n\n你先告诉我：提 2 之后括号里变成什么？"},
		{"role": "user", "content": "是不是 y = 2(x^2 + 4x) + 1？"},
		{"role": "assistant", "content": "对，第一步完全正确。\n\n接下来我们配括号里的：x^2 + 4x。\n- 这里 p=4，所以 (p/2)^2 = 2^2 = 4。\n因此：\n  x^2 + 4x = (x+2)^2 - 4\n\n代回原式：\n  y = 2[(x+2)^2 - 4] + 1\n  y = 2(x+2)^2 - 8 + 1\n  y = 2(x+2)^2 - 7\n\n所以顶点式是：\n> y = 2(x+2)^2 - 7\n顶点：(-2, -7)\n\n为了确认你真的掌握了：\n**问题**：你能用一句话解释为什么要先把 2 提出去吗？"},
	}
	metadata := map[string]interface{}{
		"default_user_id":        "user_education_001",
		"default_user_name":      "Student",
		"default_assistant_id":   "assistant_education_001",
		"default_assistant_name": "Education Assistant",
		"time":                   nowTs,
	}
	return collection.AddSession(context.Background(), mmodel.AddSessionRequest{
		SessionID: "education_session_001",
		Messages:  messages,
		Metadata:  metadata,
	})
}

func searchEducationMemories(collection *memory.CollectionClient) (map[string]*mmodel.Response, error) {
	result := map[string]*mmodel.Response{}
	if len(eventTemplates) > 0 {
		eventFilter := map[string]interface{}{
			"user_id":      "user_education_001",
			"assistant_id": "assistant_education_001",
			"memory_type":  eventTemplates,
		}
		eventResp, err := collection.SearchEventMemory(context.Background(), mmodel.SearchEventMemoryRequest{
			Query:  "二次函数 配方法",
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
			"user_id":      "user_education_001",
			"assistant_id": "assistant_education_001",
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

	fmt.Println("=== Education Assistant Best Practice ===")
	fmt.Println("Step 1: Adding education session...")
	addResult, err := addEducationSession(collection)
	if err != nil {
		panic(err)
	}
	fmt.Println("Add session result:", addResult)

	fmt.Println("Step 2: Waiting 30 seconds for data processing...")
	time.Sleep(30 * time.Second)
	fmt.Println("Wait completed")

	fmt.Println("Step 3: Searching education memories...")
	searchResult, err := searchEducationMemories(collection)
	if err != nil {
		panic(err)
	}
	fmt.Println("Search result:", searchResult)
	fmt.Println("Education Assistant workflow completed")
}
