package handler

import (
	"Tiktok/biz/model/dto"
	"Tiktok/biz/service"
	"Tiktok/pkg/consts"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func CommentPublish(ctx context.Context, c *app.RequestContext) {
	var comment dto.Comment
	err := c.Bind(&comment)
	if err != nil {
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: consts.CodeCommentError,
			Msg:  "CommentPublish comment Bind error:",
		}})
	}
	code, msg := service.CommentPublish(comment.VideoId, comment.UserId, comment.Content)
	c.JSON(200, dto.Response{Base: dto.Base{
		Code: code,
		Msg:  msg,
	}})

}
func CommentList(ctx context.Context, c *app.RequestContext) {
	videoId := c.PostForm("video_id")
	pageSize := c.PostForm("page_size")
	pageNum := c.PostForm("page_num")
	code, msg, comments, ok := service.CommentList(videoId, pageSize, pageNum)
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
		Data: dto.Items{
			Comment: comments,
			Total:   len(comments),
		}})
}

func CommentDelete(ctx context.Context, c *app.RequestContext) {
	videoId := c.PostForm("video_id")
	commentId := c.PostForm("comment_id")
	code, msg := service.CommentDelete(commentId, videoId)
	c.JSON(200, dto.Response{Base: dto.Base{
		Code: code,
		Msg:  msg,
	}})
}
