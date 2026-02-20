package service

import (
	"Tiktok/biz/dao/db"
	"Tiktok/biz/model/dto"
	"Tiktok/biz/model/entity"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/utils"
	"io"
	"mime/multipart"
	"os"
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
	err = db.CreatVideo(videoEntity)
	if err != nil {
		return consts.CodeDBOperationError, "VideoPublish db.Create err"
	}
	return consts.CodeSuccess, "success"
}

func VideoList(userId string, pageSize string, pageNum string) (int, string, []entity.VideoEntity) {
	videoList, err := db.GetVideoByUserID(userId, pageSize, pageNum)
	if err != nil {
		return consts.CodeDBSelectError, "VideoList GetVideoByUserID error", []entity.VideoEntity{}
	}
	return consts.CodeSuccess, "success", videoList
}
