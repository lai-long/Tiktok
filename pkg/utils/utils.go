// Package utils 放一些公共函数
package utils

import (
	"Tiktok/biz/model/user"
	"Tiktok/pkg/config"
	"Tiktok/pkg/consts"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

// CreateID 生成websocket的clientID
func CreateID(uid, toUID string) string {
	return uid + "->" + toUID
}

// GetID 将websocket的clientID拆成两个用户id
func GetID(id string) (string, string) {
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

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 检测密码
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
	return errors.Wrap(err, "CheckPasswordHash error")
}

// GenerateTokens generate accessToken and refreshToken
func GenerateTokens(userDto *user.UserInfo) (string, string, error) {
	refreshTime := 288 * time.Hour
	accessTime := 24 * time.Hour
	refreshToken, err := GetToken(userDto.Username, userDto.ID, refreshTime, config.Cfg.Jwt.RefreshSecret)
	if err != nil {
		return "生成refreshToken错误", "", err
	}
	accessToken, err := GetToken(userDto.Username, userDto.ID, accessTime, config.Cfg.Jwt.AccessSecret)
	if err != nil {
		return "生成accessToken错误", "", err
	}
	return refreshToken, accessToken, nil
}

// GetToken 生成单个ai
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

// IDGenerate 生成id
func IDGenerate() string {
	return xid.New().String()
}

// IsImage 检测文件是否是图片
func IsImage(file multipart.File) (bool, error) {
	head := make([]byte, 512)
	_, err := file.Read(head)
	if err != nil {
		return false, errors.Wrap(err, "->isImage read file header error")
	}
	mime := http.DetectContentType(head)
	switch mime {
	case "image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp", "image/bmp":
		return true, nil
	default:
		return false, nil
	}
}

// ai关键词
var triggerKeywords = []string{
	"@AI",
	"111",
}

// CheckAiKeyWord 检测用户是否要用ai
func CheckAiKeyWord(message string) (bool, string) {
	for _, keyword := range triggerKeywords {
		if strings.Contains(message, keyword) {
			question := strings.TrimSpace(strings.Replace(message, keyword, "", 1))
			return true, question
		}
	}
	return false, ""
}

// SaveUploadFile 保存上传的文件（视频、头像）
func SaveUploadFile(dataFile multipart.File, dir string, filename string) (int32, error) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Println(dir)
		return consts.IOOsError, errors.Wrap(err, "saveUploadFile os mkdir错误")
	}
	file, err := os.Create(dir + filename)
	if err != nil {
		return consts.IOOsError, errors.Wrap(err, "saveUploadFile creat failed")
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println("saveUploadFile close err", err)
		}
	}()
	_, err = io.Copy(file, dataFile)
	if err != nil {
		return consts.IOOsError, errors.Wrap(err, "saveUploadFile io copy error")
	}
	return consts.Success, nil
}
