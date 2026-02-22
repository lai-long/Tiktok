package handler

import (
	"Tiktok/biz/model/dto"
	"Tiktok/biz/service"
	"Tiktok/pkg/consts"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func RelationAction(ctx context.Context, c *app.RequestContext) {
	toUserId := c.PostForm("to_user_id")
	actionType := c.PostForm("action_type")
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(200, dto.Response{
			Base: dto.Base{
				Code: consts.CodeError,
				Msg:  "RelationAction userId exist error",
			},
		})
	}
	code, msg := service.RelationAction(toUserId, actionType, userId.(string))
	c.JSON(200, dto.Response{
		Base: dto.Base{
			Code: code,
			Msg:  msg,
		},
	})
}
func FollowingList(ctx context.Context, c *app.RequestContext) {
	userId := c.Query("user_id")
	pageNum := c.Query("page_num")
	pageSize := c.Query("page_size")
	code, msg, users, ok := service.FollowingList(userId, pageNum, pageSize)
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
			Items: dto.Items{
				User: users,
			},
			Total: dto.Total{
				Total: len(users),
			},
		},
	})
}
func FollowerList(ctx context.Context, c *app.RequestContext) {
	userId := c.Query("user_id")
	pageNum := c.Query("page_num")
	pageSize := c.Query("page_size")
	code, msg, followers, ok := service.FollowerList(userId, pageNum, pageSize)
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
			Items: dto.Items{
				User: followers,
			},
			Total: dto.Total{
				Total: len(followers),
			},
		},
	})
}
func FriendList(ctx context.Context, c *app.RequestContext) {
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
	}
	code, msg, friend, ok := service.FriendList(userId.(string), pageNum, pageSize)
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
			Items: dto.Items{
				User: friend,
			},
			Total: dto.Total{
				Total: len(friend),
			},
		},
	})
}
