package command

import (
	"errors"
	"fmt"
	"time"

	"apu/pkg/schema"
	"apu/pkg/source/weixin"
	"apu/pkg/store/mysql"
	"apu/pkg/store/mysql/model"
	"apu/pkg/store/mysql/query"
)

func syncWeixinAuthor(biz, mpName string) (author *model.Author, err error) {
	if biz == "" || mpName == "" {
		return nil, errors.New("biz或mpName均不可为空")
	}

	mysql.Init()
	q := query.Use(mysql.DB)
	err = q.Transaction(func(tx *query.Query) error {
		// 维护创作者
		mpUID := weixin.Biz2GhID(biz)
		if mpUID == 0 {
			return fmt.Errorf("biz2ghid 失败：%s", biz)
		}
		author, err = tx.Author.Where(
			tx.Author.Source.Eq(schema.SourceWeixin.Int()),
			tx.Author.UID.Eq(mpUID),
		).FirstOrInit()
		if err != nil {
			return err
		}
		if mpName != "" {
			author.Nickname = mpName
		}
		if err := tx.Author.Save(author); err != nil {
			return err
		}

		// 维护公众号
		mp, err := tx.WeixinMp.Where(tx.WeixinMp.Biz.Eq(biz)).FirstOrInit()
		if err != nil {
			return err
		}
		mp.UID = mpUID
		mp.AuthorID = author.ID
		if mpName != "" && mp.Name != mpName {
			mp.Name = mpName
		}
		if mp.ID == 0 {
			mp.LastPublishTime = time.Unix(0, 0)
		}
		if err := tx.WeixinMp.Save(mp); err != nil {
			return err
		}

		return nil
	})
	if author == nil {
		return nil, errors.New("未能获取到作者信息")
	}
	return
}
