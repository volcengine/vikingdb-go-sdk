# Volc-VikingDB Golang SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/volcengine/vikingdb-go-sdk.svg)](https://pkg.go.dev/github.com/volcengine/vikingdb-go-sdk)

The Volc-VikingDB Golang SDK provides a comprehensive suite of tools for interacting with the Volc-VikingDB service. It is designed to be intuitive, flexible, and efficient, enabling developers to easily integrate their applications with Volc-VikingDB.

## Features

- **Comprehensive API Coverage**: The SDK provides a complete interface to all the features of Volc-VikingDB, including collections, indexes, and data manipulation.
- **Specialized Clients**: The SDK offers specialized clients for different functional modules, such as `CollectionClient`, `IndexClient`, and `EmbeddingClient`, to provide a more focused and intuitive development experience.
- **Flexible Configuration**: The SDK supports flexible configuration of clients, including custom settings for endpoints, regions, and credentials.
- **Automatic Retries**: The SDK automatically retries failed requests due to network issues or server-side transient errors, improving the reliability of applications.
- **Context-Awareness**: The SDK is context-aware, allowing for better control over request cancellation and timeouts.

## Installation

To install the Volc-VikingDB Golang SDK, you can use the `go get` command:

```bash
go get github.com/volcengine/vikingdb-go-sdk
```

## Quick Start

The quickest path to a working setup is to load configuration from environment variables, just like the runnable guides in `examples/vector/README.md`. This keeps credentials outside of your code and lets you swap contexts by switching shells.

### Configure the Environment

| Variable                      | Purpose                                           |
|-------------------------------|---------------------------------------------------|
| `VIKINGDB_AK` / `VIKINGDB_SK` | IAM access keys used for signing requests.        |
| `VIKINGDB_HOST`               | Fully qualified API hostname (no scheme).         |
| `VIKINGDB_REGION`             | Region used for signing and routing.              |
| `VIKINGDB_COLLECTION`         | Default collection for collection/index APIs.     |
| `VIKINGDB_INDEX`              | Default index for search-focused examples.        |

Create a `.env` file so you can source these values locally. The sample at `examples/vector/.env` is a convenient starting point.

```bash
cp examples/vector/.env .env
# edit .env with your project-specific values
export $(grep -v '^#' .env | xargs)                 # load once for the current shell
# or scope the variables to a single command
env $(grep -v '^#' .env | xargs) go test ./examples/vector -run TestScenario
```

### Initialize the Client

With the variables in place, construct the SDK client by reading from the process environment:

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

If you prefer API keys, export `VIKINGDB_API_KEY` and replace the auth option with `vector.AuthAPIKey(os.Getenv("VIKINGDB_API_KEY"))`.

### Data Operations

Once the client is configured, you can use the scoped clients (`collection`, `index`, `embedding`) to call into VikingDB. The SDK exposes operations such as `Upsert`, `Update`, `Delete`, `Fetch`, `SearchByVector`, `SearchByMultiModal`, and `SearchByKeywords`.

#### Upsert Data

The `Upsert` operation is used to insert new data or update existing data.

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
    // Handle error
}
// Process response
```

#### Vector Search

Use `SearchByVector` to perform vector retrieval by comparing a query vector against indexed vectors.

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
    // Handle error
}
// Process response
```

### Text Embedding

The `Embedding` operation is used to convert text into vectors. Assuming the `client` is created as shown in the quick start:

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
    // Handle error
}
// Process response
```

## API Reference

For a detailed API reference, please visit the [Go Reference](https://pkg.go.dev/github.com/volcengine/vikingdb-go-sdk).
