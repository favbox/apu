package article

import (
	"encoding/json"
	"errors"
	"html"
	"strings"

	"apu/internal/cookieutil"
	"apu/pkg/store/mysql/query"
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
	biz, mid, idx, sn, err := GetArticleParams(rawURL)
	if err != nil {
		return nil, err
	}
	if biz == "" || mid == "" || idx == "" || sn == "" {
		return nil, errors.New("查询参数不足")
	}

	return nil, nil
}

// GetArticle 获取公开的公众号文章详情。
func GetArticle(rawURL string) (*Info, error) {
	// http://mp.weixin.qq.com/mp/appmsg/show?__biz=MjM5ODIyMTE0MA==&amp;appmsgid=10000382&amp;itemidx=1#wechat_redirect
	// http://mp.weixin.qq.com/s?__biz=MzA5ODEzMjIyMA==&amp;mid=2247713279&amp;idx=1&amp;sn=bd67c1aba187bc8833f4aee60d8a0e90&amp;chksm=909b886ca7ec017a4c72af5460dcdfe672a74450da0cec34f99cab825d512ab37e4c8715eda6#rd

	if strings.HasPrefix(rawURL, "http://") {
		rawURL = strings.Replace(rawURL, "http://", "https://", 1)
	}
	rawURL = html.UnescapeString(rawURL)
	if err := CheckURLValid(rawURL); err != nil {
		return nil, err
	}

	resp := req.MustGet(rawURL)
	defer resp.Body.Close()
	//r.SetHeaders(map[string]string{
	//	"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	//	//"Accept-Encoding": "gzip, deflate, br, zstd",
	//	"Cache-Control": "max-age=0",
	//	//"Cookie":        "rewardsn=; wxtokenkey=777",
	//	"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
	//})
	//
	//// 设置请求参数
	////r.SetQueryParams(map[string]string{
	////	"__biz": biz,
	////	"mid":   mid,
	////	"idx":   idx,
	////	"sn":    sn,
	////})
	//
	//// 发起请求
	//resp, err := r.Get(rawURL)
	//if err != nil {
	//	return nil, err
	//}
	if resp.IsErrorState() {
		return nil, errors.New(resp.GetStatus())
	}

	body := resp.Bytes()
	if e := CheckResponseError(body); e != nil {
		return nil, e
	}

	// 清理文章
	article, err := CleanArticle(body)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()

	return article, nil
}
