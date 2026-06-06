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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%v&parseTime=%v",
		config.GetCfg().MySQL.User, config.GetCfg().MySQL.Password, config.GetCfg().MySQL.Host, config.GetCfg().MySQL.Port,
		config.GetCfg().MySQL.Database, config.GetCfg().MySQL.Charset, config.GetCfg().MySQL.ParseTime)
	logger.Info("dsn: [REDACTED]")
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		logger.Fatal("InitDb open error", zap.Error(err))
	}
	cfg := config.Cfg.MySQL
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	if err := db.Ping(); err != nil {
		logger.Fatal("InitDb ping error", zap.Error(err))
	}
	logger.Info("database connection established",
		zap.Int("max_open_conns", cfg.MaxOpenConns),
		zap.Int("max_idle_conns", cfg.MaxIdleConns),
		zap.Duration("conn_max_lifetime", cfg.ConnMaxLifetime),
		zap.Duration("conn_max_idle_time", cfg.ConnMaxIdleTime),
	)
	return db
}

type MySQLdb struct {
	db *sqlx.DB
}

func NewMySQLdb(db *sqlx.DB) *MySQLdb {
	return &MySQLdb{db: db}
}
