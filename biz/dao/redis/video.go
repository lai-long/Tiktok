package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func VideoHotSet(ctx context.Context, key string, member interface{}, score float64) error {
	if err := rdb.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: member,
	}).Err(); err != nil {
		return err
	}
	return nil
}
func VideoHotGet(ctx context.Context, key string, pageNum int, pageSize int) ([]redis.Z, error) {
	z, err := rdb.ZRevRangeWithScores(ctx, key, int64(pageSize*pageNum), int64(pageSize+pageSize*pageNum)).Result()
	if err != nil {
		return nil, err
	}
	return z, nil
}
