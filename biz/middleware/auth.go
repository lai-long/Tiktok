package middleware

import (
	"Tiktok/biz/model/dto"
	"Tiktok/pkg/conf"
	"Tiktok/pkg/consts"
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/golang-jwt/jwt/v5"
)

type myClaims struct {
	UserId   string
	UserName string
	jwt.RegisteredClaims
}

func AuthMiddleware(ctx context.Context, c *app.RequestContext) {
	tokenStringByte := c.GetHeader("Authorization")
	var res dto.Response
	if len(tokenStringByte) == 0 {
		res.Base.Msg = "tokenString is empty"
		res.Base.Code = consts.CodeTokenGetError
		c.JSON(200, res)
		return
	}
	tokenString := string(tokenStringByte)
	bearer := "Bearer"
	for i := 0; i < len(bearer); i++ {
		if tokenString[i] != bearer[i] {
			res.Base.Msg = "tokenString is invalid"
			res.Base.Code = consts.CodeTokenGetError
			c.JSON(200, res)
			return
		}
	}
	tokenString = tokenString[7:]
	claims := myClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return conf.JwtSecret, nil
	})
	if err != nil {
		res.Base.Code = consts.CodeTokenGetError
		res.Base.Msg = "token parse error"
		c.JSON(200, res)
		c.Abort()
		return
	}
	if !token.Valid {
		res.Base.Code = consts.CodeTokenGetError
		res.Base.Msg = "token.valid error"
		c.JSON(200, res)
		c.Abort()
		return
	}
	c.Set("username", claims.UserName)
	c.Set("id", claims.UserId)
}
