package dao

import (
	"Tiktok/pkg/config"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func InitDb() *sqlx.DB {
	var db *sqlx.DB
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset%v&parseTime=%v", config.Cfg.MySQL.User, config.Cfg.MySQL.Password, config.Cfg.MySQL.Host, config.Cfg.MySQL.Port, config.Cfg.MySQL.Database, config.Cfg.MySQL.Charset, config.Cfg.MySQL.ParseTime)
	//	DSN       = "root:root@tcp(localhost:3306)/tiktok?charset=utf8&parseTime=True&loc=Local"
	log.Println("dsn:", dsn)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("InitDb err: %v", err)
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
