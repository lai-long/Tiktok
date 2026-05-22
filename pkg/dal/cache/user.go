package cache

import (
	"Tiktok/pkg/entity"
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

func (rdb *Redis) UserTokenSet(ctx context.Context, refreshToken string, userID string) error {
	duration := 168*time.Hour + time.Duration(rand.Intn(168))*time.Hour
	err := rdb.redis.Set(ctx, "refresh:"+refreshToken, userID, duration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rdb *Redis) UserGetByRefreshToken(ctx context.Context, refreshToken string) (userID string, err error) {
	userID, err = rdb.redis.Get(ctx, "refresh:"+refreshToken).Result()
	if err != nil {
		return userID, err
	}
	return userID, nil
}
func (rdb *Redis) UserTokenDelete(ctx context.Context, refreshToken string) error {
	err := rdb.redis.Del(ctx, "refresh:"+refreshToken).Err()
	if err != nil {
		return err
	}
	return nil
}
func (rdb *Redis) SetCachedUserInfo(ctx context.Context, userId string, user *entity.UserEntity) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	expiration := 30*time.Minute + time.Duration(rand.Intn(5))*time.Minute
	err = rdb.redis.Set(ctx, "user:info:"+userId, data, expiration).Err()
	if err != nil {
		return errors.Wrap(err, "set cache user info")
	}
	return nil
}

func (rdb *Redis) GetCachedUserInfo(ctx context.Context, userId string) (*entity.UserEntity, error) {
	data, err := rdb.redis.Get(ctx, "user:info:"+userId).Bytes()
	if err != nil {
		return nil, err
	}
	var info entity.UserEntity
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	}
	return &info, nil
}

func (rdb *Redis) DelCachedUserInfo(ctx context.Context, userId string) error {
	err := rdb.redis.Del(ctx, "user:info:"+userId).Err()
	if err != nil {
		return errors.Wrap(err, "del cache user info")
	}
	return nil
}
