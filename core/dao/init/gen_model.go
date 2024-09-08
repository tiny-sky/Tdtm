package main

import (
	"github.com/tiny-sky/Tdtm/core/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	outPath := "./dao/mysqlx/query"

	g := gen.NewGenerator(gen.Config{
		OutPath: outPath,
		OutFile: outPath + "/query.go",
	})

	db, err := gorm.Open(mysql.Open("root:tdtm@(127.0.0.1:3306)/Tdtm?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic("gorm err")
	}
	g.UseDB(db)

	g.ApplyBasic(entity.Global{}, entity.Branch{})

	g.Execute()
}
