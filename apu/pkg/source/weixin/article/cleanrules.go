package article

import (
	"regexp"
	"strings"

	"apu/pkg/utils/stringx"
)

// RuleBreakTextsMap 公众号中断行文本
var RuleBreakTextsMap = map[string][]string{
	"DEFAULT":         {"内容策划/Presented", ".END."},
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

var RuleBreakImages = []string{
	"icyksg9whhyvcIb5Dz2Zia2lxuwmELLQ1oPGpOYWoFjR1MaVsiabb78ZloJ9eRyeVDL3mxIRoegwnyiblXeiaHice1tw",
}

func isTextMatched(ruleText string, text string) bool {
	text = normalizeTextForMatch(text)
	if strings.HasPrefix(ruleText, "$$$") {
		ruleText = ruleText[3:]
		if regexp.MustCompile(ruleText).MatchString(text) {
			return true
		}
	}
	return ruleText == text
}

func normalizeTextForMatch(str string) string {
	str = stringx.RemoveAllSpace(str)
	str = strings.ReplaceAll(str, "｜", "|")
	str = strings.ReplaceAll(str, "：", ":")
	str = strings.ReplaceAll(str, "——", "-")
	str = strings.ReplaceAll(str, "—", "-")
	str = strings.ReplaceAll(str, "•", ".")
	return str
}

var RuleRemoveTextsMap = map[string][]string{
	"DEFAULT": {`点击视频`, `对gooood的分享`, `点击蓝字关注我们`, `大家好，我是`, `和邦哥`, `编辑:`, `校对:`, `编辑:`, `校对:`},

	"Architizer Awards": {
		"往期回顾",
		"编辑:",
		"校对:",
	},
	"舶乐汇ONBOX":    {"戳一下关注舶乐汇"},
	"拓者设计吧":       {"延伸阅读"},
	"设计日本":        {"设计日本", "我是专注发现全球"},
	"设计日本+":       {"设计日本", "大家好，我是"},
	"意大利室内设计中文版+": {"意大利室内设计中文版", "我是专注发现全球"},
	"DesignPro":   {"大家好，我是", "PRO哥。"},
	"德国室内设计中文版":   {"德国室内设计中文版", "TopDESIGN", "爱挖掘国内外优秀设计作品的Dave。"},
	"德国室内设计网":     {"DINZ"},
	"DesiDaiIy":   {"DesiDaiIy"},
	"gooood谷德设计网": {"对gooood的分享"},
	"T5OP设计网":     {"微信公众号T5OP设计网"},
	"设计邦": {
		"长的好看，又爱设计的",
		"从LXD离职后，并经历过设计创业失败后",
		"同小伙伴一起经营了一些自媒体账号",
		"分享全球优秀设计工作室的作品",
		"良心干货，超前理念，你想看的全都有",
		"如果你是设计师，就一定要关注哦，",
		"绝对不会让你失望！",
		"点击上方卡片，一起关注",
		"，别犹豫",
	},
	"上海陈壹周室内设计": {"视频讲解"},
	"共合设":       {"空间视频"},
	"TopDESlGN": {"TopDESIGN"},
	"LJ看设计":     {"独家专访", `$$$P.(\d+)`, "$$$LJID-"},
	"建E室内设计网":   {"左右滑动查看平面图"},
}
