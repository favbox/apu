package note

import (
	"context"
	"log"

	"apu/app/note/db"
	"apu/pkg/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ UseCase = (*Service)(nil)

// Service 实现笔记的用例服务。
type Service struct {
	service.Base
}

// ImportNotes 批量导入一批笔记，场景如采集器的分页列表。
func (s *Service) ImportNotes(ctx context.Context, params []db.CreateNotesParams) (ids []int64, err error) {
	s.Queries.CreateNotes(ctx, params).QueryRow(func(i int, id int64, e error) {
		if e != nil {
			err = e
			return
		}
		log.Println(i, id)
		ids = append(ids, id)
	})
	return
}

// NewService 创建一个新的笔记服务。
func NewService(p *pgxpool.Pool) *Service {
	q := db.New(p)
	return &Service{
		Base: service.Base{
			Pool:    p,
			Queries: q,
		},
	}
}
