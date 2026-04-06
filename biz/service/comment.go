package service

import (
	"Tiktok/biz/entity"
	"Tiktok/biz/model/react"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/utils"

	"github.com/pkg/errors"
)

type CommentDatabase interface {
	GetComments(videoId string, pageNum int64, pageSize int64) (error, []entity.CommentEntity)
	CommentDelete(commentId string) error
	GetCommentById(commentId string) (entity.CommentEntity, error)
	VideoCommentCountUp(videoId string) error
	CommentCommentCountUp(commentId string) error
	VideoCommentCountDown(videoId string) error
	CommentCommentCountDown(commentId string) error
	CreateComment(commentId string, videoId string, userId string, content string, targetType string) error
}

type CommentService struct {
	db CommentDatabase
}

func NewCommentService(db CommentDatabase) *CommentService {
	return &CommentService{db: db}
}

func (s *CommentService) CommentPublish(targetId, userId, content, targetType string) (int32, error) {
	if targetType == "1" {
		commentId := utils.IdGenerate()
		err := s.db.CreateComment(commentId, targetId, userId, content, targetType)
		if err != nil {
			return consts.ReactDBInsertError, errors.Wrap(err, "->CommentPublish Create comment error ")
		}
		err = s.db.VideoCommentCountUp(targetId)
		if err != nil {
			return consts.ReactDBUpdateError, errors.Wrap(err, "->CommentPublish Update comment count error ")
		}
		return consts.Success, nil
	} else if targetType == "2" {
		commentId := utils.IdGenerate()
		err := s.db.CreateComment(commentId, targetId, userId, content, targetType)
		if err != nil {
			return consts.ReactDBInsertError, errors.Wrap(err, "->CommentPublish Create comment error ")
		}
		err = s.db.CommentCommentCountUp(targetId)
		if err != nil {
			return consts.ReactDBUpdateError, errors.Wrap(err, "->CommentPublish update comment count error ")
		}
	}
	return consts.ReactReqValueError, nil
}

func (s *CommentService) CommentList(targetId string, pageSize int64, pageNum int64) (int32, error, []*react.CommentInfo) {

	err, commentEntity := s.db.GetComments(targetId, pageNum, pageSize)
	if err != nil {
		return consts.ReactDBSelectError, errors.Wrap(err, "->CommentList select comment err"), nil
	}
	var comments []*react.CommentInfo
	for i, _ := range commentEntity {
		comments = append(comments, commentEntity[i].ToCommentInfo())
	}
	return consts.Success, nil, comments
}

func (s *CommentService) CommentDelete(commentId string, targetId string, userId string, targetType string) (int32, error) {
	comment, err := s.db.GetCommentById(commentId)
	if err != nil {
		return consts.ReactDBSelectError, errors.Wrap(err, "->CommentDelete select comment err")
	}
	if comment.UserId != userId {
		return consts.ReactReqValueError, nil
	}
	err = s.db.CommentDelete(commentId)
	if err != nil {
		return consts.ReactDBDeleteError, errors.Wrap(err, "->CommentDelete delete comment err")
	}
	if targetType == "1" {
		err = s.db.VideoCommentCountDown(targetId)
		if err != nil {
			return consts.ReactDBUpdateError, errors.Wrap(err, "->CommentDelete update comment count error ")
		}
		return consts.Success, nil
	} else if targetType == "2" {
		err = s.db.CommentCommentCountDown(targetId)
		if err != nil {
			return consts.ReactDBUpdateError, errors.Wrap(err, "->CommentDelete update comment count error ")
		}
		return consts.Success, nil
	}
	return consts.ReactReqValueError, nil
}
