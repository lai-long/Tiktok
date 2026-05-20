package service

import (
	"Tiktok/kitex_gen/mfa/mfaservice"

	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func NewMfaRpcClient() (mfaservice.Client, error) {
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		return nil, err
	}
	return mfaservice.NewClient("mfaService", client.WithResolver(r))
}
