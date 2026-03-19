package main

import (
	"context"
	"fmt"
	"os"
	"time"

	memory "github.com/volcengine/vikingdb-go-sdk/memory"
	mmodel "github.com/volcengine/vikingdb-go-sdk/memory/model"
)

var eventTemplates = []string{"event_v1"}
var profileTemplates = []string{"profile_v1"}

func initClient() (*memory.Client, error) {
	ak := os.Getenv("VIKINGDB_AK")
	sk := os.Getenv("VIKINGDB_SK")
	return memory.New(
		memory.AuthIAM(ak, sk),
		memory.WithEndpoint("http://api-knowledgebase.mlp.cn-beijing.volces.com"),
		memory.WithRegion("cn-beijing"),
	)
}

func addCompanionSession(collection *memory.CollectionClient) (*mmodel.Response, error) {
	nowTs := time.Now().UnixNano() / int64(time.Millisecond)
	messages := []map[string]interface{}{
		{"role": "user", "content": "我刚吃完晚饭，在沙发上摊着刷手机。"},
		{"role": "assistant", "content": "这种时刻很真实：身体终于坐下了，脑子还在滑来滑去。\n\n你晚饭吃了啥？我想听点生活细节。"},
		{"role": "user", "content": "点了个外卖，番茄牛腩饭，还行但有点咸。"},
		{"role": "assistant", "content": "番茄牛腩饭还行但偏咸这评价很到位。\n\n你吃饭的时候有配点啥吗？比如可乐/豆奶/汤？还是就硬吃？"},
		{"role": "user", "content": "就喝了点水。今天也没怎么运动。"},
		{"role": "assistant", "content": "那你现在如果愿意，我们做个超日常、超轻的：在家走 3 分钟。\n不用运动，不用出汗，就当把咸味散一散。\n\n你要是懒得起来也行，我换个问题：你今晚更想继续刷到困，还是想找个东西随便看两集/听点音乐？"},
		{"role": "user", "content": "继续刷到困吧，我最近脑子不想用力。"},
		{"role": "assistant", "content": "懂，今天的你想要的是省电模式。\n\n那我帮你把刷手机也刷得更舒服一点：\n- 你愿意把亮度调低一点吗？\n- 还有一个：定个小目标，比如再刷 15 分钟就去洗漱，别让它把你拖到太晚。\n\n我不催你，你就告诉我。你一般几点睡比较舒服？"},
		{"role": "user", "content": "12点左右吧，但经常拖到1点。"},
		{"role": "assistant", "content": "那我们不追求立刻变自律，就做一个更容易成功的版本：\n\n今晚的目标：**12:40 去洗漱**，洗完继续躺也可以。\n这样你不会被必须马上睡压到，但又能把入睡门槛降下来。\n\n你要不要我陪你做个小提醒？你回复我一个好，我就当作我们达成约定：到 12:40 你去洗漱。"},
	}
	metadata := map[string]interface{}{
		"default_user_id":        "user_companion_001",
		"default_user_name":      "CompanionUser",
		"default_assistant_id":   "assistant_companion_001",
		"default_assistant_name": "AI Companion",
		"time":                   nowTs,
	}
	return collection.AddSession(context.Background(), mmodel.AddSessionRequest{
		SessionID: "companion_session_001",
		Messages:  messages,
		Metadata:  metadata,
	})
}

func searchCompanionMemories(collection *memory.CollectionClient) (map[string]*mmodel.Response, error) {
	result := map[string]*mmodel.Response{}
	if len(eventTemplates) > 0 {
		eventFilter := map[string]interface{}{
			"user_id":      "user_companion_001",
			"assistant_id": "assistant_companion_001",
			"memory_type":  eventTemplates,
		}
		eventResp, err := collection.SearchEventMemory(context.Background(), mmodel.SearchEventMemoryRequest{
			Query:  "晚饭 睡觉",
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
			"user_id":      "user_companion_001",
			"assistant_id": "assistant_companion_001",
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

	fmt.Println("=== AI Companion Best Practice ===")
	fmt.Println("Step 1: Adding companion session...")
	addResult, err := addCompanionSession(collection)
	if err != nil {
		panic(err)
	}
	fmt.Println("Add session result:", addResult)

	fmt.Println("Step 2: Waiting 30 seconds for data processing...")
	time.Sleep(30 * time.Second)
	fmt.Println("Wait completed")

	fmt.Println("Step 3: Searching companion memories...")
	searchResult, err := searchCompanionMemories(collection)
	if err != nil {
		panic(err)
	}
	fmt.Println("Search result:", searchResult)
	fmt.Println("AI Companion workflow completed")
}
