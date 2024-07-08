package weixin

import (
	"fmt"

	"apu/pkg/datasource/weixin/article"
)

// GetArticles 获取公众号下的文章列表。
func GetArticles(biz string, count, offset, syncKey int) ([]*article.BookArticle, int, error) {
	bookId := Biz2BookId(biz)
	return article.GetArticles(bookId, count, offset, syncKey)
}

// GetArticleStat 获取文章统计信息。
func GetArticleStat(biz, mid, idx, sn string) (*article.Stat, error) {
	return article.GetStat(biz, mid, idx, sn)
}

// GetArticleStatByURL 获取指定文章网址的统计信息。
func GetArticleStatByURL(rawURL string) (*article.Stat, error) {
	// https://mp.weixin.qq.com/s/UQiQocQ2MWkwImZJFMMWQw
	return article.GetStatByURL(rawURL)
}

// GetArticle 获取文章信息。
func GetArticle(biz, mid, idx string) (*article.Info, error) {
	rawURL := fmt.Sprintf(article.DetailURLPattern, biz, mid, idx)
	return GetArticleByURL(rawURL)
}

// GetArticleByURL 获取文章信息。
func GetArticleByURL(rawURL string) (*article.Info, error) {
	return article.GetArticle(rawURL)
}
