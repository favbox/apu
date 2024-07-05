package addon

import (
	"bytes"
	"fmt"
	"log"

	"apu/internal/cookieutil"
	"apu/pkg/store/mysql"
	"apu/pkg/store/mysql/model"
	"apu/pkg/store/mysql/query"
	"github.com/lqqyt2423/go-mitmproxy/proxy"
	"gorm.io/gorm/clause"
)

const reloadWindowJavascript = `
	<script type="text/javascript">
		setInterval(() => {
			window.location.href=window.location.href;
		}, 3000);
	</script>
</body>
`

// WechatAddon 微信代理插件。
// 用于拦截文章阅读量请求头。
type WechatAddon struct {
	proxy.BaseAddon
}

func (a *WechatAddon) Response(f *proxy.Flow) {
	fmt.Println(f.Request.URL.Path, f.Request.URL.RawQuery)

	// 该插件只拦截文章详情页
	if f.Request.URL.Hostname() != "mp.weixin.qq.com" ||
		f.Request.URL.Path != "/s" {
		return
	}

	f.Response.ReplaceToDecodedBody()

	// 注入定时刷新脚本 TODO 未生效
	//go injectRefreshJavascript(f)

	// 保存请求头（包含了 可请求阅读量的 appmsg_token 令牌）
	cookie := f.Request.Header.Get("cookie")
	cookieMap := cookieutil.StrToMap(cookie)
	if _, exists := cookieMap["appmsg_token"]; exists {
		mysql.Init()
		wxuin := cookieMap["wxuin"]
		err := query.WechatCookie.
			Clauses(clause.OnConflict{
				UpdateAll: true,
			}).
			Create(&model.WechatCookie{
				Wxuin:  wxuin,
				Type:   "wechat",
				Cookie: cookie,
				Status: "valid",
			})
		if err != nil {
			log.Println(err)
		}
		log.Println("已捕获新的微信可用会话", wxuin)
	}
}

func injectRefreshJavascript(f *proxy.Flow) {
	body := f.Response.Body
	newBody := bytes.Replace(body, []byte("</body>"), []byte(reloadWindowJavascript), 1)
	f.Response.Body = newBody

	log.Println("刷新窗口的脚本已注入")
}
