package schema

type Image struct {
	Source
	Key uint64

	OriginalUrl string
	Width       int
	Height      int
}
