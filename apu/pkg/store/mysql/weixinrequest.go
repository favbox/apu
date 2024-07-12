package mysql

import (
	"apu/pkg/store/mysql/model"
	"apu/pkg/store/mysql/query"
)

func FetchWexinRequest(reqType, status string) (*model.WexinRequest, error) {
	return query.WexinRequest.Where(
		query.WexinRequest.Type.Eq(reqType),
		query.WexinRequest.Status.Eq(status),
	).First()
}
