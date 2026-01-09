package dao

import (
	"aurora-im/config"
	"aurora-im/model"
)

func SaveMessageToDB(msg model.Message) error {
	return config.DB.Create(&msg).Error
}

func LoadHistory(uid, peerID int64) ([]model.Message, error) {
	var messages []model.Message
	err := config.DB.
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", uid, peerID, peerID, uid).
		Order("timestamp ASC").
		Limit(100).
		Find(&messages).Error
	return messages, err
}