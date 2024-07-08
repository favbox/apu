package article

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/imroc/req/v3"
	"golang.org/x/net/html"
)

var ErrWrongOgURL = errors.New("未能获取短网址对应的 og:url")

var (
	ErrParamsTag = []byte(`<div class="weui-msg__title warn">参数错误</div>`)
	ErrParams    = errors.New("参数错误")

	ErrEnvTag = []byte(`<p class="weui-msg__desc">当前环境异常，完成验证后即可继续访问。</p>`)
	ErrEnv    = errors.New("环境异常需验证")
)

// GetArticleParams 获取公众号文章的图文信息四元组。
func GetArticleParams(rawURL string) (biz, mid, idx, sn string, err error) {
	rawURL = html.UnescapeString(rawURL)
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", "", "", "", err
	}
	if parsedURL.Hostname() != "mp.weixin.qq.com" {
		return "", "", "", "", errors.New("主机不是 mp.weixin.qq.com")
	}

	var urlValues url.Values
	if strings.HasPrefix(parsedURL.Path, "/s/") {
		urlValues, err = GetShortURLValues(rawURL)
		if err != nil {
			return "", "", "", "", err
		}
	} else if parsedURL.Path == "/s" {
		urlValues = parsedURL.Query()
	} else {
		return "", "", "", "", errors.New("文章路径错误")
	}

	biz = strings.TrimSpace(urlValues.Get("__biz"))
	mid = strings.TrimSpace(urlValues.Get("mid"))
	idx = strings.TrimSpace(urlValues.Get("idx"))
	sn = strings.TrimSpace(urlValues.Get("sn"))
	if biz == "" || mid == "" || idx == "" || sn == "" {
		return "", "", "", "", errors.New("查询参数不足")
	}

	return biz, mid, idx, sn, nil
}

func GetShortURLValues(shortURL string) (values url.Values, err error) {
	resp, err := req.C().R().Get(shortURL)
	if err != nil {
		return nil, err
	}
	if resp.IsErrorState() {
		return nil, errors.New(resp.GetStatus())
	}

	body := resp.Bytes()

	// 检测响应异常
	if err := CheckResponseError(body); err != nil {
		return nil, err
	}

	ogURLTag := []byte(`<meta property="og:url" content="`)
	if !bytes.Contains(body, ogURLTag) {
		return nil, ErrWrongOgURL
	}

	startIndex := bytes.Index(body, ogURLTag)
	if startIndex == -1 {
		return nil, ErrWrongOgURL
	}
	startIndex += len(ogURLTag)
	endIndex := bytes.Index(body[startIndex:], []byte(`"`))
	if endIndex == -1 {
		return nil, ErrWrongOgURL
	}
	ogURL := string(body[startIndex : startIndex+endIndex])

	return GetURLValues(ogURL)
}

func GetURLValues(rawURL string) (values url.Values, err error) {
	rawURL = html.UnescapeString(rawURL)
	if err != nil {
		return nil, err
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	return parsedURL.Query(), nil
}

func CheckResponseError(body []byte) error {
	if bytes.Contains(body, ErrEnvTag) {
		return ErrEnv
	}
	if bytes.Contains(body, ErrParamsTag) {
		return ErrParams
	}

	return nil
}

func CheckURLValid(rawURL string) error {
	rawURL = html.UnescapeString(rawURL)
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	if parsedURL.Hostname() != "mp.weixin.qq.com" {
		return errors.New("主机不是 mp.weixin.qq.com")
	}

	if !strings.HasPrefix(parsedURL.Path, "/s/") && parsedURL.Path != "/s" {
		return errors.New("公众号文章网址只支持短链、长链的合法网站")
	}

	return nil
}

// RemoveNode Searching node siblings (and child siblings and so on) and after successfull found - remove it
func RemoveNode(rootNode *html.Node, removeMe *html.Node) {
	foundNode := false
	checkNodes := make(map[int]*html.Node)
	i := 0

	// loop through siblings
	for n := rootNode.FirstChild; n != nil; n = n.NextSibling {
		if n == removeMe {
			foundNode = true
			n.Parent.RemoveChild(n)
		}

		checkNodes[i] = n
		i++
	}

	// check if removing node is found
	// if yes no need to check children returning
	// if no continue loop through children and so on
	if foundNode == false {
		for _, item := range checkNodes {
			RemoveNode(item, removeMe)
		}
	}
}

func TraverseNodes(node *html.Node, fn func(*html.Node)) {
	if node == nil {
		return
	}

	fn(node)

	cur := node.FirstChild

	for cur != nil {
		next := cur.NextSibling
		TraverseNodes(cur, fn)
		cur = next
	}
}

func TraverseParentNodes(node *html.Node, targetData string, fn func(*html.Node)) {
	if node == nil {
		return
	}

	cur := node.Parent
	if cur.DataAtom.String() == targetData {
		fn(cur)
		return
	}

	for cur != nil {
		next := cur.Parent
		TraverseParentNodes(cur, targetData, fn)
		cur = next
	}
}

// GetPropValue 从样式字符串中提取宽度的像素值。
func GetPropValue(style, prop string, unit ...string) (px float64) {
	style = strings.TrimSpace(style)
	if style == "" {
		return 0
	}
	u := "px"
	if len(unit) > 0 && unit[0] != "" {
		u = unit[0]
	}

	pattern := fmt.Sprintf(`\s*([\d.]+)%s`, u)
	re := regexp.MustCompile(pattern)

	ss := strings.Split(style, ";")
	for _, s := range ss {
		vv := strings.SplitN(s, ":", 2)
		if len(vv) == 2 && vv[0] == prop {
			matches := re.FindStringSubmatch(vv[1])
			if len(matches) == 2 {
				if v, _ := strconv.ParseFloat(matches[1], 32); v > 0 {
					return v
				}
			}
		}
	}

	return 0
}
