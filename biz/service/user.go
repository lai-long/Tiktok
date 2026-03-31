package service

import (
	"Tiktok/biz/model/dto"
	"Tiktok/biz/model/entity"
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
func (s *UserService) Register(userinfo dto.User) (int, string) {
	var userEntity entity.UserEntity
	var err error
	exists, err := s.IsUsernameExists(userinfo.Username)
	if err != nil {
		log.Println(err)
		return consts.CodeUserError, "register  IsUsernameExists error"
	}
	if exists {
		return consts.CodeUserError, "用户名已存在"
	}
	userEntity.Id = utils.IdGenerate()
	userEntity.Username = userinfo.Username
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

func (s *UserService) Login(userDto dto.User, mfaCode string, ctx context.Context) (int, string, dto.User, string, string) {
	userEntity, err := s.userDb.GetUserByUsername(userDto.Username)
	if err != nil {
		log.Println("get user entity error", err)
		return consts.CodeUserError, "GetUserByUsername Error", dto.User{}, "", ""
	}
	ok := utils.CheckPasswordHash(userEntity.Password, userDto.Password)
	if !ok {
		return consts.CodeUserError, "密码错误", dto.User{}, "", ""
	}
	userDto.AvatarURL = userEntity.Avatar_url
	userDto.ID = userEntity.Id
	userDto.Username = userEntity.Username
	userDto.CreatedAt = userEntity.Created_at.String()
	userDto.UpdatedAt = userEntity.Updated_at.String()
	err, enable := s.mfaDb.CheckMfaBind(userDto.ID)
	if err != nil {
		log.Println(err)
		return consts.CodeUserError, "CheckMfaBind error", dto.User{}, "", ""
	}
	if enable != 0 {
		if mfaCode == "" {
			return consts.CodeMfaError, "GetMfaCode error 请输入mfa code", dto.User{}, "", ""
		}
		mfaSecret, err := s.mfaDb.GetMfaSecret(userDto.ID)
		if err != nil {
			log.Println(err)
			return consts.CodeDBSelectError, "GetMfaSecret from db error", dto.User{}, "", ""
		}
		if !totp.Validate(mfaCode, mfaSecret) {
			return consts.CodeMfaError, "totp.Validate error", dto.User{}, "", ""
		}
	}
	reToken, acToken, ok := utils.GenerateTokens(userDto)
	if ok == false {
		return consts.CodeTokenError, "生成token错误", userDto, reToken, acToken
	}
	err = s.redis.UserTokenSet(ctx, reToken, userDto.ID)
	if err != nil {
		log.Println(err)
		return consts.CodeUserError, "db create user refresh token error", dto.User{}, "", ""
	}
	return consts.CodeSuccess, "success", userDto, reToken, acToken
}

func (s *UserService) UserInfo(userId string) (dto.User, int, string, bool) {
	userEntity, err := s.userDb.GetUserByUserId(userId)
	if err != nil {
		log.Printf("GetUserByUserIdError : %v", err)
		return dto.User{}, consts.CodeDBSelectError, "GetUserByUserIdError", false
	}
	var user dto.User
	user.Username = userEntity.Username
	user.AvatarURL = userEntity.Avatar_url
	user.ID = userEntity.Id
	user.CreatedAt = userEntity.Created_at.String()
	user.UpdatedAt = userEntity.Updated_at.String()
	return user, consts.CodeSuccess, "Get UserInfo success", true
}

func (s *UserService) UserAvatar(data *multipart.FileHeader, userId interface{}) (int, string, bool, dto.User) {
	dataFile, err := data.Open()
	if err != nil {
		log.Printf("data.Open error: %v", err)
		return consts.CodeUserError, "data.Open Error", false, dto.User{}
	}
	defer dataFile.Close()
	ok, err := utils.IsImage(dataFile)
	if err != nil {
		log.Printf("IsImage error: %v", err)
		return consts.CodeUserError, "utils.IsImage Error", false, dto.User{}
	}
	if !ok {
		return consts.CodeIOError, "IsImage false,文件不是图片", false, dto.User{}
	}
	if _, err := dataFile.Seek(0, io.SeekStart); err != nil {
		return consts.CodeIOError, "a dataFile.Seek 重置文件指针失败", false, dto.User{}
	}
	filename := utils.IdGenerate()
	err = os.MkdirAll("/home/lai-long/Tiktok/a", os.ModePerm)
	if err != nil {
		log.Println(err)
		return consts.CodeUserError, "a user avatar MkdirAll Error", false, dto.User{}
	}
	file, err := os.Create("/home/lai-long/Tiktok/a/" + filename + filepath.Ext(data.Filename))
	if err != nil {
		log.Printf("os.Create error: %v", err)
		return consts.CodeUserError, "user a upload os.Create Error", false, dto.User{}
	}
	defer file.Close()
	_, err = io.Copy(file, dataFile)
	if err != nil {
		log.Printf("io.Copy error: %v", err)
		return consts.CodeIOError, "a io.copy error", false, dto.User{}
	}
	err = s.userDb.UpdateUserAvatar("/home/lai-long/Tiktok/a/"+filename+filepath.Ext(data.Filename), userId)
	if err != nil {
		log.Printf("db.UpdateUserAvatar error: %v", err)
		return consts.CodeDBUpdateError, "a db.UpdateUserAvatar error", false, dto.User{}
	}
	userEntity, err := s.userDb.GetUserByUserId(userId.(string))
	if err != nil {
		log.Printf("db.GetUserByUserIderror: %v", err)
		return consts.CodeDBSelectError, "a db.GetUserByUserId error", false, dto.User{}
	}
	var user dto.User
	user.Username = userEntity.Username
	user.AvatarURL = userEntity.Avatar_url
	user.ID = userEntity.Id
	return consts.CodeSuccess, "a change success", true, user
}
func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (int, string, string, string, bool) {
	userId, err := s.redis.UserGetByRefreshToken(ctx, refreshToken)
	if err != nil {
		log.Printf("redis.UserGetByRefreshToken error: %v", err)
		return consts.CodeDBSelectError, "token 错误，s.redis.UserGetByRefreshToken err", "", "", false
	}
	user, err := s.userDb.GetUserByUserId(userId)
	if err != nil {
		log.Printf("s.userDb.GetUserByUserIderror: %v", err)
		return consts.CodeDBSelectError, "RefreshToken userDb.GetUserByUserId err ", "", "", false
	}
	var userDto dto.User
	userDto.Username = user.Username
	userDto.ID = user.Id
	refreshToken2, accessToken, ok := utils.GenerateTokens(userDto)
	if !ok {
		return consts.CodeUserError, "RefreshToken utils.GenerateTokens err", "", "", false
	}
	err = s.redis.UserTokenDelete(ctx, refreshToken)
	if err != nil {
		log.Printf("redis.UserTokenDelete error: %v", err)
		return consts.CodeUserError, "RefreshToken s.redis.UserTokenDelete err", "", "", false
	}
	err = s.redis.UserTokenSet(ctx, refreshToken2, userDto.ID)
	if err != nil {
		log.Printf("redis.UserTokenSet error: %v", err)
		return consts.CodeUserError, "RefreshToken s.redis.UserTokenSet err", "", "", false
	}
	return consts.CodeSuccess, "refreshToken s.redis.UserTokenSet success", refreshToken2, accessToken, true
}
