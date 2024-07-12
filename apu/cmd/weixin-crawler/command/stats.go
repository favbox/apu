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
	Use:   "stats",
	Short: "抓取文章阅读量",
	Run:   crawStats,
}

func init() {
	rootCmd.AddCommand(statCmd)
}

func crawStats(cmd *cobra.Command, args []string) {
	mysql.Init()

	// 批量获取文章统计值
	var notes []*model.Note
	batchSize := 10
	today, _ := strconv.ParseInt(time.Now().Format("20060102"), 10, 32)
	day := int32(today)
	err := query.Note.
		Where(query.Note.State.Eq(int32(schema.NoteStateInit))).
		FindInBatches(&notes, batchSize, func(tx gen.Dao, batch int) error {
			// 获取阅读量
			var interacts []*model.Interact
			for i, note := range notes {
				fmt.Println(i, note.OriginalURL)
				stat, err := weixin.GetArticleStatByURL(note.OriginalURL)
				if err != nil {
					return err
				}
				fmt.Println(note.Title, stat.ReadNum)
				interacts = append(interacts, &model.Interact{
					NoteID:         note.ID,
					Day:            day,
					ReadCount:      stat.ReadNum,
					LikedCount:     stat.LikeNum + stat.OldLikeNum,
					CollectedCount: stat.CollectNum,
					ShareCount:     stat.ShareNum,
				})

				note.ReadCount = stat.ReadNum
				note.LikedCount = stat.LikeNum + stat.OldLikeNum
				note.CollectedCount = stat.CollectNum
				note.ShareCount = stat.ShareNum

				// 更新文档
				note.State = schema.DocumentStageStated
				err = tx.Save(note)
				if err != nil {
					return err
				}
			}

			// 批量更新文档统计
			err := query.Interact.
				Clauses(clause.OnConflict{UpdateAll: true}).
				Create(interacts...)
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
