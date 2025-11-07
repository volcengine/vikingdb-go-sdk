# VikingDB Example Guides

The `examples/vector` package hosts executable guides that double as integration smoke tests for the VikingDB Go SDK. Each `E*` test demonstrates a distinct workflow while sharing a common setup so you can copy/paste working snippets straight into your own project.

## Quickstart Walkthrough

The snippet below mirrors the `TestScenarioConnectivity` flow without chasing helper functions. It shows the minimal pieces you need to talk to VikingDB.

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/volcengine/vikingdb-go-sdk/vector"
	"github.com/volcengine/vikingdb-go-sdk/vector/model"
)

func main() {
	client, err := vector.New(
		vector.AuthIAM(os.Getenv("VIKINGDB_AK"), os.Getenv("VIKINGDB_SK")),
		vector.WithEndpoint("https://" + os.Getenv("VIKINGDB_HOST")),
		vector.WithRegion(os.Getenv("VIKINGDB_REGION")),
	)
	if err != nil {
		log.Fatal(err)
	}

	index := client.Index(model.IndexLocator{
		CollectionLocator: model.CollectionLocator{
			CollectionName: os.Getenv("VIKINGDB_COLLECTION"),
		},
		IndexName: os.Getenv("VIKINGDB_INDEX"),
	})

	limit := 1
	resp, err := index.SearchByRandom(context.Background(), model.SearchByRandomRequest{
		SearchBase: model.SearchBase{Limit: &limit},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("request_id=%s hits=%d", resp.RequestID, len(resp.Result.Data))
}
```

Once this quickstart works, explore the richer scenarios below.

## Running the Guides

```bash
# prepare .env file and run with asserts
go test ./examples/vector -run TestScenario
# run more readable snippets
env $(grep -v '^#' ./examples/vector/.env | xargs) go test ./examples/vector -run 
```

| Variable                        | Purpose                                           |
|---------------------------------|---------------------------------------------------|
| `VIKINGDB_AK` / `VIKINGDB_SK`   | IAM access keys used for signing requests.        |
| `VIKINGDB_HOST`                 | Fully qualified API hostname (no scheme).         |
| `VIKINGDB_REGION`               | Region used for signing and routing.              |
| `VIKINGDB_COLLECTION`           | Default collection for collection/index APIs.     |
| `VIKINGDB_INDEX`                | Default index for search-focused guides.          |

Populate them in your shell or a `.env` file before running `go test`.

## Scenario Overview

| Guide     | Test (file)                                               | What it demonstrates                                                                 | Key SDK calls                                                                                           |
|-----------|-----------------------------------------------------------|----------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------|
| 1   | `TestScenarioConnectivity` (`E1_connectivity_test.go`)       | Bootstrap SDK clients with shared options and validate connectivity via a lightweight random search. | `vector.New`, `Client.Collection`, `Client.Index`, `Client.Embedding`, `IndexClient.SearchByRandom`     |
| 2   | `TestScenarioCollectionLifecycle` (`E2_collection_lifecycle_test.go`) | Full CRUD lifecycle for Atlas "chapter" documents, including ID hydration through search. | `CollectionClient.Upsert`, `IndexClient.SearchByMultiModal`, `CollectionClient.Update`, `CollectionClient.Fetch`, `CollectionClient.Delete` |
| 3.1 | `TestScenarioIndexSearchMultiModal` (`E3_1_index_search_multimodal_test.go`) | Multi-modal narrative search combined with scalar filters to focus on relevant chapters. | `CollectionClient.Upsert`, `IndexClient.SearchByMultiModal`                                             |
| 3.2 | `TestScenarioIndexSearchVector` (`E3_2_index_search_vector_test.go`)       | Embedding-assisted vector retrieval with score filtering and rerank validation.         | `CollectionClient.Upsert`, `EmbeddingClient.Embedding`, `IndexClient.SearchByVector`                    |
| 3.3 | `TestScenarioSearchKeywords` (`E3_3_search_by_keyword_test.go`)            | Keyword-focused retrieval with session filters to surface tagged content.               | `CollectionClient.Upsert`, `IndexClient.SearchByKeywords`                                               |
| 4   | `TestScenarioSearchExtensionsAndAnalytics` (`E4_search_aggregate_test.go`) | Aggregate score analytics over the current session's chapters.                          | `CollectionClient.Upsert`, `IndexClient.Aggregate`                                                      |
| 5   | `TestScenarioEmbeddingMultiModalPipeline` / `TestScenarioEmbeddingDSPipeline` (`E5_embedding_test.go`) | Dense and sparse embedding retrieval, including multimodal sequences.                    | `EmbeddingClient.Embedding`                                                                             |

## SDK API Coverage

`X` indicates the API is exercised by the corresponding guide.

| SDK API                       | Client     | E1 | E2 | E3.1 | E3.2 | E3.3 | E4 | E5 |
|-------------------------------|------------|----|----|------|------|------|----|----|
| `vector.New`                  | vector     | X  | X  | X    | X    | X    | X  | X  |
| `Client.Collection`           | vector     | X  | X  | X    | X    | X    | X  |    |
| `Client.Index`                | vector     | X  | X  | X    | X    | X    | X  |    |
| `Client.Embedding`            | vector     | X  |    |      | X    |      |    | X  |
| `CollectionClient.Upsert`     | Collection |    | X  | X    | X    | X    | X  |    |
| `CollectionClient.Update`     | Collection |    | X  |      |      |      |    |    |
| `CollectionClient.Delete`     | Collection |    | X  |      |      |      |    |    |
| `CollectionClient.Fetch`      | Collection |    | X  |      |      |      |    |    |
| `IndexClient.Fetch`           | Index      |    |    |      |      |      |    |    |
| `IndexClient.SearchByVector`  | Index      |    |    |      | X    |      |    |    |
| `IndexClient.SearchByMultiModal` | Index   |    | X  | X    |      |      |    |    |
| `IndexClient.SearchByID`      | Index      |    |    |      |      |      |    |    |
| `IndexClient.SearchByScalar`  | Index      |    |    |      |      |      |    |    |
| `IndexClient.SearchByKeywords`| Index      |    |    |      |      | X    |    |    |
| `IndexClient.SearchByRandom`  | Index      | X  |    |      |      |      |    |    |
| `IndexClient.Aggregate`       | Index      |    |    |      |      |      | X  |    |
| `EmbeddingClient.Embedding`   | Embedding  |    |    |      | X    |      |    | X  |

### Uncovered Areas

- Index-level fetch, ID lookup, and scalar-only search (`Fetch`, `SearchByID`, `SearchByScalar`) are not currently represented in the guides.
- API-key based constructors are unused; all examples authenticate with AK/SK credentials.
