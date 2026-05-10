package rpc

import (
	"Tiktok/kitex_gen/user"
	"Tiktok/kitex_gen/user/userservice"
	"context"

	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func InitUserRpc() {
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
	}
	cli, err := userservice.NewClient("userService", client.WithResolver(r))
	userClient = cli
}
func RegisterRpc(ctx context.Context, req *user.RegisterReq) (int32, error) {
	resp, err := userClient.UserRegister(ctx, req)
	return resp.Code, err
}
func LoginRpc(ctx context.Context, req *user.LoginReq) (code int32, data *user.UserInfo, reToken string, acToken string, err error) {
	resp, err := userClient.UserLogin(ctx, req)
	return resp.Code, resp.Data, resp.RefreshToken, resp.AccessToken, err
}
func UserInfoRpc(ctx context.Context, req *user.UserInfoReq) (int32, *user.UserInfo, error) {
	resp, err := userClient.UserInfo(ctx, req)
	return resp.Code, resp.Data, err
}
func UserAvatarRpc(ctx context.Context, req *user.UserAvatarReq) (int32, *user.UserInfo, error) {
	resp, err := userClient.UserAvatar(ctx, req)
	return resp.Code, resp.Data, err
}
func RefreshTokenRpc(ctx context.Context, req *user.RefreshTokenReq) (int32, string, string, error) {
	resp, err := userClient.RefreshToken(ctx, req)
	return resp.Code, resp.RefreshToken, resp.AccessToken, err
}
