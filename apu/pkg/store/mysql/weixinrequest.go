package mysql

import (
	"apu/pkg/store/mysql/model"
	"apu/pkg/store/mysql/query"
)

// FetchWeixinRequest 获取指定微信类型[wechat|weread]、指定状态[valid|invalid]请求信息。
func FetchWeixinRequest(reqType, status string) (*model.WeixinRequest, error) {
	return query.WeixinRequest.Where(
		query.WeixinRequest.Type.Eq(reqType),
		query.WeixinRequest.Status.Eq(status),
	).First()
}
