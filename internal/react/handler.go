package handler

import (
	react "Tiktok/kitex_gen/react"
	"context"
)

// commentService 是 CommentServiceImpl 依赖的最小接口。
type commentService interface {
	CommentPublish(ctx context.Context, targetId, userId, content, targetType string) (int32, error)
	CommentList(ctx context.Context, targetId string, pageSize int64, pageNum int64) (int32, []*react.CommentInfo, error)
	CommentDelete(ctx context.Context, commentId string, targetId string, userId string, targetType string) (int32, error)
}

// likeService 是 LikeServiceImpl 依赖的最小接口。
type likeService interface {
	LikeAction(ctx context.Context, userId string, targetId string, action string, targetType string) (int32, error)
	LikeList(ctx context.Context, userId string, pageNum int64, pageSize int64) (int32, []string, int64, error)
}

type CommentServiceImpl struct {
	commentRepo commentService
}

func NewCommentService(commentRepo commentService) *CommentServiceImpl {
	return &CommentServiceImpl{commentRepo: commentRepo}
}

func (s *CommentServiceImpl) CommentPublish(ctx context.Context, req *react.CommentPublishReq) (resp *react.CommentPublishResp, err error) {
	resp = &react.CommentPublishResp{}
	code, err := s.commentRepo.CommentPublish(ctx, req.TargetAt, req.UserID, req.Content, req.TargetType)
	if err != nil {
		resp.Code = code
		return resp, err
	}
	resp.Code = code
	return resp, nil
}

func (s *CommentServiceImpl) CommentList(ctx context.Context, req *react.CommentListReq) (resp *react.CommentListResp, err error) {
	code, comments, err := s.commentRepo.CommentList(ctx, req.TargetAt, req.PageSize, req.PageNum)
	if err != nil {
		resp = &react.CommentListResp{
			Code: code,
			Data: &react.CommentData{
				Items: comments,
				Total: int64(len(comments)),
			},
		}
		return resp, err
	}
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
	code, err := s.commentRepo.CommentDelete(ctx, req.CommentId, req.TargetAt, req.UserID, req.TargetType)
	if err != nil {
		resp.Code = code
		return resp, err
	}
	resp.Code = code
	return resp, nil
}

type LikeServiceImpl struct {
	likeRepo likeService
}

func NewLikeService(likeRepo likeService) *LikeServiceImpl {
	return &LikeServiceImpl{likeRepo: likeRepo}
}

func (s *LikeServiceImpl) LikeAction(ctx context.Context, req *react.LikeActionReq) (resp *react.LikeActionResp, err error) {
	resp = &react.LikeActionResp{}
	code, err := s.likeRepo.LikeAction(ctx, req.UserID, req.TargetAt, req.ActionType, req.TargetType)
	if err != nil {
		resp.Code = code
		return resp, err
	}
	resp.Code = code
	return resp, nil
}

func (s *LikeServiceImpl) LikeList(ctx context.Context, req *react.LikeListReq) (resp *react.LikeListResp, err error) {
	code, videoIds, total, err := s.likeRepo.LikeList(ctx, req.UserId, req.PageNum, req.PageSize)
	if err != nil {
		resp = &react.LikeListResp{
			Code: code,
			Data: &react.LikeVideoData{
				VideoIds: videoIds,
				Total:    total,
			},
		}
		return resp, err
	}
	resp = &react.LikeListResp{
		Code: code,
		Data: &react.LikeVideoData{
			VideoIds: videoIds,
			Total:    total,
		},
	}
	return resp, nil
}
