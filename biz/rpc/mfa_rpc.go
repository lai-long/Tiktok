package rpc

import (
	"Tiktok/internal/config"
	"Tiktok/kitex_gen/mfa"
	"Tiktok/kitex_gen/mfa/mfaservice"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/logger"
	"context"

	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
)

func InitMfaRpc() {
	r, err := etcd.NewEtcdResolver([]string{config.Cfg.EtcdAddr})
	if err != nil {
		logger.Error("etcd resolver error", zap.Error(err))
	}
	cli, err := mfaservice.NewClient("mfaService", client.WithResolver(r))
	if err != nil {
		logger.Error("mfa service client error", zap.Error(err))
	}
	mfaClient = cli
}

func MfaQrCodeRpc(ctx context.Context, req *mfa.MfaQrcodeReq) (int32, string, string, error) {
	resp, err := mfaClient.MfaQrcode(ctx, req)
	if err != nil || resp == nil || resp.Data == nil {
		return consts.MfaCodeFalse, "", "", err
	}
	return resp.Code, resp.Data.Secret, resp.Data.Qrcode, nil
}
func MfaBindRpc(ctx context.Context, req *mfa.MfaBindReq) (int32, error) {
	resp, err := mfaClient.MfaBind(ctx, req)
	if err != nil || resp == nil {
		return consts.MfaCodeFalse, err
	}
	return resp.Code, nil
}
