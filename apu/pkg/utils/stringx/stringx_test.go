package stringx

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	n, err := Parse[int32]("123")
	assert.Nil(t, err)
	assert.Equal(t, int32(123), n)

	b, err := Parse[bool]("1")
	assert.Nil(t, err)
	assert.Equal(t, true, b)

	n2 := MustNumber[int]("123")
	assert.Equal(t, 123, n2)

	n3 := MustNumber[int]("138.688px")
	assert.Equal(t, 138, n3)
}

// 提取 style 属性中的 width 值
func ExtractWidthFromStyle(style string) (string, bool) {
	// 正则表达式匹配 width 属性，考虑可能存在的空格
	re := regexp.MustCompile(`(?i)width\s*:\s*(\d+)(px)?`)
	match := re.FindStringSubmatch(style)

	if len(match) < 2 {
		return "", false
	}

	// 返回第一个捕获组，即数字部分
	return match[1], true
}

func TestExtractWidth(t *testing.T) {
	styles := []string{
		"width: 200px; height: 100px;",
		"width :150px; height: 100px;",
		"width:150; height: 80;",
		"display:block; width: 300px;",
		"WIDTH: 150px;", // 测试大小写不敏感
	}

	for _, style := range styles {
		width, exists := ExtractWidthFromStyle(style)
		if exists {
			fmt.Printf("Width extracted from style '%s': %s\n", style, width)
		} else {
			fmt.Printf("Error extracting width from style: %s\n", style)
		}
	}
}
