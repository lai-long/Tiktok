package cache

import (
	"Tiktok/pkg/entity"
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func (rdb *Redis) VideoHotSet(ctx context.Context, key string, member interface{}, score float64) error {
	if err := rdb.redis.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: member,
	}).Err(); err != nil {
		return err
	}
	return nil
}
func (rdb *Redis) VideoHotGet(ctx context.Context, key string, pageNum int64, pageSize int64) ([]redis.Z, error) {
	start, end := videoHotRange(pageNum, pageSize)
	z, err := rdb.redis.ZRevRangeWithScores(ctx, key, start, end).Result()
	if err != nil {
		return nil, err
	}
	return z, nil
}

func videoHotRange(pageNum, pageSize int64) (int64, int64) {
	start := pageNum * pageSize
	end := start + pageSize - 1
	return start, end
}

func (rdb *Redis) VideoHotIncrBy(ctx context.Context, key string, videoID string, delta float64) error {
	return rdb.redis.ZIncrBy(ctx, key, delta, videoID).Err()
}

func (rdb *Redis) VideoInfoSet(ctx context.Context, videoID string, video *entity.VideoEntity) error {
	data, err := json.Marshal(video)
	if err != nil {
		return errors.Wrap(err, "json marshal")
	}
	duration := 30*time.Minute + time.Duration(rand.Intn(5))*time.Second
	key := "video:info:" + videoID
	err = rdb.redis.Set(ctx, key, data, duration).Err()
	if err != nil {
		return errors.Wrap(err, "set cache video")
	}
	return nil
}
func (rdb *Redis) VideoInfoGet(ctx context.Context, videoID string) (*entity.VideoEntity, error) {
	key := "video:info:" + videoID
	data, err := rdb.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "get cache video")
	}
	var videoEntity entity.VideoEntity
	err = json.Unmarshal(data, &videoEntity)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	}
	return &videoEntity, nil
}

func (rdb *Redis) VideoInfoDelete(ctx context.Context, videoID string) error {
	return rdb.redis.Del(ctx, "video:info:"+videoID).Err()
}
