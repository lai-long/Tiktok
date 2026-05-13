package rpc

import (
	"Tiktok/kitex_gen/mfa/mfaservice"
	"Tiktok/kitex_gen/user/userservice"
	"Tiktok/kitex_gen/video/videoservice"
)

var (
	userClient  userservice.Client
	mfaClient   mfaservice.Client
	videoClient videoservice.Client
)

func Init() {
	InitUserRpc()
	InitMfaRpc()
}
