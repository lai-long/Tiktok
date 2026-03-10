package handler

import (
	"Tiktok/biz/model/dto"
	"Tiktok/biz/service"
	"context"
	"mime/multipart"
)

type LikeHandler struct {
	service LikeSever
}

func NewLikeHandler(service *service.Service) *LikeHandler {
	return &LikeHandler{service: service}
}

type UserHandler struct {
	userService UserSever
	MfaServer   MfaServer
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService,
		MfaServer: userService}
}

type VideoHandler struct {
	videoService VideoSever
}

func NewVideoHandler(videoService VideoSever) *VideoHandler {
	return &VideoHandler{videoService: videoService}
}

type MfaServer interface {
	GenerateMfa(username string, userId string) (bool, string, string, int, string)
	MfaBindByCode(code string, userId string) (int, string)
	MfaBindBySecret(secret string, userId string) (int, string)
}
type UserSever interface {
	Register(userinfo dto.User) (int, string)
	Login(userDto dto.User, mfaCode string) (int, string, dto.User, string, string)
	UserInfo(userId string) (dto.User, int, string, bool)
	UserAvatar(data *multipart.FileHeader, userId interface{}) (int, string, bool, dto.User)
}
type VideoSever interface {
	VideoPublish(video dto.Video, data *multipart.FileHeader, ctx context.Context) (int, string)
	VideoList(userId string, pageSize string, pageNum string) (int, string, []dto.Video, bool)
	VideoSearch(keyword string, pageNum string, pageSize string) (int, string, []dto.Video, bool)
	VideoPopular(ctx context.Context, pageNum string, pageSize string) (int, string, []dto.Video, bool)
}
type LikeSever interface {
	VideoLikeAction(userId string, videoId string, action string) (int, string)
	CommentLikeAction(userId string, commentId string, action string) (int, string)
	LikeList(userId string, pageNum string, pageSize string) (int, string, []dto.Video, bool)
	CommentPublish(videoId string, userId string, content string) (int, string)
	CommentList(videoId string, pageSize string, pageNum string) (int, string, []dto.Comment, bool)
	CommentDelete(commentId string, videoId string, userId string) (int, string)
}
type SocialSever interface {
	RelationAction(toUserId string, actionType string, userId string) (int, string)
	FollowingList(userId string, pageNum string, pageSize string) (int, string, []dto.User, bool)
	FollowerList(userId string, pageNum string, pageSize string) (int, string, []dto.User, bool)
	FriendList(userId string, pageNum string, pageSize string) (int, string, []dto.User, bool)
}
type SocialHandler struct {
	socialService SocialSever
}

func NewSocialHandler(service SocialSever) *SocialHandler {
	return &SocialHandler{socialService: service}
}
