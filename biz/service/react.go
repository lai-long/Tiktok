package service

import (
	"Tiktok/biz/dao/db"
	"Tiktok/biz/model/dto"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/utils"
	"log"
	"strconv"
)

func LikeAction(userId string, videoId string, action string) (int, string) {
	if action == "1" {
		err := db.LikeCreate(userId, videoId)
		if err != nil {
			return consts.CodeDBCreateError, "LikeAction LikeCreate error"
		}
		err = db.LikeCountUp(videoId)
		if err != nil {
			return consts.CodeDBUpdateError, "LikeAction LikeCountUp error"
		}
		return consts.CodeSuccess, "LikeAction success"
	}
	if action == "2" {
		err := db.LikeDelete(userId, videoId)
		if err != nil {
			return consts.CodeDBDeleteError, "LikeAction LikeDelete error"
		}
		err = db.LikeCountDown(videoId)
		if err != nil {
			return consts.CodeDBUpdateError, "LikeAction LikeCountDown error"
		}
		return consts.CodeSuccess, "LikeAction success"
	}
	return consts.CodeLikeError, "LikeAction action num error"
}
func LikeList(userId string, pageNum string, pageSize string) (int, string, []dto.Video, bool) {
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Printf("LikeList pageNum strconv error : %v", err)
		return consts.CodeError, "LikeList pageNum strconv error", []dto.Video{}, false
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Printf("LikeList pageSize strconv error : %v", err)
		return consts.CodeError, "LikeList pageSize strconv error", []dto.Video{}, false
	}
	err, videoId := db.LikeVideoIds(userId, pageNumInt, pageSizeInt)
	if err != nil {
		log.Printf("LikeList err : %v", err)
		return consts.CodeDBSelectError, "LikeList db.LikeVideoIds error", []dto.Video{}, false
	}
	ok, videos := db.LikeVideos(videoId)
	if !ok {
		return consts.CodeDBSelectError, "LikeList db.LikeVideos error", []dto.Video{}, false
	}
	videoDTOs := make([]dto.Video, len(videos))
	for i, video := range videos {
		videoDTOs[i] = dto.Video{
			ID:           video.ID,
			UserID:       video.UserID,
			Title:        video.Title,
			Description:  video.Description,
			CommentCount: int64(video.CommentCount),
			CoverURL:     video.CoverURL,
			CreatedAt:    video.CreatedAt,
			LikeCount:    int64(video.LikeCount),
			UpdatedAt:    video.UpdatedAt,
			VideoURL:     video.VideoURL,
			VisitCount:   int64(video.VisitCount),
		}
	}
	return consts.CodeSuccess, "LikeList success", videoDTOs, true
}

func CommentPublish(videoId string, userId string, content string) (int, string) {
	commentId := utils.IdGenerate()
	err := db.CreateComment(commentId, videoId, userId, content)
	if err != nil {
		log.Printf("CommentPublish err : %v", err)
		return consts.CodeDBCreateError, "CommentPublish CreateComment error"
	}
	err = db.CommentCountUp(videoId)
	if err != nil {
		log.Printf("CommentPublish err : %v", err)
		return consts.CodeDBUpdateError, "CommentPublish CommentCountUp error"
	}
	return consts.CodeSuccess, "CommentPublish success"
}
func CommentList(videoId string, pageSize string, pageNum string) (int, string, []dto.Comment, bool) {
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
	err, commentEntity := db.GetComments(videoId, pageNumInt, pageSizeInt)
	if err != nil {
		log.Printf("GetComments err: ", err)
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
func CommentDelete(commentId string, videoId string, userId string) (int, string) {
	comment, err := db.GetCommentById(commentId)
	if err != nil {
		log.Printf("CommentDelete err : %v", err)
		return consts.CodeDBSelectError, "CommentDelete GetCommentById error"
	}
	if comment.UserId != userId {
		return consts.CodeError, "CommentDelete GetUserId error,comment_userId != userId"
	}
	err = db.CommentDelete(videoId, commentId)
	if err != nil {
		log.Printf("CommentDelete err : %v", err)
		return consts.CodeDBDeleteError, "CommentDelete CreateComment error"
	}
	err = db.CommentCountDown(videoId)
	return consts.CodeSuccess, commentId
}
