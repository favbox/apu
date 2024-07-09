package article

import (
	"bytes"
	"encoding/json"
	"errors"

	"apu/pkg/schema"
	"apu/pkg/source/weixin/article/extractor"
	"apu/pkg/store/mysql/query"
	"apu/pkg/utils/cookieutil"
	"apu/pkg/utils/stringx"
	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"
)

const DetailURLPattern = "http://mp.weixin.qq.com/mp/appmsg/show?__biz=%s&appmsgid=%s&itemidx=%s#wechat_redirect"

// GetArticles 利用微信读书 headers 获取公众号的文章列表。
func GetArticles(bookId string, count, offset, syncKey int) (articles []*BookArticle, nextSyncKey int, err error) {
	// 获取微信读书请求头
	//mysql.Init()
	weRequest, err := query.WeRequest.Where(
		query.WeRequest.Type.Eq("weread"),
		query.WeRequest.Status.Eq("valid"),
	).First()
	if err != nil {
		return
	}
	var headers map[string]string
	err = json.Unmarshal([]byte(weRequest.Headers), &headers)
	if err != nil {
		return
	}
	if len(headers) == 0 {
		err = errors.New("无可用请求头")
		return
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
		return
	}
	if resp.IsErrorState() {
		err = errors.New(resp.GetStatus())
		return
	}
	if result.Errmsg != "" {
		err = errors.New(result.Errmsg)
		return
	}

	nextSyncKey = result.SyncKey
	for _, review := range result.Reviews {
		articles = append(articles, review.Review.MpInfo)
	}

	return
}

// GetStat 利用微信 cookie 获取文章统计信息。 https://www.cnblogs.com/jianpansangejian/p/17970546
func GetStat(biz, mid, idx, sn string) (*Stat, error) {
	// 获取微信阅读量请求 cookie
	//mysql.Init()
	weRequest, err := query.WeRequest.Where(
		query.WeRequest.Type.Eq("wechat"),
		query.WeRequest.Status.Eq("valid"),
	).First()
	if err != nil {
		return nil, err
	}

	request := req.R()

	// 设置查询参数
	cookieMap := cookieutil.StrToMap(weRequest.Cookie)
	request.SetQueryParams(map[string]string{
		"appmsg_token": cookieMap["appmsg_token"],
		"x5":           "0",
	})

	// 设置请求头
	request.SetHeaders(map[string]string{
		"User-Agent":   "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0Chrome/57.0.2987.132 MQQBrowser/6.2 Mobile",
		"Cookie":       weRequest.Cookie,
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
	if result.ArticleStat == nil {
		return nil, errors.New("会话已过期")
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
