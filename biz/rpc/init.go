package rpc

import (
	"Tiktok/kitex_gen/mfa/mfaservice"
	"Tiktok/kitex_gen/react/commentservice"
	"Tiktok/kitex_gen/react/likeservice"
	"Tiktok/kitex_gen/social/socialservice"
	"Tiktok/kitex_gen/user/userservice"
	"Tiktok/kitex_gen/video/videoservice"
)

var (
	userClient    userservice.Client
	mfaClient     mfaservice.Client
	videoClient   videoservice.Client
	commentClient commentservice.Client
	likeClient    likeservice.Client
	socialClient  socialservice.Client
)

func Init() {
	InitUserRpc()
	InitMfaRpc()
	InitVideoRpc()
	InitReactRpc()
	InitSocialRpc()
}
