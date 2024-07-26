package note

import (
	"context"

	"apu/app/note/db"
)

type Note struct {
}

type Repository interface {
	db.DBTX
}

type UseCase interface {
	ImportNotes(ctx context.Context, params []db.CreateNotesParams) ([]int64, error)
}
