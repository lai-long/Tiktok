package handler

import (
	"Tiktok/biz/model/dto"
	"Tiktok/pkg/consts"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func (h *Handler) MfaQrcode(ctx context.Context, c *app.RequestContext) {
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
	ok, key, secret, code, msg := h.service.GenerateMfa(userName.(string), userId.(string))
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
func (h *Handler) MfaBind(ctx context.Context, c *app.RequestContext) {
	mfaCode := c.PostForm("code")
	secret := c.PostForm("secret")
	if secret != "" {
		h.service.MfaBindBySecret(mfaCode, secret)
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
	code, msg := h.service.MfaBindByCode(mfaCode, userId.(string))
	c.JSON(200, dto.Base{
		Code: code,
		Msg:  msg,
	})
}
