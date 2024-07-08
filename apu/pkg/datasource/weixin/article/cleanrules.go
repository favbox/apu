package article

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"apu/pkg/utils/stringx"
	"github.com/PuerkitoBio/goquery"
	"github.com/yuin/goldmark/util"
)

// 标题替换规则
var (
	titleMap = map[string]string{
		"DEFAULT": `首发\s*[|｜]`,
		"环球设计":    `【环球设计\d+期】|首发\s*\.`,
		"印际":      `\s*印际(首发?)\s*x\s*`,
		"拓者设计吧":   `拓者\s*[|｜]`,
		"LJ看设计":   `\s*LJID\s*•\s*`,
		"建E室内设计网": `建E首发\s*[|｜]|新作\s*[|｜]`,
	}
	reTitleMap = map[string]*regexp.Regexp{}
)

// 时间提取规则
var reCreateTime = regexp.MustCompile(`var oriCreateTime = '(\d+)';`)

var (
	breakTextsMap = map[string][]string{
		"DEFAULT":         {"内容策划/Presented"},
		"室内设计联盟网":         {"延伸阅读"},
		"拓者设计吧":           {"延伸阅读"},
		"安邸AD":            {"撰文|", "视觉设计|", "新媒体编辑|"},
		"德国室内设计网":         {"DINZ视频号"},
		"易高国际装饰":          {"更多案例推荐"},
		"DesignBest":      {"内容策划/Presented"},
		"印际":              {"内容策划/Presented"},
		"设计日本":            {"内容策划/Presented"},
		"意大利室内设计中文版+":     {"内容策划/Presented"},
		"T5OP设计网":         {"排版:微信公众号T5OP设计网"},
		"设计邦":             {"图源:骑驴小宅", "主题参考:新微设计", "设计类公众号"},
		"Hi设计":            {"|战略合作伙伴|", "|特别策划|"},
		"南京观享际SKH室内设计":    {"MOREABOUTINFORMATION"},
		"LJ看设计":           {"REVIEW"},
		"上海陈壹周室内设计":       {"转载请注明出处，盗用必究"},
		"时尚办公网popoffices": {"这是POP君分享的"},
		"环球设计":            {"推荐阅读"},
		"ELLEDECO家居廊":     {"本文为《ELLEDECORATION家居廊》版权所有"},
		"gooood谷德设计网":     {"禁止以gooood编辑版本进行任何形式转载"},
		"建E室内设计网":         {"Previousprojects"},
		"环球商业设计":          {"-END-"},
	}
)

func init() {
	for mpName, pattern := range titleMap {
		reTitleMap[mpName] = regexp.MustCompile(pattern)
	}
}

func CleanTitle(info *Info, title string) {
	title = stringx.Trim(title)
	if re, ok := reTitleMap[info.MpName]; ok {
		title = stringx.Replace(title, re, "")
	} else if re, ok := reTitleMap["DEFAULT"]; ok {
		title = stringx.Replace(title, re, "")
	}
	info.Title = title
}

func CleanPublishTime(info *Info, body []byte) {
	if v := reCreateTime.FindSubmatch(body); len(v) == 2 {
		timestamp := util.BytesToReadOnlyString(v[1])
		if sec, err := strconv.ParseFloat(timestamp, 64); err == nil {
			info.PublishTime = time.Unix(int64(sec), 0)
		}
	}
}

func CleanJsContent(info *Info, body []byte, doc *goquery.Document) {
	breakIndex := -1 // 裁剪内容的中断标志位

	jsContent := doc.Find("#js_content")
	jsContent.Find("*").Each(func(i int, s *goquery.Selection) {
		// 移除中断标志位后的所有节点
		if i > breakIndex && breakIndex > 0 {
			s.Remove()
			return
		}

		//fmt.Println("s.Nodes[0].Type", s.Nodes[0].Data)
		switch s.Nodes[0].Data {
		case "a":
			// 移除链接的 href 和 target 属性
			s.RemoveAttr("href").RemoveAttr("target")
			return
		case "svg":
			// 移除 svg
			s.Remove()
			return
		case "img":
			fmt.Println("移除冗余图")
			if src, exists := s.Attr("data-src"); exists {
				s.SetAttr("src", src).RemoveAttr("data-src")
			}
			return
		case "section":
			fmt.Println("合并连续 span")
		case "p":
			fmt.Println("合并连续 p")
		default:
			fmt.Println("其他节点")
			text := cleanTextForMatch(s.Text())
			if len(text) == 0 {
				return
			}
			// 移除超短文本行
			if utf8.RuneCountInString(text) == 1 {
				s.Remove()
				return
			}

			// 处理50个字符以内的文本是否包含中断标志位
			if isBreakText(text, info.MpName) {
				breakIndex = i
				s.Remove()
				RemoveNode(jsContent.Get(0), s.Get(0))
				return
			}

		}
	})

	fmt.Println(jsContent.Text())
}

func isBreakText(text, mpName string) bool {
	if breakTexts, ok := breakTextsMap[mpName]; ok {
		for _, breakText := range breakTexts {
			return strings.Contains(text, breakText)
		}
	} else if breakTexts, ok := breakTextsMap["DEFAULT"]; ok {
		for _, breakText := range breakTexts {
			return strings.Contains(text, breakText)
		}
	}

	return false
}

// cleanTextForMatch 去除所有空格。
func cleanTextForMatch(str string) string {
	str = html.UnescapeString(str)
	str = strings.ReplaceAll(str, "｜", "|")
	str = strings.ReplaceAll(str, "：", ":")
	str = stringx.RemoveAllSpace(str)
	str = stringx.Trim(str)

	return str
}
