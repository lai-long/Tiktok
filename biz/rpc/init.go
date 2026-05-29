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

// Test mock client
func SetUserClient(c userservice.Client)       { userClient = c }
func GetUserClient() userservice.Client        { return userClient }
func SetVideoClient(c videoservice.Client)     { videoClient = c }
func GetVideoClient() videoservice.Client      { return videoClient }
func SetMfaClient(c mfaservice.Client)         { mfaClient = c }
func GetMfaClient() mfaservice.Client          { return mfaClient }
func SetCommentClient(c commentservice.Client) { commentClient = c }
func GetCommentClient() commentservice.Client  { return commentClient }
func SetLikeClient(c likeservice.Client)       { likeClient = c }
func GetLikeClient() likeservice.Client        { return likeClient }
func SetSocialClient(c socialservice.Client)   { socialClient = c }
func GetSocialClient() socialservice.Client    { return socialClient }
