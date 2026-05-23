package middleware

import (
	"context"
	"reflect"

	"Tiktok/pkg/consts"

	"github.com/alibaba/sentinel-golang/api"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

func SentinelMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) error {
		ri := rpcinfo.GetRPCInfo(ctx)
		if ri == nil || ri.To() == nil {
			return next(ctx, req, resp)
		}

		entry, blockErr := api.Entry(ri.To().Method())
		if blockErr != nil {
			if v := reflect.ValueOf(resp); v.Kind() == reflect.Ptr && !v.IsNil() {
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
