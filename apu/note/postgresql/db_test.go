package postgresql_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"apu/note/postgresql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestAuthors(t *testing.T) {
	ctx := context.Background()
	uri := "postgresql://zs@127.0.0.1:5432/zs?connect_timeout=10"
	db, err := pgx.Connect(ctx, uri)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)

	err = db.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	q := postgresql.New(db)

	// begin tx
	tx, err := db.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	qtx := q.WithTx(tx)

	// create an author
	a, err := qtx.CreateAuthor(ctx, "Unknown Master")
	if err != nil {
		t.Fatal(err)
	}

	// batch insert new books
	now := pgtype.Timestamptz{
		Time:  time.Now(),
		Valid: true,
	}
	newBookParams := []postgresql.CreateBookParams{
		{
			AuthorID:  a.AuthorID,
			Isbn:      "1",
			BookType:  postgresql.BookTypeFICTION,
			Title:     "my book title",
			Year:      2016,
			Available: now,
			Tags:      []string{},
		},
		{
			AuthorID:  a.AuthorID,
			Isbn:      "2",
			BookType:  postgresql.BookTypeFICTION,
			Title:     "the second book title",
			Year:      2016,
			Available: now,
			Tags:      []string{"cool", "unique"},
		},
		{
			AuthorID:  a.AuthorID,
			Isbn:      "3",
			BookType:  postgresql.BookTypeFICTION,
			Title:     "the third book title",
			Year:      2001,
			Available: now,
			Tags:      []string{"cool"},
		},
		{
			AuthorID:  a.AuthorID,
			Isbn:      "4",
			BookType:  postgresql.BookTypeNONFICTION,
			Title:     "4th place finisher",
			Year:      2011,
			Available: now,
			Tags:      []string{"other"},
		},
	}
	newBooks := make([]postgresql.Book, len(newBookParams))
	var cnt int
	qtx.CreateBook(ctx, newBookParams).QueryRow(func(i int, book postgresql.Book, err error) {
		if err != nil {
			t.Fatalf("failed inserting book (#%d): %s", i, err)
		}
		newBooks[i] = book
		cnt = i
		fmt.Println(book.BookID, book.AuthorID, book.Isbn, book.Available)
	})
	// first i was 0, so add 1
	cnt++
	newBooksExpected := len(newBooks)
	assert.Equalf(t, newBooksExpected, cnt, "expected to insert %d books; got %d", newBooksExpected, cnt)

	// batch update the title and tags
	updateBookParams := []postgresql.UpdateBookParams{
		{
			Title:  "changed second title",
			Tags:   []string{"cool", "disastor"},
			BookID: newBooks[1].BookID,
		},
	}
	qtx.UpdateBook(ctx, updateBookParams).Exec(func(i int, err error) {
		if err != nil {
			t.Fatalf("failed updating book (%d): %s", updateBookParams[i].BookID, err)
		}
	})

	// batch many to retrieve books by year
	searchBooksByTitleYearParams := []int32{2001, 2016}
	var books0 []postgresql.Book
	qtx.BooksByYear(ctx, searchBooksByTitleYearParams).Query(func(i int, books []postgresql.Book, err error) {
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("num books for %d: %d", searchBooksByTitleYearParams[i], len(books))
		books0 = append(books0, books...)
	})

	for _, book := range books0 {
		t.Logf("Book %d (%s): %s available: %s", book.BookID, book.BookType, book.Title, book.Available.Time.Format(time.DateTime))
		author, err := qtx.GetAuthor(ctx, book.AuthorID)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Book %d author: %s", book.BookID, author.Name)
	}

	// batch delete books
	deleteBooksParams := make([]int32, len(newBooks))
	for i, b := range newBooks {
		deleteBooksParams[i] = b.BookID
	}
	batchDelete := q.DeleteBook(ctx, deleteBooksParams)
	numDeletesProcessed := 0
	wantNumDeletesProcessed := 2
	batchDelete.Exec(func(i int, err error) {
		if err != nil && err.Error() != "batch already closed" {
			t.Fatalf("error deleteing book %d: %s", deleteBooksParams[i], err)
		}

		if err == nil {
			numDeletesProcessed++
		}

		if i == wantNumDeletesProcessed-1 {
			//	close batch operation before processing all errors from delete operation
			if err := batchDelete.Close(); err != nil {
				t.Fatalf("failed to close batch delete: %s", err)
			}
		}
	})
	if numDeletesProcessed != wantNumDeletesProcessed {
		t.Fatalf("expected Close to short-circuit record processing (expected %d, got %d)", wantNumDeletesProcessed, numDeletesProcessed)
	}

	tx.Commit(ctx)
}

func TestQueries_CreateAuthors(t *testing.T) {
	ctx := context.Background()
	uri := "postgresql://zs@127.0.0.1:5432/zs?connect_timeout=10"
	db, err := pgx.Connect(ctx, uri)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)

	err = db.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	q := postgresql.New(db)

	createAuthorsParams := postgresql.CreateAuthorsParams{
		Ids:   []int64{1, 2, 3},
		Names: []string{"John", "Smith", "Mike"},
	}
	err = q.CreateAuthors(ctx, createAuthorsParams)
	if err != nil {
		t.Fatal(err)
	}
}
