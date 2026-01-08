package websocket

import (
	"aurora-im/dao"
	"aurora-im/model"
	"sync"
)

type Hub struct {
	Clients map[int64]*Client
	Lock    sync.RWMutex
}

var hub *Hub

// Initialize hub
func InitHub() {
	hub = &Hub{
		Clients: make(map[int64]*Client),
	}
}

// RegisterClient adds a client and sends offline messages
func RegisterClient(client *Client) {
	hub.Lock.Lock()
	defer hub.Lock.Unlock()
	hub.Clients[client.UID] = client

	// Fetch offline messages asynchronously
	go func() {
		msgs, _ := dao.GetOfflineMessages(client.UID)
		for _, msg := range msgs {
			client.Send <- msg
		}
	}()
}

// UnregisterClient removes client (EXPORT with uppercase!)
func UnregisterClient(uid int64) {
	hub.Lock.Lock()
	defer hub.Lock.Unlock()
	delete(hub.Clients, uid)
}

// SendMessage delivers message or stores in Redis if offline
func SendMessage(msg model.Message) {
	hub.Lock.RLock()
	receiver, online := hub.Clients[msg.ReceiverID]
	hub.Lock.RUnlock()

	if online {
		receiver.Send <- msg
	} else {
		_ = dao.PushOfflineMessage(msg)
	}
}
