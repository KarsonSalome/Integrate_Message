package dao

import (
	"aurora-im/config"
	"aurora-im/model"
	"encoding/json"
)

// Redis key: "offline:<receiver_id>"

// Store message to Redis
func PushOfflineMessage(msg model.Message) error {
	key := "offline:" + string(rune(msg.ReceiverID))
	data, _ := json.Marshal(msg)
	return config.RedisClient.RPush(config.Ctx, key, data).Err()
}

// Get all offline messages
func GetOfflineMessages(receiverID int64) ([]model.Message, error) {
	key := "offline:" + string(rune(receiverID))
	list, err := config.RedisClient.LRange(config.Ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var msgs []model.Message
	for _, item := range list {
		var msg model.Message
		json.Unmarshal([]byte(item), &msg)
		msgs = append(msgs, msg)
	}

	// Delete after fetching
	config.RedisClient.Del(config.Ctx, key)

	return msgs, nil
}
