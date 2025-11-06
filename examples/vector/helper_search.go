// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/volcengine/vikingdb-go-sdk/vector"
	"github.com/volcengine/vikingdb-go-sdk/vector/model"
)

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
		t.Logf("SearchByMultiModal request_id=%s chapter_key=%s id=%v title=%s score=%v",
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
			t.Logf("SearchByMultiModal attempt %d failed: %v", attempt+1, err)
		} else if resp != nil && resp.Result != nil && len(resp.Result.Data) > 0 {
			return resp.Result.Data[0], resp.RequestID
		} else {
			lastErr = nil
			t.Logf("SearchByMultiModal attempt %d returned no matches", attempt+1)
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
