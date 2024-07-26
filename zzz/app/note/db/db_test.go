package db_test

import (
	"context"
	"testing"

	"apu/app/note/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestNote(t *testing.T) {
	ctx := context.Background()
	d, err := pgxpool.New(ctx, "host=localhost port=5432 user=zs password=zs sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	err = d.Ping(ctx)
	if err != nil {
		t.Fatal(err)
	}

	q := db.New(d)

	// 批量创建笔记
	q.CreateNotes(ctx, []db.CreateNotesParams{
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
			Title:       "标题",
			Description: "",
			TagIds:      []int64{1, 2, 3},
			Type:        db.NoteTypeNormal,
			IsPrivacy:   true,
			SourceType:  db.SourceTypeUnset,
			SourceUrl:   "",
		},
		{
			Title:       "标题",
			Description: "",
			TagIds:      []int64{1, 2, 3},
			Type:        db.NoteTypeNormal,
			IsPrivacy:   true,
			SourceType:  db.SourceTypeUnset,
			SourceUrl:   "1",
		},
		{
			Title:       "标题",
			Description: "",
			TagIds:      []int64{1, 2, 3},
			Type:        db.NoteTypeNormal,
			IsPrivacy:   true,
			SourceType:  db.SourceTypeUnset,
			SourceUrl:   "2",
		},
	}).QueryRow(func(i int, noteID int64, err error) {
		if err != nil {
			t.Fatalf("QueryRow %d (id=%d): %v", i, noteID, err)
		}
		t.Logf("#%d note id %d", i, noteID)
	})

	//n, err := q.GetNote(ctx, id)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//assert.Equal(t, "标题", n.Title)
	//assert.Equal(t, []int64{1, 2, 3}, n.TagIds)
	//assert.False(t, n.PostTime.Valid)
	//assert.Equal(t, db.NoteTypeNormal, n.Type)
	//assert.True(t, n.IsPrivacy)
	//assert.False(t, n.IsEnabled)
}
