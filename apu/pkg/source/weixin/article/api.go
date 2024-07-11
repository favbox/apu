package article

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"apu/pkg/schema"
	"apu/pkg/source/weixin/article/extractor"
	"apu/pkg/store/mysql/query"
	"apu/pkg/utils/cookieutil"
	"apu/pkg/utils/stringx"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/bytedance/gopkg/util/xxhash3"
	"github.com/imroc/req/v3"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	"go.uber.org/ratelimit"
)

var limiterBySecond1 = ratelimit.New(1)

// GetArticles 利用微信读书 headers 获取公众号的文章列表。
func GetArticles(bookId string, count, offset, syncKey int) ([]*schema.Document, int, error) {
	limiterBySecond1.Take()

	// 获取微信读书请求头
	//mysql.Init()
	wexinRequest, err := query.WexinRequest.Where(
		query.WexinRequest.Type.Eq("weread"),
		query.WexinRequest.Status.Eq("valid"),
	).First()
	if err != nil {
		return nil, 0, err
	}
	var headers map[string]string
	err = json.Unmarshal([]byte(wexinRequest.Headers), &headers)
	if err != nil {
		return nil, 0, err
	}
	if len(headers) == 0 {
		return nil, 0, errors.New("无可用请求头")
	}

	request := req.R()

	// 设置查询参数
	request.SetQueryParamsAnyType(map[string]any{
		"bookId":  bookId,
		"count":   count,
		"offset":  offset,
		"synckey": syncKey,
	})

	// 设置请求头
	request.SetHeadersNonCanonical(headers)

	// 设置响应结果
	var result ArticlesResult
	request.SetSuccessResult(&result)

	// 发起请求
	resp, err := request.Get("https://i.weread.qq.com/book/articles")
	if err != nil {
		return nil, 0, err
	}
	if resp.IsErrorState() {
		return nil, 0, errors.New(resp.GetStatus())
	}
	if result.Errmsg != "" {
		return nil, 0, errors.New(result.Errmsg)
	}

	var articles []*schema.Document
	for _, review := range result.Reviews {
		a := review.Review.MpInfo
		keyInfo, err := GetKeyInfo(a.DocUrl)
		if err != nil {
			return nil, 0, err
		}
		articles = append(articles, &schema.Document{
			Source:      schema.Weixin,
			Key:         keyInfo.Key,
			Author:      keyInfo.Biz,
			PublishTime: time.Unix(a.Time, 0),
			OriginalUrl: a.DocUrl,
			Title:       a.Title,
			Content:     a.Content,
		})
	}

	return articles, result.SyncKey, nil
}

// var lastStatTime time.Time
var limiterBySecond3 = ratelimit.New(1, ratelimit.Per(3*time.Second))

// GetStat 利用微信 cookie 获取文章统计信息。 https://www.cnblogs.com/jianpansangejian/p/17970546
func GetStat(biz, mid, idx, sn string) (*Stat, error) {
	//if !lastStatTime.IsZero() {
	//	duration := time.Since(lastStatTime)
	//	minDuration := 3 * time.Second
	//	if duration < minDuration {
	//		time.Sleep(minDuration - duration)
	//	}
	//}
	//lastStatTime = time.Now()
	limiterBySecond3.Take()

	// 获取微信阅读量请求 cookie
	//mysql.Init()
	wexinRequest, err := query.WexinRequest.Where(
		query.WexinRequest.Type.Eq("wechat"),
		query.WexinRequest.Status.Eq("valid"),
	).First()
	if err != nil {
		return nil, err
	}

	request := req.R()

	// 设置查询参数
	cookieMap := cookieutil.StrToMap(wexinRequest.Cookie)
	request.SetQueryParams(map[string]string{
		"appmsg_token": cookieMap["appmsg_token"],
		"x5":           "0",
	})

	// 设置请求头
	request.SetHeaders(map[string]string{
		"User-Agent":   "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0Chrome/57.0.2987.132 MQQBrowser/6.2 Mobile",
		"Cookie":       wexinRequest.Cookie,
		"Origin":       "https://mp.weixin.qq.com",
		"Content-Type": "application/x-www-form-urlencoded; charset=utf-8",
		"Host":         "mp.weixin.qq.com",
	})

	// 设置表单数据
	request.SetFormData(map[string]string{
		"is_only_read": "1",
		"is_temp_url":  "0",
		"appmsg_type":  "9",
		"__biz":        biz,
		"mid":          mid,
		"idx":          idx,
		"sn":           sn,
	})

	// 设置响应结果
	var result StatResult
	request.SetSuccessResult(&result)

	// 发起请求
	resp, err := request.Post("https://mp.weixin.qq.com/mp/getappmsgext")
	if err != nil {
		return nil, err
	}
	if resp.IsErrorState() {
		return nil, errors.New(resp.GetStatus())
	}
	if result.BaseResp.Ret == 301 {
		return nil, errors.New("基础响应码为 301")
	}
	if result.ArticleStat == nil {
		return nil, errors.New("当前无法获取阅读量")
	}

	return result.ArticleStat, nil
}

func GetStatByURL(rawURL string) (*Stat, error) {
	k, err := GetKeyInfo(rawURL)
	if err != nil {
		return nil, err
	}

	return GetStat(k.Biz, k.Mid, k.Idx, k.Sn)
}

// GetArticle 获取公开的公众号文章详情。
func GetArticle(rawURL string) (*schema.Document, error) {
	var err error

	// 检测网址异常
	rawURL, err = HasURLError(rawURL)
	if err != nil {
		return nil, err
	}

	resp, err := req.Get(rawURL)

	// 检测请求异常
	if err != nil {
		return nil, err
	}
	if resp.IsErrorState() {
		return nil, errors.New(resp.GetStatus())
	}

	// 检测响应体异常
	body, err := resp.ToBytes()
	if err != nil {
		return nil, err
	}
	if err := HasResponseError(body); err != nil {
		return nil, err
	}

	// 创建为 goquery 文档
	gq, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// 提取关键参数
	link := gq.Find("meta[property='og:url']").AttrOr("content", "")
	if link == "" {
		return nil, errors.New("无法找到 og:url")
	}
	keyInfo, err := GetKeyInfo(link)
	if err != nil {
		return nil, err
	}

	// 构建初始文档
	article := &schema.Document{
		Source:      schema.Weixin,
		Key:         keyInfo.Key,
		Author:      keyInfo.Biz,
		OriginalUrl: keyInfo.Url,
	}

	// 提取发布时间
	if publishTime, ok := extractor.ExtractPublishTime(body); ok {
		article.PublishTime = publishTime
	}

	// 提取标题
	mpName := stringx.Trim(gq.Find("#js_name").Text())
	article.Title = extractor.ExtractTitle(
		mpName,
		gq.Find("meta[property='og:title']").AttrOr("content", ""),
	)

	// 提取图片列表
	images, imageSizeMap, err := extractor.ExtractImages(body)
	if err != nil {
		return nil, err
	}
	article.Images = images

	// 提取正文
	var (
		imageCount = 0
		breakIndex = -1 // 裁剪线的索引
		//modules    []string
	)
	jsContent := gq.Find("#js_content")
	jsContent.Find("*").Each(func(i int, s *goquery.Selection) {
		// 移除裁剪线以后的所有元素
		if breakIndex > 0 && i > breakIndex {
			s.Remove()
			return
		}

		// 移除所有链接🔗的 href 和 target
		if s.Is("a") {
			// 移除内部链接
			if v, exists := s.Attr("tab"); exists && v == "innerlink" {
				s.Remove()
				return
			}
			s.RemoveAttr("target").RemoveAttr("href")
			return
		}

		// 移除 svg
		if s.Is("svg") {
			s.Remove()
			return
		}

		// 移除冗余图片
		if s.Is("img") {
			// 移除跳转链接的图片
			if _, exists := s.Parent().Attr("js_jump_icon"); exists {
				s.Remove()
				return
			}

			// 重置 src 以便 markdown 正确解析
			imgSrc := s.AttrOr("data-src", s.AttrOr("src", ""))
			s.SetAttr("src", imgSrc).RemoveAttr("data-src")

			if strings.Contains(imgSrc, "icyksg9whhyvcIb5Dz2Zia2lxuwmELLQ1oPGpOYWoFjR1MaVsiabb78ZloJ9eRyeVDL3mxIRoegwnyiblXeiaHice1tw") {
				fmt.Println()
			}
			// 判断是否为图片中断标志位
			if isBreakImage(imgSrc) {
				breakIndex = i
				s.Remove()
				return
			}

			// 删除过小的图片
			imgKey := xxhash3.HashString(imgSrc)
			imgSize := imageSizeMap[imgKey]
			if isSmallImage(s, imgSize) {
				s.Remove()
				return
			}

			//modules = append(modules, fmt.Sprintf(`<img src="%s" />`, imgSrc))

			imageCount++
			return
		}

		text := stringx.Trim(s.Text())

		// 移除单个字符的文本行
		textNum := utf8.RuneCountInString(text)
		if textNum == 1 {
			s.Remove()
			return
		}

		// 判断是否为中断标志位
		if textNum > 1 && textNum < 50 {
			if isBreakTextLine(mpName, text) {
				breakIndex = i
				s.Remove()
				return
			}
			if isRemovableTextLine(mpName, text) {
				s.Remove()
				return
			}
		}
	})

	// goquery -> markdown
	start := time.Now()
	converter := md.NewConverter("", true, nil)
	mdContent := converter.Convert(jsContent)
	log.Debug().Dur("耗时", time.Since(start)).Msg("goquery -> markdown")
	//article.Content = jsContent.Text()

	// markdown -> html buffer
	start = time.Now()
	buf := bytes.NewBuffer(nil)
	err = goldmark.Convert([]byte(mdContent), buf)
	log.Debug().Dur("耗时", time.Since(start)).Msg("markdown -> html")
	if err != nil {
		return nil, err
	}

	article.Content = buf.String()
	//article.Content, _ = jsContent.Html()
	//article.Content = strings.Join(modules, "<br />")

	// 保存测试文件
	saveTestHtml(article)

	return article, nil
}

func isBreakImage(src string) bool {
	for _, breakImageKey := range RuleBreakImages {
		if strings.Contains(src, breakImageKey) {
			return true
		}
	}
	return false
}

func isRemovableTextLine(mpName, text string) bool {
	removableTexts := RuleRemoveTextsMap["DEFAULT"]
	if vs, exists := RuleRemoveTextsMap[mpName]; exists {
		removableTexts = append(removableTexts, vs...)
	}

	for _, removableText := range removableTexts {
		if isTextMatched(removableText, text) {
			return true
		}
	}

	return false
}

func isBreakTextLine(mpName, text string) bool {
	breakTexts := RuleBreakTextsMap["DEFAULT"]
	if vs, exists := RuleBreakTextsMap[mpName]; exists {
		breakTexts = append(breakTexts, vs...)
	}

	for _, breakText := range breakTexts {
		if isTextMatched(breakText, text) {
			return true
		}
	}

	return false
}

func isSmallImage(s *goquery.Selection, size [2]int) bool {
	minSize := 300
	minOriginalSize := 600
	var w, h int

	// 原图宽度【或者】高度过小
	w, h = size[0], size[1]
	if w > 0 && w < minOriginalSize || h > 0 && h < minOriginalSize {
		return true
	}

	// 展示宽度过小
	w = extractWidthFromStyle(s.AttrOr("style", ""))
	if w > 0 && w < minSize {
		return true
	}

	// 候补宽度【或者】高度过小
	w = stringx.MustNumber[int](s.AttrOr("data-backw", "0"))
	h = stringx.MustNumber[int](s.AttrOr("data-backh", "0"))
	if w > 0 && w < minSize || h > 0 && h < minSize {
		return true
	}

	return false
}

// 提取 style 属性中的 width 值
func extractWidthFromStyle(style string) int {
	// 正则表达式匹配 width 属性，考虑可能存在的空格
	re := regexp.MustCompile(`(?i)width\s*:\s*(\d+)(px|%)?`)
	match := re.FindStringSubmatch(style)

	if strings.Contains(style, "%s") {
		fmt.Println()
	}
	if len(match) < 2 {
		return 0
	}

	if hasPercent := len(match) > 2 && match[2] == "%"; hasPercent {
		return 0
	}

	// 返回第一个捕获组，即数字部分
	return stringx.MustNumber[int](match[1])
}

func saveTestHtml(a *schema.Document) {
	html := fmt.Sprintf(`
<style>
* {
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    box-sizing: border-box;
}
body {
    width: 700px;
	margin: 0 auto;
    padding: 20px;
    background-color: antiquewhite;
}
.title {
    font-size: 22px;
    line-height: 1.4;
    margin-bottom: 14px;
    font-weight: 500;
}
.text, .image {
    position: relative;
    margin-top: 40px;
    font-size: 16px;
    font-family: SF Pro Display,-apple-system,BlinkMacSystemFont,Segoe UI,PingFang SC,Hiragino Sans GB,Microsoft YaHei,Helvetica Neue,Helvetica,Arial,sans-serif;
    font-weight: 400;
    color: #333;
    line-height: 40px;
}
.image img {
	max-width: 700px;
}
</style>
<h1 class="title">%s</h1>
%s`, a.Title, a.Content)
	_ = os.WriteFile("test.html", []byte(html), 0666)
}
