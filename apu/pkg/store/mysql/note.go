package mysql

import (
	"apu/pkg/store/mysql/model"
	"apu/pkg/store/mysql/query"
	"gorm.io/gorm/clause"
)

// CreateNotesOrSkip 批量创建一批笔记，如果已存在则跳过。
func CreateNotesOrSkip(notes []*model.Note) error {
	return query.Note.
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Create(notes...)
}
