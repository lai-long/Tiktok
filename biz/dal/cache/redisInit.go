// cache use redis

package cache

import (
	"Tiktok/pkg/config"
	"context"
	"fmt"
	"log"

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
	redisAddr := fmt.Sprintf("%s:%d", config.Cfg.Redis.Host, config.Cfg.Redis.Port)
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: config.Cfg.Redis.Password,
		DB:       config.Cfg.Redis.Database,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("re 连接失败 错误: %v", err)
	}
	log.Println("re 连接成功")
	return rdb
}
