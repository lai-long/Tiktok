// cache use redis

package cache

import (
	"Tiktok/internal/config"
	"Tiktok/pkg/logger"
	"context"
	"fmt"
	"go.uber.org/zap"

	"github.com/redis/go-redis/v9"
)

// Redis is a struct encapsulation redis client
type Redis struct {
	redis *redis.Client
}

// NewRedis creat new Redis struct
func NewRedis(client *redis.Client) *Redis {
	return &Redis{redis: client}
}
func InitRedis() *redis.Client {
	cfg := config.Cfg.Redis
	redisAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:            redisAddr,
		Password:        cfg.Password,
		DB:              cfg.Database,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		PoolTimeout:     cfg.PoolTimeout,
		ConnMaxIdleTime: cfg.ConnMaxIdleTime,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Fatal("redis connection failed", zap.Error(err))
	}
	logger.Info("redis connection established",
		zap.Int("pool_size", cfg.PoolSize),
		zap.Int("min_idle_conns", cfg.MinIdleConns),
		zap.Int("max_idle_conns", cfg.MaxIdleConns),
		zap.Duration("pool_timeout", cfg.PoolTimeout),
		zap.Duration("conn_max_idle_time", cfg.ConnMaxIdleTime),
		zap.Duration("conn_max_lifetime", cfg.ConnMaxLifetime),
	)
	return rdb
}
