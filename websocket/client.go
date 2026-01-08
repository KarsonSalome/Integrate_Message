package websocket

import (
	"sync"
	"aurora-im/model"
	"github.com/gorilla/websocket"
)

// Client represents a connected user
type Client struct {
	UID  int64
	Conn *websocket.Conn
	Send chan model.Message
	Lock sync.Mutex
}

// NewClient creates a new client
func NewClient(uid int64, conn *websocket.Conn) *Client {
	return &Client{
		UID:  uid,
		Conn: conn,
		Send: make(chan model.Message, 256), // buffered channel
	}
}

// Write safely sends message to WebSocket
func (c *Client) Write(msg model.Message) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.Conn.WriteJSON(msg)
}
