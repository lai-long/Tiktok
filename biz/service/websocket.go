package service

import (
	"Tiktok/biz/cache"
	"Tiktok/biz/dao"
	"Tiktok/biz/model/dto"
	"Tiktok/pkg/consts"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/encoding/protojson"
)

var Manager = ClientManager{
	Clients:        make(map[string]*Client),
	Broadcast:      make(chan *Broadcast),
	Groups:         make(map[string][]*Client),
	GroupBroadcast: make(chan *GroupBroadcast),
	Register:       make(chan *Client),
	Unregister:     make(chan *Client),
}

type Client struct {
	ID      string
	GroupId string
	SendID  string
	Socket  *websocket.Conn
	Send    chan []byte
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
type ClientManager struct {
	Clients        map[string]*Client
	Groups         map[string][]*Client
	Broadcast      chan *Broadcast
	GroupBroadcast chan *GroupBroadcast
	Register       chan *Client
	Unregister     chan *Client
	mu             sync.RWMutex
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for {
		c.Socket.PongHandler()
		sendMsg := new(dto.SendMsg)
		err := c.Socket.ReadJSON(sendMsg)
		if err != nil {
			log.Println("client ReadJSON err", err)
			Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}
		//type 1 一对一聊天
		//type 2 获取未在线时的消息
		//type 3 获取历史消息
		//type 4 群聊
		if sendMsg.Type == "1" {
			Manager.Broadcast <- &Broadcast{
				Type:    "1",
				Clients: c,
				Message: []byte(sendMsg.Content),
			}
		} else if sendMsg.Type == "2" {
			Manager.Broadcast <- &Broadcast{
				Type:    "2",
				Clients: c,
			}
		} else if sendMsg.Type == "3" {
			Manager.Broadcast <- &Broadcast{
				Type:     "3",
				Clients:  c,
				PageNum:  sendMsg.PageNum,
				PageSize: sendMsg.PageSize,
			}
		} else if sendMsg.Type == "4" {
			Manager.mu.Lock()
			members, _ := Manager.Groups[c.GroupId]
			Manager.GroupBroadcast <- &GroupBroadcast{
				Clients: members,
				Message: []byte(sendMsg.Content),
				Type:    sendMsg.Type,
			}
			Manager.mu.Unlock()
		}
	}
}
func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "websocket connection closed")
				err := c.Socket.WriteMessage(websocket.CloseMessage, closeMsg)
				if err != nil {
					log.Println("write closeMsg err", err)
				}
				log.Println("client Write err")
				return
			}
			_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
func (manager *ClientManager) Start(m *dao.MySQLdb, re *cache.Redis) {

	for {
		select {
		case client := <-manager.Register:
			log.Println("建立websocket连接", client.ID)
			if client.GroupId != "" {
				manager.mu.Lock()
				manager.Groups[client.GroupId] = append(manager.Groups[client.GroupId], client)
				manager.mu.Unlock()
			}
			Manager.Clients[client.ID] = client
			replyMSg := dto.ReplyMsg{
				From:    client.ID,
				Code:    consts.CodeSuccess,
				Content: "连接成功",
			}
			msg, _ := protojson.Marshal(&replyMSg)
			client.Send <- msg
		case client := <-manager.Unregister:
			log.Println("断开websocket连接", client.ID)
			if client.GroupId != "" {
				manager.mu.Lock()
				for i, v := range manager.Groups[client.GroupId] {
					if v.ID == client.ID {
						manager.Groups[client.GroupId] = append(manager.Groups[client.GroupId][:i], manager.Groups[client.GroupId][i+1:]...)
						break
					}
				}
				manager.mu.Unlock()
			}
			if _, ok := Manager.Clients[client.ID]; ok {
				replyMSg := dto.ReplyMsg{
					From:    client.ID,
					Code:    consts.CodeSuccess,
					Content: "连接中断",
				}
				msg, _ := protojson.Marshal(&replyMSg)
				client.Send <- msg
				close(client.Send)
				delete(manager.Clients, client.ID)
			}
		case broadcast := <-manager.Broadcast:
			if broadcast.Type == "1" {
				message := broadcast.Message
				sendId := broadcast.Clients.SendID
				flag := false
				manager.mu.RLock()
				for id, client := range Manager.Clients {
					if id != sendId {
						continue
					}
					replyMSg := dto.ReplyMsg{
						From:    client.ID,
						Code:    consts.CodeSuccess,
						Content: string(message),
					}
					msg, _ := protojson.Marshal(&replyMSg)
					select {
					case client.Send <- msg:
						flag = true
					default:
						close(client.Send)
						delete(manager.Clients, client.ID)
					}
				}
				manager.mu.RUnlock()
				id := broadcast.Clients.ID
				if flag {
					replyMSg := dto.ReplyMsg{
						From:    broadcast.Clients.ID,
						Code:    consts.CodeSuccess,
						Content: "对方在线",
					}
					msg, _ := protojson.Marshal(&replyMSg)
					broadcast.Clients.Send <- msg
					m.InsertMsg(id, string(message))
				} else {
					replyMSg := dto.ReplyMsg{
						From:    broadcast.Clients.ID,
						Code:    consts.CodeSuccess,
						Content: "对方不在线",
					}
					m.InsertMsg(id, string(message))
					re.SaveOfflineMsg(broadcast.Clients.SendID, string(message))
					msg, _ := protojson.Marshal(&replyMSg)
					broadcast.Clients.Send <- msg
				}
			} else if broadcast.Type == "2" {
				message, err := re.FetchOfflineMsg(broadcast.Clients.SendID)
				if err != nil {
					log.Println(err)
				}
				str := strings.Join(message, ",\n ")
				finalInfo := str + fmt.Sprintf("\ntotal:%d", len(str))
				replyMSg := dto.ReplyMsg{
					From:    "未在线时收到消息",
					Code:    consts.CodeSuccess,
					Content: finalInfo,
				}
				msg, _ := protojson.Marshal(&replyMSg)
				broadcast.Clients.Send <- msg
			} else if broadcast.Type == "3" {
				pageNum := 0
				pageSize := 10
				pageNum, err := strconv.Atoi(broadcast.PageNum)
				if err != nil {
					log.Println(err)
				}
				pageSize, err = strconv.Atoi(broadcast.PageSize)
				if err != nil {
					log.Println(err)
				}

				msgs := m.GetWebsocketHistory(broadcast.Clients.ID, broadcast.Clients.SendID, pageNum, pageSize)
				str := strings.Join(msgs, ",\n ")
				finalInfo := str + fmt.Sprintf("\ntotal:%d", len(msgs))
				if err != nil {
					log.Println(err)
				}
				replyMSg := dto.ReplyMsg{
					From:    broadcast.Clients.ID + "and" + broadcast.Clients.SendID,
					Code:    consts.CodeSuccess,
					Content: finalInfo,
				}
				msg, _ := protojson.Marshal(&replyMSg)
				broadcast.Clients.Send <- msg
			}
		case groupBroadcast := <-manager.GroupBroadcast:
			for _, client := range groupBroadcast.Clients {
				replyMSg := dto.ReplyMsg{
					From:    client.ID,
					Code:    consts.CodeSuccess,
					Content: string(groupBroadcast.Message),
				}
				msg, _ := protojson.Marshal(&replyMSg)
				client.Send <- msg
			}
		}
	}
}
