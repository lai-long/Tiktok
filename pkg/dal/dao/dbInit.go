package dao

import (
	"Tiktok/internal/config"
	"Tiktok/pkg/logger"
	"fmt"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func InitDb() *sqlx.DB {
	var db *sqlx.DB
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%v&parseTime=%v",
		config.GetCfg().MySQL.User, config.GetCfg().MySQL.Password, config.GetCfg().MySQL.Host, config.GetCfg().MySQL.Port,
		config.GetCfg().MySQL.Database, config.GetCfg().MySQL.Charset, config.GetCfg().MySQL.ParseTime)
	logger.Info("dsn: [REDACTED]")
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		logger.Fatal("InitDb error", zap.Error(err))
	}
	logger.Info("database connection established")
	return db
}

type MySQLdb struct {
	db *sqlx.DB
}

func NewMySQLdb(db *sqlx.DB) *MySQLdb {
	return &MySQLdb{db: db}
}
