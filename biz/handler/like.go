package handler

import (
	"Tiktok/biz/model/dto"

	"Tiktok/pkg/consts"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type VideoLikeSever interface {
	VideoLikeAction(userId string, videoId string, action string) (int, string)
	LikeList(userId string, pageNum string, pageSize string) (int, string, []dto.Video, bool)
}
type CommentLikeSever interface {
	CommentLikeAction(userId string, commentId string, action string) (int, string)
}
type LikesHandler struct {
	video   VideoLikeSever
	comment CommentLikeSever
}

func NewLikesHandler(video VideoLikeSever, comment CommentLikeSever) *LikesHandler {
	return &LikesHandler{
		video:   video,
		comment: comment,
	}
}

func (h *LikesHandler) LikeAction(ctx context.Context, c *app.RequestContext) {
	action := c.PostForm("action_type")
	videoId := c.PostForm("video_id")
	commentId := c.PostForm("comment_id")
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: consts.CodeLikeError,
			Msg:  "likeAction Get userId error",
		}})
		return
	}
	if commentId == "" && videoId == "" {
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: consts.CodeLikeError,
			Msg:  "likeAction Get commentId or videoId error",
		}})
	}
	if videoId != "" {
		code, msg := h.video.VideoLikeAction(userId.(string), videoId, action)
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: code,
			Msg:  msg,
		}})
	}
	if commentId != "" {
		code, msg := h.comment.CommentLikeAction(userId.(string), commentId, action)
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: code,
			Msg:  msg,
		}})
	}
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
