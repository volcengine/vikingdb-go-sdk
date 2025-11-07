package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"github.com/volcengine/vikingdb-go-sdk/vector"
	"github.com/volcengine/vikingdb-go-sdk/vector/model"
)

// Scenario 1 – Connecting to VikingDB
//
// This guide shows how to bootstrap the Go SDK with credentials sourced from the environment
// and validate connectivity via a lightweight SearchByRandom request. Treat the test as executable documentation:
//  1. Collect AK/SK, host, and region from the environment (set VIKINGDB_* before running).
//  2. Build a shared client configuration with endpoint, region, and retry tuning.
//  3. Instantiate collection/index/embedding clients through helper constructors.
//  4. Execute a small SearchByRandom call to confirm authentication and network reachability.
//  5. Reuse the helpers in other scenarios to keep configuration consistent.
func TestScenarioConnectivity(t *testing.T) {
	env := requireEnv(t)

	client := mustNewClient(t, env)
	collectionClient := client.Collection(collectionBase(env))
	indexClient := client.Index(indexBase(env))
	embeddingClient := client.Embedding()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("Checking VikingDB connectivity host=%s region=%s collection=%s index=%s", env.Host, env.Region, env.Collection, env.Index)

	limit := 1
	randomReq := model.SearchByRandomRequest{
		SearchBase: model.SearchBase{
			Limit: &limit,
		},
	}
	randomResp, err := indexClient.SearchByRandom(ctx, randomReq)
	require.NoError(t, err, "SearchByRandom health check failed")
	if randomResp != nil && randomResp.Result != nil {
		log.Printf("SearchByRandom request_id=%s hits=%d", randomResp.RequestID, len(randomResp.Result.Data))
	}

	require.NotNil(t, collectionClient)
	require.NotNil(t, embeddingClient)
}

// Scenario 2 – Collection Lifecycle
//
// Demonstrates how to manage data within a VikingDB collection:
//   - Upsert documents with scalar fields and vectors.
//   - Update existing documents to adjust scalar values or TTL.
//   - Fetch documents by primary key and inspect missing IDs.
//   - Delete documents when no longer needed.
func TestScenarioCollectionLifecycle(t *testing.T) {
	env := requireEnv(t)
	sdk := mustNewClient(t, env)
	collectionClient := sdk.Collection(collectionBase(env))
	indexClient := sdk.Index(indexBase(env))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sessionTag := newSessionTag("collection-lifecycle")
	baseParagraph := time.Now().UnixNano() % 1_000_000
	chapters := buildStoryChapters(sessionTag, baseParagraph)

	// Note: Collection with vectorize does not support upsert data more than one.
	for _, c := range chaptersToUpsert(chapters) {
		upsertReq := model.UpsertDataRequest{
			WriteDataBase: model.WriteDataBase{
				Data: []model.MapStr{c},
			},
		}
		upsertResp, err := collectionClient.Upsert(ctx, upsertReq)
		require.NoError(t, err, "upsert failed")
		require.NotNil(t, upsertResp.Result)
		log.Printf("Upsert request_id=%s token_usage=%v chapters_written=%d", upsertResp.RequestID, upsertResp.Result.TokenUsage, len(chapters))
	}

	assignChapterIDsViaSearch(ctx, t, indexClient, chapters, []string{"title", "paragraph", "score"})

	targetChapter := findChapter(t, chapters, "retrieval-lab")
	require.NotNil(t, targetChapter.ID, "retrieval lab chapter must have an id assigned")
	log.Printf("Managing lifecycle for chapter_id=%v title=%q", targetChapter.ID, targetChapter.Title)

	newScore := targetChapter.Score + 4.25
	updateReq := model.UpdateDataRequest{
		WriteDataBase: model.WriteDataBase{
			Data: []model.MapStr{
				{
					"__AUTO_ID__": targetChapter.ID,
					"text":        "updated content",
					"score":       newScore,
				},
			},
		},
	}
	updateResp, err := collectionClient.Update(ctx, updateReq)
	require.NoError(t, err, "update failed")
	require.NotNil(t, updateResp.Result)
	log.Printf("Update request_id=%s new_score=%.2f", updateResp.RequestID, newScore)

	fetchReq := model.FetchDataInCollectionRequest{
		IDs: []interface{}{targetChapter.ID},
	}
	fetchResp, err := collectionClient.Fetch(ctx, fetchReq)
	require.NoError(t, err, "fetch failed")
	require.NotNil(t, fetchResp.Result)
	require.Len(t, fetchResp.Result.Items, 1, "expected one document back")
	fetched := fetchResp.Result.Items[0]
	cFetched, _ := fetched.Fields["score"].(json.Number).Float64()
	require.EqualValues(t, newScore, cFetched, "score should reflect update")
	log.Printf("Fetch request_id=%s record=%+v missing=%v", fetchResp.RequestID, fetched, fetchResp.Result.NotFoundIDs)

	deleteReq := model.DeleteDataRequest{
		IDs: []interface{}{targetChapter.ID},
	}
	deleteResp, err := collectionClient.Delete(ctx, deleteReq)
	require.NoError(t, err, "delete failed")
	require.NotNil(t, deleteResp)
	log.Printf("Delete request_id=%s removed_id=%v", deleteResp.RequestID, targetChapter.ID)
}

// Scenario 3.1 – Multi-Modal Retrieval With Filters
//
// This chapter shows how to blend multi-modal search with scalar filters to focus on the
// most relevant sessions from our Atlas journey:
//  1. Upsert several themed chapters into the collection.
//  2. Use SearchByMultiModal with a narrative prompt to surface related chapters.
//  3. Apply score/paragraph filters so only the current session's highlights appear.
func TestScenarioIndexSearchMultiModal(t *testing.T) {
	env := requireEnv(t)
	sdk := mustNewClient(t, env)
	collectionClient := sdk.Collection(collectionBase(env))
	indexClient := sdk.Index(indexBase(env))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sessionTag := newSessionTag("index-multimodal")
	baseParagraph := time.Now().UnixNano() % 1_000_000
	chapters := buildStoryChapters(sessionTag, baseParagraph)

	// only one record at a time.
	for _, chapter := range chaptersToUpsert(chapters) {
		upsertReq := model.UpsertDataRequest{
			WriteDataBase: model.WriteDataBase{
				Data: []model.MapStr{chapter},
			},
		}
		upsertResp, err := collectionClient.Upsert(ctx, upsertReq)
		require.NoError(t, err, "upsert failed")
		require.NotNil(t, upsertResp.Result)
		log.Printf("Upsert request_id=%s token_usage=%v chapters_written=%d", upsertResp.RequestID, upsertResp.Result.TokenUsage, len(chapters))
	}

	time.Sleep(3 * time.Second)

	filter := andFilters(
		sessionParagraphBounds(baseParagraph, len(chapters)),
		scoreAtLeastFilter(85.0),
	)

	query := findChapter(t, chapters, "retrieval-lab").Text
	limit := 3
	searchReq := model.SearchByMultiModalRequest{
		Text: &query,
		SearchBase: model.SearchBase{
			RecallBase: model.RecallBase{
				Filter: filter,
			},
			Limit:        &limit,
			OutputFields: []string{"title", "score", "paragraph"},
		},
	}

	d, _ := json.Marshal(searchReq)
	log.Printf("Search query: %s", string(d))

	var searchResp *model.SearchResponse
	var lastErr error
	resp, err := indexClient.SearchByMultiModal(ctx, searchReq)
	if err != nil {
		lastErr = err
		log.Printf("SearchByMultiModal failed: %v", err)
		return
	} else if resp != nil && resp.Result != nil && len(resp.Result.Data) > 0 {
		searchResp = resp
	} else {
		lastErr = nil
		t.Errorf("SearchByMultiModal returned no matches")
		return
	}

	require.NotNilf(t, searchResp, "expected multi-modal search results (last error: %v)", lastErr)
	require.NotNil(t, searchResp.Result)

	for idx, item := range searchResp.Result.Data {
		log.Printf("SearchByMultiModal request_id=%s hit[%d]=id:%v title:%v score:%v paragraph:%v",
			searchResp.RequestID, idx, item.ID, item.Fields["title"], item.Fields["score"], item.Fields["paragraph"])
	}

	require.Equal(t, 2, len(searchResp.Result.Data), "score filter should surface the two advanced chapters")
}

// Scenario 3.2 – Vector Retrieval With Embeddings
//
// This guide mirrors the Python quickstart: write chapters that include a dense vector field,
// generate a query embedding, and retrieve matching chapters with SearchByVector.
func TestScenarioIndexSearchVector(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	env := requireEnv(t)
	env.Collection = "vector"
	env.Index = "vector_index"

	client := mustNewClient(t, env)
	collectionClient := client.Collection(collectionBase(env))
	indexClient := client.Index(indexBase(env))
	embeddingClient := client.Embedding()

	sessionTag := newSessionTag("search-vector")
	baseParagraph := currentParagraphSeed()
	chapters := buildStoryChapters(sessionTag, baseParagraph)

	// 1. Embed the chapter narratives so we can upsert them with a pre-computed dense vector.
	modelName := "bge-m3"
	modelVersion := "default"
	chapterEmbeddings := batchEmbedTexts(t, ctx, embeddingClient, chapters, modelName, modelVersion)

	upsertPayload := make([]model.MapStr, 0, len(chapters))
	for idx, chapter := range chapters {
		upsertPayload = append(upsertPayload, model.MapStr{
			"title":     chapter.Title,
			"paragraph": chapter.Paragraph,
			"score":     chapter.Score,
			"text":      chapter.Text,
			"vector":    chapterEmbeddings[idx],
		})
	}

	upsertReq := model.UpsertDataRequest{
		WriteDataBase: model.WriteDataBase{
			Data: upsertPayload,
		},
	}
	upsertResp, err := collectionClient.Upsert(ctx, upsertReq)
	require.NoError(t, err, "vector upsert failed")
	log.Printf("Upsert request_id=%s", upsertResp.RequestID)

	// 2. Allow the index a moment to surface the newly written vectors.
	time.Sleep(3 * time.Second)

	// 3. Embed the retrieval chapter again to build our query vector.
	targetChapter := findChapter(t, chapters, "retrieval-lab")
	queryVector := embedSingleText(t, ctx, embeddingClient, targetChapter.Text, modelName, modelVersion)

	// 4. Run SearchByVector over the chapter window we just created.
	filter := sessionParagraphBounds(baseParagraph, len(chapters))
	limit := 5
	searchReq := model.SearchByVectorRequest{
		SearchBase: model.SearchBase{
			RecallBase: model.RecallBase{
				Filter: filter,
			},
			Limit:        &limit,
			OutputFields: []string{"title", "score", "paragraph"},
		},
		DenseVector: queryVector,
	}

	searchResp, err := indexClient.SearchByVector(ctx, searchReq)
	require.NoError(t, err, "vector search failed")
	require.NotNil(t, searchResp.Result, "vector search should return results")
	require.Greater(t, len(searchResp.Result.Data), 0, "vector search should produce hits")

	hits := searchResp.Result.Data
	titles := make([]string, 0, len(hits))
	for _, hit := range hits {
		if title, ok := hit.Fields["title"].(string); ok {
			titles = append(titles, title)
		}
	}
	require.Containsf(t, titles, targetChapter.Title, "expected %q to appear in vector search results", targetChapter.Title)
}

// Scenario 3_3 – Search Extensions & Analytics
//
// Building on the Atlas journey, this guide strings together advanced discovery APIs:
//   - SearchByMultiModal keeps session prompts scoped with filters.
//   - SearchByKeywords surfaces content by explicit tags.
//   - SearchByRandom reminds us what's already in rotation.
//   - Aggregate captures score analytics over the current session.
//   - Sort reorders specific chapters using a query vector.
func TestScenarioSearchKeywords(t *testing.T) {
	env := requireEnv(t)
	sdk := mustNewClient(t, env)
	collectionClient := sdk.Collection(collectionBase(env))
	indexClient := sdk.Index(indexBase(env))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sessionTag := newSessionTag("index-extensions")
	baseParagraph := time.Now().UnixNano() % 1_000_000
	chapters := buildStoryChapters(sessionTag, baseParagraph)

	// only one record at a time.
	for _, chapter := range chaptersToUpsert(chapters) {
		upsertReq := model.UpsertDataRequest{
			WriteDataBase: model.WriteDataBase{
				Data: []model.MapStr{chapter},
			},
		}
		upsertResp, err := collectionClient.Upsert(ctx, upsertReq)
		require.NoError(t, err, "upsert failed")
		require.NotNil(t, upsertResp.Result)
		log.Printf("Upsert request_id=%s token_usage=%v chapters_written=%d", upsertResp.RequestID, upsertResp.Result.TokenUsage, len(chapters))
	}

	time.Sleep(5 * time.Second)

	limit := 3
	sessionFilter := sessionParagraphBounds(baseParagraph, len(chapters))

	keywordsReq := model.SearchByKeywordsRequest{
		Keywords: []string{"Signal"},
		SearchBase: model.SearchBase{
			RecallBase: model.RecallBase{
				Filter: sessionFilter,
			},
			Limit:        &limit,
			OutputFields: []string{"title", "score", "text"},
		},
	}
	keyResp, keyErr := indexClient.SearchByKeywords(ctx, keywordsReq)
	require.NoError(t, keyErr, "SearchByKeywords failed")
	require.NotNil(t, keyResp)
	require.NotNil(t, keyResp.Result)
	require.NotEmpty(t, keyResp.Result.Data, "keywords search should surface tagged chapters")
	for _, item := range keyResp.Result.Data {
		title := requireStringField(t, item.Fields, "title")
		marshal, _ := json.Marshal(item)
		log.Printf(string(marshal))
		// search_by_keyword only take effect on rank, possible returned record don't contain keywords.
		// require.Contains(t, title, "Signal", "keywords search should return Signal chapters")
		require.Contains(t, title, sessionTag, "keywords search should remain scoped to this session")
	}
	log.Printf("SearchByKeywords request_id=%s hits=%d", keyResp.RequestID, len(keyResp.Result.Data))
}

// Scenario 4 – Search Aggregations
//
// Building on the Atlas journey, discover counts, group by paragraphs:
//   - Aggregate captures score analytics over the current session.
func TestScenarioSearchExtensionsAndAnalytics(t *testing.T) {
	env := requireEnv(t)
	sdk := mustNewClient(t, env)
	collectionClient := sdk.Collection(collectionBase(env))
	indexClient := sdk.Index(indexBase(env))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sessionTag := newSessionTag("index-extensions")
	baseParagraph := time.Now().UnixNano() % 1_000_000
	chapters := buildStoryChapters(sessionTag, baseParagraph)

	// only one record at a time.
	for _, chapter := range chaptersToUpsert(chapters) {
		upsertReq := model.UpsertDataRequest{
			WriteDataBase: model.WriteDataBase{
				Data: []model.MapStr{chapter},
			},
		}
		upsertResp, err := collectionClient.Upsert(ctx, upsertReq)
		require.NoError(t, err, "upsert failed")
		require.NotNil(t, upsertResp.Result)
		log.Printf("Upsert request_id=%s token_usage=%v chapters_written=%d", upsertResp.RequestID, upsertResp.Result.TokenUsage, len(chapters))
	}

	time.Sleep(3 * time.Second)

	field := "paragraph"
	op := "count"
	aggReq := model.AggRequest{
		Field: &field,
		Op:    op,
		Cond:  model.MapStr{"gt": 1},
	}
	aggResp, aggErr := indexClient.Aggregate(ctx, aggReq)
	require.NoError(t, aggErr, "Aggregate failed")
	require.NotNil(t, aggResp)
	require.NotNil(t, aggResp.Result)
	require.NotNil(t, aggResp.Result.Agg)
	d, _ := json.Marshal(aggResp.Result.Agg)

	require.NotZero(t, len(aggResp.Result.Agg))
	log.Printf("Aggregate request_id=%s agg=%s", aggResp.RequestID, string(d))
}

// Scenario 5 – Embedding Pipeline
//
// This guide demonstrates how to obtain embeddings from VikingDB:
//   - Configure dense (and optional sparse) models.
//   - Provide text/image/video content, including multimodal sequences.
//   - Inspect the response for dense/sparse vectors and token usage metadata.
//
// As with other guides, ensure the models referenced are available in your region/account.
func TestScenarioEmbeddingMultiModalPipeline(t *testing.T) {
	env := requireEnv(t)
	client := mustNewClient(t, env).Embedding()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	text := "generate embeddings with VikingDB"
	modelName := "doubao-embedding-vision"
	modelVersion := "250615"

	bge := "bge-m3"
	req := model.EmbeddingRequest{
		DenseModel: &model.EmbeddingModelOpt{
			ModelName:    &modelName,
			ModelVersion: &modelVersion,
		},
		SparseModel: &model.EmbeddingModelOpt{
			ModelName: &bge,
		},
		Data: []*model.EmbeddingData{
			{
				//Text: &text,
				FullModalSeq: []model.FullModalData{
					{Text: &text},
				},
			},
		},
	}
	d, _ := json.Marshal(req)
	log.Printf(string(d))

	resp, err := client.Embedding(ctx, req)
	require.NoError(t, err, "embedding request failed")
	require.NotNil(t, resp.Result)
	require.NotEmpty(t, resp.Result.Data)
	require.NotEmpty(t, resp.Result.Data[0].DenseVectors)
	//require.NotEmpty(t, resp.Result.Data[0].SparseVectors)
	log.Printf("Embedding request_id=%s dense_dims=%d token_usage=%v", resp.RequestID, len(resp.Result.Data[0].DenseVectors), resp.Result.TokenUsage)
	log.Printf("Dense[:5]=%v, Sparse=%v", resp.Result.Data[0].DenseVectors[:5], resp.Result.Data[0].SparseVectors)
}

func TestScenarioEmbeddingDSPipeline(t *testing.T) {
	env := requireEnv(t)
	client := mustNewClient(t, env).Embedding()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	text := "generate text dense&sparse embeddings with VikingDB"
	modelName := "bge-m3"

	req := model.EmbeddingRequest{
		DenseModel: &model.EmbeddingModelOpt{
			ModelName: &modelName,
			//ModelVersion: &modelVersion, // use default
		},
		SparseModel: &model.EmbeddingModelOpt{
			ModelName: &modelName,
		},
		Data: []*model.EmbeddingData{
			{
				Text: &text,
			},
		},
	}

	resp, err := client.Embedding(ctx, req)
	require.NoError(t, err, "embedding request failed")
	require.NotNil(t, resp.Result)
	require.NotEmpty(t, resp.Result.Data)
	require.NotEmpty(t, resp.Result.Data[0].DenseVectors)
	//require.NotEmpty(t, resp.Result.Data[0].SparseVectors)
	log.Printf("Embedding request_id=%s dense_dims=%d token_usage=%v", resp.RequestID, len(resp.Result.Data[0].DenseVectors), resp.Result.TokenUsage)
	log.Printf("Dense[:5]=%v, Sparse=%v", resp.Result.Data[0].DenseVectors[:5], resp.Result.Data[0].SparseVectors)
}

func batchEmbedTexts(t *testing.T, ctx context.Context, embeddingClient vector.EmbeddingClient, chapters []*storyChapter, modelName, modelVersion string) [][]float64 {
	t.Helper()

	req := model.EmbeddingRequest{
		DenseModel: &model.EmbeddingModelOpt{
			ModelName:    &modelName,
			ModelVersion: &modelVersion,
		},
		Data: make([]*model.EmbeddingData, 0, len(chapters)),
	}
	for _, chapter := range chapters {
		text := chapter.Text
		req.Data = append(req.Data, &model.EmbeddingData{Text: &text})
	}

	resp, err := embeddingClient.Embedding(ctx, req)
	require.NoError(t, err, "embedding batch request failed")
	require.NotNil(t, resp.Result, "embedding batch should return a result")
	require.Len(t, resp.Result.Data, len(chapters), "embedding batch must mirror chapter count")

	out := make([][]float64, len(resp.Result.Data))
	for idx, item := range resp.Result.Data {
		require.NotEmptyf(t, item.DenseVectors, "missing dense vector for chapter %s", chapters[idx].Key)
		out[idx] = float32SliceToFloat64(item.DenseVectors)
	}
	return out
}

func embedSingleText(t *testing.T, ctx context.Context, embeddingClient vector.EmbeddingClient, text, modelName, modelVersion string) []float64 {
	t.Helper()

	req := model.EmbeddingRequest{
		DenseModel: &model.EmbeddingModelOpt{
			ModelName:    &modelName,
			ModelVersion: &modelVersion,
		},
		Data: []*model.EmbeddingData{{Text: &text}},
	}

	resp, err := embeddingClient.Embedding(ctx, req)
	require.NoError(t, err, "query embedding request failed")
	require.NotNil(t, resp.Result, "embedding response should include a result")
	require.NotEmpty(t, resp.Result.Data, "embedding response should include data")
	require.NotEmpty(t, resp.Result.Data[0].DenseVectors, "embedding response should include dense vectors")

	return float32SliceToFloat64(resp.Result.Data[0].DenseVectors)
}

func currentParagraphSeed() int64 {
	return time.Now().Unix() % 1_000_000
}

type guideEnv struct {
	AccessKey  string
	SecretKey  string
	Host       string
	Region     string
	Collection string
	Index      string
}

func requireEnv(t *testing.T) guideEnv {
	t.Helper()
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}
	if os.Getenv("VIKINGDB_AK") == "" {
		_ = godotenv.Load(".env.local")
	}

	env := guideEnv{
		AccessKey:  os.Getenv("VIKINGDB_AK"),
		SecretKey:  os.Getenv("VIKINGDB_SK"),
		Host:       os.Getenv("VIKINGDB_HOST"),
		Region:     os.Getenv("VIKINGDB_REGION"),
		Collection: os.Getenv("VIKINGDB_COLLECTION"),
		Index:      os.Getenv("VIKINGDB_INDEX"),
	}

	missing := []string{}
	if env.AccessKey == "" {
		missing = append(missing, "VIKINGDB_AK")
	}
	if env.SecretKey == "" {
		missing = append(missing, "VIKINGDB_SK")
	}
	if env.Host == "" {
		missing = append(missing, "VIKINGDB_HOST")
	}
	if env.Region == "" {
		missing = append(missing, "VIKINGDB_REGION")
	}
	if env.Collection == "" {
		missing = append(missing, "VIKINGDB_COLLECTION")
	}
	if env.Index == "" {
		missing = append(missing, "VIKINGDB_INDEX")
	}

	if len(missing) > 0 {
		t.Skipf("missing required environment variables: %v", missing)
	}

	return env
}

func mustNewClient(t *testing.T, env guideEnv) *vector.Client {
	t.Helper()

	client, err := vector.New(
		vector.AuthIAM(env.AccessKey, env.SecretKey),
		sharedClientOptions(env)...,
	)
	require.NoError(t, err)
	return client
}

func sharedClientOptions(env guideEnv) []vector.ClientOption {
	// all has default values
	return []vector.ClientOption{
		vector.WithEndpoint(fmt.Sprintf("https://%s", env.Host)),
		vector.WithRegion(env.Region),
		vector.WithMaxRetries(3),
		vector.WithTimeout(30 * time.Second),
		vector.WithUserAgent("vikingdb-go-sdk-guide"),
	}
}

func collectionBase(env guideEnv) model.CollectionLocator {
	return model.CollectionLocator{
		CollectionName: env.Collection,
		ProjectName:    "",
	}
}

func indexBase(env guideEnv) model.IndexLocator {
	return model.IndexLocator{
		CollectionLocator: model.CollectionLocator{
			CollectionName: env.Collection,
		},
		IndexName: env.Index,
	}
}

// storyChapter captures the fields we write for each guide document.
type storyChapter struct {
	Key       string
	Title     string
	Paragraph int64
	Score     float64
	Text      string
	ID        interface{}
}

// buildStoryChapters returns a repeatable set of documents anchored to a unique session tag.
func buildStoryChapters(sessionTag string, baseParagraph int64) []*storyChapter {
	return []*storyChapter{
		{
			Key:       "orientation",
			Title:     fmt.Sprintf("Atlas Orientation · %s", sessionTag),
			Paragraph: baseParagraph,
			Score:     72.5,
			Text: fmt.Sprintf("%s — Orientation covers connectivity checks and environment hygiene before diving into search.",
				sessionTag),
		},
		{
			Key:       "retrieval-lab",
			Title:     fmt.Sprintf("Atlas Retrieval Lab · %s", sessionTag),
			Paragraph: baseParagraph + 1,
			Score:     88.4,
			Text: fmt.Sprintf("%s — Retrieval lab explores multi-modal prompts and scalar filters to focus recommendations.",
				sessionTag),
		},
		{
			Key:       "signal-tuning",
			Title:     fmt.Sprintf("Atlas Signal Tuning · %s", sessionTag),
			Paragraph: baseParagraph + 2,
			Score:     93.1,
			Text: fmt.Sprintf("%s — Signal tuning stresses analytics, aggregations, and reranking to refine the journey.",
				sessionTag),
		},
	}
}

// chaptersToUpsert converts chapters into the payload expected by Upsert/Update calls.
func chaptersToUpsert(chapters []*storyChapter) []model.MapStr {
	payload := make([]model.MapStr, 0, len(chapters))
	for _, chapter := range chapters {
		payload = append(payload, model.MapStr{
			"title":     chapter.Title,
			"paragraph": chapter.Paragraph,
			"score":     chapter.Score,
			"text":      chapter.Text,
		})
	}
	return payload
}

// assignChapterIDsViaSearch resolves backend generated IDs with SearchByMultiModal.
func assignChapterIDsViaSearch(ctx context.Context, t *testing.T, indexClient vector.IndexClient, chapters []*storyChapter, outputFields []string) {
	t.Helper()

	for _, chapter := range chapters {
		hit, requestID := searchChapterByNarrative(ctx, t, indexClient, chapter.Text, outputFields)
		require.NotNilf(t, hit.ID, "SearchByMultiModal returned nil id for chapter %s", chapter.Key)

		chapter.ID = hit.ID
		log.Printf("SearchByMultiModal request_id=%s chapter_key=%s id=%v title=%s score=%v",
			requestID, chapter.Key, chapter.ID, hit.Fields["title"], hit.Fields["score"])
	}
}

// searchChapterByNarrative looks up a chapter using its narrative text and retries until it appears in the index.
func searchChapterByNarrative(ctx context.Context, t *testing.T, indexClient vector.IndexClient, query string, outputFields []string) (model.SearchItemResult, string) {
	t.Helper()

	limit := 3
	nI := false
	req := model.SearchByMultiModalRequest{
		Text:            &query,
		NeedInstruction: &nI,
		SearchBase: model.SearchBase{
			Limit:        &limit,
			OutputFields: outputFields,
		},
	}

	var lastErr error
	for attempt := 0; attempt < 10; attempt++ {
		resp, err := indexClient.SearchByMultiModal(ctx, req)
		if err != nil {
			lastErr = err
			log.Printf("SearchByMultiModal attempt %d failed: %v", attempt+1, err)
		} else if resp != nil && resp.Result != nil && len(resp.Result.Data) > 0 {
			return resp.Result.Data[0], resp.RequestID
		} else {
			lastErr = nil
			log.Printf("SearchByMultiModal attempt %d returned no matches", attempt+1)
		}
		time.Sleep(500 * time.Millisecond)
	}

	require.Failf(t, "search timeout", "unable to locate chapter for query %q (last error: %v)", query, lastErr)
	return model.SearchItemResult{}, ""
}

// newSessionTag generates a session-specific suffix to keep documents unique.
func newSessionTag(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

// findChapter selects a chapter by its logical key.
func findChapter(t *testing.T, chapters []*storyChapter, key string) *storyChapter {
	t.Helper()

	for _, chapter := range chapters {
		if chapter.Key == key {
			return chapter
		}
	}

	require.Failf(t, "chapter missing", "could not locate chapter key %q", key)
	return nil
}

// sessionParagraphBounds returns a simple paragraph range filter to scope the session's documents.
func sessionParagraphBounds(base int64, count int) model.MapStr {
	return model.MapStr{
		"op":    "range",
		"field": "paragraph",
		"gte":   base,
		"lt":    base + int64(count),
	}
}

func mustFilter(field string, conds ...interface{}) model.MapStr {
	return model.MapStr{
		"op":    "must",
		"field": field,
		"conds": conds,
	}
}

func scoreAtLeastFilter(min float64) model.MapStr {
	return model.MapStr{
		"op":    "range",
		"field": "score",
		"gt":    min,
	}
}

func andFilters(conds ...model.MapStr) model.MapStr {
	filtered := make([]model.MapStr, 0, len(conds))
	for _, cond := range conds {
		if cond == nil || len(cond) == 0 {
			continue
		}
		filtered = append(filtered, cond)
	}

	switch len(filtered) {
	case 0:
		return nil
	case 1:
		return filtered[0]
	default:
		return model.MapStr{
			"op":    "and",
			"conds": filtered,
		}
	}
}

// float32SliceToFloat64 converts embedding vectors to the dtype expected by SearchByVector.
func float32SliceToFloat64(src []float32) []float64 {
	out := make([]float64, len(src))
	for i, v := range src {
		out[i] = float64(v)
	}
	return out
}

func requireListField(t *testing.T, fields map[string]interface{}, key string) []interface{} {
	t.Helper()

	raw, ok := fields[key]
	require.Truef(t, ok, "expected field %q in search result", key)

	switch v := raw.(type) {
	case []interface{}:
		return v
	case []string:
		out := make([]interface{}, len(v))
		for i, item := range v {
			out[i] = item
		}
		return out
	case []int:
		out := make([]interface{}, len(v))
		for i, item := range v {
			out[i] = item
		}
		return out
	case []int64:
		out := make([]interface{}, len(v))
		for i, item := range v {
			out[i] = item
		}
		return out
	case []float32:
		out := make([]interface{}, len(v))
		for i, item := range v {
			out[i] = float64(item)
		}
		return out
	case []float64:
		out := make([]interface{}, len(v))
		for i, item := range v {
			out[i] = item
		}
		return out
	default:
		require.Failf(t, "unexpected field type", "field %q expected list but was %T (%v)", key, raw, raw)
		return nil
	}
}

func requireFloat64Field(t *testing.T, fields map[string]interface{}, key string) float64 {
	t.Helper()

	raw, ok := fields[key]
	require.Truef(t, ok, "expected field %q in search result", key)

	switch v := raw.(type) {
	case json.Number:
		vi, _ := v.Int64()
		return float64(vi)
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	default:
		require.Failf(t, "unexpected field type", "field %q has unsupported type %T (%v)", key, raw, raw)
		return 0
	}
}

func requireStringField(t *testing.T, fields map[string]interface{}, key string) string {
	t.Helper()

	raw, ok := fields[key]
	require.Truef(t, ok, "expected field %q in search result", key)

	str, ok := raw.(string)
	require.Truef(t, ok, "field %q expected string but was %T", key, raw)
	return str
}
