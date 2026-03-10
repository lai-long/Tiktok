package handler

import (
	"Tiktok/biz/model/dto"
	"Tiktok/pkg/consts"
	"context"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
)

func (h *Handler) UserRegister(ctx context.Context, c *app.RequestContext) {
	var userinfo dto.User
	var err error
	if err = c.BindAndValidate(&userinfo); err != nil {
		log.Println("UserRegister.BindAndValidate error:", err)
		c.JSON(200, dto.Response{Base: dto.Base{Code: consts.CodeUserError, Msg: "UserRegister BindAndValidate error"}})
		c.Abort()
		return
	}
	code, msg := h.service.Register(userinfo)
	c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}})
}

func (h *Handler) UserLogin(ctx context.Context, c *app.RequestContext) {
	var userDto dto.User
	if err := c.BindAndValidate(&userDto); err != nil {
		log.Println("UserLogin.bindAndValidate error:", err)
		c.JSON(200, dto.Response{Base: dto.Base{Code: consts.CodeUserError, Msg: "UserLogin BindAndValidate error"}})
		c.Abort()
		return
	}
	mfcCode := c.PostForm("code")
	code, msg, user, reToken, acToken := h.service.Login(userDto, mfcCode)
	res := dto.LoginResponse{
		Response: dto.Response{
			Base: dto.Base{
				Code: code,
				Msg:  msg,
			},
			Data: user,
		},
		AccessToken:  acToken,
		RefreshToken: reToken,
	}

	c.JSON(200, res)
}

func (h *Handler) UserInfo(ctx context.Context, c *app.RequestContext) {
	userId := c.Query("user_id")
	user, code, msg, ok := h.service.UserInfo(userId)
	if !ok {
		c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}})
		c.Abort()
		return
	}
	c.JSON(200, dto.Response{Base: dto.Base{Code: 10000, Msg: "success"}, Data: user})
}

func (h *Handler) UserAvatar(ctx context.Context, c *app.RequestContext) {
	data, _ := c.FormFile("data")
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(200, dto.Response{Base: dto.Base{Code: consts.CodeUserError, Msg: "用户不存在，c.Get error"}})
	}
	code, msg, ok, user := h.service.UserAvatar(data, userId)
	if !ok {
		c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}})
		return
	}
	c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}, Data: user})
}
