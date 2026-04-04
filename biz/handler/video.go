package handler

import (
	"Tiktok/biz/model/common"
	"Tiktok/biz/model/dto"
	"Tiktok/biz/model/video"
	"Tiktok/pkg/consts"
	"context"
	"log"
	"mime/multipart"

	"github.com/cloudwego/hertz/pkg/app"
)

type VideoSever interface {
	VideoPublish(video *video.VideoInfo, data *multipart.FileHeader, ctx context.Context) (int, string)
	VideoList(userId string, pageSize int64, pageNum int64) (int, string, []*video.VideoInfo, bool)
	VideoSearch(keyword string, pageNun, pageSize int64) (int, string, []*video.VideoInfo, bool)
	VideoPopular(ctx context.Context, pageNum int64, pageSize int64) (int, string, []*video.VideoInfo, bool)
	VideoStream() (int, string, []*video.VideoInfo)
}
type VideoHandler struct {
	videoService VideoSever
}

func NewVideoHandler(videoService VideoSever) *VideoHandler {
	return &VideoHandler{videoService: videoService}
}

func (h *VideoHandler) VideoPublish(ctx context.Context, c *app.RequestContext) {
	req := new(video.VideoPublishReq)
	if err := c.BindAndValidate(req); err != nil {
		log.Printf("c.Bind: %v", err)
		c.JSON(200, video.VideoPublishResp{
			Base: &common.Base{
				Code: consts.CodeVideoError,
				Msg:  "VideoPublish BindAndValidate error",
			},
		})
		return
	}
	data, err := c.FormFile("data")
	if err != nil {
		log.Printf("c.FormFile: %v", err)
		c.JSON(200, video.VideoPublishResp{
			Base: &common.Base{
				Code: consts.CodeVideoError,
				Msg:  "VideoPublish FormFile error",
			},
		})
	}
	userId, ok := c.Value("user_id").(string)
	if !ok {
		c.JSON(200, video.VideoPublishResp{
			Base: &common.Base{
				Code: consts.CodeVideoError,
				Msg:  "VideoPublish Get User Error",
			},
		})
	}
	videoInfo := &video.VideoInfo{
		UserID:      userId,
		Title:       req.Title,
		Description: req.Description,
	}
	code, msg := h.videoService.VideoPublish(videoInfo, data, ctx)
	c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}})
}

func (h *VideoHandler) VideoList(ctx context.Context, c *app.RequestContext) {
	req := new(video.VideoListReq)
	if err := c.BindAndValidate(req); err != nil {
		log.Printf("c.Bind err: %v", err)
		c.JSON(200, video.VideoListResp{
			Base: &common.Base{
				Code: consts.CodeVideoError,
				Msg:  "VideoList BindAndValidate error",
			},
		})
	}
	code, msg, videoInfos, _ := h.videoService.VideoList(req.UserId, req.PageSize, req.PageNum)
	c.JSON(200, video.VideoListResp{
		Base: &common.Base{
			Code: int32(code),
			Msg:  msg,
		},
		Data: &video.VideoData{
			Items: videoInfos,
			Total: int64(len(videoInfos)),
		},
	})
}

func (h *VideoHandler) VideoSearch(ctx context.Context, c *app.RequestContext) {
	req := new(video.VideoSearchReq)
	if err := c.BindAndValidate(req); err != nil {
		c.JSON(200, video.VideoSearchResp{
			Base: &common.Base{
				Code: consts.CodeVideoError,
				Msg:  "VideoSearch BindAndValidate error",
			},
		})
	}
	code, msg, videoInfos, _ := h.videoService.VideoSearch(req.KeyWord, req.PageNum, req.PageSize)
	c.JSON(200, video.VideoSearchResp{
		Base: &common.Base{
			Code: int32(code),
			Msg:  msg,
		},
		Data: &video.VideoData{
			Items: videoInfos,
			Total: int64(len(videoInfos)),
		},
	})
}

func (h *VideoHandler) VideoPopular(ctx context.Context, c *app.RequestContext) {
	req := new(video.VideoHotReq)
	if err := c.BindAndValidate(req); err != nil {
		c.JSON(200, video.VideoHotResp{
			Base: &common.Base{
				Code: consts.CodeVideoError,
				Msg:  "VideoPopular BindAndValidate error",
			},
		})
	}
	code, msg, videoInfos, _ := h.videoService.VideoPopular(ctx, req.PageNum, req.PageSize)
	c.JSON(200, video.VideoHotResp{
		Base: &common.Base{
			Code: int32(code),
			Msg:  msg,
		},
		Data: &video.VideoData{
			Items: videoInfos,
			Total: int64(len(videoInfos)),
		},
	})
}

func (h *VideoHandler) VideoStream(ctx context.Context, c *app.RequestContext) {
	code, msg, videoInfos := h.videoService.VideoStream()
	c.JSON(200, dto.Response{Base: dto.Base{Code: code, Msg: msg}, Data: dto.Data{Items: videoInfos, Total: len(videoInfos)}})
}
