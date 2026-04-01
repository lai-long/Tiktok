package handler

import (
	"Tiktok/biz/model/dto"

	"Tiktok/pkg/consts"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type LikeSever interface {
	LikeAction(userId string, videoId string, action string, targetType string) (int, string)
	LikeList(userId string, pageNum string, pageSize string) (int, string, []dto.Video, bool)
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
	action := c.PostForm("action_type")
	targetId := c.PostForm("target_id")
	//targetType 1、视频 2、评论
	targetType := c.PostForm("target_type")
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: consts.CodeLikeError,
			Msg:  "likeAction Get userId error",
		}})
		return
	}
	if targetId == "" {
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: consts.CodeLikeError,
			Msg:  "likeAction Get commentId or videoId error",
		}})
		return
	}
	code, msg := h.likeService.LikeAction(userId.(string), targetId, action, targetType)
	c.JSON(200, dto.Response{Base: dto.Base{
		Code: code,
		Msg:  msg,
	}})
	return

}
func (h *LikesHandler) LikeList(ctx context.Context, c *app.RequestContext) {
	pageSize := c.Query("page_size")
	pageNum := c.Query("page_num")
	userId := c.Query("user_id")
	code, msg, videos, ok := h.video.LikeList(userId, pageNum, pageSize)
	if !ok {
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: code,
			Msg:  msg,
		}})
		return
	}

	c.JSON(
		200,
		dto.Response{
			Base: dto.Base{
				Code: code,
				Msg:  msg,
			},
			Data: dto.Data{
				Items: videos,
				Total: len(videos),
			},
		},
	)
}
