package handler

import (
	"Tiktok/biz/cache"
	"Tiktok/biz/dao"
	"Tiktok/biz/model/dto"
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
	userid, exist := c.Get("user_id")
	if !exist {
		c.JSON(200, dto.Response{
			Base: dto.Base{
				Code: consts.CodeError,
				Msg:  "WebSocketHandler userid not found",
			},
		})
		return
	}
	uid := userid.(string)
	toUserId := c.Query("to_userid")
	groupId := c.Query("group_id")
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
			ID:      utils.CreateId(uid, toUserId),
			SendID:  utils.CreateId(toUserId, uid),
			GroupId: groupId,
			Socket:  conn,
			Send:    make(chan []byte, 128),
		}
		service.Manager.Register <- client
		go client.Read(m.re, m.db)
		go client.Write()
	})
	wsAdaptor := adaptor.HertzHandler(stdHandler)
	wsAdaptor(ctx, c)
}
