package websocket

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"github.com/olahol/melody"
	"aurora-im/model"
	"aurora-im/dao"
)


type Broadcast struct {
	From    	int64  		`json:"from"`
	Content 	string 		`json:"content"`
	Type    	string 		`json:"type"`
	TimeAt	    time.Time 	`json:"time_at"`
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

	// Send offline messages
	offlineMessages, err := dao.GetOfflineMessages(uid)
	if err != nil {
		fmt.Println("Error fetching offline messages:", err)
		return
	}
	for _, msg := range offlineMessages {
		data, _ := json.Marshal(msg)
		s.Write(data)
		fmt.Println("Sent offline message to", uid, ":", msg)
	}
	dao.UpdateDeliveredHistory(uid)
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

func (h *Hub) SendMessage(toUID int64, msg model.Message) {
	h.Lock.RLock()
	defer h.Lock.RUnlock()
	if sessions, ok := h.UserSessions[toUID]; ok {
		msg.State = "delivered"
		
		switch msg.Type {
		case "message":
			dao.SaveMessageToDB(msg)
			dao.UpdateContactMsg(msg)
		case "read":
			dao.UpdateReadHistory(msg)
		}

		data, _ := json.Marshal(msg)
		for s := range sessions {
			s.Write(data)
		}
	} else {
		msg.State = "sent"
		switch msg.Type {
		case "message":
			dao.SaveMessageToDB(msg)
			dao.UpdateContactMsg(msg)
		case "read":
			dao.UpdateReadHistory(msg)
		}
		fmt.Println("User offline:", toUID)
	}
}
