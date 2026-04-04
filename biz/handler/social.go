package handler

import (
	"Tiktok/biz/model/common"
	"Tiktok/biz/model/social"
	"Tiktok/biz/model/user"
	"Tiktok/pkg/consts"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type SocialSever interface {
	RelationAction(toUserId string, actionType string, userId string) (int, string)
	FollowingList(userId string, pageNum int64, pageSize int64) (int, string, []*user.UserInfo, bool)
	FollowerList(userId string, pageNum int64, pageSize int64) (int, string, []*user.UserInfo, bool)
	FriendList(userId string, pageNum int64, pageSize int64) (int, string, []*user.UserInfo, bool)
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
		Data: &social.SocialData{
			Items: userInfos,
			Total: int64(len(userInfos)),
		},
	})
}

func (h *SocialHandler) FollowerList(ctx context.Context, c *app.RequestContext) {
	req := new(social.FollowerListReq)
	err := c.BindAndValidate(req)
	if err != nil {
		c.JSON(200, social.FollowerListResp{
			Base: &common.Base{
				Code: consts.CodeError,
				Msg:  "social.FollowerListResp err",
			},
		})
	}
	code, msg, followers, _ := h.socialService.FollowerList(req.UserId, req.PageNum, req.PageSize)
	c.JSON(200, social.FollowerListResp{
		Base: &common.Base{
			Code: int32(code),
			Msg:  msg,
		},
		Data: &social.SocialData{
			Items: followers,
			Total: int64(len(followers)),
		},
	})
}

func (h *SocialHandler) FriendList(ctx context.Context, c *app.RequestContext) {
	req := new(social.FriendListReq)
	err := c.BindAndValidate(req)
	if err != nil {
		c.JSON(200, social.FriendListResp{
			Base: &common.Base{
				Code: consts.CodeError,
				Msg:  "social.FriendListResp err",
			},
		})
	}
	userId, ok := ctx.Value("user_id").(string)
	if !ok {
		c.JSON(200, social.FriendListResp{
			Base: &common.Base{
				Code: consts.CodeError,
				Msg:  "social.FriendListResp err",
			},
		})
		return
	}
	code, msg, friendsInfo, _ := h.socialService.FriendList(userId, req.PageNum, req.PageSize)
	c.JSON(200, social.FriendListResp{
		Base: &common.Base{
			Code: int32(code),
			Msg:  msg,
		},
		Data: &social.SocialData{
			Items: friendsInfo,
			Total: int64(len(friendsInfo)),
		},
	})
}
