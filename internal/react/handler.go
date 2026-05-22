package handler

import (
	"Tiktok/internal/react/service"
	react "Tiktok/kitex_gen/react"
	"Tiktok/pkg/consts"
	"context"

	"github.com/alibaba/sentinel-golang/api"
)

type CommentServiceImpl struct {
	commentRepo *service.CommentRepo
}

func NewCommentService(commentRepo *service.CommentRepo) *CommentServiceImpl {
	return &CommentServiceImpl{commentRepo: commentRepo}
}

func (s *CommentServiceImpl) CommentPublish(ctx context.Context, req *react.CommentPublishReq) (resp *react.CommentPublishResp, err error) {
	// TODO: Your code here...
	resp = &react.CommentPublishResp{}
	entry, blockErr := api.Entry("/comment/publish")
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		return resp, nil
	}
	defer entry.Exit()

	code, _ := s.commentRepo.CommentPublish(req.TargetAt, req.UserID, req.Content, req.TargetType)
	resp.Code = code
	return resp, nil
}

func (s *CommentServiceImpl) CommentList(ctx context.Context, req *react.CommentListReq) (resp *react.CommentListResp, err error) {
	// TODO: Your code here...
	resp = &react.CommentListResp{}
	entry, blockErr := api.Entry("/comment/list")
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		return resp, nil
	}
	defer entry.Exit()

	code, comments, _ := s.commentRepo.CommentList(req.TargetAt, req.PageSize, req.PageNum)
	resp = &react.CommentListResp{
		Code: code,
		Data: &react.CommentData{
			Items: comments,
			Total: int64(len(comments)),
		},
	}
	return resp, nil
}

func (s *CommentServiceImpl) CommentDelete(ctx context.Context, req *react.CommentDeleteReq) (resp *react.CommentDeleteResp, err error) {
	// TODO: Your code here...
	resp = &react.CommentDeleteResp{}
	entry, blockErr := api.Entry("/comment/delete")
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		return resp, nil
	}
	defer entry.Exit()

	code, _ := s.commentRepo.CommentDelete(req.CommentId, req.TargetAt, req.UserID, req.TargetType)
	resp.Code = code
	return resp, nil
}

type LikeServiceImpl struct {
	likeRepo *service.LikeRepo
}

func NewLikeService(likeRepo *service.LikeRepo) *LikeServiceImpl {
	return &LikeServiceImpl{likeRepo: likeRepo}
}

func (s *LikeServiceImpl) LikeAction(ctx context.Context, req *react.LikeActionReq) (resp *react.LikeActionResp, err error) {
	// TODO: Your code here...
	resp = &react.LikeActionResp{}
	entry, blockErr := api.Entry("/like/action")
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		return resp, nil
	}
	defer entry.Exit()

	code, _ := s.likeRepo.LikeAction(ctx, req.UserID, req.TargetAt, req.ActionType, req.TargetType)
	resp.Code = code
	return resp, nil
}

func (s *LikeServiceImpl) LikeList(ctx context.Context, req *react.LikeListReq) (resp *react.LikeListResp, err error) {
	// TODO: Your code here...
	resp = &react.LikeListResp{}
	entry, blockErr := api.Entry("/like/list")
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		return resp, nil
	}
	defer entry.Exit()

	code, videoIds, total, _ := s.likeRepo.LikeList(ctx, req.UserId, req.PageNum, req.PageSize)
	resp = &react.LikeListResp{
		Code: code,
		Data: &react.LikeVideoData{
			VideoIds: videoIds,
			Total:    total,
		},
	}
	return resp, nil
}
