package weixin_test

import (
	"fmt"
	"testing"
	"time"

	"apu/pkg/source/weixin"
	"apu/pkg/store/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetArticleStat(t *testing.T) {
	biz := "MzkyMDA3MDcwMQ=="
	mid := "2248142696"
	idx := "1"
	sn := "7400cde03e86b5450481cd10d4fbfbe6"

	mysql.Init()
	stat, err := weixin.GetArticleStat(biz, mid, idx, sn)
	if err != nil {
		t.Fatal(err)
	}

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

func TestGetArticleByURL(t *testing.T) {
	rawURL := "https://mp.weixin.qq.com/s?__biz=MjM5MzcxOTcyMQ==&mid=2651661917&idx=1&sn=4eb32349e238778487de133374017d25&chksm=bc897a3b5cd786861e6520c330279d4bf2e3587b315fdf579f5e332ff1e78b6b1cb163ca850d&scene=132&exptype=timeline_recommend_article_extendread_samebiz&show_related_article=1&subscene=132&scene=132#wechat_redirect"
	document, err := weixin.GetArticleByURL(rawURL)
	require.Nil(t, err)
	for _, img := range document.Images {
		fmt.Println(img.OriginalUrl)
	}
}

func TestGetArticleKey(t *testing.T) {
	rawURL := "https://mp.weixin.qq.com/s/6gGMs7dK5-ngfUMYJ2IxDA"
	keyInfo, err := weixin.GetArticleKeyInfo(rawURL)
	assert.NotNil(t, err)
	assert.Nil(t, keyInfo)

	rawURL = "https://mp.weixin.qq.com/s?__biz=MjM5MzcxOTcyMQ==&mid=2651661917&idx=1&sn=4eb32349e238778487de133374017d25&chksm=bc897a3b5cd786861e6520c330279d4bf2e3587b315fdf579f5e332ff1e78b6b1cb163ca850d&scene=132&exptype=timeline_recommend_article_extendread_samebiz&show_related_article=1&subscene=132&scene=132#wechat_redirect"
	keyInfo, err = weixin.GetArticleKeyInfo(rawURL)
	assert.Equal(t, "MjM5MzcxOTcyMQ==", keyInfo.Biz)
	assert.Equal(t, "2651661917", keyInfo.Mid)
	assert.Equal(t, "1", keyInfo.Idx)
}
