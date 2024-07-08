package weixin_test

import (
	"fmt"
	"testing"
	"time"

	"apu/pkg/datasource/weixin"
	"apu/pkg/datasource/weixin/article"
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

func TestGetArticleParams(t *testing.T) {
	rawURL := "https://mp.weixin.qq.com/s/UQiQocQ2MWkwImZJFMMWQw"
	biz, mid, idx, sn, err := article.GetArticleParams(rawURL)
	assert.Nil(t, err)
	fmt.Println(biz, mid, idx, sn)

	rawURL = "http://mp.weixin.qq.com/s?__biz=MzA5ODEzMjIyMA==&amp;mid=2247713279&amp;idx=1&amp;sn=bd67c1aba187bc8833f4aee60d8a0e90&amp;chksm=909b886ca7ec017a4c72af5460dcdfe672a74450da0cec34f99cab825d512ab37e4c8715eda6#rd"
	biz, mid, idx, sn, err = article.GetArticleParams(rawURL)
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

func TestGetArticleByURL(t *testing.T) {
	rawURL := "https://mp.weixin.qq.com/s/WnGKB9j4-CGyqdX70kl7XA"
	rawURL = "https://mp.weixin.qq.com/s?__biz=MzA4NzEwMTMyOQ==&amp;mid=2649849538&amp;idx=1&amp;sn=c07678777df1b73d68fef027544affdc&amp;chksm=883b443cbf4ccd2af975876b5e7731c43ab616cba4e3d82ce10feabe2530d03088abfc5e5655#rd"
	rawURL = "https://mp.weixin.qq.com/s?__biz=MjM5MzcxOTcyMQ==&amp;mid=2651628506&amp;idx=1&amp;sn=a5a68a7da94bd092f0cb9abb10f6d11f&amp;chksm=bd6ab71c8a1d3e0ad9d4f985548811fe0abb1364b958ad3ee7bc82ef08d4c72c0f8239bcda5c#rd"
	rawURL = "https://mp.weixin.qq.com/s?__biz=MjM5MzcxOTcyMQ==&mid=2651661917&idx=1&sn=4eb32349e238778487de133374017d25&chksm=bc897a3b5cd786861e6520c330279d4bf2e3587b315fdf579f5e332ff1e78b6b1cb163ca850d&scene=132&exptype=timeline_recommend_article_extendread_samebiz&show_related_article=1&subscene=132&scene=132#wechat_redirect"
	info, err := weixin.GetArticleByURL(rawURL)
	require.Nil(t, err)
	fmt.Printf("%#v\n", info)
}
