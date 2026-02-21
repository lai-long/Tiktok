package service

import (
	"Tiktok/biz/dao/db"
	"Tiktok/biz/model/dto"
	"Tiktok/pkg/consts"
	"strconv"
)

func LikeAction(user_id string, video_id string, action string) (int, string) {
	if action == "1" {
		err := db.LikeCountUp(user_id, video_id)
		if err != nil {
			return consts.CodeDBOperationError, "LikeAction LikeCountUp error"
		}
		err = db.LikeCreate(user_id, video_id)
		if err != nil {
			return consts.CodeDBOperationError, "LikeAction LikeCreate error"
		}
		return consts.CodeSuccess, "LikeAction success"
	}
	if action == "2" {
		err := db.LikeCountDown(user_id, video_id)
		if err != nil {
			return consts.CodeDBOperationError, "LikeAction LikeCountDown error"
		}
		err = db.LikeDelete(user_id, video_id)
		if err != nil {
			return consts.CodeDBOperationError, "LikeAction LikeDelete error"
		}
		return consts.CodeSuccess, "LikeAction success"
	}
	return consts.CodeLikeError, "LikeAction action num error"
}
func LikeList(user_id string, pageNum string, pageSize string) (int, string, []dto.Video, bool) {
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		return consts.CodeError, "LikeList pageNum strconv error", []dto.Video{}, false
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return consts.CodeError, "LikeList pageSize strconv error", []dto.Video{}, false
	}
	err, videoId := db.LikeVideoIds(user_id, pageNumInt, pageSizeInt)
	if err != nil {
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
			CommentCount: video.CommentCount,
			CoverURL:     video.CoverURL,
			CreatedAt:    video.CreatedAt,
			DeletedAt:    video.DeletedAt,
			LikeCount:    video.LikeCount,
			UpdatedAt:    video.UpdatedAt,
			VideoURL:     video.VideoURL,
			VisitCount:   video.VisitCount,
		}
	}
	return consts.CodeSuccess, "LikeList success", videoDTOs, true
}
