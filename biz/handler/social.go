package handler

import (
	"Tiktok/biz/model/dto"
	"Tiktok/pkg/consts"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func (h *SocialHandler) RelationAction(ctx context.Context, c *app.RequestContext) {
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
	code, msg := h.socialService.RelationAction(toUserId, actionType, userId.(string))
	c.JSON(200, dto.Response{
		Base: dto.Base{
			Code: code,
			Msg:  msg,
		},
	})
}
func (h *SocialHandler) FollowingList(ctx context.Context, c *app.RequestContext) {
	userId := c.Query("user_id")
	pageNum := c.Query("page_num")
	pageSize := c.Query("page_size")
	code, msg, users, ok := h.socialService.FollowingList(userId, pageNum, pageSize)
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
func (h *SocialHandler) FollowerList(ctx context.Context, c *app.RequestContext) {
	userId := c.Query("user_id")
	pageNum := c.Query("page_num")
	pageSize := c.Query("page_size")
	code, msg, followers, ok := h.socialService.FollowerList(userId, pageNum, pageSize)
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
func (h *SocialHandler) FriendList(ctx context.Context, c *app.RequestContext) {
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
	code, msg, friend, ok := h.socialService.FriendList(userId.(string), pageNum, pageSize)
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
