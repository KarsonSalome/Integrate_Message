package model

import "time"

// Message represents a chat message in Tiny IM
type Message struct {
	SenderID   int64     `json:"sender_id"`   // Who sends
	ReceiverID int64     `json:"receiver_id"` // Who receives
	Content    string    `json:"content"`     // Message body
	Type       string    `json:"type"`        // text / file / typing_start / typing_end / open 
	Timestamp  time.Time `json:"timestamp"`   // Sending time
}
