package redis

import (
	"Tiktok/pkg/conf"
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedis() *redis.Client {
	redisAddr := fmt.Sprintf("%s:%d", conf.Cfg.Redis.Host, conf.Cfg.Redis.Port)
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: conf.Cfg.Redis.Password,
		DB:       conf.Cfg.Redis.Database,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis 连接失败", err)
	}
	log.Println("redis 连接成功")
	return rdb
}
