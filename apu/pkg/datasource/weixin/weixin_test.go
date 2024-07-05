package weixin_test

import (
	"fmt"
	"testing"

	"apu/pkg/datasource/weixin"
	"apu/pkg/store/mysql"
	"github.com/stretchr/testify/assert"
)

func init() {
	mysql.Init()
}

func TestGetArticleStat(t *testing.T) {
	biz := "MzkyMDA3MDcwMQ=="
	mid := "2248142696"
	idx := "1"
	sn := "7400cde03e86b5450481cd10d4fbfbe6"

	stat, err := weixin.GetArticleStat(biz, mid, idx, sn)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v", stat)
}

func TestGetArticleStatByURL(t *testing.T) {
	rawURL := "https://mp.weixin.qq.com/s/UQiQocQ2MWkwImZJFMMWQw"
	stat, err := weixin.GetArticleStatByURL(rawURL)
	assert.Nil(t, err)
	assert.NotNil(t, stat)
	fmt.Printf("%#v", stat)
}
