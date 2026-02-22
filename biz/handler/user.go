package handler

import (
	"Tiktok/biz/model/dto"
	"Tiktok/biz/service"
	"Tiktok/pkg/consts"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func UserRegister(ctx context.Context, c *app.RequestContext) {
	var userinfo dto.User
	var err error
	if err = c.BindAndValidate(&userinfo); err != nil {
		c.JSON(200, dto.Response{Base: dto.Base{Code: consts.CodeUserError, Msg: "UserRegister BindAndValidate error"}})
		c.Abort()
		return
	}
	code, msg := service.Register(userinfo)
	c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}})
}

func UserLogin(ctx context.Context, c *app.RequestContext) {
	var userDto dto.User
	if err := c.BindAndValidate(&userDto); err != nil {
		c.JSON(200, dto.Response{Base: dto.Base{Code: consts.CodeUserError, Msg: "UserLogin BindAndValidate error"}})
		c.Abort()
		return
	}
	code, msg, user, reToken, acToken := service.Login(userDto)
	res := dto.Response{
		Base: dto.Base{
			Code: code,
			Msg:  msg,
		},
		Data:         user,
		AccessToken:  acToken,
		RefreshToken: reToken,
	}

	c.JSON(200, res)
}

func UserInfo(ctx context.Context, c *app.RequestContext) {
	userId := c.Query("user_id")
	user, code, msg, ok := service.UserInfo(userId)
	if !ok {
		c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}})
		c.Abort()
		return
	}
	c.JSON(200, dto.Response{Base: dto.Base{Code: 10000, Msg: "success"}, Data: user})
}

func UserAvatar(ctx context.Context, c *app.RequestContext) {
	data, _ := c.FormFile("data")
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(200, dto.Response{Base: dto.Base{Code: consts.CodeUserError, Msg: "用户不存在，c.Get error"}})
	}
	code, msg, ok, user := service.UserAvatar(data, userId)
	if !ok {
		c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}})
		return
	}
	c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}, Data: user})
}
