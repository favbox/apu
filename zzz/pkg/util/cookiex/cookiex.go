package cookiex

import (
	"fmt"
	"net/http"
	"strings"
)

func StrToMap(s string) (m map[string]string) {
	m = make(map[string]string)
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}
	cookies := strings.Split(s, ";")
	for _, c := range cookies {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		kv := strings.SplitN(c, "=", 2)
		if len(kv[1]) > 0 {
			m[kv[0]] = kv[1]
		}
	}
	return
}

func StrToHttpCookies(s string) []*http.Cookie {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	var ret []*http.Cookie
	cookies := strings.Split(s, ";")
	for _, c := range cookies {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		kv := strings.SplitN(c, "=", 2)
		if len(kv[1]) > 0 {
			ret = append(ret, &http.Cookie{
				Domain: ".xiaohongshu.com",
				Path:   "/",
				Name:   kv[0],
				Value:  kv[1],
			})
		}
	}
	return ret
}

func MapToStr(m map[string]string) string {
	var cookies []string
	for name, value := range m {
		cookies = append(cookies, fmt.Sprintf("%s=%s", name, value))
	}
	return strings.Join(cookies, ";")
}

func MapToHttpCookies(m map[string]string) []*http.Cookie {
	var cookies []*http.Cookie
	for name, value := range m {
		cookies = append(cookies, &http.Cookie{
			Domain: ".xiaohongshu.com",
			Name:   name,
			Value:  value,
		})
	}
	return cookies
}

func HttpCookiesToStr(cookies []*http.Cookie) string {
	var ret []string
	for _, c := range cookies {
		if strings.HasSuffix(c.Domain, ".xiaohongshu.com") && c.Valid() == nil {
			ret = append(ret, fmt.Sprintf("%s=%s", c.Name, c.Value))
		}
	}
	return strings.Join(ret, ";")
}

func HttpCookiesToMap(cookies []*http.Cookie) map[string]string {
	m := make(map[string]string)
	for _, c := range cookies {
		if strings.HasSuffix(c.Domain, ".xiaohongshu.com") && c.Valid() == nil {
			m[c.Name] = c.Value
		}
	}
	return m
}
