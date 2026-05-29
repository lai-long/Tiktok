package service

import (
	"Tiktok/biz/model/chat"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/dal/cache"
	"Tiktok/pkg/dal/dao"
	"Tiktok/pkg/logger"
	"Tiktok/pkg/utils"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

type WebsocketService struct {
	mysql   *dao.MySQLdb
	redis   *cache.Redis
	Manager *ClientManager
}

func NewWebsocketService(mysql *dao.MySQLdb, re *cache.Redis) *WebsocketService {
	return &WebsocketService{
		mysql:   mysql,
		redis:   re,
		Manager: NewClientManager(),
	}
}

type Client struct {
	ID       string
	GroupId  string
	SendID   string
	Socket   *websocket.Conn
	Send     chan []byte
	Ctx      context.Context
	agent    *Agent
	agentMu  sync.Mutex
	lastPong time.Time
	mu       sync.Mutex
}
type Broadcast struct {
	Clients  *Client
	Message  []byte
	Type     string
	PageNum  string
	PageSize string
}
type GroupBroadcast struct {
	Clients []*Client
	Message []byte
	Type    string
}

func (ws *WebsocketService) Read(c *Client) {
	defer func() {
		ws.Manager.Unregister <- c
		_ = c.Socket.Close()
		if c.agent != nil {
			c.agent.StopAction()
		}
	}()

	c.Socket.SetPongHandler(func(string) error {
		c.mu.Lock()
		c.lastPong = time.Now()
		c.mu.Unlock()
		return nil
	})
	for {
		sendMsg := new(chat.SendMsg)
		err := c.Socket.ReadJSON(sendMsg)
		if err != nil {
			logger.Error("client ReadJSON error", zap.Error(err))
			ws.Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}
		ok, question := utils.CheckAiKeyWord(sendMsg.Content)
		if ok {
			go func(q string) {
				ctx, cancel := context.WithTimeout(c.Ctx, 30*time.Second)
				defer cancel()
				c.agentMu.Lock()
				defer c.agentMu.Unlock()
				if c.agent == nil {
					c.agent = NewAgent(ctx)
				}
				if c.agent == nil {
					ws.aiReplyToClient("", c)
					return
				}
				resp := c.agent.StartAction(q)
				ws.aiReplyToClient(resp, c)
			}(question)
		}
		switch sendMsg.Type {
		case consts.MsgTypePrivate:
			ws.Manager.Broadcast <- &Broadcast{
				Type:    consts.MsgTypePrivate,
				Clients: c,
				Message: []byte(sendMsg.Content),
			}
		case consts.MsgTypeOffline:
			ws.Manager.Broadcast <- &Broadcast{
				Type:    consts.MsgTypeOffline,
				Clients: c,
			}
		case consts.MsgTypeHistory:
			ws.Manager.Broadcast <- &Broadcast{
				Type:     consts.MsgTypeHistory,
				Clients:  c,
				PageNum:  sendMsg.PageNum,
				PageSize: sendMsg.PageSize,
			}
		case consts.MsgTypeGroupMsg:
			ws.Manager.mu.Lock()
			ws.Manager.GroupBroadcast <- &GroupBroadcast{
				Clients: ws.Manager.Groups[c.GroupId],
				Message: []byte(sendMsg.Content),
				Type:    sendMsg.Type,
			}
			ws.Manager.mu.Unlock()
		}
	}
}

func (ws *WebsocketService) Write(c *Client) {
	defer func() {
		_ = c.Socket.Close()
	}()
	for message := range c.Send {
		_ = c.Socket.WriteMessage(websocket.TextMessage, message)

	}
}

func (ws *WebsocketService) Heartbeat(c *Client) {
	ticker := time.NewTicker(15 * time.Second)
	defer func() {
		ticker.Stop()
		c.mu.Lock()
		_ = c.Socket.Close()
		c.mu.Unlock()
		ws.Manager.Unregister <- c
	}()

	for {
		<-ticker.C
		c.mu.Lock()
		if time.Since(c.lastPong) > 40*time.Second {
			c.mu.Unlock()
			logger.Warn("heartbeat timeout, disconnecting", zap.String("client_id", c.ID))
			return
		}
		c.mu.Unlock()

		if err := c.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
			logger.Error("ping failed", zap.String("client_id", c.ID), zap.Error(err))
			return
		}
	}
}

func (ws *WebsocketService) Start() {
	for {
		select {
		case client := <-ws.Manager.Register:
			ws.startRegister(client)
		case client := <-ws.Manager.Unregister:
			ws.startUnregister(client)
		case broadcast := <-ws.Manager.Broadcast:
			switch broadcast.Type {
			case "1":
				ws.startBroadcastOneOnline(broadcast)
			case "2":
				ws.startBroadcastOneOffline(broadcast)
			case "3":
				ws.startBroadcastOneHistory(broadcast)
			default:
				ws.startBroadcastOneError(broadcast)
			}
		case groupBroadcast := <-ws.Manager.GroupBroadcast:
			ws.startBroadcastGroupOnline(groupBroadcast)
		}
	}
}

func (ws *WebsocketService) startRegister(client *Client) {
	logger.Info("websocket connected", zap.String("client_id", client.ID))
	if client.GroupId != "" {
		ws.Manager.mu.Lock()
		ws.Manager.Groups[client.GroupId] = append(ws.Manager.Groups[client.GroupId], client)
		ws.Manager.mu.Unlock()
	}
	ws.Manager.Clients[client.ID] = client
	ws.sendReplyMsg(client, client.ID, consts.WsConnectSuccess, consts.GetErrorCodeMsg(consts.WsConnectSuccess))
}

func (ws *WebsocketService) startUnregister(client *Client) {
	logger.Info("websocket disconnected", zap.String("client_id", client.ID))
	if client.GroupId != "" {
		ws.Manager.mu.Lock()
		for i, v := range ws.Manager.Groups[client.GroupId] {
			if v.ID == client.ID {
				ws.Manager.Groups[client.GroupId] = append(ws.Manager.Groups[client.GroupId][:i],
					ws.Manager.Groups[client.GroupId][i+1:]...)
				break
			}
		}
		ws.Manager.mu.Unlock()
	}
	if _, ok := ws.Manager.Clients[client.ID]; ok {
		ws.sendReplyMsg(client, client.ID, consts.WsDisconnect, consts.GetErrorCodeMsg(consts.WsDisconnect))
		close(client.Send)
		delete(ws.Manager.Clients, client.ID)
	}
}

func (ws *WebsocketService) startBroadcastOneOnline(broadcast *Broadcast) {
	message := broadcast.Message
	sendId := broadcast.Clients.SendID
	flag := false
	ws.Manager.mu.Lock()
	for id, client := range ws.Manager.Clients {
		if id != sendId {
			continue
		}
		replyMSg := chat.ReplyMsg{From: client.ID, Code: consts.Success, Content: string(message)}
		msg, _ := protojson.Marshal(&replyMSg)
		select {
		case client.Send <- msg:
			flag = true
		default:
			close(client.Send)
			delete(ws.Manager.Clients, client.ID)
		}
	}
	ws.Manager.mu.Unlock()
	id := broadcast.Clients.ID
	if flag {
		ws.sendReplyMsg(broadcast.Clients, broadcast.Clients.ID, consts.WsClientOnline, consts.GetErrorCodeMsg(consts.WsClientOnline))
		err := ws.mysql.InsertMsg(id, string(message))
		if err != nil {
			logger.Error("insert message error", zap.Error(err))
		}
	} else {
		ws.sendReplyMsg(broadcast.Clients, broadcast.Clients.ID, consts.WsClientNotOnline, consts.GetErrorCodeMsg(consts.WsClientNotOnline))
		err := ws.mysql.InsertMsg(id, string(message))
		if err != nil {
			logger.Error("insert message error", zap.Error(err))
		}
		err = ws.redis.SaveOfflineMsg(broadcast.Clients.SendID, string(message))
		if err != nil {
			logger.Error("save offline message error", zap.Error(err))
		}
	}
}

func (ws *WebsocketService) startBroadcastOneOffline(broadcast *Broadcast) {
	message, err := ws.redis.FetchOfflineMsg(broadcast.Clients.SendID)
	if err != nil {
		logger.Error("fetch offline message error", zap.Error(err))
		ws.sendReplyMsg(broadcast.Clients, "系统", consts.WsGetOfflineError, consts.GetErrorCodeMsg(consts.WsGetOfflineError))
		return
	}
	ws.sendReplyMsg(broadcast.Clients, "未在线时收到消息", consts.Success, formatMessageList(message))
}

func (ws *WebsocketService) startBroadcastOneHistory(broadcast *Broadcast) {
	pageNum := 0
	pageSize := 10
	if broadcast.PageNum != "" {
		pageNum, _ = strconv.Atoi(broadcast.PageNum)
	}
	if broadcast.PageSize != "" {
		pageSize, _ = strconv.Atoi(broadcast.PageSize)
	}
	msgs, err := ws.mysql.GetWebsocketHistory(broadcast.Clients.ID, broadcast.Clients.SendID, pageNum, pageSize)
	if err != nil || msgs == nil {
		ws.sendReplyMsg(broadcast.Clients, "系统", consts.WsGetHistoryError, consts.GetErrorCodeMsg(consts.WsGetHistoryError))
		return
	}
	ws.sendReplyMsg(broadcast.Clients, broadcast.Clients.ID+"and"+broadcast.Clients.SendID, consts.Success, formatMessageList(msgs))
}

func (ws *WebsocketService) startBroadcastGroupOnline(groupBroadcast *GroupBroadcast) {
	for _, client := range groupBroadcast.Clients {
		ws.sendReplyMsg(client, client.ID, consts.Success, string(groupBroadcast.Message))
	}
}

func (ws *WebsocketService) startBroadcastOneError(broadcast *Broadcast) {
	logger.Error("unknown broadcast type")
	ws.sendReplyMsg(broadcast.Clients, "system", consts.WsReqValidError, consts.GetErrorCodeMsg(consts.WsReqValidError))
}

func (ws *WebsocketService) aiReplyToClient(resp string, c *Client) {
	if resp == "" {
		ws.sendReplyMsgAndSendID(c, "AI", consts.WsAIReplyEmpty, consts.GetErrorCodeMsg(consts.WsAIReplyEmpty))
		return
	}
	ws.sendReplyMsgAndSendID(c, "AI", consts.Success, resp)
}

func (ws *WebsocketService) sendReplyMsg(client *Client, from string, code int32, content string) {
	replyMSg := chat.ReplyMsg{From: from, Code: code, Content: content}
	msg, _ := protojson.Marshal(&replyMSg)
	client.Send <- msg
}

func (ws *WebsocketService) sendReplyMsgAndSendID(c *Client, from string, code int32, content string) {
	replyMSg := chat.ReplyMsg{From: from, Code: code, Content: content}
	msg, _ := protojson.Marshal(&replyMSg)
	c.Send <- msg
	if c.SendID != "" {
		ws.Manager.mu.Lock()
		for id, client := range ws.Manager.Clients {
			if id == c.SendID {
				client.Send <- msg
				break
			}
		}
		ws.Manager.mu.Unlock()
	}
}

func formatMessageList(messages []string) string {
	if len(messages) == 0 {
		return ""
	}
	return strings.Join(messages, ",\n ") + fmt.Sprintf("\ntotal:%d", len(messages))
}
