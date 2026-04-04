package handler

import (
	"Tiktok/biz/model/common"
	"Tiktok/biz/model/dto"
	"Tiktok/biz/model/react"
	"Tiktok/pkg/consts"
	"context"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
)

type CommentSever interface {
	CommentPublish(targetId, userId, content, targetType string) (int, string)
	CommentList(targetId string, pageSize string, pageNum string) (int, string, []*react.CommentInfo, bool)
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
	commentPublishReq := new(react.CommentPublishReq)
	err := c.BindAndValidate(commentPublishReq)
	if err != nil {
		c.JSON(200, react.CommentPublishResp{
			Base: &common.Base{
				Code: consts.CodeError,
				Msg:  "CommentPublish BindAndValidate error",
			},
		})
		return
	}
	userId, ok := ctx.Value("user_id").(string)
	if !ok {
		c.JSON(200, react.CommentPublishResp{
			Base: &common.Base{
				Code: consts.CodeError,
				Msg:  "CommentPublish ctx.Value userid error",
			},
		})
		log.Printf("CommentPublish userId exists error: %v\n", err)
		return
	}
	code, msg := h.service.CommentPublish(commentPublishReq.TargetAt, userId, commentPublishReq.Content, commentPublishReq.TargetType)
	c.JSON(200, react.CommentPublishResp{
		Base: &common.Base{
			Code: int32(code),
			Msg:  msg,
		},
	})

}

func (h *CommentHandler) CommentList(ctx context.Context, c *app.RequestContext) {
	commentListReq := new(react.CommentListReq)
	err := c.BindAndValidate(commentListReq)
	if err != nil {
		c.JSON(200, react.CommentListResp{
			Base: &common.Base{
				Code: consts.CodeError,
				Msg:  "CommentList BindAndValidate error",
			},
		})
	}
	code, msg, commentInfos, _ := h.service.CommentList(commentListReq.TargetAt, commentListReq.PageSize, commentListReq.PageNum)
	c.JSON(200, react.CommentListResp{
		Base: &common.Base{
			Code: int32(code),
			Msg:  msg,
		},
		Data: &react.CommentData{Items: commentInfos},
	})
}

func (h *CommentHandler) CommentDelete(ctx context.Context, c *app.RequestContext) {
	commentDeleteReq := new(react.CommentDeleteReq)
	err := c.BindAndValidate(commentDeleteReq)
	if err != nil {
		c.JSON(200, react.CommentDeleteResp{
			Base: &common.Base{
				Code: consts.CodeError,
				Msg:  "CommentDelete BindAndValidate error",
			},
		})
	}
	userId, ok := ctx.Value("user_id").(string)
	if !ok {
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: consts.CodeCommentError,
			Msg:  "commentDelete Get userId error",
		}})
	}
	code, msg := h.service.CommentDelete(commentDeleteReq.CommentId, commentDeleteReq.TargetAt, userId, commentDeleteReq.TargetType)
	c.JSON(200, react.CommentDeleteResp{
		Base: &common.Base{
			Code: int32(code),
			Msg:  msg,
		},
	})
}
