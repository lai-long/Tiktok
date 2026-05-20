package service

import (
	"Tiktok/kitex_gen/mfa"
	"Tiktok/kitex_gen/mfa/mfaservice"
	"Tiktok/kitex_gen/user"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/entity"
	"Tiktok/pkg/utils"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

type UserRedis interface {
	UserTokenSet(ctx context.Context, refreshToken string, userId string) error
	UserGetByRefreshToken(ctx context.Context, refreshToken string) (userId string, err error)
	UserTokenDelete(ctx context.Context, refreshToken string) error
	GetCachedUserInfo(ctx context.Context, userId string) (*entity.UserEntity, error)
	SetCachedUserInfo(ctx context.Context, userId string, info *entity.UserEntity) error
	DelCachedUserInfo(ctx context.Context, userId string) error
}

type UserDatabase interface {
	CreateUser(user entity.UserEntity) error
	GetUserByUsername(username string) (entity.UserEntity, error)
	GetUserByUserId(userId string) (entity.UserEntity, error)
	UpdateUserAvatar(url string, userId interface{}) error
}
type UserRepo struct {
	userDb    UserDatabase
	mfaClient mfaservice.Client
	redis     UserRedis
}

func NewUserRepo(userDb UserDatabase, mfaClient mfaservice.Client, redis UserRedis) *UserRepo {
	return &UserRepo{userDb: userDb, mfaClient: mfaClient, redis: redis}
}

func (s *UserRepo) IsUsernameExists(username string) (bool, error) {
	_, err := s.userDb.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, errors.Wrap(err, "get user by username")
	}
	return true, nil
}

func (s *UserRepo) Register(userName string, password string) (int32, error) {
	var userEntity entity.UserEntity
	var err error
	exists, err := s.IsUsernameExists(userName)
	if err != nil {
		return consts.UserDBSelectError, errors.Wrap(err, "IsUsernameExists error")
	}
	if exists {
		return consts.UserNameExists, nil
	}
	userEntity.ID = utils.IDGenerate()
	userEntity.Username = userName
	userEntity.Password, err = utils.HashPassword(password)
	if err != nil {
		return consts.UserHashError, errors.Wrap(err, "utils.HashPassword error")
	}
	if err = s.userDb.CreateUser(userEntity); err != nil {
		return consts.UserDBInsertError, errors.Wrap(err, "CreateUser error")
	}
	return consts.Success, nil
}

func (s *UserRepo) Login(userName, password, mfaCode string, ctx context.Context) (int32, *user.UserInfo, string, string, error) {
	userEntity, err := s.userDb.GetUserByUsername(userName)
	if errors.Is(err, sql.ErrNoRows) {
		return consts.UserNotExists, &user.UserInfo{}, "", "", nil
	}
	if err != nil {
		return consts.UserDBSelectError, &user.UserInfo{}, "", "", errors.Wrap(err, "GetUserByUsername failed")
	}
	err = utils.CheckPasswordHash(userEntity.Password, password)
	if err != nil {
		return consts.UserPasswordError, &user.UserInfo{}, "", "", errors.Wrap(err, "->login: check password failed")
	}
	userInfo := userEntity.ToUserInfo()
	resp, err := s.mfaClient.MfaConfirm(ctx, &mfa.MfaConfirmReq{
		UserID: userInfo.ID,
		QrCode: mfaCode,
	})
	if err != nil || resp == nil {
		return consts.MfaDBSelectError, &user.UserInfo{}, "", "", err
	}
	if resp.Code != consts.Success {
		return resp.Code, &user.UserInfo{}, "", "", nil
	}
	reToken, acToken, err := utils.GenerateTokens(userInfo)
	if err != nil {
		return consts.GenerateTokenError, userInfo, reToken, acToken, errors.Wrap(err, "->login 生成token错误")
	}
	err = s.redis.UserTokenSet(ctx, reToken, userInfo.ID)
	if err != nil {
		return consts.UserRedisSetError, &user.UserInfo{}, "", "", errors.Wrap(err, "->login 将refresh token存入redis错误")
	}
	return consts.Success, userInfo, reToken, acToken, nil
}

func (s *UserRepo) UserInfo(ctx context.Context, userId string) (*user.UserInfo, int32, error) {
	userEntity, err := s.redis.GetCachedUserInfo(ctx, userId)
	if err == nil && userEntity != nil {
		return userEntity.ToUserInfo(), consts.Success, nil
	}
	userEntity2, err := s.userDb.GetUserByUserId(userId)
	if err != nil {
		return &user.UserInfo{}, consts.UserDBSelectError, errors.Wrap(err, "->UserInfo GetUserByUserId error")
	}
	userInfo := userEntity2.ToUserInfo()
	err = s.redis.SetCachedUserInfo(ctx, userId, &userEntity2)
	if err != nil {
		return &user.UserInfo{}, consts.UserDBSelectError, errors.Wrap(err, "->UserInfo SetCachedUserInfo error")
	}
	return userInfo, consts.Success, nil
}

func (s *UserRepo) UserAvatar(url string, userID string) (int32, *user.UserInfo, error) {
	err := s.userDb.UpdateUserAvatar(url, userID)
	if err != nil {
		return consts.UserDBUpdateError, &user.UserInfo{}, errors.Wrap(err, "->userinfo 更新头像错误")
	}
	if err := s.redis.DelCachedUserInfo(context.Background(), userID); err != nil {
		return consts.UserRedisDelError, &user.UserInfo{}, errors.Wrap(err, "->userinfo 删除缓存错误")
	}
	userEntity, err := s.userDb.GetUserByUserId(userID)
	if err != nil {
		return consts.UserDBSelectError, &user.UserInfo{}, errors.Wrap(err, "->userinfo get user by userid failed")
	}
	userInfo := userEntity.ToUserInfo()
	return consts.Success, userInfo, nil
}

func (s *UserRepo) RefreshToken(ctx context.Context, refreshToken string) (int32, string, string, error) {
	userId, err := s.redis.UserGetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return consts.UserRedisGetError, "", "", errors.Wrap(err, "->RefreshToken GetUserIDByRefreshToken error")
	}
	userEntity, err := s.userDb.GetUserByUserId(userId)
	if err != nil {
		return consts.UserDBSelectError, "", "", errors.Wrap(err, "->RefreshToken GetUserByUserId error")
	}
	userInfo := userEntity.ToUserInfo()
	refreshToken2, accessToken, err := utils.GenerateTokens(userInfo)
	if err != nil {
		return consts.GenerateTokenError, "", "", errors.Wrap(err, "->RefreshToken GenerateTokens error")
	}
	err = s.redis.UserTokenDelete(ctx, refreshToken)
	if err != nil {
		return consts.UserRedisDelError, "", "", errors.Wrap(err, "->RefreshToken DeleteToken error")
	}
	err = s.redis.UserTokenSet(ctx, refreshToken2, userInfo.ID)
	if err != nil {
		return consts.UserRedisSetError, "", "", errors.Wrap(err, "->RefreshToken SetToken error")
	}
	return consts.Success, refreshToken2, accessToken, nil
}
