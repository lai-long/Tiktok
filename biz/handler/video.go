package handler

import (
	"Tiktok/biz/model/dto"
	"Tiktok/pkg/consts"
	"context"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
)

func (h *Handler) VideoPublish(ctx context.Context, c *app.RequestContext) {
	var video dto.Video
	if err := c.Bind(&video); err != nil {
		log.Printf("c.Bind: %v", err)
		c.JSON(200, dto.Response{Base: dto.Base{Code: consts.CodeVideoError, Msg: "VideoPublish Bind Error"}})
		return
	}
	data, err := c.FormFile("data")
	if err != nil {
		log.Printf("c.FormFile: %v", err)
		c.JSON(200, dto.Response{Base: dto.Base{Code: consts.CodeVideoError, Msg: "VideoPublish FormFile Error"}})
	}
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(200, dto.Response{
			Base: dto.Base{
				Code: consts.CodeVideoError,
				Msg:  "VideoPublish Get User Error",
			},
		})
	}
	video.UserID = userId.(string)
	code, msg := h.service.VideoPublish(video, data, ctx)
	c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}})
}
func (h *Handler) VideoList(ctx context.Context, c *app.RequestContext) {
	userId := c.Query("user_id")
	pageSize := c.Query("page_size")
	pageNum := c.Query("page_num")
	code, msg, video, ok := h.service.VideoList(userId, pageSize, pageNum)
	if !ok {
		c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}})
		return
	}
	c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}, Data: dto.Data{
		Items: video,
		Total: len(video),
	}})
}
func (h *Handler) VideoSearch(ctx context.Context, c *app.RequestContext) {
	keywords := c.PostForm("keywords")
	pageSize := c.PostForm("page_size")
	pageNum := c.PostForm("page_num")
	code, msg, video, ok := h.service.VideoSearch(keywords, pageNum, pageSize)
	if !ok {
		c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}})
		return
	}
	c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}, Data: dto.Data{
		Items: video,
		Total: len(video),
	}})
}
func (h *Handler) VideoPopular(ctx context.Context, c *app.RequestContext) {
	pageNum := c.Query("page_num")
	pageSize := c.Query("page_size")
	code, msg, videos, ok := h.service.VideoPopular(ctx, pageNum, pageSize)
	if !ok {
		c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}})
		return
	}
	c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}, Data: videos})
}
