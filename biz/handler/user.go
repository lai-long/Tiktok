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
		baseResponse := dto.Base{Code: consts.CodeError, Msg: "UserRegister BindAndValidate error"}
		c.JSON(200, dto.Response{baseResponse, nil})
		c.Abort()
		return
	}
	code, msg := service.Register(userinfo)
	res := dto.Response{
		Base: dto.Base{Code: code, Msg: msg},
		Data: nil,
	}
	c.JSON(200, res)
}

func UserLogin(ctx context.Context, c *app.RequestContext) {
	var userDto dto.User
	if err := c.BindAndValidate(&userDto); err != nil {
		baseResponse := dto.Base{Code: -1, Msg: "UserLogin BindAndValidate error"}
		c.JSON(200, dto.Response{Base: baseResponse})
		c.Abort()
		return
	}
	code, msg, user, reToken, acToken := service.Login(userDto)
	res := dto.LoginResponse{
		Response: dto.Response{
			Base: dto.Base{
				Code: code,
				Msg:  msg,
			},
			Data: user,
		},
		RefreshToken: reToken,
		AccessToken:  acToken,
	}
	c.JSON(200, res)
}

func UserInfo(ctx context.Context, c *app.RequestContext) {
	//userId := c.Param("user_id")
	//userEntity, err := db.GetUserByUserId(userId)
	//if err != nil {
	//	baseResponse := dto.Base{Code: -1, Msg: "GetUserByUserId error"}
	//	c.JSON(200, baseResponse)
	//	c.Abort()
	//	return
	//}
	//c.JSON(200, dto.Response{Base: dto.Base{Code: 10000, Msg: "success"}, Data: userEntity})
}
