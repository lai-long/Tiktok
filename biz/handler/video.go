package handler

import (
	"Tiktok/biz/model/dto"
	"Tiktok/biz/service"
	"Tiktok/pkg/consts"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func VideoPublish(ctx context.Context, c *app.RequestContext) {
	var video dto.Video
	if err := c.Bind(&video); err != nil {
		c.JSON(200, dto.Response{Base: dto.Base{Code: consts.CodeVideoError, Msg: "VideoPublish Bind Error"}})
		return
	}
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(200, dto.Response{Base: dto.Base{Code: consts.CodeVideoError, Msg: "VideoPublish FormFile Error"}})
	}
	code, msg := service.VideoPublish(video, data)
	c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}})
}
func VideoList(ctx context.Context, c *app.RequestContext) {
	userId := c.Query("user_id")
	pageSize := c.Query("page_size")
	pageNum := c.Query("page_num")
	code, msg, video := service.VideoList(userId, pageSize, pageNum)
	c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}, Data: dto.Items{Video: video, Total: 0}})
}
func VideoSearch(ctx context.Context, c *app.RequestContext) {
	title := c.Query("title")
	description := c.Query("description")
	code, msg, video := service.VideoSearch(title, description)
	c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}, Data: dto.Items{Video: video}})
}
