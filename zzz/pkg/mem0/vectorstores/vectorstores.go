package vectorstores

import (
	pb "github.com/qdrant/go-client/qdrant"
)

var (
	distance = pb.Distance_Cosine
)

// VectorStorer 定义向量存储器所需实现的方法。
type VectorStorer interface {
	// CreateCol 创建一个新的向量集合。
	CreateCol(name string, vectorSize int, distance string) error
	Insert(name string, vectors []float32, playloads any, ids []int) error
	Search(name string) ([]int, error)
}
