# Volc-VikingDB Golang SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/volcengine/vikingdb-go-sdk.svg)](https://pkg.go.dev/github.com/volcengine/vikingdb-go-sdk)

Volc-VikingDB Golang SDK 为与 Volc-VikingDB 服务交互提供了一套全面的工具。它旨在实现直观、灵活和高效，使开发人员能够轻松地将其应用程序与 Volc-VikingDB 集成。

## 特性

- **全面的 API 覆盖**：SDK 为 Volc-VikingDB 的所有功能提供了完整的接口，包括集合、索引和数据操作。
- **专用客户端**：SDK 为不同的功能模块（如 `CollectionClient`、`IndexClient` 和 `EmbeddingClient`）提供专用客户端，以提供更专注、更直观的开发体验。
- **灵活的配置**：SDK 支持客户端的灵活配置，包括端点、区域和凭据的自定义设置。
- **自动重试**：SDK 会自动重试因网络问题或服务器端瞬时错误而失败的请求，从而提高了应用程序的可靠性。
- **上下文感知**：SDK 具有上下文感知能力，可以更好地控制请求取消和超时。

## 安装

要安装 Volc-VikingDB Golang SDK，您可以使用 `go get` 命令：

```bash
go get github.com/volcengine/vikingdb-go-sdk
```

## 快速入门

最直接的方式是像 `examples/vector/README.md` 中的示例一样，通过环境变量加载配置。这样可以把密钥与代码分离，并且可以随时切换不同的凭据。

### 配置环境变量

| 变量                            | 说明                                       |
|---------------------------------|--------------------------------------------|
| `VIKINGDB_AK` / `VIKINGDB_SK`   | 用于请求签名的 IAM Access Key/Secret Key。  |
| `VIKINGDB_HOST`                 | 完整的 API 域名（不含协议前缀）。           |
| `VIKINGDB_REGION`               | 用于签名与路由的地域。                     |
| `VIKINGDB_COLLECTION`           | 集合与索引客户端默认使用的集合名称。       |
| `VIKINGDB_INDEX`                | 搜索示例默认使用的索引名称。               |

建议在项目根目录准备一个 `.env` 文件，便于本地加载。`examples/vector/.env` 提供了一个可直接复制的模板。

```bash
cp examples/vector/.env .env
# 根据实际情况修改 .env 中的配置
export $(grep -v '^#' .env | xargs)                 # 加载到当前 shell
# 或者仅在单条命令中注入变量
env $(grep -v '^#' .env | xargs) go test ./examples/vector -run TestScenario
```

### 初始化客户端

环境变量准备好后，可以在代码中从进程环境读取配置并初始化 SDK：

```go
package main

import (
    "context"
    "log"
    "os"
    "time"

    "github.com/volcengine/vikingdb-go-sdk/vector"
    "github.com/volcengine/vikingdb-go-sdk/vector/model"
)

func main() {
    client, err := vector.New(
        vector.AuthIAM(os.Getenv("VIKINGDB_AK"), os.Getenv("VIKINGDB_SK")),
        vector.WithEndpoint("https://" + os.Getenv("VIKINGDB_HOST")),
        vector.WithRegion(os.Getenv("VIKINGDB_REGION")),
        vector.WithTimeout(30*time.Second),
    )
    if err != nil {
        log.Fatal(err)
    }

    collection := client.Collection(model.CollectionLocator{
        CollectionName: os.Getenv("VIKINGDB_COLLECTION"),
    })
    index := client.Index(model.IndexLocator{
        CollectionLocator: model.CollectionLocator{
            CollectionName: os.Getenv("VIKINGDB_COLLECTION"),
        },
        IndexName: os.Getenv("VIKINGDB_INDEX"),
    })

    limit := 5
    resp, err := index.SearchByRandom(context.Background(), model.SearchByRandomRequest{
        SearchBase: model.SearchBase{Limit: &limit},
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("request_id=%s hits=%d", resp.RequestID, len(resp.Result.Data))
}
```

如果需要使用 API Key，可以导出 `VIKINGDB_API_KEY` 并将认证选项替换为 `vector.AuthAPIKey(os.Getenv("VIKINGDB_API_KEY"))`。

### 数据操作

完成初始化后，就可以使用对应的客户端（`collection`、`index`、`embedding`）调用 VikingDB 的各类接口，例如 `Upsert`、`Update`、`Delete`、`Fetch`、`SearchByVector`、`SearchByMultiModal` 和 `SearchByKeywords`。

#### Upsert 数据

`Upsert` 操作用于插入新数据或更新现有数据。

```go
import (
    "context"
    "github.com/volcengine/vikingdb-go-sdk/vector/model"
)

req := model.UpsertDataRequest{
    WriteDataBase: model.WriteDataBase{
        Data: []model.MapStr{
            {
                "ID": 1,
                "vector": []float64{1.1,2.2,3.4,4.2},
            },
        },
    },
}

resp, err := collection.Upsert(context.Background(), req)
if err != nil {
    // 处理错误
}
// 处理响应
```

#### 向量检索

使用 `SearchByVector` 可以将查询向量与索引中的向量进行对比，实现向量检索。

```go
import (
    "context"
    "github.com/volcengine/vikingdb-go-sdk/vector/model"
)

limit := 5
req := model.SearchByVectorRequest{
    SearchBase: model.SearchBase{
        Limit:        &limit,
        OutputFields: []string{"title", "score"},
    },
    DenseVector: []float64{0.1, 0.5, 0.2, 0.8},
}

resp, err := index.SearchByVector(context.Background(), req)
if err != nil {
    // 处理错误
}
// 处理响应
```

### 文本嵌入

`Embedding` 操作用于将文本转换为向量。以下示例假设已经按前文构建了 `client`：

```go
import (
    "context"
    "github.com/volcengine/vikingdb-go-sdk/vector/model"
)

text := "hello world"
model_name := "doubao-embedding"
model_version := "240715"
req := model.EmbeddingRequest{
    DenseModel: &model.EmbeddingModelOpt{
        ModelName: &model_name,
        ModelVersion: &model_version,
    },
    Data: []*model.EmbeddingData{
        {
            Text: &text,
        },
    },
}

embeddingClient := client.Embedding()

resp, err := embeddingClient.Embedding(context.Background(), req)
if err != nil {
    // 处理错误
}
// 处理响应
```

## API 参考

有关详细的 API 参考，请访问 [Go Reference](https://pkg.go.dev/github.com/volcengine/vikingdb-go-sdk)。
