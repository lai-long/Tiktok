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
		return false, err
	}
	return true, nil
}

func (s *UserService) Register(userinfo *user.RegisterReq) (int, string) {
	var userEntity entity.UserEntity
	var err error
	exists, err := s.IsUsernameExists(userinfo.UserName)
	if err != nil {
		log.Println(err)
		return consts.CodeUserError, "register  IsUsernameExists error"
	}
	if exists {
		return consts.CodeUserError, "用户名已存在"
	}
	userEntity.Id = utils.IdGenerate()
	userEntity.Username = userinfo.UserName
	userEntity.Password, err = utils.HashPassword(userinfo.Password)
	if err != nil {
		log.Println(err)
		return consts.CodeHashError, "hashPassword error"
	}
	if err = s.userDb.CreateUser(userEntity); err != nil {
		log.Println(err)
		return consts.CodeDBCreateError, "db create user error"
	}
	return consts.CodeSuccess, "success"
}

func (s *UserService) Login(userName, password, mfaCode string, ctx context.Context) (int, string, *user.UserInfo, string, string) {
	userEntity, err := s.userDb.GetUserByUsername(userName)
	if err != nil {
		log.Println("get user entity error", err)
		return consts.CodeUserError, "GetUserByUsername Error", &user.UserInfo{}, "", ""
	}
	ok := utils.CheckPasswordHash(userEntity.Password, password)
	if !ok {
		return consts.CodeUserError, "密码错误", &user.UserInfo{}, "", ""
	}
	userInfo := userEntity.ToUserInfo()
	err, enable := s.mfaDb.CheckMfaBind(userInfo.ID)
	if err != nil {
		log.Println(err)
		return consts.CodeUserError, "CheckMfaBind error", &user.UserInfo{}, "", ""
	}
	if enable != 0 {
		if mfaCode == "" {
			return consts.CodeMfaError, "GetMfaCode error 请输入mfa code", &user.UserInfo{}, "", ""
		}
		mfaSecret, err := s.mfaDb.GetMfaSecret(userInfo.ID)
		if err != nil {
			log.Println(err)
			return consts.CodeDBSelectError, "GetMfaSecret from db error", &user.UserInfo{}, "", ""
		}
		if !totp.Validate(mfaCode, mfaSecret) {
			return consts.CodeMfaError, "totp.Validate error", &user.UserInfo{}, "", ""
		}
	}
	reToken, acToken, ok := utils.GenerateTokens(userInfo)
	if ok == false {
		return consts.CodeTokenError, "生成token错误", userInfo, reToken, acToken
	}
	err = s.redis.UserTokenSet(ctx, reToken, userInfo.ID)
	if err != nil {
		log.Println(err)
		return consts.CodeUserError, "db create user refresh token error", &user.UserInfo{}, "", ""
	}
	return consts.CodeSuccess, "success", userInfo, reToken, acToken
}

func (s *UserService) UserInfo(userId string) (*user.UserInfo, int, string, bool) {
	userEntity, err := s.userDb.GetUserByUserId(userId)
	if err != nil {
		log.Printf("GetUserByUserIdError : %v", err)
		return &user.UserInfo{}, consts.CodeDBSelectError, "GetUserByUserIdError", false
	}
	userInfo := userEntity.ToUserInfo()
	return userInfo, consts.CodeSuccess, "Get UserInfo success", true
}

func (s *UserService) UserAvatar(data *multipart.FileHeader, userId interface{}) (int, string, bool, *user.UserInfo) {
	dataFile, err := data.Open()
	if err != nil {
		log.Printf("data.Open error: %v", err)
		return consts.CodeUserError, "data.Open Error", false, &user.UserInfo{}
	}
	defer dataFile.Close()
	ok, err := utils.IsImage(dataFile)
	if err != nil {
		log.Printf("IsImage error: %v", err)
		return consts.CodeUserError, "utils.IsImage Error", false, &user.UserInfo{}
	}
	if !ok {
		return consts.CodeIOError, "IsImage false,文件不是图片", false, &user.UserInfo{}
	}
	if _, err := dataFile.Seek(0, io.SeekStart); err != nil {
		return consts.CodeIOError, "a dataFile.Seek 重置文件指针失败", false, &user.UserInfo{}
	}
	filename := utils.IdGenerate()
	err = os.MkdirAll("/home/lai-long/Tiktok/a", os.ModePerm)
	if err != nil {
		log.Println(err)
		return consts.CodeUserError, "a user avatar MkdirAll Error", false, &user.UserInfo{}
	}
	file, err := os.Create("/home/lai-long/Tiktok/a/" + filename + filepath.Ext(data.Filename))
	if err != nil {
		log.Printf("os.Create error: %v", err)
		return consts.CodeUserError, "user a upload os.Create Error", false, &user.UserInfo{}
	}
	defer file.Close()
	_, err = io.Copy(file, dataFile)
	if err != nil {
		log.Printf("io.Copy error: %v", err)
		return consts.CodeIOError, "a io.copy error", false, &user.UserInfo{}
	}
	err = s.userDb.UpdateUserAvatar("/home/lai-long/Tiktok/a/"+filename+filepath.Ext(data.Filename), userId)
	if err != nil {
		log.Printf("db.UpdateUserAvatar error: %v", err)
		return consts.CodeDBUpdateError, "a db.UpdateUserAvatar error", false, &user.UserInfo{}
	}
	userEntity, err := s.userDb.GetUserByUserId(userId.(string))
	if err != nil {
		log.Printf("db.GetUserByUserIderror: %v", err)
		return consts.CodeDBSelectError, "a db.GetUserByUserId error", false, &user.UserInfo{}
	}
	userInfo := userEntity.ToUserInfo()
	return consts.CodeSuccess, "a change success", true, userInfo
}

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (int, string, string, string, bool) {
	userId, err := s.redis.UserGetByRefreshToken(ctx, refreshToken)
	if err != nil {
		log.Printf("redis.UserGetByRefreshToken error: %v", err)
		return consts.CodeDBSelectError, "token 错误，s.redis.UserGetByRefreshToken err", "", "", false
	}
	userEntity, err := s.userDb.GetUserByUserId(userId)
	if err != nil {
		log.Printf("s.userDb.GetUserByUserIderror: %v", err)
		return consts.CodeDBSelectError, "RefreshToken userDb.GetUserByUserId err ", "", "", false
	}
	userInfo := userEntity.ToUserInfo()
	refreshToken2, accessToken, ok := utils.GenerateTokens(userInfo)
	if !ok {
		return consts.CodeUserError, "RefreshToken utils.GenerateTokens err", "", "", false
	}
	err = s.redis.UserTokenDelete(ctx, refreshToken)
	if err != nil {
		log.Printf("redis.UserTokenDelete error: %v", err)
		return consts.CodeUserError, "RefreshToken s.redis.UserTokenDelete err", "", "", false
	}
	err = s.redis.UserTokenSet(ctx, refreshToken2, userInfo.ID)
	if err != nil {
		log.Printf("redis.UserTokenSet error: %v", err)
		return consts.CodeUserError, "RefreshToken s.redis.UserTokenSet err", "", "", false
	}
	return consts.CodeSuccess, "refreshToken s.redis.UserTokenSet success", refreshToken2, accessToken, true
}
