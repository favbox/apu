package addon

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	"apu/pkg/store/mysql"
	"apu/pkg/store/mysql/model"
	"apu/pkg/store/mysql/query"
	"github.com/lqqyt2423/go-mitmproxy/proxy"
	"gorm.io/gorm/clause"
)

// WereadAddon 微信读书代理插件。
// 用于拦截微信读书APP[PC|Mobile]包含 vid,skey 的 header。
type WereadAddon struct {
	proxy.BaseAddon
}

func (a *WereadAddon) Response(f *proxy.Flow) {
	// 该插件只拦截书信息页
	if f.Request.URL.Hostname() != "i.weread.qq.com" ||
		f.Request.URL.Path != "/book/info" {
		return
	}

	f.Response.ReplaceToDecodedBody()

	// 保存请求头（包含了 可请求阅读量的 appmsg_token 令牌）
	var headersBuf bytes.Buffer
	err := f.Request.Header.Write(&headersBuf)
	if err != nil {
		log.Fatal(err)
		return
	}

	var vid string
	headerMap := make(map[string]string)
	lines := strings.Split(headersBuf.String(), "\n")
	for _, line := range lines {
		vs := strings.SplitN(line, ":", 2)
		if len(vs) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(vs[0]))
		value := strings.TrimSpace(vs[1])
		headerMap[key] = value
		if key == "vid" {
			vid = value
		}
	}
	headers, err := json.Marshal(headerMap)
	if err != nil {
		return
	}

	if vid != "" {
		mysql.Init()
		err := query.WexinRequest.
			Clauses(clause.OnConflict{
				UpdateAll: true,
			}).
			Create(&model.WexinRequest{
				Type:    "weread",
				UserID:  vid,
				Headers: string(headers),
				Status:  "valid",
			})
		if err != nil {
			log.Println(err)
		}
		log.Println("已捕获新的【微信读书】可用会话", vid)
	}
}
