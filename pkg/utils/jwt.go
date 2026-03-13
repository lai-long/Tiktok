package utils

import (
	"Tiktok/biz/model/dto"
	"Tiktok/pkg/conf"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokens(userDto dto.User) (string, string, bool) {
	refreshTime := 288 * time.Hour
	accessTime := 24 * time.Hour
	refreshToken, err := GetToken(userDto.Username, userDto.ID, refreshTime, conf.Cfg.Jwt.RefreshSecret)
	if err != nil {
		return "生成refreshToken错误", "", false
	}
	accessToken, err := GetToken(userDto.Username, userDto.ID, accessTime, conf.Cfg.Jwt.AccessSecret)
	if err != nil {
		return "生成accessToken错误", "", false
	}
	return refreshToken, accessToken, true
}
func GetToken(username string, userid string, t time.Duration, secret string) (string, error) {
	jwtClaims := &jwt.MapClaims{
		"username": username,
		"userid":   userid,
		"exp":      time.Now().Add(t).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "生成token错误", err
	}
	return tokenString, nil
}
