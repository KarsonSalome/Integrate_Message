package dao

import (
	"aurora-im/config"
	"aurora-im/model"
	"fmt"
)

func SaveMessageToDB(msg model.Message) error {
	return config.DB.Create(&msg).Error
}

func UpdateContactMsg(res model.Message) {
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

	result := config.DB.
		Model(&model.Contact{}).
		Where("owner_id = ? AND contact_id = ?", contact.OwnerID, contact.ContactID).
		Select("last_msg", "state", "last_sender_id").
		Updates(contact)

	if result.Error != nil {
		fmt.Println("msg: contact update error")
		return
	} else {
		fmt.Println("msg: contact updated!")
	}
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

func UpdateReadHistory(msg model.Message) {
	result := config.DB.
		Model(&model.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND readable = ?", msg.ReceiverID, msg.SenderID, "unread").
		Update("readable", "read")

	if result.Error != nil {
		fmt.Println("msg: contact update error")
		return
	} else {
		fmt.Println("msg: contact updated!")
	}
}

func UpdateTypingHistory(res model.Message) {
	var contact model.Contact

	err := config.DB.
		Where("(owner_id = ? AND contact_id = ?) OR (contact_id = ? AND owner_id = ?)", res.SenderID, res.ReceiverID, res.SenderID, res.ReceiverID).
		Find(&contact).Error

	if err != nil {
		fmt.Println("contact find error")
		return
	}

	contact.State = "typing"
	contact.LastSenderID = uint(res.SenderID)

	result := config.DB.
		Model(&model.Contact{}).
		Where("owner_id = ? AND contact_id = ?", contact.OwnerID, contact.ContactID).
		Select("state", "last_sender_id").
		Updates(contact)

	if result.Error != nil {
		fmt.Println("msg: contact update error")
		return
	} else {
		fmt.Println("msg: contact updated!")
	}
}