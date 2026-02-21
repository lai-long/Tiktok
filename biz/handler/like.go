package handler

import (
	"Tiktok/biz/model/dto"
	"Tiktok/biz/service"
	"Tiktok/pkg/consts"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func LikeAction(ctx context.Context, c *app.RequestContext) {
	action := c.PostForm("action_type")
	id := c.PostForm("video_id")
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: consts.CodeLikeError,
			Msg:  "likeAction Get userId error",
		}})
		return
	}
	code, msg := service.LikeAction(userId.(string), id, action)
	c.JSON(200, dto.Response{Base: dto.Base{
		Code: code,
		Msg:  msg,
	}})
}
func LikeList(ctx context.Context, c *app.RequestContext) {
	pageSize := c.Query("page_size")
	pageNum := c.Query("page_num")
	userId := c.Query("user_id")
	code, msg, videos, ok := service.LikeList(userId, pageSize, pageNum)
	if !ok {
		c.JSON(200, dto.Response{Base: dto.Base{
			Code: code,
			Msg:  msg,
		}})
		return
	}

	c.JSON(200, dto.Response{
		Base: dto.Base{
			Code: code,
			Msg:  msg,
		},
		Data: dto.Items{
			Video: videos,
			Total: len(videos),
		},
	})
}
