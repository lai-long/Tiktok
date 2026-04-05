package service

import (
	"Tiktok/biz/entity"
	"Tiktok/biz/model/user"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/utils"
	"context"
	"database/sql"
	"io"
	"log"
	"mime/multipart"
	"os"
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

type UserService struct {
	userDb UserDatabase
	mfaDb  MfaDatabase
	redis  UserRedis
}

func NewUserService(userDb UserDatabase, mfaDb MfaDatabase, redis UserRedis) *UserService {
	return &UserService{userDb: userDb, mfaDb: mfaDb, redis: redis}
}

func (s *UserService) IsUsernameExists(username string) (bool, error) {
	_, err := s.userDb.GetUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.Wrap(err, "get user by username")
	}
	return true, nil
}

func (s *UserService) Register(userinfo *user.RegisterReq) (int32, error) {
	var userEntity entity.UserEntity
	var err error
	exists, err := s.IsUsernameExists(userinfo.UserName)
	if err != nil {
		return consts.UserDBSelectError, errors.Wrap(err, "IsUsernameExists error")
	}
	if exists {
		return consts.UserNameExists, nil
	}
	userEntity.Id = utils.IdGenerate()
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

func (s *UserService) Login(userName, password, mfaCode string, ctx context.Context) (int32, error, *user.UserInfo, string, string) {
	userEntity, err := s.userDb.GetUserByUsername(userName)
	if err != nil && err != sql.ErrNoRows {
		return consts.UserDBSelectError, errors.Wrap(err, "->Login GetUserByUsername数据库查询错误"), &user.UserInfo{}, "", ""
	} else if err == sql.ErrNoRows {
		return consts.UserNotExists, errors.Wrap(err, "->login用户不存在"), &user.UserInfo{}, "", ""
	}
	err = utils.CheckPasswordHash(userEntity.Password, password)
	if err != nil {
		return consts.UserPasswordError, errors.Wrap(err, "->login: check password failed"), &user.UserInfo{}, "", ""
	}
	userInfo := userEntity.ToUserInfo()
	err, enable := s.mfaDb.CheckMfaBind(userInfo.ID)
	if err != nil {
		return consts.UserDBSelectError, errors.Wrap(err, "->login: check mfa bind failed"), &user.UserInfo{}, "", ""
	}
	if enable != 0 {
		if mfaCode == "" {
			return consts.MfaLack, nil, &user.UserInfo{}, "", ""
		}
		mfaSecret, err := s.mfaDb.GetMfaSecret(userInfo.ID)
		if err != nil {
			return consts.UserDBSelectError, errors.Wrap(err, "->login get mfa secret failed"), &user.UserInfo{}, "", ""
		}
		if !totp.Validate(mfaCode, mfaSecret) {
			return consts.MfaCodeFalse, nil, &user.UserInfo{}, "", ""
		}
	}
	reToken, acToken, err := utils.GenerateTokens(userInfo)
	if err != nil {
		return consts.GenerateTokenError, errors.Wrap(err, "->login 生成token错误"), userInfo, reToken, acToken
	}
	err = s.redis.UserTokenSet(ctx, reToken, userInfo.ID)
	if err != nil {
		return consts.UserRedisSetError, errors.Wrap(err, "->login 将refresh token存入redis错误"), &user.UserInfo{}, "", ""
	}
	return consts.Success, nil, userInfo, reToken, acToken
}

func (s *UserService) UserInfo(userId string) (*user.UserInfo, int32, error) {
	userEntity, err := s.userDb.GetUserByUserId(userId)
	if err != nil {
		return &user.UserInfo{}, consts.UserDBSelectError, errors.Wrap(err, "->UserInfo GetUserByUserId error")
	}
	userInfo := userEntity.ToUserInfo()
	return userInfo, consts.Success, nil
}

func (s *UserService) UserAvatar(data *multipart.FileHeader, userId interface{}) (int32, error, *user.UserInfo) {
	dataFile, err := data.Open()
	if err != nil {
		return consts.FileError, errors.Wrap(err, "->UserInfo data open 错误"), &user.UserInfo{}
	}
	defer dataFile.Close()
	ok, err := utils.IsImage(dataFile)
	if err != nil {
		log.Printf("IsImage error: %v", err)
		return consts.FileError, errors.Wrap(err, "->userInfo check image failed"), &user.UserInfo{}
	}
	if !ok {
		return consts.ImageFalse, nil, &user.UserInfo{}
	}
	if _, err := dataFile.Seek(0, io.SeekStart); err != nil {
		return consts.FileError, errors.Wrap(err, "->userInfo dataFile error"), &user.UserInfo{}
	}
	filename := utils.IdGenerate()
	err = os.MkdirAll("/home/lai-long/Tiktok/a", os.ModePerm)
	if err != nil {
		return consts.IOOsError, errors.Wrap(err, "->userInfo os mkdir错误"), &user.UserInfo{}
	}
	file, err := os.Create("/home/lai-long/Tiktok/a/" + filename + filepath.Ext(data.Filename))
	if err != nil {
		return consts.IOOsError, errors.Wrap(err, "->userInfo os creat failed"), &user.UserInfo{}
	}
	defer file.Close()
	_, err = io.Copy(file, dataFile)
	if err != nil {
		return consts.IOOsError, errors.Wrap(err, "->userInfo io copy error"), &user.UserInfo{}
	}
	err = s.userDb.UpdateUserAvatar("/home/lai-long/Tiktok/a/"+filename+filepath.Ext(data.Filename), userId)
	if err != nil {
		return consts.UserDBUpdateError, errors.Wrap(err, "->userinfo 更新头像错误"), &user.UserInfo{}
	}
	userEntity, err := s.userDb.GetUserByUserId(userId.(string))
	if err != nil {
		return consts.UserDBSelectError, errors.Wrap(err, "->userinfo get user by userid failed"), &user.UserInfo{}
	}
	userInfo := userEntity.ToUserInfo()
	return consts.Success, nil, userInfo
}

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (int32, string, string, error) {
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
