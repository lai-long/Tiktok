package cache

import (
	"context"
	"log"
	"time"
)

func (r *Redis) SaveOfflineMsg(id, content string) error {
	ctx := context.Background()
	key := "offline:" + id
	err := r.redis.RPush(ctx, key, content).Err()
	if err != nil {
		log.Println("Save offline message error:", err)
		return err
	}
	err = r.redis.Expire(ctx, key, 72*time.Hour).Err()
	if err != nil {
		log.Println("Set expire error:", err)
		return err
	}
	return nil
}
func (r *Redis) FetchOfflineMsg(id string) ([]string, error) {
	ctx := context.Background()
	key := "offline:" + id
	messages, err := r.redis.LRange(ctx, key, 0, 100).Result()
	if err != nil {
		return []string{}, err
	}
	err = r.redis.Del(ctx, key).Err()
	if err != nil {
		log.Println("redis FetchOfflineMsg del messages err:", err)
		return []string{}, err
	}
	return messages, nil
}
