package service

import (
	"Tiktok/biz/dao/db"
	"Tiktok/biz/dao/re"
	"Tiktok/biz/model/dto"
	"Tiktok/pkg/consts"
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

type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Register   chan *Client
	Unregister chan *Client
}

func (c *Client) Read(ctx context.Context) {
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()
	for {
		c.Socket.PongHandler()
		sendMsg := new(dto.SendMsg)
		err := c.Socket.ReadJSON(sendMsg)
		if err != nil {
			log.Println("client ReadJSON err", err)
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		//type 1 一对一
		//type 2 获取未在线时的消息
		if sendMsg.Type == "1" {
			Manager.Broadcast <- &Broadcast{
				Type:    "1",
				Clients: c,
				Message: []byte(sendMsg.Content),
			}
		} else if sendMsg.Type == "2" {
			msg, _ := re.FetchOfflineMsg(c.ID)
			message, err := json.Marshal(msg)
			if err != nil {
				log.Println("re.FetchOfflineMsg json.Marshal err", err)
			}
			Manager.Broadcast <- &Broadcast{
				Type:    "2",
				Clients: c,
				Message: message,
			}
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
			replyMsg := dto.ReplyMsg{
				From:    c.ID,
				Code:    consts.CodeSuccess,
				Content: fmt.Sprintf("%s", string(message)),
			}
			msg, _ := json.Marshal(replyMsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
func (manager *ClientManager) Start(m *db.MySQLdb) {
	for {
		select {
		case client := <-manager.Register:
			log.Println("建立websocket连接", client.ID)
			Manager.Clients[client.ID] = client
			replyMSg := dto.ReplyMsg{
				Code:    consts.CodeSuccess,
				Content: "连接成功",
			}
			msg, _ := json.Marshal(replyMSg)
			_ = client.Socket.WriteMessage(websocket.TextMessage, msg)
		case client := <-manager.Unregister:
			log.Println("断开websocket连接", client.ID)
			if _, ok := Manager.Clients[client.ID]; ok {
				replyMSg := dto.ReplyMsg{
					Code:    consts.CodeSuccess,
					Content: "连接中断",
				}
				msg, _ := json.Marshal(replyMSg)
				_ = client.Socket.WriteMessage(websocket.TextMessage, msg)
				close(client.Send)
				delete(manager.Clients, client.ID)
			}
		case broadcast := <-manager.Broadcast:
			if broadcast.Type == "1" {
				log.Println("进行broadcast")
				message := broadcast.Message
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
					replyMSg := dto.ReplyMsg{
						From:    broadcast.Clients.ID,
						Code:    consts.CodeSuccess,
						Content: "对方在线",
					}
					msg, _ := json.Marshal(replyMSg)
					_ = broadcast.Clients.Socket.WriteMessage(websocket.TextMessage, msg)
					sender, receiver := utils.GetId(id)
					m.InsertMsg(id, string(message), sender, receiver)
				} else {
					replyMSg := dto.ReplyMsg{
						From:    broadcast.Clients.ID,
						Code:    consts.CodeSuccess,
						Content: "对方不在线",
					}
					re.SaveOfflineMsg(broadcast.Clients.SendID, string(message))
					msg, _ := json.Marshal(replyMSg)
					_ = broadcast.Clients.Socket.WriteMessage(websocket.TextMessage, msg)
				}
			} else if broadcast.Type == "2" {
				replyMSg := dto.ReplyMsg{
					From:    "未在线时收到消息",
					Code:    consts.CodeSuccess,
					Content: string(broadcast.Message),
				}
				msg, _ := json.Marshal(replyMSg)
				_ = broadcast.Clients.Socket.WriteMessage(websocket.TextMessage, msg)
				sender, receiver := utils.GetId(broadcast.Clients.SendID)
				m.InsertMsg(broadcast.Clients.SendID, string(broadcast.Message), sender, receiver)
			}
		}
	}
}
