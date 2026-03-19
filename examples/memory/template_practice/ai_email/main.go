package main

import (
	"context"
	"fmt"
	"os"
	"time"

	memory "github.com/volcengine/vikingdb-go-sdk/memory"
	mmodel "github.com/volcengine/vikingdb-go-sdk/memory/model"
)

var eventTemplates = []string{"event_email"}
var profileTemplates = []string{"profile_email"}

func initClient() (*memory.Client, error) {
	ak := os.Getenv("VIKINGDB_AK")
	sk := os.Getenv("VIKINGDB_SK")
	return memory.New(
		memory.AuthIAM(ak, sk),
		memory.WithEndpoint("http://api-knowledgebase.mlp.cn-beijing.volces.com"),
		memory.WithRegion("cn-beijing"),
	)
}

func addEmailSession(collection *memory.CollectionClient) (*mmodel.Response, error) {
	nowTs := time.Now().UnixNano() / int64(time.Millisecond)
	messages := []map[string]interface{}{
		{"role": "user", "content": `Subject: [Action Required] 推荐页改版 V2 灰度发布计划确认（01/27）

各位好，

我们计划在 01/27（周一）开启 V2 的 10% 灰度，01/29（周三）视数据决定是否全量。请大家确认以下事项：

1) 埋点字段是否今日可冻结（exposure_id、rank_index、event_time 毫秒）
2) 客户端是否能在周一 10:00 前提供提测包
3) 后端 /v2/feed 接口提测时间与回滚开关说明
4) QA 回归范围与灰度验收口径

文档：
- 发布计划（草稿）：https://doc.example.com/reco-v2-release-plan

请在今天 18:00 前邮件回复确认/风险。

Thanks,
Alice`},
		{"role": "user", "content": `Re: [Action Required] 推荐页改版 V2 灰度发布计划确认（01/27）

Alice 好，

1) 埋点字段我今天 18:00 前可以冻结并补充字段定义（必填/可选）。
- exposure_id：必填
- rank_index：必填
- event_time：新埋点统一毫秒（旧链路保持原样），我会在数据层做兼容转换，避免口径混乱。

我会把最终版埋点文档更新到：
https://doc.example.com/reco-v2-tracking-spec

另外灰度期间我会按小时监控 CTR、完播率、负反馈率，并在群里同步异常。

Bob`},
		{"role": "user", "content": `Re: [Action Required] 推荐页改版 V2 灰度发布计划确认（01/27）

确认后端计划：
2) /v2/feed 接口今天 20:00 前提测。
- 默认 AB 分桶：若请求缺少 ab_group，会走对照组逻辑，但我会加告警（缺失率>1% 报警）。
- 回滚：保留 /v1/feed，灰度开关关闭后端仍可回退到 v1，预计无需重启。

回滚说明我会在 01/24（周五）下班前补到发布计划文档。

Charlie`},
		{"role": "user", "content": `Re: [Action Required] 推荐页改版 V2 灰度发布计划确认（01/27）

我这边有风险需要同步：
- 目前联调完成约 80%，但有 2 个疑似与图片缓存相关的 crash 在排查。
- 我可以承诺：：明天（01/24）11:00 前给定位结论。

对发布节奏建议：
- 01/24 先给内测包
- 01/27 周一 10:00 前给提测包
- 若 crash 风险较高，建议先 1% 观察 2 小时再扩到 10%

另外 event_time 毫秒口径：我会保证推荐页 V2 新埋点上报毫秒，不动公共库，避免影响其他页面。

Diana`},
		{"role": "user", "content": `Re: [Action Required] 推荐页改版 V2 灰度发布计划确认（01/27）

QA 侧确认：
- 如果周一（01/27）要灰度 10%，我最晚需要周一 10:00 拿到冻结提测包。
- 回归重点：AB 分桶是否生效、曝光埋点触发时机、灰度开关/回滚是否可用。

请补充：
1) 灰度开关配置路径/操作人
2) 关键验收指标阈值（例如 CTR 下跌多少触发回滚）

Ellen`},
		{"role": "user", "content": `Re: [Action Required] 推荐页改版 V2 灰度发布计划确认（01/27）

感谢大家，结论我汇总如下（如有遗漏请直接 reply）：

【结论/决策】
- 灰度节奏：01/27（周一）10% 灰度；若 crash 风险高则先 1% 观察 2 小时再扩；01/29（周三）视数据与稳定性决定是否全量。
- 埋点口径：V2 新埋点 event_time 统一毫秒；旧链路保持不变；数据层做兼容转换。

【待办（Owner/DDL）】
1) Bob：01/23 18:00 前冻结埋点文档（含字段定义、必填/可选）
   - https://doc.example.com/reco-v2-tracking-spec
2) Charlie：01/23 20:00 前 /v2/feed 提测；01/24（周五）下班前补回滚说明与 ab_group 缺失告警策略
3) Diana：01/24 11:00 前给 crash 定位结论；01/27 10:00 前提供提测包；补端上灰度/回滚开关说明
4) Ellen：拿到周一提测包后开始回归；灰度阶段重点验证 AB/曝光埋点/回滚

【待确认】
- 灰度开关配置路径与操作人（Charlie/Diana 请在文档补充）
- 验收阈值（CTR/负反馈触发回滚的阈值）：Bob 提个建议，我来拍板写入发布计划

下次同步：01/27（周一）16:00 复盘灰度数据与稳定性。

Alice`},
	}
	metadata := map[string]interface{}{
		"default_user_id":        "user_email_001",
		"default_user_name":      "EmailUser",
		"default_assistant_id":   "assistant_email_001",
		"default_assistant_name": "Email Summary",
		"time":                   nowTs,
	}
	return collection.AddSession(context.Background(), mmodel.AddSessionRequest{
		SessionID: "email_session_001",
		Messages:  messages,
		Metadata:  metadata,
	})
}

func searchEmailMemories(collection *memory.CollectionClient) (map[string]*mmodel.Response, error) {
	result := map[string]*mmodel.Response{}
	if len(eventTemplates) > 0 {
		eventFilter := map[string]interface{}{
			"user_id":      "user_email_001",
			"assistant_id": "assistant_email_001",
			"memory_type":  eventTemplates,
		}
		eventResp, err := collection.SearchEventMemory(context.Background(), mmodel.SearchEventMemoryRequest{
			Query:  "灰度发布 待办",
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
			"user_id":      "user_email_001",
			"assistant_id": "assistant_email_001",
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

	fmt.Println("=== Email Summary Best Practice ===")
	fmt.Println("Step 1: Adding email session...")
	addResult, err := addEmailSession(collection)
	if err != nil {
		panic(err)
	}
	fmt.Println("Add session result:", addResult)

	fmt.Println("Step 2: Waiting 30 seconds for data processing...")
	time.Sleep(30 * time.Second)
	fmt.Println("Wait completed")

	fmt.Println("Step 3: Searching email memories...")
	searchResult, err := searchEmailMemories(collection)
	if err != nil {
		panic(err)
	}
	fmt.Println("Search result:", searchResult)
	fmt.Println("Email Summary workflow completed")
}
