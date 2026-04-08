package middleware

import (
	"Tiktok/biz/model/common"
	"Tiktok/biz/model/user"
	"Tiktok/pkg/config"
	"Tiktok/pkg/consts"
	"context"
	"log"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v5"
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
		return []byte(config.Cfg.Jwt.AccessSecret), nil
	})
	if err != nil {
		log.Printf("JWT parse error: %v", err)
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
	newCtx := context.WithValue(ctx, "user_id", userid)
	newCtx = context.WithValue(newCtx, "username", username)
	c.Next(newCtx)
}
