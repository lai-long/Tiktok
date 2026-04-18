package websocket

import (
	"Tiktok/biz/dal/cache"
	"Tiktok/biz/dal/dao"
	"Tiktok/biz/model/chat"
	"Tiktok/biz/service/ai"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/utils"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/encoding/protojson"
)

type WebsocketService struct {
	mysql   *dao.MySQLdb
	redis   *cache.Redis
	Manager *clientManager
}

func NewWebsocketService(mysql *dao.MySQLdb, re *cache.Redis) *WebsocketService {
	return &WebsocketService{
		mysql:   mysql,
		redis:   re,
		Manager: NewClientManager(),
	}
}

type Client struct {
	ID      string
	GroupId string
	SendID  string
	Socket  *websocket.Conn
	Send    chan []byte
	Ctx     context.Context
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
	}()

	c.Socket.SetPongHandler(func(string) error {
		return nil
	})
	for {
		sendMsg := new(chat.SendMsg)
		err := c.Socket.ReadJSON(sendMsg)
		if err != nil {
			log.Println("client ReadJSON err", err)
			ws.Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}
		ok, question := utils.CheckAiKeyWord(sendMsg.Content)
		if ok {
			agent := ai.NewAgent(context.Background())
			go func(q string) {
				resp := agent.StartAction(question)
				ws.aiReplyToClient(resp, c)
			}(question)
		}
		switch sendMsg.Type {
		case "1":
			ws.Manager.Broadcast <- &Broadcast{
				Type:    "1",
				Clients: c,
				Message: []byte(sendMsg.Content),
			}
		case "2":
			ws.Manager.Broadcast <- &Broadcast{
				Type:    "2",
				Clients: c,
			}
		case "3":
			ws.Manager.Broadcast <- &Broadcast{
				Type:     "3",
				Clients:  c,
				PageNum:  sendMsg.PageNum,
				PageSize: sendMsg.PageSize,
			}
		case "4":
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
	log.Println("建立websocket连接", client.ID)
	if client.GroupId != "" {
		ws.Manager.mu.Lock()
		ws.Manager.Groups[client.GroupId] = append(ws.Manager.Groups[client.GroupId], client)
		ws.Manager.mu.Unlock()
	}
	ws.Manager.Clients[client.ID] = client
	replyMSg := chat.ReplyMsg{
		From:    client.ID,
		Code:    consts.Success,
		Content: "连接成功",
	}
	msg, _ := protojson.Marshal(&replyMSg)
	client.Send <- msg
}

func (ws *WebsocketService) startUnregister(client *Client) {
	log.Println("断开websocket连接", client.ID)
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
		replyMSg := chat.ReplyMsg{
			From:    client.ID,
			Code:    consts.Success,
			Content: "连接中断",
		}
		msg, _ := protojson.Marshal(&replyMSg)
		client.Send <- msg
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
		replyMSg := chat.ReplyMsg{
			From:    client.ID,
			Code:    consts.Success,
			Content: string(message),
		}
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
		replyMSg := chat.ReplyMsg{
			From:    broadcast.Clients.ID,
			Code:    consts.Success,
			Content: "对方在线",
		}
		msg, _ := protojson.Marshal(&replyMSg)
		broadcast.Clients.Send <- msg
		err := ws.mysql.InsertMsg(id, string(message))
		if err != nil {
			log.Println("Insert message error:", err)
		}
	} else {
		replyMSg := chat.ReplyMsg{
			From:    broadcast.Clients.ID,
			Code:    consts.Success,
			Content: "对方不在线",
		}
		err := ws.mysql.InsertMsg(id, string(message))
		if err != nil {
			log.Println("Insert message error:", err)
		}
		err = ws.redis.SaveOfflineMsg(broadcast.Clients.SendID, string(message))
		if err != nil {
			log.Println("Save offline message error:", err)
		}
		msg, _ := protojson.Marshal(&replyMSg)
		broadcast.Clients.Send <- msg
	}
}

func (ws *WebsocketService) startBroadcastOneOffline(broadcast *Broadcast) {
	message, err := ws.redis.FetchOfflineMsg(broadcast.Clients.SendID)
	if err != nil {
		log.Println("Fetch offline message error:", err)
		replyMSg := chat.ReplyMsg{
			From:    "系统",
			Code:    consts.Success,
			Content: "获取离线消息失败",
		}
		msg, _ := protojson.Marshal(&replyMSg)
		broadcast.Clients.Send <- msg
		return
	}
	str := strings.Join(message, ",\n ")
	finalInfo := str + fmt.Sprintf("\ntotal:%d", len(message))
	replyMSg := chat.ReplyMsg{
		From:    "未在线时收到消息",
		Code:    consts.Success,
		Content: finalInfo,
	}
	msg, _ := protojson.Marshal(&replyMSg)
	broadcast.Clients.Send <- msg
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
		replyMSg := chat.ReplyMsg{
			From:    "系统",
			Code:    consts.Success,
			Content: "获取历史消息失败",
		}
		msg, _ := protojson.Marshal(&replyMSg)
		broadcast.Clients.Send <- msg
		return
	}
	str := strings.Join(msgs, ",\n ")
	finalInfo := str + fmt.Sprintf("\ntotal:%d", len(msgs))
	replyMSg := chat.ReplyMsg{
		From:    broadcast.Clients.ID + "and" + broadcast.Clients.SendID,
		Code:    consts.Success,
		Content: finalInfo,
	}
	msg, _ := protojson.Marshal(&replyMSg)
	broadcast.Clients.Send <- msg
}

func (ws *WebsocketService) startBroadcastGroupOnline(groupBroadcast *GroupBroadcast) {
	for _, client := range groupBroadcast.Clients {
		replyMSg := chat.ReplyMsg{
			From:    client.ID,
			Code:    consts.Success,
			Content: string(groupBroadcast.Message),
		}
		msg, _ := protojson.Marshal(&replyMSg)
		client.Send <- msg
	}
}

func (ws *WebsocketService) startBroadcastOneError(broadcast *Broadcast) {
	replyMSg := chat.ReplyMsg{
		From:    "system",
		Code:    consts.Success,
		Content: "请求类型不存在",
	}
	log.Println("请求类型不存在")
	msg, _ := protojson.Marshal(&replyMSg)
	broadcast.Clients.Send <- msg
}

func (ws *WebsocketService) aiReplyToClient(resp string, c *Client) {
	if resp == "" {
		replyMSg := chat.ReplyMsg{
			From:    "AI",
			Code:    consts.Success,
			Content: "ai不理你",
		}
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
		return
	}
	replyContent := resp
	replyMSg := chat.ReplyMsg{
		From:    "AI",
		Code:    consts.Success,
		Content: replyContent,
	}
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
