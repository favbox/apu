package article

import (
	"errors"
	"fmt"
	"html"
	"net/url"
	"strings"

	"github.com/bytedance/gopkg/util/xxhash3"
)

type KeyInfo struct {
	Biz string
	Mid string
	Idx string
	Sn  string

	Key uint64
	Url string
}

// GetKeyInfo 获取指定文章网址的 KeyInfo。网址必须为经典格式。
func GetKeyInfo(canonicalURL string) (*KeyInfo, error) {
	canonicalURL = html.UnescapeString(canonicalURL)
	if strings.HasPrefix(canonicalURL, "http://") {
		canonicalURL = strings.Replace(canonicalURL, "http://", "https://", 1)
	}

	parsedURL, err := url.Parse(canonicalURL)
	if err != nil {
		return nil, err
	}

	if parsedURL.Hostname() != "mp.weixin.qq.com" {
		return nil, fmt.Errorf("不是常规的公众号域名")
	}
	if parsedURL.Path != "/s" {
		return nil, errors.New("不是常规的公众号文章网址路径，应为 /s")
	}

	query := parsedURL.Query()
	keyInfo := &KeyInfo{
		Biz: query.Get("__biz"),
		Mid: query.Get("mid"),
		Idx: query.Get("idx"),
		Sn:  query.Get("sn"),
		Url: canonicalURL,
	}
	if keyInfo.Biz == "" || keyInfo.Mid == "" || keyInfo.Idx == "" {
		return nil, errors.New("图文三元组不可为空")
	}

	// 生成唯一键
	keyInfo.Key = xxhash3.HashString(fmt.Sprintf("%s:%s:%s", keyInfo.Biz, keyInfo.Mid, keyInfo.Idx))

	return keyInfo, nil
}
