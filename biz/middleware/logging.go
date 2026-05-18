package middleware

import (
	"context"
	"time"

	"Tiktok/pkg/consts"
	"Tiktok/pkg/logger"

	"github.com/cloudwego/hertz/pkg/app"
	"go.uber.org/zap"
)

func LoggingMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		path := string(c.Request.URI().Path())
		method := string(c.Request.Method())
		clientIP := c.ClientIP()

		c.Next(ctx)

		latency := time.Since(start).Milliseconds()
		status := c.Response.StatusCode()

		userIDVal, _ := c.Get(consts.UserIDKey)
		userID := ""
		if userIDVal != nil {
			userID, _ = userIDVal.(string)
		}

		fields := []zap.Field{
			logger.WithServiceName("api"),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("client_ip", clientIP),
			zap.Int("status", status),
			zap.Int64("latency_ms", latency),
		}

		if userID != "" {
			fields = append(fields, logger.WithUserID(userID))
		}

		switch {
		case status >= 500:
			logger.Error("request completed with server error", fields...)
		case status >= 400:
			logger.Warn("request completed with client error", fields...)
		default:
			logger.Info("request completed", fields...)
		}
	}
}
