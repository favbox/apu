package weixin_test

import (
	"fmt"
	"testing"
	"time"

	"apu/pkg/datasource/weixin"
	"apu/pkg/store/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestGetArticleParams(t *testing.T) {
	rawURL := "https://mp.weixin.qq.com/s/UQiQocQ2MWkwImZJFMMWQw"
	biz, mid, idx, sn, err := weixin.GetArticleParams(rawURL)
	assert.Nil(t, err)
	fmt.Println(biz, mid, idx, sn)

	rawURL = "http://mp.weixin.qq.com/s?__biz=MzA5ODEzMjIyMA==&amp;mid=2247713279&amp;idx=1&amp;sn=bd67c1aba187bc8833f4aee60d8a0e90&amp;chksm=909b886ca7ec017a4c72af5460dcdfe672a74450da0cec34f99cab825d512ab37e4c8715eda6#rd"
	biz, mid, idx, sn, err = weixin.GetArticleParams(rawURL)
	assert.Nil(t, err)
	fmt.Println(biz, mid, idx, sn)
}

func TestGetArticleStatByURL(t *testing.T) {
	rawURL := "https://mp.weixin.qq.com/s/UQiQocQ2MWkwImZJFMMWQw"
	rawURL = "https://mp.weixin.qq.com/s/bYTmk-ZRBJ6YaJG3f00l_g"
	stat, err := weixin.GetArticleStatByURL(rawURL)
	assert.Nil(t, err)
	assert.NotNil(t, stat)
	fmt.Printf("%#v", stat)
}

func TestGetArticlesWithSyncKey(t *testing.T) {
	biz := "MzkyMDA3MDcwMQ=="

	count := 5
	offset := 0
	syncKey := 1719763200

	articles, syncKey, err := weixin.GetArticles(biz, count, offset, syncKey)
	assert.Nil(t, err)
	for i, a := range articles {
		fmt.Println(i+1, a.Time, a.Title)
	}
}

func TestGetArticles(t *testing.T) {
	biz := "MzkyMDA3MDcwMQ=="

	count := 5
	offset := 0
	syncKey := 0
	minTimestamp := 1719763200

	articles, nextKey, err := weixin.GetArticles(biz, count, offset, syncKey)
	assert.Nil(t, err)
	fmt.Println("synckey", nextKey)
	next := true
	for _, a := range articles {
		if a.Time < minTimestamp {
			next = false
			break
		}
		stat, err := weixin.GetArticleStatByURL(a.DocUrl)
		require.Nil(t, err)
		fmt.Println(time.Unix(int64(a.Time), 0).Format(time.DateTime), a.Title, stat.ReadNum, stat.LikeNum, stat.OldLikeNum, stat.CollectNum, stat.FriendLikeNum)
	}

	for next {
		offset += len(articles)
		articles, nextKey, err := weixin.GetArticles(biz, count, offset, 0)
		assert.Nil(t, err)
		if len(articles) == 0 {
			next = false
		}

		fmt.Println("synckey", nextKey)
		for _, a := range articles {
			stat, err := weixin.GetArticleStatByURL(a.DocUrl)
			require.Nil(t, err)
			fmt.Println(time.Unix(int64(a.Time), 0).Format(time.DateTime), a.Title, stat.ReadNum, stat.LikeNum, stat.OldLikeNum, stat.CollectNum, stat.FriendLikeNum)
			if a.Time < minTimestamp {
				next = false
			}
		}
	}
}
