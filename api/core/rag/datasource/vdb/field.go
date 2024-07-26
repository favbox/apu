package vdb

type Field string

const (
	FieldContentKey  Field = "page_content"
	FieldMetadataKey Field = "metadata"
	FieldGroupKey    Field = "group_id"
	FieldVector      Field = "vector"
	FieldTextKey     Field = "text"
	FieldPrimaryKey  Field = "id"
	FieldDocID       Field = "metadata.doc_id"
)
