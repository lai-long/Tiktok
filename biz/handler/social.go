package handler

import (
	"Tiktok/biz/model/common"
	"Tiktok/biz/model/dto"
	"Tiktok/biz/model/social"
	"Tiktok/biz/model/user"
	"Tiktok/pkg/consts"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type SocialSever interface {
	RelationAction(toUserId string, actionType string, userId string) (int, string)
	FollowingList(userId string, pageNum int64, pageSize int64) (int, string, []*user.UserInfo, bool)
	FollowerList(userId string, pageNum int64, pageSize int64) (int, string, []dto.User, bool)
	FriendList(userId string, pageNum int64, pageSize int64) (int, string, []dto.User, bool)
}
type SocialHandler struct {
	socialService SocialSever
}

func NewSocialHandler(service SocialSever) *SocialHandler {
	return &SocialHandler{socialService: service}
}

func (h *SocialHandler) RelationAction(ctx context.Context, c *app.RequestContext) {
	req := new(social.RelationActionReq)
	err := c.BindAndValidate(req)
	if err != nil {
		c.JSON(200, social.RelationActionResp{
			Base: &common.Base{
				Code: consts.CodeError,
				Msg:  "social.RelationActionResp err",
			},
		})
	}
	userId, ok := ctx.Value("user_id").(string)
	if !ok {
		c.JSON(200, social.RelationActionResp{
			Base: &common.Base{
				Code: consts.CodeError,
				Msg:  "social.RelationActionResp err",
			},
		})
		return
	}
	code, msg := h.socialService.RelationAction(req.ToUserId, req.ActionType, userId)
	c.JSON(200, social.RelationActionResp{
		Base: &common.Base{
			Code: int32(code),
			Msg:  msg,
		},
	})
}

func (h *SocialHandler) FollowingList(ctx context.Context, c *app.RequestContext) {
	req := new(social.FollowingListReq)
	err := c.BindAndValidate(req)
	if err != nil {
		c.JSON(200, social.FollowingListResp{
			Base: &common.Base{
				Code: consts.CodeError,
				Msg:  "social.FollowingListResp err",
			},
		})
	}
	code, msg, userInfos, _ := h.socialService.FollowingList(req.UserId, req.PageNum, req.PageSize)
	c.JSON(200, social.FollowingListResp{
		Base: &common.Base{
			Code: int32(code),
			Msg:  msg,
		},
		Data: userInfos,
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
	userId, ok := ctx.Value("user_id").(string)
	if !ok {
		c.JSON(200, dto.Response{
			Base: dto.Base{
				Code: consts.CodeRelationError,
				Msg:  "FriendList userId exist not exist",
			},
		})
		return
	}
	code, msg, friend, ok := h.socialService.FriendList(userId, pageNum, pageSize)
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
