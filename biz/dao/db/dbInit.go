package db

import (
	"Tiktok/pkg/conf"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func InitDb() *sqlx.DB {
	var err error
	db, err = sqlx.Connect("mysql", conf.DSN)
	if err != nil {
		log.Fatalf("InitDb err: %v", err)
	}
	log.Println("数据库连接成功")
	return db
}
