package mysql_test

import (
	"testing"

	"apu/pkg/store/mysql"
	"apu/pkg/store/mysql/query"
	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	mysql.Init()
	count, err := query.WechatHeader.Count()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, int64(0), count)
}
