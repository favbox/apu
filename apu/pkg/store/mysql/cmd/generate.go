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

	g.WithOpts(gen.FieldModify(func(field gen.Field) gen.Field {
		if field.Type == "int32" {
			field.Type = "int"
		}
		return field
	}))

	gormdb, _ := gorm.Open(mysql.Open("root:asdfasdf@tcp(127.0.0.1:3306)/apu?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"))
	g.UseDB(gormdb)

	g.ApplyBasic(
		g.GenerateModel("author"),
		g.GenerateModel("category"),
		g.GenerateModel("image", gen.FieldType("note_id", "uint64"), gen.FieldType("uid", "uint64")),
		g.GenerateModel("interact", gen.FieldType("note_id", "uint64")),
		g.GenerateModel("note", gen.FieldType("id", "uint64")),
		g.GenerateModel("note_category"),
		g.GenerateModel("note_content", gen.FieldType("id", "uint64")),
		g.GenerateModel("note_tag"),
		g.GenerateModel("pipeline", gen.FieldType("id", "uint64")),
		g.GenerateModel("tag"),
		g.GenerateModel("video", gen.FieldType("uid", "uint64")),
		g.GenerateModel("weixin_mp"),
		g.GenerateModel("weixin_request"),
	)
	g.Execute()
}
