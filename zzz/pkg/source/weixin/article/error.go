package article

import (
	"bytes"
	"errors"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

var (
	ErrParamsTag = []byte(`<div class="weui-msg__title warn">参数错误</div>`)
	ErrParams    = errors.New("参数错误")

	ErrEnvTag = []byte(`<p class="weui-msg__desc">当前环境异常，完成验证后即可继续访问。</p>`)
	ErrEnv    = errors.New("环境异常需验证")
)

func HasResponseError(body []byte) error {
	if bytes.Contains(body, ErrEnvTag) {
		return ErrEnv
	}
	if bytes.Contains(body, ErrParamsTag) {
		return ErrParams
	}

	return nil
}

func HasURLError(rawURL string) (string, error) {
	if strings.HasPrefix(rawURL, "http://") {
		rawURL = strings.Replace(rawURL, "http://", "https://", 1)
	}
	rawURL = html.UnescapeString(rawURL)
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if parsedURL.Hostname() != "mp.weixin.qq.com" {
		return "", errors.New("主机不是 mp.weixin.qq.com")
	}

	if !strings.HasPrefix(parsedURL.Path, "/s/") && parsedURL.Path != "/s" {
		return "", errors.New("公众号文章网址只支持短链、长链的合法网站")
	}

	return rawURL, nil
}
