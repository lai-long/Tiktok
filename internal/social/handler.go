package social

import (
	"Tiktok/internal/social/service"
	social "Tiktok/kitex_gen/social"
	"Tiktok/pkg/consts"
	"context"
)

type SocialServiceImpl struct {
	socialService *service.SocialRepo
}

func NewSocialServiceImpl(socialRepo *service.SocialRepo) *SocialServiceImpl {
	return &SocialServiceImpl{socialService: socialRepo}
}

func (s *SocialServiceImpl) RelationAction(ctx context.Context, req *social.RelationActionReq) (resp *social.RelationActionResp, err error) {
	resp = &social.RelationActionResp{}
	code, _ := s.socialService.RelationAction(ctx, req.ToUserId, req.ActionType, req.UserId)
	resp.Code = code
	resp.Msg = consts.GetErrorCodeMsg(code)
	return resp, nil
}

func (s *SocialServiceImpl) FollowingList(ctx context.Context, req *social.FollowingListReq) (resp *social.FollowingListResp, err error) {
	resp = &social.FollowingListResp{}
	code, userInfos, _ := s.socialService.FollowingList(req.UserId, req.PageNum, req.PageSize)
	resp.Code = code
	resp.Msg = consts.GetErrorCodeMsg(code)
	resp.Data = &social.SocialData{
		Items: userInfos,
		Total: int64(len(userInfos)),
	}
	return resp, nil
}

func (s *SocialServiceImpl) FollowerList(ctx context.Context, req *social.FollowerListReq) (resp *social.FollowerListResp, err error) {
	resp = &social.FollowerListResp{}
	code, userInfos, _ := s.socialService.FollowerList(req.UserId, req.PageNum, req.PageSize)
	resp.Code = code
	resp.Msg = consts.GetErrorCodeMsg(code)
	resp.Data = &social.SocialData{
		Items: userInfos,
		Total: int64(len(userInfos)),
	}
	return resp, nil
}

func (s *SocialServiceImpl) FriendList(ctx context.Context, req *social.FriendListReq) (resp *social.FriendListResp, err error) {
	resp = &social.FriendListResp{}
	code, userInfos, _ := s.socialService.FriendList(req.UserId, req.PageNum, req.PageSize)
	resp.Code = code
	resp.Msg = consts.GetErrorCodeMsg(code)
	resp.Data = &social.SocialData{
		Items: userInfos,
		Total: int64(len(userInfos)),
	}
	return resp, nil
}
