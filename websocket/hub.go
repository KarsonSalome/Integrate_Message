package websocket

import (
	"encoding/json"
	"fmt"
	"sync"
	// "time"
	"github.com/olahol/melody"
	"aurora-im/model"
	"aurora-im/dao"
)

var (
	selectMap  = make(map[int64]int64)
	selectLock sync.RWMutex
)

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
		if err := s.Write(data); err != nil {
			fmt.Println("Offline send failed:", err)
			return
		}
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

			selectLock.Lock()
			delete(selectMap, uid)
			selectLock.Unlock()

			fmt.Println("Client unregistered:", uid)
		}
	}
}

func isReceiverViewingSender(msg model.Message) bool {
	selectLock.RLock()
	defer selectLock.RUnlock()

	fmt.Println(selectMap)

	return selectMap[msg.ReceiverID] == msg.SenderID
}

func (h *Hub) SendMessage(toUID int64, msg model.Message) {

	h.Lock.RLock()
	sessions, online := h.UserSessions[toUID]
	returnSession, exists := h.UserSessions[msg.SenderID]
	h.Lock.RUnlock()

	/* ---------- READ EVENT ---------- */
	if msg.Type == "read" {
		selectLock.Lock()
		selectMap[msg.SenderID] = toUID
		selectLock.Unlock()

		dao.UpdateReadState(msg)
		dao.UpdateUnreadCount(msg)
		return
	}

	/* ---------- MESSAGE EVENT ---------- */
	if msg.Type == "message" {
		isActive := isReceiverViewingSender(msg)
		fmt.Println("isActive:", isActive)
		if online && isActive{
			msg.State = "read"
			dao.SaveMessageToDB(msg)
			dao.UpdateContactReadMsg(msg)

		} else if online {
			msg.State = "delivered"
			dao.SaveMessageToDB(msg)
			dao.UpdateContactMsg(msg)

		} else {
			msg.State = "sent"
			dao.SaveMessageToDB(msg)
			dao.UpdateContactMsg(msg)
		}		
	}

	/* ---------- PUSH TO CLIENT ---------- */
	if online {
		data, _ := json.Marshal(msg)
		for s := range sessions {
			s.Write(data)
		}
	} else {
		fmt.Println("User offline:", toUID)
	}

	if exists {
		data, _ := json.Marshal(msg)
		for s := range returnSession {
			s.Write(data)
		}
	}
}