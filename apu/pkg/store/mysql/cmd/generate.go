package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "pkg/store/mysql/query",
		ModelPkgPath:  "pkg/store/mysql/model",
		FieldNullable: false,
		Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	gormdb, _ := gorm.Open(mysql.Open("root:asdfasdf@tcp(127.0.0.1:3306)/apu?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"))
	g.UseDB(gormdb)

	g.ApplyBasic(
		g.GenerateModel("wexin_request"),
		g.GenerateModel("document", gen.FieldType("key", "uint64")),
		g.GenerateModel("image", gen.FieldType("key", "uint64")),
		g.GenerateModel("interaction", gen.FieldType("key", "uint64")),
	)
	g.Execute()
}
