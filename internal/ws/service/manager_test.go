package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientManager(t *testing.T) {
	t.Run("NewClientManager_initializes_maps_and_channels", func(t *testing.T) {
		m := NewClientManager()
		assert.NotNil(t, m.clients)
		assert.NotNil(t, m.groups)
		assert.NotNil(t, m.Broadcast)
		assert.NotNil(t, m.GroupBroadcast)
		assert.NotNil(t, m.Register)
		assert.NotNil(t, m.Unregister)
	})

	t.Run("AddClient_and_GetClient", func(t *testing.T) {
		m := NewClientManager()
		c := &Client{ID: "client1"}

		m.AddClient(c)

		got, ok := m.GetClient("client1")
		assert.True(t, ok)
		assert.Equal(t, c, got)

		_, ok = m.GetClient("missing")
		assert.False(t, ok)
	})

	t.Run("RemoveClient", func(t *testing.T) {
		m := NewClientManager()
		m.AddClient(&Client{ID: "client1"})

		assert.False(t, m.RemoveClient("missing"))
		assert.True(t, m.RemoveClient("client1"))
		assert.False(t, m.RemoveClient("client1"))

		_, ok := m.GetClient("client1")
		assert.False(t, ok)
	})

	t.Run("AddToGroup_and_GetGroup", func(t *testing.T) {
		m := NewClientManager()
		c1 := &Client{ID: "client1"}
		c2 := &Client{ID: "client2"}

		m.AddToGroup("group1", c1)
		m.AddToGroup("group1", c2)

		group := m.GetGroup("group1")
		assert.Len(t, group, 2)
		assert.Equal(t, c1, group[0])
		assert.Equal(t, c2, group[1])

		// GetGroup should return a copy; mutating the returned slice
		// must not affect the manager's internal state.
		group = append(group, &Client{ID: "client3"})
		assert.Len(t, m.GetGroup("group1"), 2)

		assert.Nil(t, m.GetGroup("empty"))
	})

	t.Run("RemoveFromGroup", func(t *testing.T) {
		m := NewClientManager()
		c1 := &Client{ID: "client1"}
		c2 := &Client{ID: "client2"}
		c3 := &Client{ID: "client3"}

		m.AddToGroup("group1", c1)
		m.AddToGroup("group1", c2)
		m.AddToGroup("group1", c3)

		m.RemoveFromGroup("group1", "client2")
		group := m.GetGroup("group1")
		assert.Len(t, group, 2)
		assert.Equal(t, "client1", group[0].ID)
		assert.Equal(t, "client3", group[1].ID)

		// Removing a missing client should not panic.
		m.RemoveFromGroup("group1", "missing")
		m.RemoveFromGroup("missing_group", "client1")
	})
}
