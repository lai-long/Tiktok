package ws

import (
	"Tiktok/biz/model/chat"
	"Tiktok/biz/model/common"
	"Tiktok/internal/ws/service"
	"Tiktok/pkg/dal/cache"
	"Tiktok/pkg/dal/dao"
	"Tiktok/pkg/logger"
	"Tiktok/pkg/utils"
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WebsocketServer struct {
	db        *dao.MySQLdb
	re        *cache.Redis
	websocket *service.WebsocketService
}

func NewWebsocketServer(db *dao.MySQLdb, re *cache.Redis, ws *service.WebsocketService) *WebsocketServer {
	return &WebsocketServer{
		db:        db,
		re:        re,
		websocket: ws,
	}
}

func (m *WebsocketServer) WebSocketHandler(ctx context.Context, c *app.RequestContext) {
	userid := utils.GetUserID(c)
	if userid == "" {
		c.JSON(200, chat.WebsocketResp{Base: &common.Base{
			Code: 200,
			Msg:  "Unauthorized: user_id not found",
		}})
		return
	}
	logger.Info("user connected", logger.WithUserID(userid))
	req := new(chat.WebsocketReq)
	err := c.BindAndValidate(req)
	if err != nil {
		logger.Error("bindAndValidate error", zap.Error(err))
		c.JSON(200, chat.WebsocketResp{Base: &common.Base{
			Code: 0,
			Msg:  "Invalid request parameters",
		}})
		return
	}
	stdHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}).Upgrade(w, r, nil)
		if err != nil {
			logger.Error("websocket upgrade error", zap.Error(err))
			http.Error(w, "Could not upgrade to WebSocket", http.StatusInternalServerError)
			return
		}
		client := &service.Client{
			ID:      utils.CreateID(userid, req.ToUserId),
			SendID:  utils.CreateID(req.ToUserId, userid),
			GroupId: req.GroupId,
			Socket:  conn,
			Send:    make(chan []byte, 128),
			Ctx:     ctx,
		}
		logger.Info("websocket client registered", zap.String("client_id", client.ID))
		m.websocket.Manager.Register <- client
		go m.websocket.Read(client)
		go m.websocket.Write(client)
		go m.websocket.Heartbeat(client)
	})
	wsAdaptor := adaptor.HertzHandler(stdHandler)
	wsAdaptor(ctx, c)
}
