package schema

import (
	"time"
)

// Document 是与文档交互的接口。
type Document struct {
	Source
	Key uint64

	Author      string
	PublishTime time.Time
	OriginalUrl string
	Title       string
	Content     string

	Images []*Image
}

var (
	DocumentStageInit       int32 = 0
	DocumentStageStated     int32 = 1
	DocumentStageDetailed   int32 = 2
	DocumentStageStructured int32 = 3
	DocumentStageEmbedded   int32 = 4
)
