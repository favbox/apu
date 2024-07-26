package embedding

// Embedding 是嵌入模型的接口。
type Embedding interface {
	// EmbedDocuments 嵌入搜索文档。
	EmbedDocuments(texts []string) [][]float32
	// EmbedQuery 嵌入查询文本。
	EmbedQuery(text string) []float32
	// AEmbedDocuments 异步嵌入搜索文档。
	AEmbedDocuments(texts []string) [][]float32
	// AEmbedQuery 异步嵌入查询文本。
	AEmbedQuery(text string) []float32
}

type CacheEmbedding struct {
}
