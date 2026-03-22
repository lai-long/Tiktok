package service

import (
	"Tiktok/biz/dao/db"
	"Tiktok/biz/dao/re"
	"Tiktok/biz/model/websocketModel"
	"Tiktok/pkg/utils"
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/gorilla/websocket"
)

var Manager = ClientManager{
	Clients:    make(map[string]*Client),
	Broadcast:  make(chan *Broadcast),
	Reply:      make(chan *Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan []byte
}
type Broadcast struct {
	Clients *Client
	Message []byte
	Type    string
}

//新用户连接>>创建client>>放入register channel>>加入client map
//用户发送消息>>解析消息>>放入Broadcast Channel>>从broadcast得到目标用户>>发送消息
//用户断开>>放入unregister通道>>从clients map中删除

type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

func (c *Client) Read(ctx context.Context) {
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()
	for {
		//检测活跃
		c.Socket.PongHandler()
		sendMsg := new(websocketModel.SendMsg)
		err := c.Socket.ReadJSON(sendMsg)
		if err != nil {
			log.Println("client ReadJSON err", err)
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		//type 1 1->2
		if sendMsg.Type == "1" {
			//r1 := re.GetMsgCountsBYClientID(ctx, c.ID)
			//r2 := re.GetMsgCountsBYClientID(ctx, c.SendID)
			re.MSgCountIncr(ctx, c.ID)
		}
		Manager.Broadcast <- &Broadcast{
			Clients: c,
			Message: []byte(sendMsg.Content),
		}
	}
}
func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("client Write err")
				return
			}
			replyMsg := websocketModel.ReplyMsg{
				Code:    "100",
				Content: fmt.Sprintf("%s", string(message)),
			}
			msg, _ := json.Marshal(replyMsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
func (manager *ClientManager) Start(m *db.MySQLdb) {
	for {
		log.Println("websocket 启动")
		select {
		case client := <-manager.Register:
			log.Println("建立websocket连接", client.ID)
			Manager.Clients[client.ID] = client
			replyMSg := websocketModel.ReplyMsg{
				Code:    "100",
				Content: "连接成功",
			}
			msg, _ := json.Marshal(replyMSg)
			_ = client.Socket.WriteMessage(websocket.TextMessage, msg)
		case client := <-manager.Unregister:
			log.Println("断开websocket连接", client.ID)
			if _, ok := Manager.Clients[client.ID]; ok {
				replyMSg := websocketModel.ReplyMsg{
					Code:    "100",
					Content: "连接中断",
				}
				msg, _ := json.Marshal(replyMSg)
				_ = client.Socket.WriteMessage(websocket.TextMessage, msg)
				close(client.Send)
				delete(manager.Clients, client.ID)
			}
		case broadcast := <-manager.Broadcast:
			log.Println("进行broadcast")
			message := broadcast.Message
			log.Println("[]byte message", message)
			log.Println("string message", string(message))
			sendId := broadcast.Clients.SendID
			flag := false
			for id, client := range Manager.Clients {
				if id != sendId {
					continue
				}
				select {
				case client.Send <- message:
					flag = true
				default:
					close(client.Send)
					delete(manager.Clients, client.ID)
				}
			}
			id := broadcast.Clients.ID
			if flag {
				replyMSg := websocketModel.ReplyMsg{
					Code:    "100",
					Content: "在线",
				}
				msg, _ := json.Marshal(replyMSg)
				_ = broadcast.Clients.Socket.WriteMessage(websocket.TextMessage, msg)
				sender, receiver := utils.GetId(id)
				m.InsertMsg(id, string(message), sender, receiver)
			} else {
				replyMSg := websocketModel.ReplyMsg{
					Code:    "100",
					Content: "对方不在线",
				}
				msg, _ := json.Marshal(replyMSg)
				_ = broadcast.Clients.Socket.WriteMessage(websocket.TextMessage, msg)
				sender, receiver := utils.GetId(id)
				log.Println(sender, receiver, string(message), id)
				m.InsertMsg(id, string(message), sender, receiver)
			}
		}

	}
}
