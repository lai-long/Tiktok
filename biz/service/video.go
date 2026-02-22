package service

import (
	"Tiktok/biz/dao/db"
	"Tiktok/biz/model/dto"
	"Tiktok/biz/model/entity"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/utils"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"strconv"
)

func VideoPublish(video dto.Video, data *multipart.FileHeader) (int, string) {
	dataFile, err := data.Open()
	if err != nil {
		return consts.CodeIOError, "VideoPublish data.Open err"
	}
	defer dataFile.Close()
	file, err := os.Create("/home/lai/project/video" + data.Filename)
	if err != nil {
		return consts.CodeIOError, "VideoPublish os.Create err"
	}
	defer file.Close()
	if _, err := io.Copy(file, dataFile); err != nil {
		return consts.CodeIOError, "VideoPublish io.copy err"
	}
	var videoEntity entity.VideoEntity
	videoEntity.Title = video.Title
	videoEntity.Description = video.Description
	videoEntity.VideoURL = video.VideoURL
	videoEntity.UserID = video.UserID
	videoEntity.ID = utils.IdGenerate()
	videoEntity.VisitCount = rand.Intn(100)
	err = db.CreatVideo(videoEntity)
	if err != nil {
		return consts.CodeDBOperationError, "VideoPublish db.Create err"
	}
	return consts.CodeSuccess, "success"
}

func VideoList(userId string, pageSize string, pageNum string) (int, string, []dto.Video, bool) {
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return consts.CodeError, "VideoList pageSize strconv error", []dto.Video{}, false
	}
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		return consts.CodeError, "VideoList pageNum error", []dto.Video{}, false
	}
	videoList, err := db.GetVideoByUserID(userId, pageSizeInt, pageNumInt)
	if err != nil {
		return consts.CodeDBSelectError, "VideoList GetVideoByUserID error", []dto.Video{}, false
	}
	videoDTOs := make([]dto.Video, len(videoList))
	for i := 0; i < len(videoList); i++ {
		videoDTOs[i] = dto.Video{
			ID:           videoList[i].ID,
			UserID:       videoList[i].UserID,
			Title:        videoList[i].Title,
			Description:  videoList[i].Description,
			CommentCount: int64(videoList[i].CommentCount),
			CoverURL:     videoList[i].CoverURL,
			CreatedAt:    videoList[i].CreatedAt,
			DeletedAt:    videoList[i].DeletedAt,
			LikeCount:    int64(videoList[i].LikeCount),
			UpdatedAt:    videoList[i].UpdatedAt,
			VideoURL:     videoList[i].VideoURL,
			VisitCount:   int64(videoList[i].VisitCount),
		}
	}
	return consts.CodeSuccess, "success", videoDTOs, true
}

func VideoSearch(keyword string, pageNum string, pageSize string) (int, string, []dto.Video, bool) {
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return consts.CodeError, "VideoSearch pageSize strconv error", []dto.Video{}, false
	}
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		return consts.CodeError, "VideoSearch pageNum error", []dto.Video{}, false
	}
	video, err := db.GetVideoByKeyWord(keyword, pageNumInt, pageSizeInt)
	if err != nil {
		return consts.CodeVideoError, "GetVideoByVideoTitleOrDescription error", []dto.Video{}, false
	}
	videoDTOs := make([]dto.Video, len(video))
	for i := 0; i < len(video); i++ {
		videoDTOs[i].ID = video[i].ID
		videoDTOs[i].Title = video[i].Title
		videoDTOs[i].Description = video[i].Description
		videoDTOs[i].VideoURL = video[i].VideoURL
		videoDTOs[i].CreatedAt = video[i].CreatedAt
		videoDTOs[i].DeletedAt = video[i].DeletedAt
		videoDTOs[i].LikeCount = int64(video[i].LikeCount)
		videoDTOs[i].UpdatedAt = video[i].UpdatedAt
		videoDTOs[i].VideoURL = video[i].VideoURL
		videoDTOs[i].CoverURL = video[i].CoverURL
		videoDTOs[i].CommentCount = int64(video[i].CommentCount)
		videoDTOs[i].CreatedAt = video[i].CreatedAt
		videoDTOs[i].DeletedAt = video[i].DeletedAt
	}
	return consts.CodeSuccess, "success", videoDTOs, true
}
