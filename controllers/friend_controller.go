package controllers

import (
	"github.com/gin-gonic/gin"

	"aurora-im/config"
	"aurora-im/model"

	"strconv"
	"errors"
	"gorm.io/gorm"
)

func AddContact(c *gin.Context) {
	var req model.ContactReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("uid") // From JWT middleware
	ownerID, err := strconv.ParseUint(userID, 10, 64)

	if err != nil {
		c.JSON(400, gin.H{"error": "invalid uid"})
		return
	}

	contact := model.Contact{
		OwnerID:   uint(ownerID),
		ContactID: req.ContactID,
		LastMsg:   "",
		State: "not_typing",
		LastSenderID: 0,
		UnreadCount: 0,
	}

	var oldContact model.Contact

	dberro := config.DB.Where("(owner_id = ? AND contact_id = ?) OR (contact_id = ? AND owner_id = ?)", req.ContactID, uint(ownerID), req.ContactID, uint(ownerID)).First(&oldContact).Error
	if dberro == nil || req.ContactID == uint(ownerID) {
    	c.JSON(400, gin.H{"error": "contact already exists"})
    	return
	}

	if !errors.Is(dberro, gorm.ErrRecordNotFound) {
    	c.JSON(500, gin.H{"error": err.Error()})
    	return
	}

	if err := config.DB.Create(&contact).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"msg": "contact added"})
}

func GetContacts(c *gin.Context) {
	userID := c.GetString("uid") // From JWT middleware
	ownerID, err := strconv.ParseUint(userID, 10, 64)

	if err != nil {
		c.JSON(400, gin.H{"error": "invalid uid"})
		return
	}

	var contacts []model.Contact
	config.DB.
		Preload("Contact").
		Preload("Owner").
		Where("owner_id = ? OR contact_id = ?", uint(ownerID), uint(ownerID)).
		Find(&contacts)

	var res []gin.H
	for _, c := range contacts {
		if uint(ownerID) == c.OwnerID {
			res = append(res, gin.H{
				"contact_id": c.ContactID,
				"username":   c.Contact.Username,
				"phone":      c.Contact.Phone,
				"avatar":     c.Contact.Avatar,
				"last_msg":   c.LastMsg,
				"state":	  c.State,
				"last_sender":c.LastSenderID,
				"last_time":   c.UpdatedAt,
				"unread_count": c.UnreadCount,
			})
		} else if uint(ownerID) == c.ContactID {
			res = append(res, gin.H{
				"contact_id": c.OwnerID,
				"username":   c.Owner.Username,
				"phone":      c.Owner.Phone,
				"avatar":     c.Owner.Avatar,
				"last_msg":    c.LastMsg,
				"state":	  c.State,
				"last_sender":c.LastSenderID,
				"last_time":   c.UpdatedAt,
				"unread_count": c.UnreadCount,
			})
		}
		
	}

	c.JSON(200, res)
}

func SearchUser(c *gin.Context) {
	var req model.SearchReq
	var user model.User
	var res model.SearchRes

	c.ShouldBindJSON(&req)

	if err := config.DB.Where("phone = ?", req.Phone).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"msg": "invalid credentials"})
		return
	}

	res.UserID = user.ID
	res.Username = user.Username
	res.Phone = user.Phone
	res.Avatar = user.Avatar

	c.JSON(200, res)
}

func DeleteContact(c *gin.Context) {
	var req model.ContactReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("uid") // From JWT middleware
	ownerID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid uid"})
		return
	}
	result := config.DB.
		Where("(owner_id = ? AND contact_id = ?) OR (contact_id = ? AND owner_id = ?)", uint(ownerID), req.ContactID, uint(ownerID), req.ContactID).
		Delete(&model.Contact{})
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "contact deleted"})
}
