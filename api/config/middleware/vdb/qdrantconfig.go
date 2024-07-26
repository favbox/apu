package vdb

type QdrantConfig struct {
	QdrantUrl           string
	QdrantApiKey        string
	QdrantClientTimeout int
	QdrantGrpcEnabled   bool
	QdrantGrpcPort      int
}

func NewQdrantConfig() *QdrantConfig {
	return &QdrantConfig{
		QdrantUrl:           "",
		QdrantApiKey:        "",
		QdrantClientTimeout: 20,
		QdrantGrpcEnabled:   false,
		QdrantGrpcPort:      6334,
	}
}
