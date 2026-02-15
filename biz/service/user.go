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
		return consts.CodeError, "register exists error"
	}
	if exists {
		return consts.CodeError, "用户名已存在"
	}
	userEntity.Id = utils.IdGenerate()
	userEntity.Username = userinfo.Username
	userEntity.Password, err = utils.HashPassword(userinfo.Password)
	if err != nil {
		log.Println(err)
		return consts.CodeError, "hashPassword error"
	}
	if err = db.CreateUser(userEntity); err != nil {
		return consts.CodeError, "create user error"
	}
	return consts.CodeSuccess, "success"
}
