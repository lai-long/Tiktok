package middleware

import (
	"context"
	"runtime/debug"
	"time"

	"Tiktok/pkg/consts"
	"Tiktok/pkg/logger"
	"Tiktok/pkg/utils"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"go.uber.org/zap"
)

func TracingMiddleware(serviceName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			method := "default"
			if ri := rpcinfo.GetRPCInfo(ctx); ri != nil && ri.To() != nil {
				method = ri.To().Method()
			}

			traceID := utils.GetTraceID(ctx)
			if traceID == "" {
				traceID = utils.IDGenerate()
				ctx = metainfo.WithValue(ctx, consts.TraceIDKey, traceID)
			}

			defer func() {
				if r := recover(); r != nil {
					logger.Error("svc_panic",
						logger.WithServiceName(serviceName),
						logger.WithTraceID(traceID),
						zap.String("method", method),
						zap.Any("panic", r),
						zap.String("stack", string(debug.Stack())),
					)
					panic(r)
				}
			}()

			start := time.Now()
			logger.Info("svc_start",
				logger.WithServiceName(serviceName),
				logger.WithTraceID(traceID),
				zap.String("method", method),
			)

			err = next(ctx, req, resp)

			duration := time.Since(start).Milliseconds()
			status := "success"
			if err != nil {
				status = "error"
			}
			fields := []zap.Field{
				logger.WithServiceName(serviceName),
				logger.WithTraceID(traceID),
				zap.String("method", method),
				zap.String("status", status),
				zap.Int64("duration_ms", duration),
			}
			if err != nil {
				fields = append(fields, zap.Error(err))
				logger.Error("svc_done", fields...)
			} else {
				logger.Info("svc_done", fields...)
			}

			return err
		}
	}
}
