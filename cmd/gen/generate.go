package main

import (
	"github.com/tiny-sky/Tdtm/core/dao/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	outPath := "../../dao/mysqlx/query"

	g := gen.NewGenerator(gen.Config{
		OutPath: outPath,
		OutFile: outPath + "/query.go",
	})

	db, err := gorm.Open(mysql.Open("root:tdtm@(127.0.0.1:3306)/Tdtm?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic("gorm err")
	}

	//复用已有的SQL连接配置db(*gorm.DB)
	g.UseDB(db)

	// 生成对应 CRUD 代码
	g.ApplyBasic(entity.Global{}, entity.Branch{})

	g.Execute()
}
