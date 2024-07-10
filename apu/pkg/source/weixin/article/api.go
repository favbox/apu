package article

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"apu/pkg/schema"
	"apu/pkg/source/weixin/article/extractor"
	"apu/pkg/store/mysql/query"
	"apu/pkg/utils/cookieutil"
	"apu/pkg/utils/stringx"
	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"
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

	// 提取正文
	article.Content = gq.Find("#js_content").Text()

	// 提取图片列表
	images, err := extractor.ExtractImages(body)
	if err != nil {
		return nil, err
	}
	article.Images = images

	return article, nil
}
