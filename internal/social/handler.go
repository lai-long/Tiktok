package social

import (
	social "Tiktok/kitex_gen/social"
	"Tiktok/pkg/consts"
	"context"
)

// socialService 是 SocialServiceImpl 依赖的最小接口。
type socialService interface {
	RelationAction(ctx context.Context, toUserId string, actionType string, userId string) (int32, error)
	FollowingList(ctx context.Context, userId string, pageNum int64, pageSize int64) (int32, []*social.UserInfo, error)
	FollowerList(ctx context.Context, userId string, pageNum int64, pageSize int64) (int32, []*social.UserInfo, error)
	FriendList(ctx context.Context, userId string, pageNum int64, pageSize int64) (int32, []*social.UserInfo, error)
}

type SocialServiceImpl struct {
	socialService socialService
}

func NewSocialServiceImpl(socialRepo socialService) *SocialServiceImpl {
	return &SocialServiceImpl{socialService: socialRepo}
}

func (s *SocialServiceImpl) RelationAction(ctx context.Context, req *social.RelationActionReq) (resp *social.RelationActionResp, err error) {
	resp = &social.RelationActionResp{}
	code, err := s.socialService.RelationAction(ctx, req.ToUserId, req.ActionType, req.UserId)
	if err != nil {
		resp.Code = code
		resp.Msg = consts.GetErrorCodeMsg(code)
		return resp, err
	}
	resp.Code = code
	resp.Msg = consts.GetErrorCodeMsg(code)
	return resp, nil
}

func (s *SocialServiceImpl) FollowingList(ctx context.Context, req *social.FollowingListReq) (resp *social.FollowingListResp, err error) {
	resp = &social.FollowingListResp{}
	code, userInfos, err := s.socialService.FollowingList(ctx, req.UserId, req.PageNum, req.PageSize)
	if err != nil {
		resp.Code = code
		resp.Msg = consts.GetErrorCodeMsg(code)
		resp.Data = &social.SocialData{
			Items: userInfos,
			Total: int64(len(userInfos)),
		}
		return resp, err
	}
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
	code, userInfos, err := s.socialService.FollowerList(ctx, req.UserId, req.PageNum, req.PageSize)
	if err != nil {
		resp.Code = code
		resp.Msg = consts.GetErrorCodeMsg(code)
		resp.Data = &social.SocialData{
			Items: userInfos,
			Total: int64(len(userInfos)),
		}
		return resp, err
	}
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
	code, userInfos, err := s.socialService.FriendList(ctx, req.UserId, req.PageNum, req.PageSize)
	if err != nil {
		resp.Code = code
		resp.Msg = consts.GetErrorCodeMsg(code)
		resp.Data = &social.SocialData{
			Items: userInfos,
			Total: int64(len(userInfos)),
		}
		return resp, err
	}
	resp.Code = code
	resp.Msg = consts.GetErrorCodeMsg(code)
	resp.Data = &social.SocialData{
		Items: userInfos,
		Total: int64(len(userInfos)),
	}
	return resp, nil
}
