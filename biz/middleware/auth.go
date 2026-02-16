package middleware

import (
	"Tiktok/biz/model/dto"
	"Tiktok/pkg/conf"
	"Tiktok/pkg/consts"
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v5"
)

type myClaims struct {
	UserId   string `json:"user_id"`
	UserName string `json:"username"`
	jwt.RegisteredClaims
}

func AuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		authHeader := string(c.GetHeader("Access-Token"))
		if authHeader == "" {
			c.JSON(200, dto.Response{
				Base: struct {
					Code int    `json:"code"`
					Msg  string `json:"msg"`
				}{
					Code: consts.CodeTokenGetError,
					Msg:  "Get access_token error",
				},
			})
			c.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.JSON(200, dto.Response{
				Base: struct {
					Code int    `json:"code"`
					Msg  string `json:"msg"`
				}{
					Code: consts.CodeTokenGetError,
					Msg:  "token hasprefix error",
				},
			})
			c.Abort()
			return
		}
		tokenString := strings.TrimSpace(authHeader[len(bearerPrefix):])
		if tokenString == "" {
			c.JSON(200, dto.Response{
				Base: struct {
					Code int    `json:"code"`
					Msg  string `json:"msg"`
				}{
					Code: consts.CodeTokenGetError,
					Msg:  "empty token",
				},
			})
			c.Abort()
			return
		}
		claims := &myClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return conf.JwtSecret, nil
		})
		if err != nil {
			c.JSON(200, dto.Response{
				Base: struct {
					Code int    `json:"code"`
					Msg  string `json:"msg"`
				}{
					Code: consts.CodeTokenGetError,
					Msg:  "token parse error",
				},
			})
			c.Abort()
			return
		}
		if !token.Valid {
			c.JSON(200, dto.Response{
				Base: struct {
					Code int    `json:"code"`
					Msg  string `json:"msg"`
				}{
					Code: consts.CodeTokenGetError,
					Msg:  "token valid error",
				},
			})
			c.Abort()
			return
		}

		c.Set("username", claims.UserName)
		c.Set("user_id", claims.UserId)

		c.Next(ctx)
	}
}
