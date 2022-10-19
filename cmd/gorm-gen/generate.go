package main

import (
	"beneburg/pkg/database"
	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:           "pkg/database/query",
		FieldNullable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})
	g.ApplyBasic(database.Models...)
	g.Execute()
}
