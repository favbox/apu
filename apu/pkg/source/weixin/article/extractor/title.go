package extractor

import (
	"regexp"

	"apu/pkg/util/stringx"
)

// 标题替换规则
var (
	titleMap = map[string]string{
		"DEFAULT": `首发\s*[|｜]`,
		"环球设计":    `【环球设计\d+期】|首发\s*\.`,
		"印际":      `\s*印际(首发?)\s*x\s*`,
		"拓者设计吧":   `拓者\s*[|｜]`,
		"LJ看设计":   `\s*LJID\s*•\s*|LJID首发\s*[|｜]`,
		"建E室内设计网": `建E首发\s*[|｜]|新作\s*[|｜]`,
	}
	reTitleMap = map[string]*regexp.Regexp{}
)

func init() {
	for mpName, pattern := range titleMap {
		reTitleMap[mpName] = regexp.MustCompile(pattern)
	}
}

// ExtractTitle 按公众号名称清理标题中的冗余文本。
func ExtractTitle(mpName, title string) string {
	title = stringx.Trim(title)
	if re, ok := reTitleMap[mpName]; ok {
		title = stringx.Replace(title, re, "")
	} else if re, ok := reTitleMap["DEFAULT"]; ok {
		title = stringx.Replace(title, re, "")
	}
	return title
}
