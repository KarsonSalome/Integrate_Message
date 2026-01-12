package model

import (
	"time"
	// "gorm.io/gorm"
)

type Message struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SenderID   int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id"` // for one-on-one chat
	Content    string    `json:"content"`
	Timestamp  time.Time `json:"timestamp"`
	Type       string    `json:"type"` // text / file / typing_start / typing_end / open
	State      string    `json:"state"` // sent / delivered / read
}

type HistoryReq struct {
	PeerID int64 `json:"peerID"`
}
