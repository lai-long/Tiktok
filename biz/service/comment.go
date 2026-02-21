package service

import (
	"Tiktok/biz/dao/db"
	"Tiktok/biz/model/dto"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/utils"
	"strconv"
)

func CommentPublish(videoId string, userId string, content string) (int, string) {
	commentId := utils.IdGenerate()
	err := db.CreateComment(commentId, videoId, userId, content)
	if err != nil {
		return consts.CodeDBOperationError, "CommentPublish CreateComment error"
	}
	return consts.CodeSuccess, commentId
}
func CommentList(videoId string, pageSize string, pageNum string) (int, string, []dto.Comment, bool) {
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		return consts.CodeError, "CommentList pageNumInt strconv error", []dto.Comment{}, false
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return consts.CodeError, "CommentList pageSize strconv error", []dto.Comment{}, false
	}
	err, commentEntity := db.GetComments(videoId, pageNumInt, pageSizeInt)
	if err != nil {
		return consts.CodeDBSelectError, "service CommentList GetComments error", []dto.Comment{}, false
	}
	comments := make([]dto.Comment, len(commentEntity))
	for i, _ := range commentEntity {
		comments[i].CommentId = commentEntity[i].CommentId
		comments[i].UserId = commentEntity[i].UserId
		comments[i].Content = commentEntity[i].Content
		comments[i].VideoId = commentEntity[i].VideoId
		comments[i].CreatedAt = commentEntity[i].CreatedAt
	}
	return consts.CodeSuccess, "CommentList success", comments, true
}
func CommentDelete(commentId string, videoId string) (int, string) {
	err := db.CommentDelete(videoId, commentId)
	if err != nil {
		return consts.CodeDBOperationError, "CommentDelete CreateComment error"
	}
	return consts.CodeSuccess, commentId
}
