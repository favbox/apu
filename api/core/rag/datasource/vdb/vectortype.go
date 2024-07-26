package vdb

type VectorType string

const (
	VectorTypeAnalyticDB VectorType = "analyticdb"
	VectorTypeChroma     VectorType = "chroma"
	VectorTypeMilvus     VectorType = "milvus"
	VectorTypeMyScale    VectorType = "myscale"
	VectorTypePgVector   VectorType = "pgvector"
	VectorTypePgVectoRs  VectorType = "pgvecto-rs"
	VectorTypeQdrant     VectorType = "qdrant"
	VectorTypeRelyt      VectorType = "relyt"
	VectorTypeTiDBVector VectorType = "tidb_vector"
	VectorTypeWeaviate   VectorType = "weaviate"
	VectorTypeOpenSearch VectorType = "opensearch"
	VectorTypeTencent    VectorType = "tencent"
	VectorTypeOracle     VectorType = "oracle"
)
