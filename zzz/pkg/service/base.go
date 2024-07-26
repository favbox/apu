package service

import (
	"context"

	"apu/app/note/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Base 是一个服务基类，为服务提供事务型数据操作的支持。
type Base struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
}

// Tx 执行一个事务型数据服务。
func (s *Base) Tx(ctx context.Context, fn func(qtx *db.Queries)) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.Queries.WithTx(tx)
	fn(qtx)
	return tx.Commit(ctx)
}
