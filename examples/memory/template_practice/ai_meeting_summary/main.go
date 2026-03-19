package main

import (
	"context"
	"fmt"
	"os"
	"time"

	memory "github.com/volcengine/vikingdb-go-sdk/memory"
	mmodel "github.com/volcengine/vikingdb-go-sdk/memory/model"
)

var eventTemplates = []string{"meeting_summary"}
var profileTemplates = []string{"profile_meeting"}

func initClient() (*memory.Client, error) {
	ak := os.Getenv("VIKINGDB_AK")
	sk := os.Getenv("VIKINGDB_SK")
	return memory.New(
		memory.AuthIAM(ak, sk),
		memory.WithEndpoint("http://api-knowledgebase.mlp.cn-beijing.volces.com"),
		memory.WithRegion("cn-beijing"),
	)
}

func addMeetingSession(collection *memory.CollectionClient) (*mmodel.Response, error) {
	nowTs := time.Now().UnixNano() / int64(time.Millisecond)
	messages := []map[string]interface{}{
		{"role": "user", "content": "大家好，开始今天的项目例会，主题：推荐页改版 V2 的发布评审。目标是对齐范围、风险、上线时间、以及各自 action items。"},
		{"role": "user", "content": "我先同步后端：接口改造基本完成，新增了 /v2/feed 接口，支持AB参数。现在卡点是埋点字段还没最终定，怕上线后数据口径变动。"},
		{"role": "user", "content": "埋点字段我这边建议本周内冻结。V2 需要新增 exposure_id 和 rank_index，否则后面归因会很难做。另外老字段 event_time 要统一成毫秒。"},
		{"role": "user", "content": "客户端这边 UI 已经联调到 80%。但如果 event_time 要改成毫秒，我们端上要改一处公共库，会影响别的页面，风险有点大。"},
		{"role": "user", "content": "先记一下争议点：event_time 口径要不要在本次上线强制统一。王五你觉得不统一会有什么后果？"},
		{"role": "user", "content": "不统一的话，会出现同一个指标在不同报表里差 1000 倍的问题，后面修成本更高。我们可以允许旧字段继续上报，但新埋点必须毫秒，同时在数据层做兼容转换。"},
		{"role": "user", "content": "数据层兼容我支持。后端接口里我也可以直接返回毫秒时间戳，端上只要不做二次转换就行。"},
		{"role": "user", "content": "如果只要求新埋点毫秒，但不改公共库，只在推荐页埋点里保证毫秒，那我这边可控。公共库改动就先不动。"},
		{"role": "user", "content": "OK，这里形成一个决策：本次上线推荐页 V2 新埋点使用毫秒，旧链路保持不变；数据层做兼容转换，避免口径混乱。"},
		{"role": "user", "content": "测试侧有两个风险：1）AB 分桶参数如果缺失会走默认逻辑，可能导致灰度失效；2）曝光埋点触发时机变了，需要重新跑一遍回归。上线时间如果是下周三，我需要最晚周一中午拿到冻结包。"},
		{"role": "user", "content": "排期我们原来是下周三（01/29）全量，上周五（01/24）开始 10% 灰度。赵六、李四，这个能保证吗？"},
		{"role": "user", "content": "我这边灰度 01/24 有点紧。UI 还有两个 crash 问题在排查，比较像是图片缓存引起的。保守点：01/24 做内测包，01/27 周一再灰度。"},
		{"role": "user", "content": "后端没问题，接口今天可以提测。AB 参数缺失默认逻辑我再加一个告警，如果请求里没有 ab_group 就打点并报警。"},
		{"role": "user", "content": "那我们调整发布节奏：01/27（周一）10% 灰度，01/29（周三）视数据与稳定性再全量。小周，按这个节奏你需要什么？"},
		{"role": "user", "content": "我需要：1）周一上午 10 点前给到提测包；2）埋点字段列表冻结，最好今天发邮件/飞书文档；3）灰度开关和回滚滚方案写清楚。"},
		{"role": "user", "content": "埋点字段我今天 18:00 前给出最终版文档，并在文档里标注必填/可选。灰度期间我会盯核心指标：CTR、完播率、负反馈率，按小时出监控。"},
		{"role": "user", "content": "补充一个行动项：赵六你提到的 crash 问题，最迟什么时候能定位？如果周一还没解决，我们灰度要不要延后？"},
		{"role": "user", "content": "我今晚先把图片缓存策略降级，明天上午 11 点前给定位结论。如果还是高风险，我建议灰度推迟到周二，并先在 1% 观察 2 小时。"},
		{"role": "user", "content": "好，记录：客户端 crash 明天 11:00 前出结论；如高风险，先 1% 观察再扩到 10%。回滚方案由李四提供接口回退开关说明，赵六提供端上开关说明。"},
		{"role": "user", "content": "最后确认一下行动项：\n1）王五：今日 18:00 前冻结埋点文档（含 exposure_id、rank_index、毫秒口径说明）。\n2）李四：今日提测 /v2/feed；补 ab_group 缺失告警；周五前给回滚说明。\n3）赵六：明天 11:00 前 crash 定位；周一 10:00 前给测试包；端上灰度/回滚开关说明。\n4）小周：基于周一包开始回归，灰度阶段重点验证曝光埋点与AB分桶。\n\n下次会议：周一 16:00 复盘灰度数据与稳定性，决定周三是否全量。"},
	}
	metadata := map[string]interface{}{
		"default_user_id":        "user_meeting_001",
		"default_user_name":      "MeetingUser",
		"default_assistant_id":   "assistant_meeting_001",
		"default_assistant_name": "Meeting Summary",
		"time":                   nowTs,
	}
	return collection.AddSession(context.Background(), mmodel.AddSessionRequest{
		SessionID: "meeting_session_001",
		Messages:  messages,
		Metadata:  metadata,
	})
}

func searchMeetingMemories(collection *memory.CollectionClient) (map[string]*mmodel.Response, error) {
	result := map[string]*mmodel.Response{}
	if len(eventTemplates) > 0 {
		eventFilter := map[string]interface{}{
			"user_id":      "user_meeting_001",
			"assistant_id": "assistant_meeting_001",
			"memory_type":  eventTemplates,
		}
		eventResp, err := collection.SearchEventMemory(context.Background(), mmodel.SearchEventMemoryRequest{
			Query:  "灰度 发布 待办",
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
			"user_id":      "user_meeting_001",
			"assistant_id": "assistant_meeting_001",
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

	fmt.Println("=== Meeting Summary Best Practice ===")
	fmt.Println("Step 1: Adding meeting session...")
	addResult, err := addMeetingSession(collection)
	if err != nil {
		panic(err)
	}
	fmt.Println("Add session result:", addResult)

	fmt.Println("Step 2: Waiting 30 seconds for data processing...")
	time.Sleep(30 * time.Second)
	fmt.Println("Wait completed")

	fmt.Println("Step 3: Searching meeting memories...")
	searchResult, err := searchMeetingMemories(collection)
	if err != nil {
		panic(err)
	}
	fmt.Println("Search result:", searchResult)
	fmt.Println("Meeting Summary workflow completed")
}
