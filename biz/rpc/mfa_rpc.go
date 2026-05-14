package rpc

import (
	"Tiktok/kitex_gen/mfa"
	"Tiktok/kitex_gen/mfa/mfaservice"
	"context"
	"log"

	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func InitMfaRpc() {
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Println(err)
	}
	cli, err := mfaservice.NewClient("mfaService", client.WithResolver(r))
	if err != nil {
		log.Println(err)
	}
	mfaClient = cli
}

func MfaQrCodeRpc(ctx context.Context, req *mfa.MfaQrcodeReq) (int32, string, string, error) {
	resp, err := mfaClient.MfaQrcode(ctx, req)
	return resp.Code, resp.Data.Secret, resp.Data.Qrcode, err
}
func MfaBindRpc(ctx context.Context, req *mfa.MfaBindReq) (int32, error) {
	resp, err := mfaClient.MfaBind(ctx, req)
	return resp.Code, err
}
