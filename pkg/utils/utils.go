// Package utils 放一些公共函数
package utils

import (
	"Tiktok/internal/config"
	user "Tiktok/kitex_gen/user"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/logger"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// CreateID 生成websocket的clientID
func CreateID(uid, toUID string) string {
	return uid + "->" + toUID
}

// GetID 将websocket的clientID拆成两个用户id
func GetID(id string) (string, string) {
	parts := strings.Split(id, "->")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", ""
}

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 检测密码
func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return errors.Wrap(err, "CheckPasswordHash error")
}

// GenerateTokens generate accessToken and refreshToken
func GenerateTokens(userDto *user.UserInfo) (string, string, error) {
	refreshTime := 288 * time.Hour
	accessTime := 24 * time.Hour
	refreshToken, err := GetToken(userDto.Username, userDto.ID, refreshTime, config.GetCfg().Jwt.RefreshSecret)
	if err != nil {
		return "生成refreshToken错误", "", err
	}
	accessToken, err := GetToken(userDto.Username, userDto.ID, accessTime, config.GetCfg().Jwt.AccessSecret)
	if err != nil {
		return "生成accessToken错误", "", err
	}
	return refreshToken, accessToken, nil
}

// GetToken 生成单个token
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
		logger.Error("saveUploadFile mkdir error", zap.String("dir", dir), zap.Error(err))
		return consts.IOOsError, errors.Wrap(err, "saveUploadFile os mkdir错误")
	}
	file, err := os.Create(dir + filename)
	if err != nil {
		return consts.IOOsError, errors.Wrap(err, "saveUploadFile creat failed")
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error("saveUploadFile close error", zap.Error(err))
		}
	}()
	_, err = io.Copy(file, dataFile)
	if err != nil {
		return consts.IOOsError, errors.Wrap(err, "saveUploadFile io copy error")
	}
	return consts.Success, nil
}

// GetUserID 从请求上下文中获取用户ID
func GetUserID(c *app.RequestContext) string {
	if userID, ok := c.Get(consts.UserIDKey); ok {
		if uid, ok := userID.(string); ok {
			return uid
		}
	}
	return ""
}

// GetTraceID 从 context 中获取 trace_id
func GetTraceID(ctx context.Context) string {
	if v, ok := metainfo.GetValue(ctx, consts.TraceIDKey); ok {
		return v
	}
	return ""
}

// TrackTime 返回一个 defer 闭包,记录 op_done 日志,包含 duration_ms 和 ctx 中的 trace_id
func TrackTime(ctx context.Context, opName string) func() {
	start := time.Now()
	traceID := GetTraceID(ctx)
	return func() {
		logger.Info("op_done",
			logger.WithTraceID(traceID),
			zap.String("op", opName),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
	}
}

// TrackRPC 返回一个 defer 闭包,记录 rpc_call 日志,包含 duration_ms、目标 service、method 和 trace_id
func TrackRPC(ctx context.Context, serviceName, method string) func() {
	start := time.Now()
	traceID := GetTraceID(ctx)
	return func() {
		logger.Info("rpc_call",
			logger.WithTraceID(traceID),
			zap.String("service", serviceName),
			zap.String("method", method),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
	}
}

// GetUserIDAndName 从请求上下文中获取用户ID和用户名
func GetUserIDAndName(c *app.RequestContext) (string, string) {
	uid := GetUserID(c)
	if userName, ok := c.Get(consts.UsernameKey); ok {
		if name, ok := userName.(string); ok {
			return uid, name
		}
	}
	return uid, ""
}
