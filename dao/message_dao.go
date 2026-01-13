package dao

import (
	"aurora-im/config"
	"aurora-im/model"
)

// Get all offline messages
func GetOfflineMessages(receiverID int64) ([]model.Message, error) {
	var messages []model.Message
	err := config.DB.
		Where("receiver_id = ? AND state = ?", receiverID, "sent").
		Order("timestamp ASC").
		Find(&messages).Error
	return messages, err
}
