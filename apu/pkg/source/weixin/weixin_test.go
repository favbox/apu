package weixin_test

import (
	"fmt"
	"testing"
	"time"
	"unicode/utf8"

	"apu/pkg/source/weixin"
	"apu/pkg/store/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/ratelimit"
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
		fmt.Println(i+1, a.PublishTime, a.Title)
	}
}

func TestGetArticles(t *testing.T) {
	biz := "MzkyMDA3MDcwMQ=="

	count := 5
	offset := 0
	syncKey := 0
	var minTimestamp int64 = 1719763200

	articles, nextKey, err := weixin.GetArticles(biz, count, offset, syncKey)
	assert.Nil(t, err)
	fmt.Println("synckey", nextKey)
	next := true
	for _, a := range articles {
		if a.PublishTime.Unix() < minTimestamp {
			next = false
			break
		}
		stat, err := weixin.GetArticleStatByURL(a.OriginalUrl)
		require.Nil(t, err)
		fmt.Println(a.PublishTime.Format(time.DateTime), a.Title, stat.ReadNum, stat.LikeNum, stat.OldLikeNum, stat.CollectNum, stat.FriendLikeNum)
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
			stat, err := weixin.GetArticleStatByURL(a.OriginalUrl)
			require.Nil(t, err)
			fmt.Println(a.PublishTime.Format(time.DateTime), a.Title, stat.ReadNum, stat.LikeNum, stat.OldLikeNum, stat.CollectNum, stat.FriendLikeNum)
			if a.PublishTime.Unix() < minTimestamp {
				next = false
			}
		}
	}
}

func TestGetArticleKey(t *testing.T) {
	rawURL := "https://mp.weixin.qq.com/s/6gGMs7dK5-ngfUMYJ2IxDA"
	keyInfo, err := weixin.GetArticleKeyInfo(rawURL)
	assert.NotNil(t, err)
	assert.Nil(t, keyInfo)

	rawURL = "https://mp.weixin.qq.com/s?__biz=MjM5MzcxOTcyMQ==&mid=2651661917&idx=1&sn=4eb32349e238778487de133374017d25&chksm=bc897a3b5cd786861e6520c330279d4bf2e3587b315fdf579f5e332ff1e78b6b1cb163ca850d&scene=132&exptype=timeline_recommend_article_extendread_samebiz&show_related_article=1&subscene=132&scene=132#wechat_redirect"
	keyInfo, err = weixin.GetArticleKeyInfo(rawURL)
	assert.Nil(t, err)
	assert.Equal(t, "MjM5MzcxOTcyMQ==", keyInfo.Biz)
	assert.Equal(t, "2651661917", keyInfo.Mid)
	assert.Equal(t, "1", keyInfo.Idx)
}

func TestGetArticleByURL(t *testing.T) {
	rawURL := "https://mp.weixin.qq.com/s?__biz=MjM5MzcxOTcyMQ==&mid=2651661917&idx=1&sn=4eb32349e238778487de133374017d25&chksm=bc897a3b5cd786861e6520c330279d4bf2e3587b315fdf579f5e332ff1e78b6b1cb163ca850d&scene=132&exptype=timeline_recommend_article_extendread_samebiz&show_related_article=1&subscene=132&scene=132#wechat_redirect"
	//rawURL = "https://mp.weixin.qq.com/s?__biz=Mzg5MTIxNjQ3NQ==&mid=2247497762&idx=1&sn=03cba3a0de7c845611e438cb2b767f01&chksm=cfd20f16f8a58600baeae765d9aceb4aff469ceca560c9e3ea1fed36a6bf2064c4091070691f#rd"
	//rawURL = "https://mp.weixin.qq.com/s/ya2rYHfsRu0DbOSE3Qb4Ew"
	//rawURL = "https://mp.weixin.qq.com/s/DUpXEV9dCU2DlqpanCHBug"
	//rawURL = "https://mp.weixin.qq.com/s/KXupZGje7CLda_7mdduURQ"
	//rawURL = "https://mp.weixin.qq.com/s/fnnS3j1zKqiTQe3nUcXrvg"
	//rawURL = "https://mp.weixin.qq.com/s/DJC3aejOlmX9llKi4H5Oag"
	//rawURL = "https://mp.weixin.qq.com/s?__biz=MzA4NjAyNjAzMw==&mid=2651132177&idx=1&sn=8c0d49e60ebf05ee1c180286922371f0&chksm=843f434fb348ca59e66c6458a80531d207409a495338fc762cdf50534b3e765f916947e6a53c&scene=58&subscene=0#rd"
	rawURL = "https://mp.weixin.qq.com/s/r2TqMPZnY0wO7Cv07l8QfA"
	rawURL = "https://mp.weixin.qq.com/s/WnGKB9j4-CGyqdX70kl7XA"
	rawURL = "https://mp.weixin.qq.com/s?__biz=MjM5MzcxOTcyMQ==&mid=2651654837&idx=1&sn=077aea1dac9f3a25ae3c0b006c47e82e&chksm=bd6b1e738a1c9765115c478c4c0fc8bb6f554a0ad2247cd417e56ec6062d7bdd7bfd8c6a0170&scene=58&subscene=0#rd"
	a, err := weixin.GetArticleByURL(rawURL)
	require.Nil(t, err)
	for i, img := range a.Images {
		fmt.Println(i+1, img.UID, img.Width, img.Height, img.OriginalUrl)
	}
	fmt.Println(a.UID, a.Source, a.PublishTime, a.Title, utf8.RuneCountInString(a.Content), a.OriginalUrl)
}

func TestLimiter(t *testing.T) {
	rl := ratelimit.New(1, ratelimit.Per(3*time.Second))
	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := rl.Take()
		fmt.Println(i, now.Sub(prev))
		prev = now
	}
}
