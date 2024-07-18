package mysql

import (
	"context"
	"fmt"
	"log"

	"apu/note/mysql/ent"
)

// NewEntClient 使用默认配置创建一个 ent 客户端。
func NewEntClient() *ent.Client {
	entClient, err := ent.Open(
		"mysql",
		fmt.Sprintf("root:asdfasdf@tcp(127.0.0.1:3306)/apu?charset=utf8mb4&parseTime=true"),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err = entClient.Schema.Create(context.TODO()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return entClient
}
