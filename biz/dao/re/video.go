package re

import (
	"context"
	"time"

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
func (rdb *Redis) VideoHotGet(ctx context.Context, key string, pageNum int, pageSize int) ([]redis.Z, error) {
	z, err := rdb.redis.ZRevRangeWithScores(ctx, key, int64(pageSize*pageNum), int64(pageSize+pageSize*pageNum)).Result()
	if err != nil {
		return nil, err
	}
	return z, nil
}
func (rdb *Redis) UserTokenSet(ctx context.Context, refreshToken string, userId string) error {
	err := rdb.redis.Set(ctx, "refresh:"+refreshToken, userId, 7*24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

//func (rdb *Redis) UserGetByRefreshToken(ctx context.Context, refreshToken string) (userId string, err error) {
//	rdb.redis.Get(ctx, "refresh:"+refreshToken).Result()
//}
