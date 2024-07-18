package note

import (
	"apu/note/mysql"
	"github.com/google/wire"
)

var ServiceProvider = wire.NewSet(NewService)

var _ UseCase = (*Service)(nil)

type Service struct {
	queries mysql.Queries
}

func NewService(queries mysql.Queries) UseCase {
	return &Service{
		queries: queries,
	}
}
