package re

import (
	"context"
	"time"
)

func (rdb *Redis) UserTokenSet(ctx context.Context, refreshToken string, userId string) error {
	err := rdb.redis.Set(ctx, "refresh:"+refreshToken, userId, 7*24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rdb *Redis) UserGetByRefreshToken(ctx context.Context, refreshToken string) (userId string, err error) {
	userId, err = rdb.redis.Get(ctx, "refresh:"+refreshToken).Result()
	if err != nil {
		return userId, err
	}
	return userId, nil
}
func (rdb *Redis) UserTokenDelete(ctx context.Context, refreshToken string) error {
	err := rdb.redis.Del(ctx, "refresh:"+refreshToken).Err()
	if err != nil {
		return err
	}
	return nil
}
