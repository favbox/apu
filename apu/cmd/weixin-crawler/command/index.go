package command

import (
	"time"

	"apu/pkg/schema"
	"apu/pkg/source/weixin"
	"apu/pkg/store/mysql"
	"apu/pkg/store/mysql/model"
	"apu/pkg/store/mysql/query"
	"github.com/spf13/cobra"
)

// var biz string
var minTimestamp int64

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "索引指定公众号的文章列表",
	Run:   indexBiz,
}

func init() {
	// 默认为当年1月1日
	minTimestamp = time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local).Unix()

	//indexCmd.Flags().StringVar(&biz, "biz", "", "biz")
	//indexCmd.Flags().Int64Var(&minTimestamp, "minTimestamp", minTimestamp, "最晚日期，格式形如 1704038400")
	//err := indexCmd.MarkFlagRequired("biz")
	//if err != nil {
	//	log.Fatal().Err(err).Msg("参数错误")
	//}

	rootCmd.AddCommand(indexCmd)
}

func indexBiz(cmd *cobra.Command, args []string) {
	mysql.Init()

	// 获取公众号上一次同步后采集到的最后发布时间
	biz := "MzI5NTYyMTgwMA=="
	mp, err := query.WeixinMp.Where(query.WeixinMp.Biz.Eq(biz)).First()
	if err != nil || mp == nil {
		log.Fatal().Err(err).Msg("该 biz 未录入在 DB 中")
	}

	syncArticlesByBiz(mp)

	// 全量更新
	//weixinMps, err := query.WeixinMp.Find()
	//if err != nil {
	//	log.Fatal().Err(err).Msg("读取DB失败")
	//}
	//for _, mp := range weixinMps {
	//	log.Info().Str("公众号", mp.Name).Send()
	//	biz = mp.Biz
	//	lastPublishTime := mp.LastPublishTime
	//	//syncKey = stringx.MustNumber[int](mp.LastPublishTime.Format("20060102"))
	//	syncArticlesByBiz(biz, lastPublishTime)
	//	time.Sleep(10 * time.Second)
	//}
}

func syncArticlesByBiz(mp *model.WeixinMp) {
	biz := mp.Biz
	count := 20
	offset := 0
	next := true

	// 初始化作者信息
	author, err := syncWeixinAuthor(mp.Biz, mp.Name)
	if err != nil {
		log.Fatal().Err(err).Msg("无法同步作者信息")
	}

	for next {
		// 获取文章列表
		articles, _, err := weixin.GetArticles(biz, count, offset, 0)
		if err != nil {
			log.Fatal().Err(err).Str("biz", biz).Msg("无法获取线上公众号文章列表")
		} else if len(articles) == 0 {
			next = false
		}
		offset += len(articles)

		// 批量保存为笔记
		var notes []*model.Note
		for i, a := range articles {
			// 第一页每次都更新，其他页看发布时间
			if i > 0 && a.PublishTime.Unix() < minTimestamp {
				next = false
				break
			}
			//content := a.Content
			//if len(content) > 1000 {
			//	content = stringx.Cut(content, 1000)
			//}
			log.Debug().Int("序号", i+1).Time("发布事件", a.PublishTime).Str("标题", a.Title).Str("网址", a.OriginalUrl).Send()
			notes = append(notes, &model.Note{
				ID:          a.UID,
				Source:      schema.SourceWeixin.Int(),
				State:       schema.NoteStateInit,
				Type:        schema.NoteTypeNormal,
				PublishTime: a.PublishTime,
				Title:       a.Title,
				//Description: content,
				AuthorID:    author.ID,
				OriginalURL: a.OriginalUrl,
			})
		}
		if len(notes) > 0 {
			err := mysql.CreateNotesOrSkip(notes)
			//err := query.Note.Save(notes...)
			if err != nil {
				log.Fatal().Err(err).Msg("批量入库失败")
			}
		}
	}
}
