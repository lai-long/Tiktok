package service

import (
	"Tiktok/kitex_gen/user"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/entity"
	"Tiktok/pkg/utils"
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/pquerna/otp/totp"
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
	CreateUser(ctx context.Context, user entity.UserEntity) error
	GetUserByUsername(ctx context.Context, username string) (entity.UserEntity, error)
	GetUserByUserId(ctx context.Context, userId string) (entity.UserEntity, error)
	UpdateUserAvatar(ctx context.Context, url string, userId interface{}) error
	CheckMfaBind(ctx context.Context, userId string) (int, error)
	GetMfaSecret(ctx context.Context, userId string) (string, error)
}
type UserRepo struct {
	userDb UserDatabase
	redis  UserRedis
}

func NewUserRepo(userDb UserDatabase, redis UserRedis) *UserRepo {
	return &UserRepo{userDb: userDb, redis: redis}
}

func toUserInfo(e entity.UserEntity) *user.UserInfo {
	return &user.UserInfo{
		ID:        e.ID,
		Username:  e.Username,
		AvatarURL: utils.SignQiNiuURL(e.Avatar_url),
		CreatedAt: e.Created_at.String(),
		UpdatedAt: e.Updated_at.String(),
	}
}

func (s *UserRepo) IsUsernameExists(ctx context.Context, username string) (bool, error) {
	_, err := s.userDb.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, errors.Wrap(err, "get user by username")
	}
	return true, nil
}

func (s *UserRepo) Register(ctx context.Context, userName string, password string) (int32, error) {
	defer utils.TrackTime(ctx, "UserRegister")()
	var userEntity entity.UserEntity
	var err error
	exists, err := s.IsUsernameExists(ctx, userName)
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
	if err = s.userDb.CreateUser(ctx, userEntity); err != nil {
		return consts.UserDBInsertError, errors.Wrap(err, "CreateUser error")
	}
	return consts.Success, nil
}

func (s *UserRepo) Login(ctx context.Context, userName, password, mfaCode string) (int32, *user.UserInfo, string, string, error) {
	defer utils.TrackTime(ctx, "UserLogin")()
	userEntity, err := s.userDb.GetUserByUsername(ctx, userName)
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
	userInfo := toUserInfo(userEntity)
	code, err := s.mfaConfirm(ctx, userInfo.ID, mfaCode)
	if err != nil {
		return consts.MfaDBSelectError, &user.UserInfo{}, "", "", err
	}
	if code != consts.Success {
		return code, &user.UserInfo{}, "", "", nil
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
	defer utils.TrackTime(ctx, "UserInfo")()
	userEntity, err := s.redis.GetCachedUserInfo(ctx, userId)
	if err == nil && userEntity != nil {
		return toUserInfo(*userEntity), consts.Success, nil
	}
	userEntity2, err := s.userDb.GetUserByUserId(ctx, userId)
	if err != nil {
		return &user.UserInfo{}, consts.UserDBSelectError, errors.Wrap(err, "->UserInfo GetUserByUserId error")
	}
	userInfo := toUserInfo(userEntity2)
	err = s.redis.SetCachedUserInfo(ctx, userId, &userEntity2)
	if err != nil {
		return &user.UserInfo{}, consts.UserDBSelectError, errors.Wrap(err, "->UserInfo SetCachedUserInfo error")
	}
	return userInfo, consts.Success, nil
}

func (s *UserRepo) UserAvatar(ctx context.Context, url string, userID string) (int32, *user.UserInfo, error) {
	defer utils.TrackTime(ctx, "UserAvatar")()
	err := s.userDb.UpdateUserAvatar(ctx, url, userID)
	if err != nil {
		return consts.UserDBUpdateError, &user.UserInfo{}, errors.Wrap(err, "->userinfo 更新头像错误")
	}
	if err := s.redis.DelCachedUserInfo(ctx, userID); err != nil {
		return consts.UserRedisDelError, &user.UserInfo{}, errors.Wrap(err, "->userinfo 删除缓存错误")
	}
	userEntity, err := s.userDb.GetUserByUserId(ctx, userID)
	if err != nil {
		return consts.UserDBSelectError, &user.UserInfo{}, errors.Wrap(err, "->userinfo get user by userid failed")
	}
	userInfo := toUserInfo(userEntity)
	return consts.Success, userInfo, nil
}

func (s *UserRepo) RefreshToken(ctx context.Context, refreshToken string) (int32, string, string, error) {
	defer utils.TrackTime(ctx, "RefreshToken")()
	userId, err := s.redis.UserGetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return consts.UserRedisGetError, "", "", errors.Wrap(err, "->RefreshToken GetUserIDByRefreshToken error")
	}
	userEntity, err := s.userDb.GetUserByUserId(ctx, userId)
	if err != nil {
		return consts.UserDBSelectError, "", "", errors.Wrap(err, "->RefreshToken GetUserByUserId error")
	}
	userInfo := toUserInfo(userEntity)
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

func (s *UserRepo) mfaConfirm(ctx context.Context, userID string, mfaCode string) (int32, error) {
	isBind, err := s.userDb.CheckMfaBind(ctx, userID)
	if err != nil {
		return consts.MfaDBSelectError, errors.Wrap(err, "->check mfa bind error")
	}
	if isBind != 0 {
		if mfaCode == "" {
			return consts.MfaReqValidError, nil
		}
		mfaSecret, err := s.userDb.GetMfaSecret(ctx, userID)
		if err != nil {
			return consts.MfaDBSelectError, errors.Wrap(err, "->mfa confirm mfa secret error")
		}
		if !totp.Validate(mfaCode, mfaSecret) {
			return consts.MfaCodeFalse, nil
		}
	}
	return consts.Success, nil
}
