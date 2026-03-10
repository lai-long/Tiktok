package db

import (
	"Tiktok/pkg/conf"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func InitDb() *sqlx.DB {
	var db *sqlx.DB
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset%v&parseTime=%v", conf.Cfg.MySQL.User, conf.Cfg.MySQL.Password, conf.Cfg.MySQL.Host, conf.Cfg.MySQL.Port, conf.Cfg.MySQL.Database, conf.Cfg.MySQL.Charset, conf.Cfg.MySQL.ParseTime)
	//	DSN       = "root:root@tcp(localhost:3306)/tiktok?charset=utf8&parseTime=True&loc=Local"
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Printf("InitDb err: %v", err)
	}
	log.Println("数据库连接成功")
	return db
}

type MySQLdb struct {
	db *sqlx.DB
}

func NewMySQLdb(db *sqlx.DB) *MySQLdb {
	return &MySQLdb{db: db}
}
