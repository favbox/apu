package schema

type Source int

var (
	SourceWeixin      Source = 1
	SourceZhimo       Source = 2
	SourceXiaohongshu Source = 3
	SourceBehance     Source = 4
)

func (s Source) Int() int {
	return int(s)
}
