package handler

import (
	"Tiktok/biz/model/dto"
	"Tiktok/biz/service"
	"Tiktok/pkg/consts"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func LikeAction(ctx context.Context, c *app.RequestContext) {
	action := c.PostForm("action_type")
	id := c.PostForm("video_id")
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: consts.CodeLikeError,
			Msg:  "likeAction Get userId error",
		}})
		return
	}
	code, msg := service.LikeAction(userId.(string), id, action)
	c.JSON(200, dto.Response{Base: dto.Base{
		Code: code,
		Msg:  msg,
	}})
}
func LikeList(ctx context.Context, c *app.RequestContext) {
	pageSize := c.Query("page_size")
	pageNum := c.Query("page_num")
	userId := c.Query("user_id")
	code, msg, videos, ok := service.LikeList(userId, pageSize, pageNum)
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
				Items: dto.Items{Video: videos},
				Total: dto.Total{
					Total: len(videos),
				},
			},
		},
	)
}
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
		Data: dto.Data{
			dto.Items{Comment: comments},
			dto.Total{
				Total: len(comments),
			},
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
