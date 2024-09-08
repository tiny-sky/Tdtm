package mysqlx

import (
	"context"
	"time"

	"github.com/tiny-sky/Tdtm/core/dao"
	"github.com/tiny-sky/Tdtm/core/dao/sql"
	"github.com/tiny-sky/Tdtm/tools"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Settings struct {
	DbURL        string `json:"dbURL" yaml:"dbURL"`
	MaxLifetime  int    `json:"maxLifetime" yaml:"maxLifetime"`
	MaxIdleConns int    `json:"maxIdleConns" yaml:"maxIdleConns"`
	MaxOpenConns int    `json:"maxOpenConns" yaml:"maxOpenConns"`
}

func (m *Settings) Init() {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               m.DbURL,
		DefaultStringSize: 256,
	}), &gorm.Config{SkipDefaultTransaction: true})
	tools.WrapDbErr(err)

	d, err := db.DB()
	tools.ErrToPanic(err)

	if m.MaxLifetime > 0 {
		d.SetConnMaxLifetime(time.Duration(m.MaxLifetime) * time.Second)
	}
	if m.MaxOpenConns > 0 {
		d.SetMaxOpenConns(m.MaxOpenConns)
	}

	if m.MaxIdleConns > 0 {
		d.SetMaxIdleConns(m.MaxIdleConns)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	for _, tableSql := range sql.Sql() {
		_, err := d.ExecContext(ctx, tableSql)
		tools.ErrToPanic(err)
	}

	dao.SetTransaction(NewDao(db))
}
