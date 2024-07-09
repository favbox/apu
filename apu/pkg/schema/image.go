package schema

type Image struct {
	Source
	Key uint64
	
	OriginalUrl string
	Format      string
	Width       int
	Height      int
}
