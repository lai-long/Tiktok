package websocket

import "sync"

type clientManager struct {
	Clients        map[string]*Client
	Groups         map[string][]*Client
	Broadcast      chan *Broadcast
	GroupBroadcast chan *GroupBroadcast
	Register       chan *Client
	Unregister     chan *Client
	mu             sync.RWMutex
}

func NewClientManager() *clientManager {
	return &clientManager{
		Clients:        make(map[string]*Client),
		Groups:         make(map[string][]*Client),
		Broadcast:      make(chan *Broadcast),
		GroupBroadcast: make(chan *GroupBroadcast),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		mu:             sync.RWMutex{},
	}
}
