package service

import (
	"Tiktok/biz/dao/db"
	"Tiktok/biz/model/dto"
	"Tiktok/biz/model/entity"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/utils"
	"io"
	"log"
	"mime/multipart"
	"os"
)

func Register(userinfo dto.User) (int, string) {
	var userEntity entity.UserEntity
	var err error
	exists, err := utils.IsUsernameExists(userinfo.Username)
	if err != nil {
		log.Fatal(err)
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
	if err = db.CreateUser(userEntity); err != nil {
		return consts.CodeDBCreateError, "create user error"
	}
	return consts.CodeSuccess, "success"
}

func Login(userDto dto.User) (int, string, dto.User, string, string) {
	userEntity, err := db.GetUserByUsername(userDto.Username)
	if err != nil {
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
	reToken, acToken, ok := utils.GenerateTokens(userDto)
	if ok == false {
		return consts.CodeTokenError, "生成token错误", userDto, reToken, acToken
	}
	return consts.CodeSuccess, "success", userDto, reToken, acToken
}

func UserInfo(userId string) (dto.User, int, string, bool) {
	userEntity, err := db.GetUserByUserId(userId)
	if err != nil {
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

func UserAvatar(data *multipart.FileHeader, userId interface{}) (int, string, bool, dto.User) {
	dataFile, err := data.Open()
	if err != nil {
		return consts.CodeUserError, "data.Open Error", false, dto.User{}
	}
	defer dataFile.Close()
	ok, err := utils.IsImage(dataFile)
	if err != nil {
		return consts.CodeUserError, "utils.IsImage Error", false, dto.User{}
	}
	if !ok {
		return consts.CodeIOError, "IsImage false,文件不是图片", false, dto.User{}
	}
	if _, err := dataFile.Seek(0, io.SeekStart); err != nil {
		return consts.CodeIOError, "avatar dataFile.Seek 重置文件指针失败", false, dto.User{}
	}
	file, err := os.Create("/home/lai-long/Tiktok/avatar/" + data.Filename)
	if err != nil {
		return consts.CodeUserError, "user avatar upload os.Create Error", false, dto.User{}
	}
	defer file.Close()
	_, err = io.Copy(file, dataFile)
	if err != nil {
		return consts.CodeIOError, "avatar io.copy error", false, dto.User{}
	}
	err = db.UpdateUserAvatar("/home/lai-long/Tiktok/avatar/"+data.Filename, userId)
	if err != nil {
		return consts.CodeDBUpdateError, "avatar db.UpdateUserAvatar error", false, dto.User{}
	}
	userEntity, err := db.GetUserByUserId(userId.(string))
	if err != nil {
		return consts.CodeDBSelectError, "avatar db.GetUserByUserId error", false, dto.User{}
	}
	var user dto.User
	user.Username = userEntity.Username
	user.AvatarURL = userEntity.Avatar_url
	user.ID = userEntity.Id
	return consts.CodeSuccess, "avatar change success", true, user
}
