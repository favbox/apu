package note_test

import (
	"context"
	"fmt"
	"testing"

	"apu/app/note"
	"apu/app/note/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool
var ctx = context.Background()
var err error

func init() {
	pool, err = db.NewPool()
	if err != nil {
		panic(err)
	}
}

func TestService_ImportNotes(t *testing.T) {
	s := note.NewService(pool)
	params := []db.CreateNotesParams{
		{
			Title:       "标题",
			Description: "",
			TagIds:      []int64{1, 2, 3},
			Type:        db.NoteTypeNormal,
			IsPrivacy:   true,
			SourceType:  db.SourceTypeUnset,
			SourceUrl:   "",
		},
		{
			Title:       "标题2",
			Description: "",
			TagIds:      []int64{1, 2},
			Type:        db.NoteTypeVideo,
			IsPrivacy:   true,
			SourceType:  db.SourceTypeWeixin,
			SourceUrl:   "https://mp.weixin.qq.com/s/WnGKB9j4-CGyqdX70kl7XA",
		},
	}
	ids, err := s.ImportNotes(ctx, params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ids)
}
