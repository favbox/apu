package schema

import (
	"time"
)

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
