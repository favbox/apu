package mysql_test

import (
	"fmt"
	"testing"
	"time"

	"apu/pkg/store/mysql"
	"apu/pkg/store/mysql/query"
	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	mysql.Init()
	count, err := query.WexinRequest.Count()
	if err != nil {
		t.Fatal(err)
	}
	assert.GreaterOrEqual(t, count, int64(0))
}

func TestBefore(t *testing.T) {
	fmt.Print(Year())
}

func Year() int64 {
	return time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local).Unix()
}
