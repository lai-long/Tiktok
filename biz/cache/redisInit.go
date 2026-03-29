package cache

import (
	"Tiktok/pkg/config"
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	redis *redis.Client
}

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
		log.Fatalf("re 连接失败", err)
	}
	log.Println("re 连接成功")
	return rdb
}
