package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool 创建一个可复用的数据库连接池。
func NewPool() (*pgxpool.Pool, error) {
	var ctx = context.Background()
	pool, err := pgxpool.New(ctx, "host=localhost port=5432 user=zs password=zs sslmode=disable")
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
