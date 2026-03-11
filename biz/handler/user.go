package handler

import (
	"Tiktok/biz/model/dto"
	"Tiktok/biz/service"
	"Tiktok/pkg/consts"
	"context"
	"log"
	"mime/multipart"

	"github.com/cloudwego/hertz/pkg/app"
)

type UserSever interface {
	Register(userinfo dto.User) (int, string)
	Login(userDto dto.User, mfaCode string) (int, string, dto.User, string, string)
	UserInfo(userId string) (dto.User, int, string, bool)
	UserAvatar(data *multipart.FileHeader, userId interface{}) (int, string, bool, dto.User)
}
type UserHandler struct {
	userService UserSever
	MfaServer   MfaServer
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService,
		MfaServer: userService}
}

func (h *UserHandler) UserRegister(ctx context.Context, c *app.RequestContext) {
	var userinfo dto.User
	var err error
	if err = c.BindAndValidate(&userinfo); err != nil {
		log.Println("UserRegister.BindAndValidate error:", err)
		c.JSON(200, dto.Response{Base: dto.Base{Code: consts.CodeUserError, Msg: "UserRegister BindAndValidate error"}})
		c.Abort()
		return
	}
	code, msg := h.userService.Register(userinfo)
	c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}})
}

func (h *UserHandler) UserLogin(ctx context.Context, c *app.RequestContext) {
	var userDto dto.User
	if err := c.BindAndValidate(&userDto); err != nil {
		log.Println("UserLogin.bindAndValidate error:", err)
		c.JSON(200, dto.Response{Base: dto.Base{Code: consts.CodeUserError, Msg: "UserLogin BindAndValidate error"}})
		c.Abort()
		return
	}
	mfcCode := c.PostForm("code")
	code, msg, user, reToken, acToken := h.userService.Login(userDto, mfcCode)
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

func (h *UserHandler) UserInfo(ctx context.Context, c *app.RequestContext) {
	userId := c.Query("user_id")
	user, code, msg, ok := h.userService.UserInfo(userId)
	if !ok {
		c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}})
		c.Abort()
		return
	}
	c.JSON(200, dto.Response{Base: dto.Base{Code: 10000, Msg: "success"}, Data: user})
}

func (h *UserHandler) UserAvatar(ctx context.Context, c *app.RequestContext) {
	data, _ := c.FormFile("data")
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(200, dto.Response{Base: dto.Base{Code: consts.CodeUserError, Msg: "用户不存在，c.Get error"}})
	}
	code, msg, ok, user := h.userService.UserAvatar(data, userId)
	if !ok {
		c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}})
		return
	}
	c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}, Data: user})
}
