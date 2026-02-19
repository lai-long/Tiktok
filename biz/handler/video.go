package handler

import (
	"Tiktok/biz/dao/db"
	"Tiktok/biz/model/dto"
	"Tiktok/biz/model/entity"
	"Tiktok/pkg/utils"
	"context"
	"io"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
)

func VideoPublish(ctx context.Context, c *app.RequestContext) {
	var video dto.Video
	if err := c.Bind(&video); err != nil {
	}
	var videoEntity entity.VideoEntity
	videoEntity.Title = video.Title
	videoEntity.Description = video.Description
	videoEntity.VideoURL = video.VideoURL
	videoEntity.UserID = video.UserID
	videoEntity.ID = utils.IdGenerate()
	db.CreatVideo(videoEntity)
	data, _ := c.FormFile("data")
	dataFile, _ := data.Open()
	defer dataFile.Close()
	file, _ := os.Create(data.Filename)
	defer file.Close()
	if _, err := io.Copy(file, dataFile); err != nil {
		c.JSON(200, dto.Response{
			Base: dto.Base{
				Code: 0,
				Msg:  err.Error(),
			},
		})
	}

}
func VideoList(ctx context.Context, c *app.RequestContext) {
	userId := c.Query("user_id")
	pageSize := c.Query("page_size")
	pageNum := c.Query("page_num")
	db.GetVideoByUserID(userId, pageSize, pageNum)
}
func VideoSearch(ctx context.Context, c *app.RequestContext) {
	title := c.Query("title")
	description := c.Query("description")
	db.GetVideoByVideoTitleOrDescription(title, description)
}
