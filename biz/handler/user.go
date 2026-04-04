package handler

import (
	"Tiktok/biz/model/common"
	"Tiktok/biz/model/user"

	"Tiktok/biz/model/dto"
	"Tiktok/biz/service"
	"Tiktok/pkg/consts"
	"context"
	"log"
	"mime/multipart"

	"github.com/cloudwego/hertz/pkg/app"
)

type UserSever interface {
	Register(registerReq *user.RegisterReq) (int, string)
	Login(loginReq *user.LoginReq, mfaCode string, ctx context.Context) (int, string, *user.UserInfo, string, string)
	UserInfo(userInfoReq *user.UserInfoReq) (*user.UserInfo, int, string, bool)
	UserAvatar(data *multipart.FileHeader, userId interface{}) (int, string, bool, dto.User)
	RefreshToken(ctx context.Context, refreshToken string) (int, string, string, string, bool)
}
type UserHandler struct {
	userService UserSever
	MfaServer   MfaServer
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService,
		MfaServer: userService,
	}
}

func (h *UserHandler) UserRegister(ctx context.Context, c *app.RequestContext) {
	registerReq := new(user.RegisterReq)
	var err error
	if err = c.BindAndValidate(registerReq); err != nil {
		log.Println("UserRegister.BindAndValidate error:", err)
		c.JSON(200, dto.Response{Base: dto.Base{Code: consts.CodeUserError, Msg: "UserRegister BindAndValidate error"}})
		c.Abort()
		return
	}
	code, msg := h.userService.Register(registerReq)
	RegisterResp := &user.RegisterResp{
		Base: &common.Base{
			Code: int32(code),
			Msg:  msg,
		},
	}
	c.JSON(200, RegisterResp)
}

func (h *UserHandler) UserLogin(ctx context.Context, c *app.RequestContext) {
	loginReq := new(user.LoginReq)
	if err := c.BindAndValidate(loginReq); err != nil {
		log.Println("UserLogin.bindAndValidate error:", err)
		loginResp := &user.LoginResp{
			Base: &common.Base{
				Code: consts.CodeUserError,
				Msg:  "UserLogin.bindAndValidate error",
			},
		}
		c.JSON(200, loginResp)
		c.Abort()
		return
	}
	code, msg, userInfo, reToken, acToken := h.userService.Login(loginReq, loginReq.Code, ctx)
	loginResp := &user.LoginResp{
		Base: &common.Base{
			Code: int32(code),
			Msg:  msg,
		},
		Data:         userInfo,
		AccessToken:  acToken,
		RefreshToken: reToken,
	}
	c.JSON(200, loginResp)
}

func (h *UserHandler) UserInfo(ctx context.Context, c *app.RequestContext) {
	userInfoReq := new(user.UserInfoReq)
	err := c.BindAndValidate(userInfoReq)
	if err != nil {
		log.Println("UserInfoReq.bindAndValidate error:", err)
		c.JSON(200, map[string]interface{}{
			"code": consts.CodeUserError,
			"msg":  "UserInfoReq.bindAndValidate error",
		})
	}
	userInfo, code, msg, _ := h.userService.UserInfo(userInfoReq)
	userInfoResp := &user.UserInfoResp{
		Base: &common.Base{
			Code: int32(code),
			Msg:  msg,
		},
		Data: userInfo,
	}
	c.JSON(200, userInfoResp)
}

func (h *UserHandler) UserAvatar(ctx context.Context, c *app.RequestContext) {
	data, _ := c.FormFile("data")
	userId, ok := ctx.Value("user_id").(string)
	if !ok {
		c.JSON(200, dto.Response{Base: dto.Base{Code: consts.CodeUserError, Msg: "用户不存在，c.Get error"}})
	}
	code, msg, ok, user := h.userService.UserAvatar(data, userId)
	if !ok {
		c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}})
		return
	}
	c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}, Data: user})
}

func (h *UserHandler) RefreshToken(ctx context.Context, c *app.RequestContext) {
	refreshToken := c.PostForm("refresh_token")
	code, msg, reToken, acToken, ok := h.userService.RefreshToken(ctx, refreshToken)
	if !ok {
		c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}})
		return
	}
	c.JSON(200, dto.LoginResponse{
		Response: dto.Response{
			Base: dto.Base{Code: code, Msg: msg},
		},
		RefreshToken: reToken,
		AccessToken:  acToken,
	})
}
