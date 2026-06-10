package service

import "sync"

type ClientManager struct {
	clients        map[string]*Client
	groups         map[string][]*Client
	Broadcast      chan *Broadcast
	GroupBroadcast chan *GroupBroadcast
	Register       chan *Client
	Unregister     chan *Client
	mu             sync.RWMutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients:        make(map[string]*Client),
		groups:         make(map[string][]*Client),
		Broadcast:      make(chan *Broadcast),
		GroupBroadcast: make(chan *GroupBroadcast),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
	}
}

func (m *ClientManager) AddClient(c *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[c.ID] = c
}

func (m *ClientManager) RemoveClient(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.clients[id]; !ok {
		return false
	}
	delete(m.clients, id)
	return true
}

func (m *ClientManager) GetClient(id string) (*Client, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	client, ok := m.clients[id]
	return client, ok
}

func (m *ClientManager) AddToGroup(groupId string, c *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.groups[groupId] = append(m.groups[groupId], c)
}

func (m *ClientManager) RemoveFromGroup(groupId, clientId string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	group := m.groups[groupId]
	for i, v := range group {
		if v.ID == clientId {
			m.groups[groupId] = append(group[:i], group[i+1:]...)
			return
		}
	}
}

func (m *ClientManager) GetGroup(groupId string) []*Client {
	m.mu.RLock()
	defer m.mu.RUnlock()
	groupClients := m.groups[groupId]
	if len(groupClients) == 0 {
		return nil
	}
	members := make([]*Client, len(groupClients))
	copy(members, groupClients)
	return members
}
