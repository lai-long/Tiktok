package handler

import (
	"Tiktok/biz/model/dto"
	"Tiktok/pkg/consts"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type MfaServer interface {
	GenerateMfa(username string, userId string) (bool, string, string, int, string)
	MfaBindByCode(code string, userId string) (int, string)
	MfaBindBySecret(secret string, userId string) (int, string)
}

func (h *UserHandler) MfaQrcode(ctx context.Context, c *app.RequestContext) {
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(200, dto.Base{
			Code: consts.CodeError,
			Msg:  "GET USER ID not found",
		})
		return
	}
	userName, exist := c.Get("username")
	if !exist {
		c.JSON(200, dto.Base{
			Code: consts.CodeError,
			Msg:  "GET USER USERNAME not found",
		})
		return
	}
	ok, key, secret, code, msg := h.MfaServer.GenerateMfa(userName.(string), userId.(string))
	if !ok {
		c.JSON(200, dto.Base{
			Code: code,
			Msg:  msg,
		})
		return
	}
	c.JSON(200, dto.Response{
		Base: dto.Base{
			Code: code,
			Msg:  msg,
		},
		Data: map[string]string{
			"secret": secret,
			"qrcode": key,
		},
	})
}
func (h *UserHandler) MfaBind(ctx context.Context, c *app.RequestContext) {
	mfaCode := c.PostForm("code")
	secret := c.PostForm("secret")
	if secret != "" {
		h.MfaServer.MfaBindBySecret(mfaCode, secret)
		c.JSON(200, dto.Base{
			Code: consts.CodeSuccess,
			Msg:  "mfa bind success",
		})
		return
	}
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(200, dto.Base{
			Code: consts.CodeError,
			Msg:  "GET USER ID not found",
		})
	}
	code, msg := h.MfaServer.MfaBindByCode(mfaCode, userId.(string))
	c.JSON(200, dto.Base{
		Code: code,
		Msg:  msg,
	})
}
