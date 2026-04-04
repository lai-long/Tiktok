package handler

import (
	"Tiktok/biz/model/dto"
	"Tiktok/pkg/consts"
	"context"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
)

type CommentSever interface {
	CommentPublish(targetId string, userId string, content string, targetType string) (int, string)
	CommentList(targetId string, pageSize string, pageNum string) (int, string, []dto.Comment, bool)
	CommentDelete(commentId string, target string, userId string, targetType string) (int, string)
}

type CommentHandler struct {
	service CommentSever
}

func NewCommentHandler(service CommentSever) *CommentHandler {
	return &CommentHandler{
		service: service,
	}
}

func (h *CommentHandler) CommentPublish(ctx context.Context, c *app.RequestContext) {
	var comment dto.Comment
	err := c.Bind(&comment)
	if err != nil {
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: consts.CodeCommentError,
			Msg:  "CommentPublish comment Bind error:",
		}})
		return
	}
	userId, ok := ctx.Value("user_id").(string)
	if !ok {
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: consts.CodeCommentError,
			Msg:  "CommentPublish userId exists error",
		}})
		log.Println("CommentPublish userId exists error: %v", err)
		return
	}
	comment.UserId = userId
	code, msg := h.service.CommentPublish(comment.TargetId, comment.UserId, comment.Content, comment.TargetType)
	c.JSON(200, dto.Response{Base: dto.Base{
		Code: code,
		Msg:  msg,
	}})

}

func (h *CommentHandler) CommentList(ctx context.Context, c *app.RequestContext) {
	targetId := c.Query("target_id")
	pageSize := c.Query("page_size")
	pageNum := c.Query("page_num")
	code, msg, comments, ok := h.service.CommentList(targetId, pageSize, pageNum)
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

func (h *CommentHandler) CommentDelete(ctx context.Context, c *app.RequestContext) {
	commentId := c.PostForm("comment_id")
	targetType := c.PostForm("target_type")
	targetId := c.PostForm("target_id")
	userId, ok := ctx.Value("user_id").(string)
	if !ok {
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: consts.CodeCommentError,
			Msg:  "commentDelete Get userId error",
		}})
	}
	code, msg := h.service.CommentDelete(commentId, targetId, userId, targetType)
	c.JSON(200, dto.Response{Base: dto.Base{
		Code: code,
		Msg:  msg,
	}})
}
