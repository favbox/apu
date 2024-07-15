package command

import (
	"apu/pkg/schema"
	"apu/pkg/source/weixin"
	"apu/pkg/store/mysql"
	"apu/pkg/store/mysql/model"
	"apu/pkg/store/mysql/query"
	"github.com/spf13/cobra"
)

var url string

var detailCmd = &cobra.Command{
	Use:   "detail",
	Short: "抓取文章详情",
	Run:   crawDetail,
}

func init() {
	detailCmd.Flags().StringVar(&url, "url", "", "公众号文章网址")
	if err := detailCmd.MarkFlagRequired("url"); err != nil {
		log.Fatal().Err(err).Msg("参数错误")
	}

	rootCmd.AddCommand(detailCmd)
}

func crawDetail(cmd *cobra.Command, args []string) {
	// 获取文章信息
	a, err := weixin.GetArticleByURL(url)
	if err != nil {
		log.Fatal().Err(err).Msg("获取文章信息失败")
	}

	// 确保作者事先存在
	metadata := a.Metadata
	var author *model.Author
	author, err = syncWeixinAuthor(metadata["biz"].(string), metadata["mpName"].(string))
	if err != nil {
		log.Fatal().Err(err).Msg("无法事先准备作者")
	}

	// 使用事务入库
	mysql.Init()
	q := query.Use(mysql.DB)
	err = q.Transaction(func(tx *query.Query) error {
		// 保存为笔记
		note, err := tx.Note.Where(tx.Note.ID.Eq(a.UID)).FirstOrInit()
		if err != nil {
			return err
		}
		note.Source = schema.SourceWeixin.Int()
		note.State = schema.NoteStateDetailed
		note.Type = schema.NoteTypeNormal
		note.PublishTime = a.PublishTime
		note.Title = a.Title
		note.OriginalURL = a.OriginalUrl
		if author != nil {
			note.AuthorID = author.ID
		}
		if err := tx.Note.Save(note); err != nil {
			return err
		}

		// 保存笔记内容长文本
		if err := tx.NoteContent.Save(&model.NoteContent{
			ID:   note.ID,
			Text: a.Content,
		}); err != nil {
			return err
		}

		// 批量保存为笔记图片
		var images []*model.Image
		for i, img := range a.Images {
			images = append(images, &model.Image{
				UID:         img.UID,
				NoteID:      note.ID,
				OriginalURL: img.OriginalUrl,
				Width:       img.Width,
				Height:      img.Height,
				Sort:        i,
			})
		}
		if err := tx.Image.Save(images...); err != nil {
			return err
		}

		// 更新管道交互量已采集
		err = mysql.UpdatePipeline(note.ID, mysql.PipelineOptions{IsDetailed: true})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Msg("入库失败")
	}
}
