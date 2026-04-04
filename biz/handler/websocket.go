package handler

import (
	"Tiktok/biz/cache"
	"Tiktok/biz/dao"
	"Tiktok/biz/model/chat"
	"Tiktok/biz/model/common"
	"Tiktok/biz/service"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/utils"
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	"github.com/gorilla/websocket"
)

type WebsocketSever struct {
	db *dao.MySQLdb
	re *cache.Redis
}

func NewWebsocketSever(db *dao.MySQLdb, re *cache.Redis) *WebsocketSever {
	return &WebsocketSever{
		db: db,
		re: re,
	}
}
func (m *WebsocketSever) WebSocketHandler(ctx context.Context, c *app.RequestContext) {
	userid, ok := ctx.Value("user_id").(string)
	if !ok {
		c.JSON(200, chat.WebsocketResp{Base: &common.Base{
			Code: consts.CodeError,
			Msg:  "WebSocketHandler userid not found",
		}})
		return
	}
	uid := userid
	req := new(chat.WebsocketReq)
	err := c.BindAndValidate(req)
	if err != nil {
		c.JSON(200, chat.WebsocketResp{Base: &common.Base{
			Code: consts.CodeError,
			Msg:  "c.BindAndValidate(req) error",
		}})
	}

	stdHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}).Upgrade(w, r, nil)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		client := &service.Client{
			ID:      utils.CreateId(uid, req.ToUserId),
			SendID:  utils.CreateId(req.ToUserId, uid),
			GroupId: req.GroupId,
			Socket:  conn,
			Send:    make(chan []byte, 128),
		}
		service.Manager.Register <- client
		go client.Read()
		go client.Write()
	})
	wsAdaptor := adaptor.HertzHandler(stdHandler)
	wsAdaptor(ctx, c)
}
