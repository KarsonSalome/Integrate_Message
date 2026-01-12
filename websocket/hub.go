package websocket

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/olahol/melody"
)


type Broadcast struct {
	From    int64  `json:"from"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

type Hub struct {
	UserSessions map[int64]map[*melody.Session]bool
	Lock         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		UserSessions: make(map[int64]map[*melody.Session]bool),
	}
}

func (h *Hub) Register(uid int64, s *melody.Session) {
	h.Lock.Lock()
	
	if _, ok := h.UserSessions[uid]; !ok {
		h.UserSessions[uid] = make(map[*melody.Session]bool)
	}
	h.UserSessions[uid][s] = true
	fmt.Println("Client registered:", uid)

	h.Lock.Unlock()
}

func (h *Hub) Unregister(uid int64, s *melody.Session) {
	h.Lock.Lock()
	defer h.Lock.Unlock()
	if sessions, ok := h.UserSessions[uid]; ok {
		delete(sessions, s)
		if len(sessions) == 0 {
			delete(h.UserSessions, uid)
			fmt.Println("Client unregistered:", uid)
		}
	}
}

func (h *Hub) SendMessage(toUID int64, msg Broadcast) {
	h.Lock.RLock()
	defer h.Lock.RUnlock()
	if sessions, ok := h.UserSessions[toUID]; ok {
		data, _ := json.Marshal(msg)
		for s := range sessions {
			s.Write(data)
		}
	} else {
		fmt.Println("User offline:", toUID)
	}
}
