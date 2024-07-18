package mysql

import (
	"errors"
	"time"

	"apu/pkg/store/mysql/model"
	"apu/pkg/store/mysql/query"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/clause"
)

// CreateNotesOrSkip 批量创建一批笔记，如果已存在则跳过。
func CreateNotesOrSkip(notes []*model.Note) error {
	//q := query.Use(DB)
	//err := q.Transaction(func(tx *query.Query) error {
	//
	//
	//	return nil
	//})

	// 保存或跳过笔记
	err := query.Note.
		Clauses(clause.OnConflict{
			//UpdateAll: true,
			//Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"update_time": time.Now()}),
		}).
		Create(notes...)
	if err != nil {
		return err
	}

	// 初始化任务管道
	var pipelines []*model.NotePipeline
	for _, n := range notes {
		if n.ID == 0 {
			log.Fatal().Err(errors.New("note id should not be zero"))
		}
		pipelines = append(pipelines, &model.NotePipeline{
			ID: n.ID,
		})
	}
	if err = query.NotePipeline.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}, {Name: "note_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"update_time": time.Now()}),
		}).
		Create(pipelines...); err != nil {
		return err
	}

	return err
}
