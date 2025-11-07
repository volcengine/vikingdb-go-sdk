package main

func main() {
	Connectivity()
	CollectionLifecycle()
	IndexSearchMultiModal()
	IndexSearchVector("vector", "vector_index")
	IndexSearchKeywords()
	IndexSearchAggregate()
	EmbeddingMultiModal()
	EmbeddingDenseSparse()
}

func intPtr(v int) *int {
	return &v
}

func boolPtr(v bool) *bool {
	return &v
}

func stringPtr(v string) *string {
	return &v
}

func float64Ptr(v float64) *float64 {
	return &v
}
