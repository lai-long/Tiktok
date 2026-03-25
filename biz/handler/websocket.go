package handler

import (
	"Tiktok/biz/dao/db"
	"Tiktok/biz/dao/re"
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
	db *db.MySQLdb
	re *re.Redis
}

func NewWebsocketSever(db *db.MySQLdb, re *re.Redis) *WebsocketSever {
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
			ID:     utils.CreateId(uid, toUserId),
			SendID: utils.CreateId(toUserId, uid),
			Socket: conn,
			Send:   make(chan []byte),
		}
		service.Manager.Register <- client
		go client.Read(m.re, m.db)
		go client.Write()
	})
	wsAdaptor := adaptor.HertzHandler(stdHandler)
	wsAdaptor(ctx, c)
}
