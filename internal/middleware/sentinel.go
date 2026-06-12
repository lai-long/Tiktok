package middleware

import (
	"context"
	"reflect"

	"Tiktok/pkg/consts"
	"Tiktok/pkg/logger"
	"Tiktok/pkg/utils"

	"github.com/alibaba/sentinel-golang/api"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"go.uber.org/zap"
)

func SentinelMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) error {
		ri := rpcinfo.GetRPCInfo(ctx)
		if ri == nil || ri.To() == nil {
			return next(ctx, req, resp)
		}

		entry, blockErr := api.Entry(ri.To().Method())
		if blockErr != nil {
			traceID := utils.GetTraceID(ctx)
			logger.Warn("request blocked by sentinel",
				zap.String("method", ri.To().Method()),
				logger.WithTraceID(traceID),
				zap.String("block_reason", blockErr.Error()),
			)
			if v := reflect.ValueOf(resp); v.Kind() == reflect.Pointer && !v.IsNil() {
				if f := v.Elem().FieldByName("Code"); f.IsValid() && f.CanSet() {
					f.SetInt(int64(consts.SentinelBlock))
				}
			}
			return nil
		}
		defer entry.Exit()

		return next(ctx, req, resp)
	}
}
