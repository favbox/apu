package mysql

import (
	"log"
	"sync"

	"apu/pkg/store/mysql/query"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dsn = "root:asdfasdf@tcp(127.0.0.1:3306)/apu?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"

var DB *gorm.DB
var once sync.Once

func Init() {
	once.Do(func() {
		var err error
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			Logger:                 logger.Default.LogMode(logger.Warn),
		})
		if err != nil {
			panic(err)
		}

		query.SetDefault(DB)
		log.Println("mysql 已初始化")
	})
}
