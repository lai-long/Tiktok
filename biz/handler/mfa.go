package handler

import (
	"Tiktok/biz/model/common"
	"Tiktok/biz/model/mfa"
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
	userId, ok := ctx.Value("user_id").(string)
	if !ok {
		c.JSON(200, mfa.MfaQrcodeResp{Base: &common.Base{
			Code: consts.CodeError,
			Msg:  "mfa qrcode get user name err",
		}})
		return
	}
	userName, ok := ctx.Value("username").(string)
	if !ok {
		c.JSON(200, mfa.MfaQrcodeResp{
			Base: &common.Base{
				Code: consts.CodeError,
				Msg:  "GET USER USERNAME not found",
			},
		})
		return
	}
	_, key, secret, code, msg := h.MfaServer.GenerateMfa(userName, userId)
	c.JSON(200, mfa.MfaQrcodeResp{
		Base: &common.Base{
			Code: int32(code),
			Msg:  msg,
		},
		Data: &mfa.MfaData{
			Secret: secret,
			Qrcode: key,
		},
	})
}
func (h *UserHandler) MfaBind(ctx context.Context, c *app.RequestContext) {
	req := new(mfa.MfaBindReq)
	err := c.BindAndValidate(req)
	if err != nil {
		c.JSON(200, mfa.MfaBindResp{Base: &common.Base{
			Code: consts.CodeError,
			Msg:  "mfa BindAndValidate err",
		}})
		return
	}
	if req.Secret != "" {
		h.MfaServer.MfaBindBySecret(req.Code, req.Secret)
		c.JSON(200, mfa.MfaBindResp{Base: &common.Base{
			Code: consts.CodeSuccess,
			Msg:  "mfa BindBySecret success",
		}})
		return
	}
	userId, ok := ctx.Value("user_id").(string)
	if !ok {
		c.JSON(200, mfa.MfaBindResp{Base: &common.Base{
			Code: consts.CodeError,
			Msg:  "mfa BindBySecret get user id err",
		}})
	}
	code, msg := h.MfaServer.MfaBindByCode(req.Code, userId)
	c.JSON(200, mfa.MfaBindResp{
		Base: &common.Base{
			Code: int32(code),
			Msg:  msg,
		},
	})
}
