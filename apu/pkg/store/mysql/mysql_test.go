package mysql_test

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"apu/pkg/store/mysql"
	"apu/pkg/store/mysql/model"
	"apu/pkg/store/mysql/query"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/clause"
)

func TestQuery(t *testing.T) {
	mysql.Init()
	count, err := query.WeixinRequest.Count()
	if err != nil {
		t.Fatal(err)
	}
	assert.GreaterOrEqual(t, count, int64(0))
}

func TestBefore(t *testing.T) {
	fmt.Print(Year())
}

func Year() int64 {
	return time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local).Unix()
}

func TestBulkInsert(t *testing.T) {
	mysql.Init()

	urls := []string{"123", "1", "2", "3"}
	existingOriginalUrls, err := query.OriginalURL.
		Where(query.OriginalURL.URL.In(urls...)).
		Select(query.OriginalURL.URL).
		Find()
	assert.Nil(t, err)
	var existingUrls []string
	for _, existsUrl := range existingOriginalUrls {
		existingUrls = append(existingUrls, existsUrl.URL)
	}

	var newUrls []*model.OriginalURL
	for _, url := range urls {
		if !slices.Contains(existingUrls, url) {
			newUrls = append(newUrls, &model.OriginalURL{URL: url})
		}
	}

	if len(newUrls) > 0 {
		err = query.OriginalURL.
			Create(newUrls...)
		assert.Nil(t, err)
	}

	for _, u := range newUrls {
		fmt.Println(u.ID, u.URL)
	}
}

func TestBulkInsert2(t *testing.T) {
	mysql.Init()

	urls := []string{"123", "1", "2", "3"}
	us := make([]*model.OriginalURL, 0, len(urls))
	for _, u := range urls {
		us = append(us, &model.OriginalURL{
			URL: u,
		})
	}

	var newIDs []int64
	q := query.Use(mysql.DB)
	err := q.Transaction(func(tx *query.Query) error {
		err := tx.OriginalURL.
			Clauses(clause.OnConflict{
				DoNothing: true,
			}).
			Create(us...)
		if err != nil {
			return err
		}

		one, err := tx.OriginalURL.Where(tx.OriginalURL.URL.In(urls...)).Select(tx.OriginalURL.ID.Max()).First()
		if err != nil {
			return err
		}
		two, err := tx.OriginalURL.Select(tx.OriginalURL.ID.Max()).First()
		if err != nil {
			return err
		}

		newURLs, err := tx.OriginalURL.Where(
			tx.OriginalURL.URL.In(urls...),
			tx.OriginalURL.ID.Gt(one.ID),
			tx.OriginalURL.ID.Lt(two.ID),
		).Find()
		if err != nil {
			return err
		}
		for _, n := range newURLs {
			newIDs = append(newIDs, n.ID)
		}

		return nil
	})
	if err != nil {
		return
	}

	for _, newID := range newIDs {
		fmt.Println(newID)
	}
}
