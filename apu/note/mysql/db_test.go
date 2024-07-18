package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestQueries_CreateOriginalURL(t *testing.T) {
	dsn := "root:asdfasdf@tcp(127.0.0.1:3306)/apu?charset=utf8mb4&parseTime=true"
	db, err := sql.Open("mysql", dsn)
	assert.Nil(t, err)

	q := New(db)
	id, err := q.CreateOriginalURL(context.Background(), "123")
	assert.Nil(t, err)
	fmt.Println(id)
}
