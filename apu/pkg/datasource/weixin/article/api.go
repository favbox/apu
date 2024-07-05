package article

import (
	"apu/internal/cookieutil"
	"apu/pkg/store/mysql/query"
	"github.com/imroc/req/v3"
)

func GetList(biz string) {

}

// GetStat 获取文章统计信息。 https://www.cnblogs.com/jianpansangejian/p/17970546
func GetStat(biz, mid, idx, sn string) (*Stat, error) {
	//mysql.Init()
	header, err := query.WechatCookie.Where(
		query.WechatCookie.Type.Eq("wechat"),
		query.WechatCookie.Status.Eq("valid"),
	).First()
	if err != nil {
		return nil, err
	}

	request := req.R()

	// 设置查询参数
	cookieMap := cookieutil.StrToMap(header.Cookie)
	request.SetQueryParams(map[string]string{
		"appmsg_token": cookieMap["appmsg_token"],
		"x5":           "0",
	})

	// 设置请求头
	request.SetHeaders(map[string]string{
		"User-Agent":   "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0Chrome/57.0.2987.132 MQQBrowser/6.2 Mobile",
		"Cookie":       header.Cookie,
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
	_, err = request.Post("https://mp.weixin.qq.com/mp/getappmsgext")
	if err != nil {
		return nil, err
	}

	return result.ArticleStat, nil
}

func GetContent() {

}
