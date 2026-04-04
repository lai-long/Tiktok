package handler

import (
	"Tiktok/biz/model/common"
	"Tiktok/biz/model/react"
	"Tiktok/biz/model/video"

	"Tiktok/pkg/consts"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type LikeSever interface {
	LikeAction(userId string, videoId string, action string, targetType string) (int, string)
	LikeList(userId string, pageNum string, pageSize string) (int, string, []*video.VideoInfo, bool)
}

type LikesHandler struct {
	likeService LikeSever
}

func NewLikesHandler(like LikeSever) *LikesHandler {
	return &LikesHandler{
		likeService: like,
	}
}

func (h *LikesHandler) LikeAction(ctx context.Context, c *app.RequestContext) {
	//targetType 1、视频 2、评论
	likeActionReq := new(react.LikeActionReq)
	err := c.BindAndValidate(likeActionReq)
	if err != nil {
		c.JSON(200, react.LikeActionResp{Base: &common.Base{
			Code: consts.CodeLikeError,
			Msg:  "LikeAction BindAndValidate error",
		}})
		return
	}
	userId, ok := ctx.Value("user_id").(string)
	if !ok {
		c.JSON(200, react.LikeActionResp{Base: &common.Base{
			Code: consts.CodeLikeError,
			Msg:  "UserId not found in context",
		}})
		return
	}
	if likeActionReq.TargetType == "" {
		c.JSON(200, react.LikeActionResp{Base: &common.Base{
			Code: consts.CodeLikeError,
			Msg:  "likeActionReq.TargetType err is null",
		}})
		return
	}
	code, msg := h.likeService.LikeAction(userId, likeActionReq.TargetAt, likeActionReq.ActionType, likeActionReq.TargetType)
	c.JSON(200, react.LikeActionResp{Base: &common.Base{
		Code: int32(code),
		Msg:  msg,
	}})
	return

}

func (h *LikesHandler) LikeList(ctx context.Context, c *app.RequestContext) {
	likeListReq := new(react.LikeListReq)
	err := c.BindAndValidate(likeListReq)
	if err != nil {
		c.JSON(200, react.LikeListResp{Base: &common.Base{
			Code: consts.CodeLikeError,
			Msg:  "likeListReq BindAndValidate error",
		}})
	}
	code, msg, videos, _ := h.likeService.LikeList(likeListReq.UserId, likeListReq.PageNum, likeListReq.PageSize)
	c.JSON(200, react.LikeListResp{Base: &common.Base{
		Code: int32(code),
		Msg:  msg,
	}, Data: &react.LikeVideoData{Items: videos}})
}
