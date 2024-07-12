package command

import (
	"fmt"
	"time"

	"apu/pkg/schema"
	"apu/pkg/source/weixin"
	"apu/pkg/store/mysql"
	"apu/pkg/store/mysql/model"
	"github.com/spf13/cobra"
)

var biz string
var minTimestamp int64

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "索引指定公众号的文章列表",
	Run:   indexBiz,
}

func init() {
	// 默认为当年1月1日
	minTimestamp = time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local).Unix()

	indexCmd.Flags().StringVar(&biz, "biz", "", "biz")
	indexCmd.Flags().Int64Var(&minTimestamp, "minTimestamp", minTimestamp, "最晚日期，格式形如 1704038400")
	err := indexCmd.MarkFlagRequired("biz")
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	rootCmd.AddCommand(indexCmd)
}
func indexBiz(cmd *cobra.Command, args []string) {
	mysql.Init()

	count := 20
	offset := 0
	syncKey := 0
	next := true

	for next {
		// 获取文章列表
		articles, _, err := weixin.GetArticles(biz, count, offset, syncKey)
		if err != nil {
			log.Fatal().Err(err).Str("biz", biz).Msg("无法获取线上公众号文章列表")
		} else if len(articles) == 0 {
			next = false
		}
		offset += len(articles)

		// 批量保存为笔记
		var notes []*model.Note
		for i, a := range articles {
			if a.PublishTime.Unix() < minTimestamp {
				next = false
				break
			}
			fmt.Println(i+1, a.PublishTime, a.Author, a.Title, a.OriginalUrl)
			notes = append(notes, &model.Note{
				UID:         a.Key,
				Source:      int32(schema.Weixin),
				State:       int32(schema.NoteStateInit),
				Type:        string(schema.NoteTypeNormal),
				PublishTime: a.PublishTime,
				Title:       a.Title,
				Content:     a.Content,
				AuthorID:    0,
				OriginalURL: a.OriginalUrl,
			})
		}
		if len(notes) > 0 {
			err := mysql.CreateNotesOrSkip(notes)
			if err != nil {
				log.Fatal().Err(err).Msg("批量入库失败")
			}
		}

	}
}
