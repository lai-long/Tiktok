package handler

import (
	"Tiktok/biz/model/dto"
	"Tiktok/pkg/consts"
	"context"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
)

func (h *Handler) LikeAction(ctx context.Context, c *app.RequestContext) {
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
		code, msg := h.service.VideoLikeAction(userId.(string), videoId, action)
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: code,
			Msg:  msg,
		}})
	}
	if commentId != "" {
		code, msg := h.service.CommentLikeAction(userId.(string), commentId, action)
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: code,
			Msg:  msg,
		}})
	}
}
func (h *Handler) LikeList(ctx context.Context, c *app.RequestContext) {
	pageSize := c.Query("page_size")
	pageNum := c.Query("page_num")
	userId := c.Query("user_id")
	code, msg, videos, ok := h.service.LikeList(userId, pageNum, pageSize)
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
func (h *Handler) CommentPublish(ctx context.Context, c *app.RequestContext) {
	var comment dto.Comment
	err := c.Bind(&comment)
	if err != nil {
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: consts.CodeCommentError,
			Msg:  "CommentPublish comment Bind error:",
		}})
		return
	}
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: consts.CodeCommentError,
			Msg:  "CommentPublish userId exists error",
		}})
		log.Fatalf("CommentPublish userId exists error: %v", err)
		return
	}
	comment.UserId = userId.(string)
	code, msg := h.service.CommentPublish(comment.VideoId, comment.UserId, comment.Content)
	c.JSON(200, dto.Response{Base: dto.Base{
		Code: code,
		Msg:  msg,
	}})

}
func (h *Handler) CommentList(ctx context.Context, c *app.RequestContext) {
	videoId := c.Query("video_id")
	pageSize := c.Query("page_size")
	pageNum := c.Query("page_num")
	code, msg, comments, ok := h.service.CommentList(videoId, pageSize, pageNum)
	if !ok {
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: code,
			Msg:  msg,
		}})
		return
	}
	c.JSON(200, dto.Response{Base: dto.Base{
		Code: code,
		Msg:  msg,
	},
		Data: dto.Data{
			Items: comments,
			Total: len(comments),
		}})
}

func (h *Handler) CommentDelete(ctx context.Context, c *app.RequestContext) {
	videoId := c.PostForm("video_id")
	commentId := c.PostForm("comment_id")
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: consts.CodeCommentError,
			Msg:  "commentDelete Get userId error",
		}})
	}
	code, msg := h.service.CommentDelete(commentId, videoId, userId.(string))
	c.JSON(200, dto.Response{Base: dto.Base{
		Code: code,
		Msg:  msg,
	}})
}
