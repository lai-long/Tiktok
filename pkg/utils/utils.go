package utils

import (
	"Tiktok/biz/model/user"
	"Tiktok/pkg/config"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

// websocket
func CreateId(uid, toUid string) string {
	return uid + "->" + toUid
}
func GetId(id string) (string, string) {
	log.Println(id)
	parts := strings.Split(id, "->")
	log.Println("begin")
	if len(parts) == 2 {
		log.Println("part[0]", parts[0], "part[1]", parts[1])
		return parts[0], parts[1]
	}
	log.Println("false")
	return "", ""
}

// password hash
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
	return err == nil
}

// jwt generate token
func GenerateTokens(userDto *user.UserInfo) (string, string, bool) {
	refreshTime := 288 * time.Hour
	accessTime := 24 * time.Hour
	refreshToken, err := GetToken(userDto.Username, userDto.ID, refreshTime, config.Cfg.Jwt.RefreshSecret)
	if err != nil {
		log.Println(err)
		return "生成refreshToken错误", "", false
	}
	accessToken, err := GetToken(userDto.Username, userDto.ID, accessTime, config.Cfg.Jwt.AccessSecret)
	if err != nil {
		log.Println(err)
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

// IdGenerate
func IdGenerate() string {
	return xid.New().String()
}

// check image
func IsImage(file multipart.File) (bool, error) {
	head := make([]byte, 512)
	_, err := file.Read(head)
	if err != nil {
		return false, err
	}
	mime := http.DetectContentType(head)
	switch mime {
	case "image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp", "image/bmp":
		return true, nil
	default:
		return false, nil
	}
}
