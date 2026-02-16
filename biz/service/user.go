package service

import (
	"Tiktok/biz/dao/db"
	"Tiktok/biz/model/dto"
	"Tiktok/biz/model/entity"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/utils"
	"log"
)

func Register(userinfo dto.User) (int, string) {
	var userEntity entity.UserEntity
	var err error
	exists, err := utils.IsUsernameExists(userinfo.Username)
	if err != nil {
		log.Fatal(err)
		return consts.CodeCheckUserNameExistError, "register exists error"
	}
	if exists {
		return consts.CodeUserNameExist, "用户名已存在"
	}
	userEntity.Id = utils.IdGenerate()
	userEntity.Username = userinfo.Username
	userEntity.Password, err = utils.HashPassword(userinfo.Password)
	if err != nil {
		log.Println(err)
		return consts.CodeHashError, "hashPassword error"
	}
	if err = db.CreateUser(userEntity); err != nil {
		return consts.CodeDBCreateUserError, "create user error"
	}
	return consts.CodeSuccess, "success"
}

func Login(userDto dto.User) (int, string, dto.User, string, string) {
	userEntity, err := db.GetUserByUsername(userDto.Username)
	if err != nil {
		log.Printf("查询失败: %v", err)
		return consts.CodeGetUserByUserNameError, "用户名错误或用户不存在", dto.User{}, "", ""
	}
	ok := utils.CheckPasswordHash(userEntity.Password, userDto.Password)
	if !ok {
		return consts.CodeCheckPasswordError, "密码错误", dto.User{}, "", ""
	}
	userDto.AvatarURL = userEntity.Avatar_url
	userDto.ID = userEntity.Id
	userDto.Username = userEntity.Username
	userDto.CreatedAt = userEntity.Created_at.String()
	userDto.UpdatedAt = userEntity.Updated_at.String()
	reToken, acToken, ok := utils.GenerateTokens(userDto)
	if ok == false {
		return consts.CodeTokenGenerateError, "生成token错误", userDto, reToken, acToken
	}
	return consts.CodeSuccess, "success", userDto, reToken, acToken
}
