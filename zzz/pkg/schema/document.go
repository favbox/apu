package schema

import (
	"time"
)

// Document 是与文档交互的接口。
type Document struct {
	Source
	UID uint64

	Metadata map[string]any

	PublishTime time.Time
	OriginalUrl string
	Title       string

	Content string
	Images  []*Image
}

var (
	DocumentStageInit       = 0
	DocumentStageStated     = 1
	DocumentStageDetailed   = 2
	DocumentStageStructured = 3
	DocumentStageEmbedded   = 4
)
