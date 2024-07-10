package command

import (
	"fmt"
	"strconv"
	"time"

	"apu/pkg/schema"
	"apu/pkg/source/weixin"
	"apu/pkg/store/mysql"
	"apu/pkg/store/mysql/model"
	"apu/pkg/store/mysql/query"
	"github.com/spf13/cobra"
	"gorm.io/gen"
	"gorm.io/gorm/clause"
)

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "抓取文章阅读量",
	Run:   crawStat,
}

func init() {
	rootCmd.AddCommand(statCmd)
}

func crawStat(cmd *cobra.Command, args []string) {
	mysql.Init()

	// 批量获取文章统计值
	var documents []*model.Document
	batchSize := 10
	today, _ := strconv.ParseInt(time.Now().Format("20060102"), 10, 32)
	day := int32(today)
	err := query.Document.
		Where(query.Document.Stage.Eq(schema.DocumentStageInit)).
		FindInBatches(&documents, batchSize, func(tx gen.Dao, batch int) error {
			// 获取阅读量
			var interactions []*model.Interaction
			for i, doc := range documents {
				fmt.Println(i, doc.OriginalURL)
				stat, err := weixin.GetArticleStatByURL(doc.OriginalURL)
				if err != nil {
					return err
				}
				fmt.Println(doc.Title, stat.ReadNum)
				interactions = append(interactions, &model.Interaction{
					DocID:      doc.ID,
					Day:        day,
					ReadNum:    stat.ReadNum,
					LikeNum:    stat.LikeNum + stat.OldLikeNum,
					CollectNum: stat.CollectNum,
				})

				// 更新文档
				doc.ReadNum = stat.ReadNum
				doc.LikeNum = stat.LikeNum + stat.OldLikeNum
				doc.CollectNum = stat.CollectNum
				doc.Stage = schema.DocumentStageStated
				err = tx.Save(doc)
				if err != nil {
					return err
				}
			}

			// 批量更新文档统计
			err := query.Interaction.
				Clauses(clause.OnConflict{UpdateAll: true}).
				Create(interactions...)
			if err != nil {
				log.Fatal().Err(err).Msg("无法批量入库")
			}

			return nil
		})
	if err != nil {
		log.Fatal().Err(err).Msg("无法批量采集文章阅读量")
	}

	// 批量更新文档状态

}
