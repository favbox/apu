package weixin

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"net/url"
	"strings"

	"apu/pkg/datasource/weixin/article"
	"github.com/imroc/req/v3"
)

var ErrWrongURL = errors.New("网址异常")

// GetArticleStat 获取文章统计信息。
func GetArticleStat(biz, mid, idx, sn string) (*article.Stat, error) {
	return article.GetStat(biz, mid, idx, sn)
}

// GetArticleStatByURL 获取指定文章网址的统计信息。
func GetArticleStatByURL(rawURL string) (*article.Stat, error) {
	// https://mp.weixin.qq.com/s/UQiQocQ2MWkwImZJFMMWQw
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if parsedURL.Hostname() != "mp.weixin.qq.com" {
		return nil, errors.New("主机不是 mp.weixin.qq.com")
	}

	var urlValues url.Values
	if strings.HasPrefix(parsedURL.Path, "/s/") {
		urlValues, err = getShortURLValues(rawURL)
		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(parsedURL.Path, "/s?") {
		urlValues = parsedURL.Query()
	}

	var biz, mid, idx, sn string
	biz = strings.TrimSpace(urlValues.Get("__biz"))
	mid = strings.TrimSpace(urlValues.Get("mid"))
	idx = strings.TrimSpace(urlValues.Get("idx"))
	sn = strings.TrimSpace(urlValues.Get("sn"))
	if biz == "" || mid == "" || idx == "" || sn == "" {
		return nil, errors.New("查询参数不足")
	}

	fmt.Println(biz, mid, idx, sn)
	return nil, nil
}

func getShortURLValues(shortURL string) (values url.Values, err error) {
	resp, err := req.Get(shortURL)
	if err != nil {
		return nil, err
	}
	if resp.IsErrorState() {
		return nil, errors.New(resp.GetStatus())
	}

	body := resp.Bytes()
	ogURLTag := []byte(`<meta property="og:url" content="`)
	if !bytes.Contains(body, ogURLTag) {
		return nil, ErrWrongURL
	}

	startIndex := bytes.Index(body, ogURLTag)
	if startIndex == -1 {
		return nil, ErrWrongURL
	}
	startIndex += len(ogURLTag)
	endIndex := bytes.Index(body[startIndex:], []byte(`"`))
	if endIndex == -1 {
		return nil, ErrWrongURL
	}
	ogURL := string(body[startIndex : startIndex+endIndex])
	ogURL = html.UnescapeString(ogURL)
	if err != nil {
		return nil, err
	}
	parsedURL, err := url.Parse(ogURL)
	if err != nil {
		return nil, err
	}

	return parsedURL.Query(), nil
}
