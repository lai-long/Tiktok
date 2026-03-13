package re

import (
	"Tiktok/pkg/conf"
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

type Redis struct {
	redis *redis.Client
}

func NewRedis(client *redis.Client) *Redis {
	return &Redis{redis: client}
}
func InitRedis() *redis.Client {
	redisAddr := fmt.Sprintf("%s:%d", conf.Cfg.Redis.Host, conf.Cfg.Redis.Port)
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: conf.Cfg.Redis.Password,
		DB:       conf.Cfg.Redis.Database,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("re 连接失败", err)
	}
	log.Println("re 连接成功")
	return rdb
}
