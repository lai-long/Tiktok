package social

import (
	"Tiktok/internal/social/service"
	social "Tiktok/kitex_gen/social"
	"Tiktok/pkg/consts"
	"context"

	"github.com/alibaba/sentinel-golang/api"
)

// SocialServiceImpl implements the last service interface defined in the IDL.
type SocialServiceImpl struct {
	socialService *service.SocialRepo
}

func NewSocialServiceImpl(socialRepo *service.SocialRepo) *SocialServiceImpl {
	return &SocialServiceImpl{socialService: socialRepo}
}

// RelationAction implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) RelationAction(ctx context.Context, req *social.RelationActionReq) (resp *social.RelationActionResp, err error) {
	// TODO: Your code here...
	resp = &social.RelationActionResp{}
	entry, blockErr := api.Entry("/relation/action")
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		resp.Msg = consts.GetErrorCodeMsg(consts.SentinelBlock)
		return resp, nil
	}
	defer entry.Exit()

	code, err := s.socialService.RelationAction(ctx, req.ToUserId, req.ActionType, req.UserId)
	resp.Code = code
	resp.Msg = consts.GetErrorCodeMsg(code)
	return resp, nil
}

// FollowingList implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) FollowingList(ctx context.Context, req *social.FollowingListReq) (resp *social.FollowingListResp, err error) {
	// TODO: Your code here...
	resp = &social.FollowingListResp{}
	entry, blockErr := api.Entry("/following/list")
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		resp.Msg = consts.GetErrorCodeMsg(consts.SentinelBlock)
		return resp, nil
	}
	defer entry.Exit()

	code, userInfos, err := s.socialService.FollowingList(req.UserId, req.PageNum, req.PageSize)
	resp.Code = code
	resp.Msg = consts.GetErrorCodeMsg(code)
	resp.Data = &social.SocialData{
		Items: userInfos,
		Total: int64(len(userInfos)),
	}
	return resp, nil
}

// FollowerList implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) FollowerList(ctx context.Context, req *social.FollowerListReq) (resp *social.FollowerListResp, err error) {
	// TODO: Your code here...
	resp = &social.FollowerListResp{}
	entry, blockErr := api.Entry("/follower/list")
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		resp.Msg = consts.GetErrorCodeMsg(consts.SentinelBlock)
		return resp, nil
	}
	defer entry.Exit()

	code, userInfos, err := s.socialService.FollowerList(req.UserId, req.PageNum, req.PageSize)
	resp.Code = code
	resp.Msg = consts.GetErrorCodeMsg(code)
	resp.Data = &social.SocialData{
		Items: userInfos,
		Total: int64(len(userInfos)),
	}
	return resp, nil
}

// FriendList implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) FriendList(ctx context.Context, req *social.FriendListReq) (resp *social.FriendListResp, err error) {
	// TODO: Your code here...
	resp = &social.FriendListResp{}
	entry, blockErr := api.Entry("/friend/list")
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		resp.Msg = consts.GetErrorCodeMsg(consts.SentinelBlock)
		return resp, nil
	}
	defer entry.Exit()

	code, userInfos, err := s.socialService.FriendList(req.UserId, req.PageNum, req.PageSize)
	resp.Code = code
	resp.Msg = consts.GetErrorCodeMsg(code)
	resp.Data = &social.SocialData{
		Items: userInfos,
		Total: int64(len(userInfos)),
	}
	return resp, nil
}
