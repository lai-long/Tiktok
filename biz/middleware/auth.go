package middleware

import (
	"Tiktok/biz/model/common"
	"Tiktok/biz/model/user"
	"Tiktok/internal/config"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/logger"
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func AuthMiddleware(ctx context.Context, c *app.RequestContext) {
	path := string(c.Request.URI().Path())
	switch path {
	case "/user/login",
		"/user/register",
		"/user/refresh":
		c.Next(ctx)
		return
	}
	req := new(user.AuthReq)
	err := c.BindAndValidate(req)
	if err != nil {
		logger.Error("AuthMiddleware BindAndValidate error", zap.Error(err))
		return

	}
	if req.AccessToken == "" {
		c.JSON(200, user.AuthResp{Base: &common.Base{
			Code: consts.UserReqValidError,
			Msg:  "AccessToken 为空",
		}})
		c.Abort()
		return
	}
	tokenString := strings.TrimSpace(req.AccessToken)
	if tokenString == "" {
		c.JSON(200, user.AuthResp{Base: &common.Base{
			Code: consts.UserPasswordError,
			Msg:  "tokenString TrimSpace error",
		}})
		c.Abort()
		return
	}
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.GetCfg().Jwt.AccessSecret), nil
	})
	if err != nil {
		logger.Error("JWT parse error", zap.Error(err))
		c.JSON(200, user.AuthResp{Base: &common.Base{
			Code: consts.UserPasswordError,
			Msg:  "JWT parse error",
		}})
		c.Abort()
		return
	}
	if !token.Valid {
		c.JSON(200, user.AuthResp{Base: &common.Base{
			Code: consts.UserPasswordError,
			Msg:  "JWT Valid error",
		}})
		c.Abort()
		return
	}
	userid, _ := (*claims)["userid"].(string)
	username, _ := (*claims)["username"].(string)
	c.Set(consts.UserIDKey, userid)
	c.Set(consts.UsernameKey, username)
	c.Next(ctx)
}
