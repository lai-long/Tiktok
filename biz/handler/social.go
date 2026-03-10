package handler

import (
	"Tiktok/biz/model/dto"
	"Tiktok/pkg/consts"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func (h *Handler) RelationAction(ctx context.Context, c *app.RequestContext) {
	toUserId := c.PostForm("to_user_id")
	actionType := c.PostForm("action_type")
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(200, dto.Response{
			Base: dto.Base{
				Code: consts.CodeRelationError,
				Msg:  "RelationAction userId exist error",
			},
		})
	}
	code, msg := h.service.RelationAction(toUserId, actionType, userId.(string))
	c.JSON(200, dto.Response{
		Base: dto.Base{
			Code: code,
			Msg:  msg,
		},
	})
}
func (h *Handler) FollowingList(ctx context.Context, c *app.RequestContext) {
	userId := c.Query("user_id")
	pageNum := c.Query("page_num")
	pageSize := c.Query("page_size")
	code, msg, users, ok := h.service.FollowingList(userId, pageNum, pageSize)
	if !ok {
		c.JSON(200, dto.Response{
			Base: dto.Base{
				Code: code,
				Msg:  msg,
			},
		})
		return
	}
	c.JSON(200, dto.Response{
		Base: dto.Base{
			Code: code,
			Msg:  msg,
		},
		Data: dto.Data{
			Items: users,
			Total: len(users),
		},
	})
}
func (h *Handler) FollowerList(ctx context.Context, c *app.RequestContext) {
	userId := c.Query("user_id")
	pageNum := c.Query("page_num")
	pageSize := c.Query("page_size")
	code, msg, followers, ok := h.service.FollowerList(userId, pageNum, pageSize)
	if !ok {
		c.JSON(200, dto.Response{
			Base: dto.Base{
				Code: code,
				Msg:  msg,
			},
		})
		return
	}
	c.JSON(200, dto.Response{
		Base: dto.Base{
			Code: code,
			Msg:  msg,
		},
		Data: dto.Data{
			Items: followers,
			Total: len(followers),
		},
	})
}
func (h *Handler) FriendList(ctx context.Context, c *app.RequestContext) {
	pageNum := c.Query("page_num")
	pageSize := c.Query("page_size")
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(200, dto.Response{
			Base: dto.Base{
				Code: consts.CodeRelationError,
				Msg:  "FriendList userId exist not exist",
			},
		})
		return
	}
	code, msg, friend, ok := h.service.FriendList(userId.(string), pageNum, pageSize)
	if !ok {
		c.JSON(200, dto.Response{
			Base: dto.Base{
				Code: code,
				Msg:  msg,
			},
		})
		return
	}
	c.JSON(200, dto.Response{
		Base: dto.Base{
			Code: code,
			Msg:  msg,
		},
		Data: dto.Data{
			Items: friend,
			Total: len(friend),
		},
	})
}
