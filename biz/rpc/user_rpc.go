package rpc

import (
	"context"

	model "Tiktok/biz/model/user"
	"Tiktok/internal/config"
	"Tiktok/kitex_gen/user"
	"Tiktok/kitex_gen/user/userservice"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/logger"
	"Tiktok/pkg/utils"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/transmeta"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
)

func InitUserRpc() {
	r, err := etcd.NewEtcdResolver([]string{config.GetCfg().EtcdAddr})
	if err != nil {
		logger.Fatal("etcd resolver error", zap.Error(err))
	}
	cli, err := userservice.NewClient("userService", client.WithResolver(r), client.WithMetaHandler(transmeta.MetainfoClientHandler))
	if err != nil {
		logger.Fatal("user service client error", zap.Error(err))
	}
	userClient = cli
}
func RegisterRpc(ctx context.Context, req *user.RegisterReq) (int32, error) {
	defer utils.TrackRPC(ctx, "user", "UserRegister")()
	resp, err := userClient.UserRegister(ctx, req)
	if err != nil || resp == nil {
		return consts.UserReqValidError, err
	}
	return resp.Code, nil
}
func LoginRpc(ctx context.Context, req *user.LoginReq) (code int32, data *model.UserInfo, reToken string, acToken string, err error) {
	defer utils.TrackRPC(ctx, "user", "UserLogin")()
	resp, err := userClient.UserLogin(ctx, req)
	if err != nil || resp == nil {
		return consts.UserReqValidError, &model.UserInfo{}, "", "", err
	}
	return resp.Code, utils.ToBizUserInfo(resp.Data), resp.RefreshToken, resp.AccessToken, nil
}
func UserInfoRpc(ctx context.Context, req *user.UserInfoReq) (int32, *model.UserInfo, error) {
	defer utils.TrackRPC(ctx, "user", "UserInfo")()
	resp, err := userClient.UserInfo(ctx, req)
	if err != nil || resp == nil {
		return consts.UserReqValidError, &model.UserInfo{}, err
	}
	return resp.Code, utils.ToBizUserInfo(resp.Data), nil
}
func UserAvatarRpc(ctx context.Context, req *user.UserAvatarReq) (int32, *model.UserInfo, error) {
	defer utils.TrackRPC(ctx, "user", "UserAvatar")()
	resp, err := userClient.UserAvatar(ctx, req)
	if err != nil || resp == nil {
		return consts.UserReqValidError, &model.UserInfo{}, err
	}
	return resp.Code, utils.ToBizUserInfo(resp.Data), nil
}
func RefreshTokenRpc(ctx context.Context, req *user.RefreshTokenReq) (int32, string, string, error) {
	defer utils.TrackRPC(ctx, "user", "RefreshToken")()
	resp, err := userClient.RefreshToken(ctx, req)
	if err != nil || resp == nil {
		return consts.UserReqValidError, "", "", err
	}
	return resp.Code, resp.RefreshToken, resp.AccessToken, nil
}

func MfaQrCodeRpc(ctx context.Context, req *user.MfaQrcodeReq) (int32, string, string, error) {
	defer utils.TrackRPC(ctx, "user", "MfaQrcode")()
	resp, err := userClient.MfaQrcode(ctx, req)
	if err != nil || resp == nil || resp.Data == nil {
		return consts.MfaCodeFalse, "", "", err
	}
	return resp.Code, resp.Data.Secret, resp.Data.Qrcode, nil
}

func MfaBindRpc(ctx context.Context, req *user.MfaBindReq) (int32, error) {
	defer utils.TrackRPC(ctx, "user", "MfaBind")()
	resp, err := userClient.MfaBind(ctx, req)
	if err != nil || resp == nil {
		return consts.MfaCodeFalse, err
	}
	return resp.Code, nil
}
