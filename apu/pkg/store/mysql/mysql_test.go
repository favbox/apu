package mysql_test

import (
	"testing"

	"apu/pkg/store/mysql"
	"apu/pkg/store/mysql/query"
	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	mysql.Init()
	count, err := query.WechatCookie.Count()
	if err != nil {
		t.Fatal(err)
	}
	assert.GreaterOrEqual(t, count, int64(0))
}
