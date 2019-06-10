package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	DB *sqlx.DB
)

func Init(dns string) error {
	var err error
	DB, err = sqlx.Open("mysql", dns)
	if err != nil {
		return err
	}

	err = DB.Ping() //防止数据库账号密码不对没有报错
	if err != nil {
		return err
	}

	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(16)
	return nil
}
