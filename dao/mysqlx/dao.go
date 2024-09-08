package mysqlx

import (
	"github.com/tiny-sky/Tdtm/dao"
	"gorm.io/gorm"
)

var db *gorm.DB

type Dao struct {
	dao.BranchDao
	dao.GlobalDao
}

func NewDao(gorm *gorm.DB) Dao {
	db = gorm
	return Dao{
		BranchDao: NewBranchImpl(),
		GlobalDao: NewGlobalImpl(),
	}
}
