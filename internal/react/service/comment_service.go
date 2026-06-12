package service

import (
	"Tiktok/kitex_gen/react"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/entity"
	"Tiktok/pkg/utils"
	"context"

	"github.com/pkg/errors"
)

type CommentDatabase interface {
	GetComments(ctx context.Context, targetId string, pageNum int64, pageSize int64) ([]entity.CommentEntity, error)
	CommentDelete(ctx context.Context, commentId string) error
	GetCommentById(ctx context.Context, commentId string) (entity.CommentEntity, error)
	VideoCommentCountUp(ctx context.Context, videoId string) error
	CommentCommentCountUp(ctx context.Context, commentId string) error
	VideoCommentCountDown(ctx context.Context, videoId string) error
	CommentCommentCountDown(ctx context.Context, commentId string) error
	CreateComment(ctx context.Context, commentId string, videoId string, userId string, content string, targetType string) error
}

type CommentRepo struct {
	db CommentDatabase
}

func NewCommentService(db CommentDatabase) *CommentRepo {
	return &CommentRepo{db: db}
}

func toCommentInfo(e entity.CommentEntity) *react.CommentInfo {
	return &react.CommentInfo{
		UserId:       e.UserID,
		TargetId:     e.TargetID,
		CommentId:    e.CommentID,
		Content:      e.Content,
		LikeCount:    e.LikeCount,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
		TargetType:   e.TargetType,
		CommentCount: e.CommentCount,
	}
}

func (s *CommentRepo) CommentPublish(ctx context.Context, targetId, userId, content, targetType string) (int32, error) {
	defer utils.TrackTime(ctx, "CommentPublish")()
	switch targetType {
	case "1":
		commentId := utils.IDGenerate()
		err := s.db.CreateComment(ctx, commentId, targetId, userId, content, targetType)
		if err != nil {
			return consts.ReactDBInsertError, errors.Wrap(err, "->CommentPublish Create comment error ")
		}
		err = s.db.VideoCommentCountUp(ctx, targetId)
		if err != nil {
			return consts.ReactDBUpdateError, errors.Wrap(err, "->CommentPublish Update comment count error ")
		}
		return consts.Success, nil
	case "2":
		commentId := utils.IDGenerate()
		err := s.db.CreateComment(ctx, commentId, targetId, userId, content, targetType)
		if err != nil {
			return consts.ReactDBInsertError, errors.Wrap(err, "->CommentPublish Create comment error ")
		}
		err = s.db.CommentCommentCountUp(ctx, targetId)
		if err != nil {
			return consts.ReactDBUpdateError, errors.Wrap(err, "->CommentPublish update comment count error ")
		}
		return consts.Success, nil
	}
	return consts.ReactReqValueError, nil
}

func (s *CommentRepo) CommentList(ctx context.Context, targetId string, pageSize int64, pageNum int64) (int32, []*react.CommentInfo, error) {
	defer utils.TrackTime(ctx, "CommentList")()
	commentEntity, err := s.db.GetComments(ctx, targetId, pageNum, pageSize)
	if err != nil {
		return consts.ReactDBSelectError, nil, errors.Wrap(err, "->CommentList select comment err")
	}
	var comments []*react.CommentInfo
	for i := range commentEntity {
		comments = append(comments, toCommentInfo(commentEntity[i]))
	}
	return consts.Success, comments, nil
}

func (s *CommentRepo) CommentDelete(ctx context.Context, commentId string, targetId string, userId string, targetType string) (int32, error) {
	defer utils.TrackTime(ctx, "CommentDelete")()
	comment, err := s.db.GetCommentById(ctx, commentId)
	if err != nil {
		return consts.ReactDBSelectError, errors.Wrap(err, "->CommentDelete select comment err")
	}
	if comment.UserID != userId {
		return consts.ReactReqValueError, nil
	}
	err = s.db.CommentDelete(ctx, commentId)
	if err != nil {
		return consts.ReactDBDeleteError, errors.Wrap(err, "->CommentDelete delete comment err")
	}
	switch targetType {
	case "1":
		err = s.db.VideoCommentCountDown(ctx, targetId)
		if err != nil {
			return consts.ReactDBUpdateError, errors.Wrap(err, "->CommentDelete update comment count error ")
		}
		return consts.Success, nil
	case "2":
		err = s.db.CommentCommentCountDown(ctx, targetId)
		if err != nil {
			return consts.ReactDBUpdateError, errors.Wrap(err, "->CommentDelete update comment count error ")
		}
	}
	return consts.ReactReqValueError, nil
}
