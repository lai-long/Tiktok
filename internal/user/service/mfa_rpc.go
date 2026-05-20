package service

import (
	"Tiktok/kitex_gen/mfa"
	"Tiktok/kitex_gen/mfa/mfaservice"
	"Tiktok/pkg/consts"
	"context"

	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

type MfaService interface {
	MfaConfirm(ctx context.Context, userID string, mfaCode string) (int32, error)
}

type mfaRpcClient struct {
	cli mfaservice.Client
}

func NewMfaService() (MfaService, error) {
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		return nil, err

	}
	cli, err := mfaservice.NewClient("mfaService", client.WithResolver(r))
	if err != nil {
		return nil, err
	}
	return &mfaRpcClient{cli: cli}, nil
}

func (c *mfaRpcClient) MfaConfirm(ctx context.Context, userID string, mfaCode string) (int32, error) {
	resp, err := c.cli.MfaConfirm(ctx, &mfa.MfaConfirmReq{
		UserID: userID,
		QrCode: mfaCode,
	})
	if err != nil || resp == nil {
		return consts.MfaCodeFalse, err
	}
	return resp.Code, nil
}
