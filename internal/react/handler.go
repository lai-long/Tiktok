package handler

import (
	"Tiktok/internal/react/service"
	react "Tiktok/kitex_gen/react"
	"context"
)

type CommentServiceImpl struct {
	commentRepo *service.CommentRepo
}

func NewCommentService(commentRepo *service.CommentRepo) *CommentServiceImpl {
	return &CommentServiceImpl{commentRepo: commentRepo}
}

func (s *CommentServiceImpl) CommentPublish(ctx context.Context, req *react.CommentPublishReq) (resp *react.CommentPublishResp, err error) {
	resp = &react.CommentPublishResp{}
	code, _ := s.commentRepo.CommentPublish(req.TargetAt, req.UserID, req.Content, req.TargetType)
	resp.Code = code
	return resp, nil
}

func (s *CommentServiceImpl) CommentList(ctx context.Context, req *react.CommentListReq) (resp *react.CommentListResp, err error) {
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
	resp = &react.CommentDeleteResp{}
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
	resp = &react.LikeActionResp{}
	code, _ := s.likeRepo.LikeAction(ctx, req.UserID, req.TargetAt, req.ActionType, req.TargetType)
	resp.Code = code
	return resp, nil
}

func (s *LikeServiceImpl) LikeList(ctx context.Context, req *react.LikeListReq) (resp *react.LikeListResp, err error) {
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
