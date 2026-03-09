package middleware

import (
	"Tiktok/biz/model/dto"
	"Tiktok/pkg/conf"
	"Tiktok/pkg/consts"
	"context"
	"log"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(ctx context.Context, c *app.RequestContext) {
	authHeader := c.Request.Header.Get("Access-Token")
	log.Printf("%v", authHeader)
	if authHeader == "" {
		c.JSON(200, dto.Response{
			Base: dto.Base{
				Code: consts.CodeTokenError,
				Msg:  "get tokenHeader failed",
			},
		})
		c.Abort()
		return
	}
	tokenString := strings.TrimSpace(authHeader)
	if tokenString == "" {
		c.JSON(200, dto.Response{
			Base: dto.Base{
				Code: consts.CodeTokenError,
				Msg:  "get tokenString failed",
			},
		})
		c.Abort()
		return
	}
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(conf.Cfg.Jwt.Secret), nil
	})
	if err != nil {
		log.Printf("JWT parse error: %v", err)
		c.JSON(200, dto.Response{
			Base: dto.Base{
				Code: consts.CodeTokenError,
				Msg:  "token ParseWithClaims failed",
			},
		})
		c.Abort()
		return
	}
	if !token.Valid {
		c.JSON(200, dto.Response{
			Base: dto.Base{
				Code: consts.CodeTokenError,
				Msg:  "token invalid",
			},
		})
		c.Abort()
		return
	}
	username, _ := (*claims)["username"].(string)
	userid, _ := (*claims)["userid"].(string)
	c.Set("username", username)
	c.Set("user_id", userid)
	c.Next(ctx)
}
