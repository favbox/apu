package mysql

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	ctx := context.Background()

	client := NewEntClient()
	url, err := client.OriginalUrl.Query().Count(ctx)
	assert.Nil(t, err)
	fmt.Println(url)

	id, err := client.OriginalUrl.
		Create().
		SetURL("123").
		OnConflictColumns("url").
		UpdateUpdatedAt().
		ID(ctx)
	fmt.Println(id, err)

	err = client.OriginalUrl.
		CreateBulk(
			client.OriginalUrl.Create().SetURL("123"),
			client.OriginalUrl.Create().SetURL("456"),
			client.OriginalUrl.Create().SetURL("123"),
		).
		OnConflict().UpdateUpdatedAt().Exec(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
