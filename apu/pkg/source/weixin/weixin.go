package weixin

import (
	"apu/pkg/schema"
	"apu/pkg/source/weixin/article"
)

// GetArticles 获取公众号下的文章列表。
func GetArticles(biz string, count, offset, syncKey int) ([]*schema.Document, int, error) {
	bookId := Biz2BookId(biz)
	return article.GetArticles(bookId, count, offset, syncKey)
}

// GetArticleStat 获取文章统计信息。
func GetArticleStat(biz, mid, idx, sn string) (*article.Stat, error) {
	return article.GetStat(biz, mid, idx, sn)
}

// GetArticleStatByURL 获取文章统计信息。
func GetArticleStatByURL(canonicalURL string) (*article.Stat, error) {
	k, err := article.GetKeyInfo(canonicalURL)
	if err != nil {
		return nil, err
	}
	return article.GetStat(k.Biz, k.Mid, k.Idx, k.Sn)
}

// GetArticleByURL 获取文章信息。
func GetArticleByURL(canonicalURL string) (*schema.Document, error) {
	return article.GetArticle(canonicalURL)
}

// GetArticleKeyInfo 获取微信公众号文章的键信息。
func GetArticleKeyInfo(canonicalURL string) (*article.KeyInfo, error) {
	return article.GetKeyInfo(canonicalURL)
}
