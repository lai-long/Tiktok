package redis

import (
	"Tiktok/pkg/conf"
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis,
		Password: "",
		DB:       0,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis 连接失败", err)
	}
	log.Println("redis 连接成功")
	return rdb
}
