package db

import (
	"Tiktok/pkg/conf"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func InitDb() error {
	var err error
	db, err = sqlx.Connect("mysql", conf.DSN)
	if err != nil {
		return fmt.Errorf("InitDb err: %v", err)
	}
	fmt.Println("数据库连接成功")
	return nil
}
