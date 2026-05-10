package cache

import (
	"context"

	"github.com/pkg/errors"
)

func (rdb *Redis) VideoLikeSAdd(ctx context.Context, userId string, videoId string) error {
	key := "user:like:videos:" + userId
	if err := rdb.redis.SAdd(ctx, key, videoId).Err(); err != nil {
		return errors.Wrap(err, "sAdd video like")
	}
	return nil
}

func (rdb *Redis) VideoDislikeSRem(ctx context.Context, userId string, videoId string) error {
	key := "user:like:videos:" + userId
	if err := rdb.redis.SRem(ctx, key, videoId).Err(); err != nil {
		return errors.Wrap(err, "sRem video like")
	}
	return nil
}

func (rdb *Redis) VideoLikeGet(ctx context.Context, userId string) ([]string, error) {
	key := "user:like:videos:" + userId
	results, err := rdb.redis.SMembers(ctx, key).Result()
	return results, errors.Wrap(err, "get video likes")
}
