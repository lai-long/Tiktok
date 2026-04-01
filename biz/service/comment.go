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
	GetComments(videoId string, pageNum int, pageSize int) (error, []entity.CommentEntity)
	CommentDelete(commentId string) error
	GetCommentById(commentId string) (entity.CommentEntity, error)
	VideoCommentCountUp(videoId string) error
	CommentCommentCountUp(commentId string) error
	CommentCountDown(videoId string) error
	CreateComment(commentId string, videoId string, userId string, content string, targetType string) error
}

type CommentService struct {
	db CommentDatabase
}

func NewCommentService(db CommentDatabase) *CommentService {
	return &CommentService{db: db}
}

func (s *CommentService) CommentPublish(targetId, userId, content, targetType string) (int, string) {
	if targetType == "1" {
		commentId := utils.IdGenerate()
		err := s.db.CreateComment(commentId, targetId, userId, content, targetType)
		if err != nil {
			log.Println("Video db CreateComment err", err)
			return consts.CodeDBCreateError, "Video CommentPublish CreateComment error"
		}
		err = s.db.VideoCommentCountUp(targetId)
		if err != nil {
			log.Println("db CommentCountUp err", err)
			return consts.CodeDBUpdateError, "Video CommentPublish CommentCountUp error"
		}
		return consts.CodeSuccess, "Video CommentPublish success"
	} else if targetType == "2" {
		commentId := utils.IdGenerate()
		err := s.db.CreateComment(commentId, targetId, userId, content, targetType)
		if err != nil {
			log.Println("Comment db CreateComment err", err)
			return consts.CodeDBCreateError, "Comment CommentPublish CreateComment error"
		}
		err = s.db.CommentCommentCountUp(targetId)
	}
	return consts.CodeDBCreateError, "CommentPublish CreateComment error"
}

func (s *CommentService) CommentList(targetId string, pageSize string, pageNum string) (int, string, []dto.Comment, bool) {
	pageNumInt := 0
	pageSizeInt := 10
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Printf("pageNumInt, err := strconv.Atoi(pageNum) error: %v", err)
		return consts.CodeError, "CommentList pageNumInt strconv error", []dto.Comment{}, false
	}
	pageSizeInt, err = strconv.Atoi(pageSize)
	if err != nil {
		log.Printf("pageSizeInt, err := strconv.Atoi(pageSize) error: %v", err)
		return consts.CodeError, "CommentList pageSize strconv error", []dto.Comment{}, false
	}
	err, commentEntity := s.db.GetComments(targetId, pageNumInt, pageSizeInt)
	if err != nil {
		log.Printf("db GetComments err: %v", err)
		return consts.CodeDBSelectError, "service CommentList GetComments error", []dto.Comment{}, false
	}
	comments := make([]dto.Comment, len(commentEntity))
	for i := range commentEntity {
		comments[i].CommentId = commentEntity[i].CommentId
		comments[i].UserId = commentEntity[i].UserId
		comments[i].Content = commentEntity[i].Content
		comments[i].TargetId = commentEntity[i].TargetId
		comments[i].CreatedAt = commentEntity[i].CreatedAt
	}
	return consts.CodeSuccess, "CommentList success", comments, true
}

func (s *CommentService) CommentDelete(commentId string, videoId string, userId string) (int, string) {
	comment, err := s.db.GetCommentById(commentId)
	if err != nil {
		log.Println("db GetComment err", err)
		return consts.CodeDBSelectError, "CommentDelete GetCommentById error"
	}
	if comment.UserId != userId {
		return consts.CodeError, "CommentDelete GetUserId error,comment_userId != userId"
	}
	err = s.db.CommentDelete(commentId)
	if err != nil {
		log.Println("db CommentDelete err", err)
		return consts.CodeDBDeleteError, "CommentDelete CreateComment error"
	}
	err = s.db.CommentCountDown(videoId)
	if err != nil {
		log.Println("db CommentCountDown err", err)
		return consts.CodeDBDeleteError, "CommentDelete CommentCountDown error"
	}
	return consts.CodeSuccess, commentId
}
