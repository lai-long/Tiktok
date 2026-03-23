package re

import (
	"context"
	"log"
	"time"
)

func SaveOfflineMsg(id, content string) {
	ctx := context.Background()
	key := "offline:" + id
	rdb.RPush(ctx, key, content)
	rdb.Expire(ctx, key, 72*time.Hour)
}
func FetchOfflineMsg(id string) ([]string, int) {
	ctx := context.Background()
	key := "offline:" + id
	messages, err := rdb.LRange(ctx, key, 0, 100).Result()
	if err != nil {
		log.Println("redis FetchOfflineMsg err:", err)
		return []string{}, 0
	}
	err = rdb.Del(ctx, key).Err()
	if err != nil {
		log.Println("redis FetchOfflineMsg del messages err:", err)
		return []string{}, 0
	}
	return messages, len(messages)
}
