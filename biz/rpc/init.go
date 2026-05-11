package rpc

import (
	"Tiktok/kitex_gen/mfa/mfaservice"
	"Tiktok/kitex_gen/user/userservice"
)

var (
	userClient userservice.Client
	mfaClient  mfaservice.Client
)

func Init() {
	InitUserRpc()
	InitMfaRpc()
}
