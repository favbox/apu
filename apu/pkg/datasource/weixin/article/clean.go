package article

import (
	"bytes"
	"errors"
	"fmt"

	"apu/pkg/utils/stringx"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

func CleanArticle(body []byte) (*Info, error) {
	// 创建为 goquery 文档
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	info := &Info{}
	// 获取公众号昵称、微信ID和介绍
	info.MpName = stringx.Trim(doc.Find("#js_name").Text())
	doc.Find(".profile_meta_value").Each(func(i int, selection *goquery.Selection) {
		switch i {
		case 0:
			info.MpWeixinID = selection.Text()
		case 1:
			info.MpIntro = selection.Text()
		}
	})

	// 获取文章基本参数
	link := doc.Find("meta[property='og:url']").AttrOr("content", "")
	if link == "" {
		return nil, errors.New("未能在响应中找到 og:url")
	}
	info.Biz, info.Mid, info.Idx, info.Sn, err = GetArticleParams(link)
	if err != nil {
		return nil, err
	}

	// 清洗文章TDK相关信息
	CleanTitle(info, doc.Find("meta[property='og:title']").AttrOr("content", ""))
	info.Image = stringx.Trim(doc.Find("meta[property='og:image']").AttrOr("content", ""))
	info.Description = stringx.Trim(doc.Find("meta[property='og:description']").AttrOr("content", ""))

	// 清洗文章创建时间
	CleanPublishTime(info, body)

	// 清洗文章正文
	CleanJsContent(info, body, doc)

	// 正文转 markdown
	jsContent := doc.Find("#js_content")
	converter := md.NewConverter("", true, nil)
	mdContent := converter.Convert(jsContent)
	fmt.Println(mdContent)

	return info, nil
}
