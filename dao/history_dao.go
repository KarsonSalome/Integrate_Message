package dao

import (
	"aurora-im/config"
	"aurora-im/model"
	"fmt"

	// "gorm.io/gorm"
)

func SaveMessageToDB(msg model.Message) error {
	return config.DB.Create(&msg).Error
}

func UpdateContactMsg(res model.Message) {
	 var contact model.Contact

	err := config.DB.
		Where("(owner_id = ? AND contact_id = ?) OR (contact_id = ? AND owner_id = ?)", res.SenderID, res.ReceiverID, res.SenderID, res.ReceiverID).
		//Where("owner_id = ? AND contact_id = ?", res.ReceiverID, res.SenderID).
		Find(&contact).Error

	if err != nil {
		fmt.Println("contact find error")
		return
	}

	contact.LastMsg = res.Content
	contact.State = "not_typing"
	contact.LastSenderID = uint(res.SenderID)
	contact.UnreadCount += 1

	if err := config.DB.
		Model(&model.Contact{}).
		Where("id = ?", contact.ID).
		Updates(map[string]interface{}{
			"last_msg":       contact.LastMsg,
			"state":          contact.State,
			"last_sender_id": contact.LastSenderID,
			"unread_count":   contact.UnreadCount,
		}).Error; err != nil {

		fmt.Println("contact update error:", err)
	}
}

func UpdateContactReadMsg(res model.Message) {
	 var contact model.Contact

	err := config.DB.
		Where("(owner_id = ? AND contact_id = ?) OR (contact_id = ? AND owner_id = ?)", res.SenderID, res.ReceiverID, res.SenderID, res.ReceiverID).
		Find(&contact).Error

	if err != nil {
		fmt.Println("contact find error")
		return
	}

	contact.LastMsg = res.Content
	contact.State = "not_typing"
	contact.LastSenderID = uint(res.SenderID)
	contact.UnreadCount = 0

	if err := config.DB.
		Model(&model.Contact{}).
		Where("id = ?", contact.ID).
		Updates(map[string]interface{}{
			"last_msg":       contact.LastMsg,
			"state":          contact.State,
			"last_sender_id": contact.LastSenderID,
			"unread_count":   contact.UnreadCount,
		}).Error; err != nil {

		fmt.Println("contact update error:", err)
	}
}

func LoadHistory(uid, peerID int64) ([]model.Message, error) {
	var messages []model.Message
	err := config.DB.
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", uid, peerID, peerID, uid).
		Order("timestamp ASC").
		Find(&messages).Error
	return messages, err
}

func UpdateReadState(msg model.Message) {
	readRes := config.DB.
		Model(&model.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND state = ?", msg.ReceiverID, msg.SenderID, "delivered").
		Update("state", "read")

	if readRes.Error != nil {
		fmt.Println("msg: Read State update error")
		return
	} else {
		fmt.Println("msg: Read State updated!")
	}

	// result := config.DB.
	// 	Model(&model.Contact{}).
	// 	Where("owner_id = ? AND contact_id = ?", msg.SenderID, msg.ReceiverID).
	// 	Update("unread_count", 0)

	// if result.Error != nil {
	// 	fmt.Println("msg: contact update error")
	// 	return
	// } else {
	// 	fmt.Println("msg: contact updated!")
	// }
}

func UpdateUnreadCount(msg model.Message) {
	result := config.DB.
		Model(&model.Contact{}).
		Where("((owner_id = ? AND contact_id = ?) OR (contact_id = ? AND owner_id = ?)) AND last_sender_id = ?", uint(msg.ReceiverID), uint(msg.SenderID), uint(msg.ReceiverID), uint(msg.SenderID), uint(msg.ReceiverID)).
		Update("unread_count", 0)

	if result.Error != nil {
		fmt.Println("msg: Unread Count update error")
		return
	} else {
		fmt.Println("msg: Unread Count updated!")
	}
}

func UpdateReadCount(msg model.Message) {
	result := config.DB.
		Model(&model.Contact{}).
		Where("(owner_id = ? AND contact_id = ?) OR (contact_id = ? AND owner_id = ?)", uint(msg.ReceiverID), uint(msg.SenderID), uint(msg.ReceiverID), uint(msg.SenderID)).
		Update("unread_count", 0)

	if result.Error != nil {
		fmt.Println("msg: Read Count update error")
		return
	} else {
		fmt.Println("msg: Read Count updated!")
	}
}

func UpdateDeliveredHistory(receiverID int64) {
	result := config.DB.
		Model(&model.Message{}).
		Where("receiver_id = ? AND state = ?", receiverID, "sent").
		Update("state", "delivered")

	if result.Error != nil {
		fmt.Println("msg: Delivered update error")
		return
	} else {
		fmt.Println("msg: Delivered updated!")
	}
}

// func UpdateTypingHistory(res model.Message) {
// 	var contact model.Contact

// 	err := config.DB.
// 		Where("(owner_id = ? AND contact_id = ?) OR (contact_id = ? AND owner_id = ?)", res.SenderID, res.ReceiverID, res.SenderID, res.ReceiverID).
// 		Find(&contact).Error

// 	if err != nil {
// 		fmt.Println("contact find error")
// 		return
// 	}

// 	contact.State = "typing"
// 	contact.LastSenderID = uint(res.SenderID)

// 	result := config.DB.
// 		Model(&model.Contact{}).
// 		Where("owner_id = ? AND contact_id = ?", contact.OwnerID, contact.ContactID).
// 		Select("state", "last_sender_id").
// 		Updates(contact)

// 	if result.Error != nil {
// 		fmt.Println("msg: contact update error")
// 		return
// 	} else {
// 		fmt.Println("msg: contact updated!")
// 	}
// }