package model

// Document 用于存储一段文本和相关的元数据。
type Document struct {
	PageContent string
	// 关于页面内容的任意元数据(例如，来源、与其他文档的关系等)。
	Metadata map[string]any
}

// DocumentTransformer 定义文档转换系统的操作。
//
// 文档转换系统接受一系列文档并返回一个转换后的文档序列。
//
// Example:
//
//	type EmbeddingsRedundantFilter struct {
//		Embeddings          Embeddings
//		SimilarityFn        Callable
//		SimilarityThreshold float64
//	}
//
//	type Config struct {
//		ArbitraryTypesAllowed bool
//	}
//
//	func (erf EmbeddingsRedundantFilter) TransformDocuments(documents []Document, opts...interface{}) []Document {
//		statefulDocuments := getStatefulDocuments(documents)
//		embeddedDocuments := _getEmbeddingsFromStatefulDocs(erf.Embeddings, statefulDocuments)
//		includedIdxs := _filterSimilarEmbeddings(embeddedDocuments, erf.SimilarityFn, erf.SimilarityThreshold)
//		result := make([]Document, 0)
//		for _, i := range sorted(includedIdxs) {
//			result = append(result, statefulDocuments[i])
//		}
//		return result
//	}
type DocumentTransformer interface {
	// TransformDocuments 转换文档列表。
	TransformDocuments(documents []Document, opts ...any) []Document
}
