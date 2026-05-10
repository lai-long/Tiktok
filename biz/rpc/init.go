package rpc

import "Tiktok/kitex_gen/user/userservice"

var (
	userClient userservice.Client
)

func Init() {
	InitUserRpc()
}
