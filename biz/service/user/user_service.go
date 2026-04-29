package user

import (
	"Tiktok/biz/entity"
	"Tiktok/biz/model/user"
	"Tiktok/biz/service/mfa"
	"Tiktok/pkg/config"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/utils"
	"context"
	"database/sql"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/pquerna/otp/totp"
)

type UserRedis interface {
	UserTokenSet(ctx context.Context, refreshToken string, userId string) error
	UserGetByRefreshToken(ctx context.Context, refreshToken string) (userId string, err error)
	UserTokenDelete(ctx context.Context, refreshToken string) error
}

type UserDatabase interface {
	CreateUser(user entity.UserEntity) error
	GetUserByUsername(username string) (entity.UserEntity, error)
	GetUserByUserId(userId string) (entity.UserEntity, error)
	UpdateUserAvatar(url string, userId interface{}) error
}

type UserRepo struct {
	userDb UserDatabase
	mfaDb  mfa.MfaDatabase
	redis  UserRedis
}

func NewUserRepo(userDb UserDatabase, mfaDb mfa.MfaDatabase, redis UserRedis) *UserRepo {
	return &UserRepo{userDb: userDb, mfaDb: mfaDb, redis: redis}
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

func (s *UserRepo) Register(userinfo *user.RegisterReq) (int32, error) {
	var userEntity entity.UserEntity
	var err error
	exists, err := s.IsUsernameExists(userinfo.UserName)
	if err != nil {
		return consts.UserDBSelectError, errors.Wrap(err, "IsUsernameExists error")
	}
	if exists {
		return consts.UserNameExists, nil
	}
	userEntity.ID = utils.IDGenerate()
	userEntity.Username = userinfo.UserName
	userEntity.Password, err = utils.HashPassword(userinfo.Password)
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
	enable, err := s.mfaDb.CheckMfaBind(userInfo.ID)
	if err != nil {
		return consts.UserDBSelectError, &user.UserInfo{}, "", "", errors.Wrap(err, "->login: check mfa bind failed")
	}
	if enable != 0 {
		if mfaCode == "" {
			return consts.MfaLack, &user.UserInfo{}, "", "", nil
		}
		mfaSecret, err := s.mfaDb.GetMfaSecret(userInfo.ID)
		if err != nil {
			return consts.UserDBSelectError, &user.UserInfo{}, "", "", errors.Wrap(err, "->login get mfa secret failed")
		}
		if !totp.Validate(mfaCode, mfaSecret) {
			return consts.MfaCodeFalse, &user.UserInfo{}, "", "", nil
		}
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

func (s *UserRepo) UserInfo(userId string) (*user.UserInfo, int32, error) {
	userEntity, err := s.userDb.GetUserByUserId(userId)
	if err != nil {
		return &user.UserInfo{}, consts.UserDBSelectError, errors.Wrap(err, "->UserInfo GetUserByUserId error")
	}
	userInfo := userEntity.ToUserInfo()
	return userInfo, consts.Success, nil
}

func (s *UserRepo) UserAvatar(data *multipart.FileHeader, userId interface{}) (int32, *user.UserInfo, error) {
	dataFile, err := data.Open()
	if err != nil {
		return consts.IOOsError, &user.UserInfo{}, errors.Wrap(err, "->UserInfo data open 错误")
	}
	defer func() {
		err := dataFile.Close()
		if err != nil {
			log.Println(errors.Wrap(err, "-UserInfo data close"))
		}
	}()
	ok, err := utils.IsImage(dataFile)
	if err != nil {
		return consts.FileError, &user.UserInfo{}, errors.Wrap(err, "->userInfo check image failed")
	}
	if !ok {
		return consts.ImageFalse, nil, nil
	}
	if _, err := dataFile.Seek(0, io.SeekStart); err != nil {
		return consts.IOOsError, &user.UserInfo{}, errors.Wrap(err, "->userInfo dataFile error")
	}
	filename := utils.IDGenerate()
	code, err := utils.SaveUploadFile(dataFile, config.Cfg.Path.AvatarPath, filename+filepath.Ext(data.Filename))
	if err != nil {
		return code, &user.UserInfo{}, errors.Wrap(err, "->userAvatar")
	}
	err = s.userDb.UpdateUserAvatar(config.Cfg.Path.AvatarPath+filename, userId)
	if err != nil {
		return consts.UserDBUpdateError, &user.UserInfo{}, errors.Wrap(err, "->userinfo 更新头像错误")
	}
	userEntity, err := s.userDb.GetUserByUserId(userId.(string))
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
