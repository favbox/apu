package vdb

import (
	"apu/core/rag/model"
)

type Interface interface {
	GetType() string
	Create(texts []model.Document, embeddings [][]float32, opts ...any)
	AddTexts(text []model.Document, embeddings [][]float32, opts ...any)
	TextExists(id string) bool
	DeleteByIDs(ids []string)
	GetIDsByMetadataField(key, value string) []string
	SearchByVector(queryVector []float32, opts ...any) []model.Document
	SearchByFullText(query string, opts ...any) []model.Document
	Delete()
}

type BaseVector struct {
	collectionName string
}

func (b *BaseVector) GetType() string {
	//TODO implement me
	panic("implement me")
}

func (b *BaseVector) Create(texts []model.Document, embeddings [][]float32, opts ...any) {
	//TODO implement me
	panic("implement me")
}

func (b *BaseVector) AddTexts(text []model.Document, embeddings [][]float32, opts ...any) {
	//TODO implement me
	panic("implement me")
}

func (b *BaseVector) TextExists(id string) bool {
	//TODO implement me
	panic("implement me")
}

func (b *BaseVector) DeleteByIDs(ids []string) {
	//TODO implement me
	panic("implement me")
}

func (b *BaseVector) GetIDsByMetadataField(key, value string) []string {
	//TODO implement me
	panic("implement me")
}

func (b *BaseVector) SearchByVector(queryVector []float32, opts ...any) []model.Document {
	//TODO implement me
	panic("implement me")
}

func (b *BaseVector) SearchByFullText(query string, opts ...any) []model.Document {
	//TODO implement me
	panic("implement me")
}

func (b *BaseVector) Delete() {
	//TODO implement me
	panic("implement me")
}

// FilterDuplicateTexts 过滤重复文档。
func (b *BaseVector) filterDuplicateTexts(texts []model.Document) []model.Document {
	result := make([]model.Document, 0)
	for _, text := range texts {
		docID := text.Metadata["doc_id"].(string)
		if !b.TextExists(docID) {
			result = append(result, text)
		}
	}
	return result
}

func (b *BaseVector) getUUIDs(texts []model.Document) []string {
	uuids := make([]string, 0)
	for _, text := range texts {
		uuids = append(uuids, text.Metadata["doc_id"].(string))
	}
	return uuids
}

func (b *BaseVector) CollectionName() string {
	return b.collectionName
}
