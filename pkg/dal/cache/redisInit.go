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
	var rdb *redis.Client
	redisAddr := fmt.Sprintf("%s:%d", config.GetCfg().Redis.Host, config.GetCfg().Redis.Port)
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: config.GetCfg().Redis.Password,
		DB:       config.GetCfg().Redis.Database,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Fatal("redis connection failed", zap.Error(err))
	}
	logger.Info("redis connection established")
	return rdb
}
