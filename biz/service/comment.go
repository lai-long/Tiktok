package service

import (
	"Tiktok/biz/model/dto"
	"Tiktok/biz/model/entity"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/utils"
	"log"
	"strconv"
)

type CommentDatabase interface {
	GetComments(videoId string, pageNum int, pageSize int) (bool, []entity.CommentEntity)
	CommentDelete(videoId string, commentId string) bool
	GetCommentById(commentId string) (entity.CommentEntity, bool)
	CommentCountUp(videoId string) bool
	CommentCountDown(videoId string) bool
	CreateComment(commentId string, videoId string, userId string, content string) bool
}

type CommentService struct {
	db CommentDatabase
}

func NewCommentService(db CommentDatabase) *CommentService {
	return &CommentService{db: db}
}
func (s *CommentService) CommentPublish(videoId string, userId string, content string) (int, string) {
	commentId := utils.IdGenerate()
	ok := s.db.CreateComment(commentId, videoId, userId, content)
	if !ok {
		return consts.CodeDBCreateError, "CommentPublish CreateComment error"
	}
	ok = s.db.CommentCountUp(videoId)
	if !ok {
		return consts.CodeDBUpdateError, "CommentPublish CommentCountUp error"
	}
	return consts.CodeSuccess, "CommentPublish success"
}
func (s *CommentService) CommentList(videoId string, pageSize string, pageNum string) (int, string, []dto.Comment, bool) {
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Printf("pageNumInt, err := strconv.Atoi(pageNum) error: %v", err)
		return consts.CodeError, "CommentList pageNumInt strconv error", []dto.Comment{}, false
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Printf("pageSizeInt, err := strconv.Atoi(pageSize) error: %v", err)
		return consts.CodeError, "CommentList pageSize strconv error", []dto.Comment{}, false
	}
	ok, commentEntity := s.db.GetComments(videoId, pageNumInt, pageSizeInt)
	if !ok {
		return consts.CodeDBSelectError, "service CommentList GetComments error", []dto.Comment{}, false
	}
	comments := make([]dto.Comment, len(commentEntity))
	for i := range commentEntity {
		comments[i].CommentId = commentEntity[i].CommentId
		comments[i].UserId = commentEntity[i].UserId
		comments[i].Content = commentEntity[i].Content
		comments[i].VideoId = commentEntity[i].VideoId
		comments[i].CreatedAt = commentEntity[i].CreatedAt
	}
	return consts.CodeSuccess, "CommentList success", comments, true
}
func (s *CommentService) CommentDelete(commentId string, videoId string, userId string) (int, string) {
	comment, ok := s.db.GetCommentById(commentId)
	if !ok {
		return consts.CodeDBSelectError, "CommentDelete GetCommentById error"
	}
	if comment.UserId != userId {
		return consts.CodeError, "CommentDelete GetUserId error,comment_userId != userId"
	}
	ok = s.db.CommentDelete(videoId, commentId)
	if !ok {
		return consts.CodeDBDeleteError, "CommentDelete CreateComment error"
	}
	ok = s.db.CommentCountDown(videoId)
	if !ok {
		return consts.CodeDBDeleteError, "CommentDelete CommentCountDown error"
	}
	return consts.CodeSuccess, commentId
}
