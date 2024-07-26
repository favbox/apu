package shortid_test

import (
	"fmt"
	"testing"

	"apu/pkg/schema"
	"apu/pkg/shortid"
	"apu/pkg/source/weixin"
	"github.com/stretchr/testify/assert"
)

func TestEncodeWeixinArticle(t *testing.T) {
	url := "https://mp.weixin.qq.com/s?__biz=MjM5MzcxOTcyMQ==&mid=2651661917&idx=1&sn=4eb32349e238778487de133374017d25&chksm=bc897a3b5cd786861e6520c330279d4bf2e3587b315fdf579f5e332ff1e78b6b1cb163ca850d&scene=132&exptype=timeline_recommend_article_extendread_samebiz&show_related_article=1&subscene=132&scene=132#wechat_redirect"
	keyInfo, err := weixin.GetArticleKeyInfo(url)
	assert.Nil(t, err)

	id, err := shortid.EncodeWeixinArticleID(keyInfo.Biz, keyInfo.Mid, keyInfo.Idx)
	assert.Nil(t, err)
	assert.Equal(t, "1a30b2dc886fc8fb7f9c65ce", id)

	biz, mid, idx, err := shortid.DecodeWeixinArticleID(id)
	assert.Nil(t, err)
	assert.Equal(t, "MjM5MzcxOTcyMQ==", biz)
	assert.Equal(t, "2651661917", mid)
	assert.Equal(t, "1", idx)
}

func TestDecode(t *testing.T) {
	source, nums, err := shortid.Decode("1a30b2dc886fc8fb7f9c65ce")
	assert.Equal(t, schema.SourceWeixin, source)
	assert.Nil(t, err)
	assert.Len(t, nums, 4)
	fmt.Println(nums)
}

func TestName(t *testing.T) {
	url := "https://mp.weixin.qq.com/s?__biz=MjM5MzcxOTcyMQ==&mid=2651661917&idx=1&sn=4eb32349e238778487de133374017d25&chksm=bc897a3b5cd786861e6520c330279d4bf2e3587b315fdf579f5e332ff1e78b6b1cb163ca850d&scene=132&exptype=timeline_recommend_article_extendread_samebiz&show_related_article=1&subscene=132&scene=132#wechat_redirect"
	url := "https://mp.weixin.qq.com/s?__biz=MjM5MzcxOTcyMQ==&mid=2651661917&idx=1&sn=4eb32349e238778487de133374017d25&chksm=bc897a3b5cd786861e6520c330279d4bf2e3587b315fdf579f5e332ff1e78b6b1cb163ca850d&scene=132&exptype=timeline_recommend_article_extendread_samebiz&show_related_article=1&subscene=132&scene=132#wechat_redirect"
	keyInfo, err := weixin.GetArticleKeyInfo(url)
	assert.Nil(t, err)
	//weixinKey := xxhash3.HashString(fmt.Sprintf("%s:%s:%s", keyInfo.Biz, keyInfo.Mid, keyInfo.Idx))

	// 唯一、无序、可逆、快

}

func BenchmarkEncodeWeixinArticle(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		shortid.EncodeWeixinArticleID("MjM5MzcxOTcyMQ==", "2651661917", "1")
	}
}

func BenchmarkDecodeWeixinArticle(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		shortid.DecodeWeixinArticleID("1a30b2dc886fc8fb7f9c65ce")
	}
}
