package controllers

import (
	"github.com/gin-gonic/gin"

	"aurora-im/config"
	"aurora-im/model"

	"strconv"
)

func AddContact(c *gin.Context) {
	var req model.AddContactReq
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
	config.DB.Preload("Contact").
		Where("owner_id = ?", uint(ownerID)).
		Find(&contacts)

	var res []gin.H
	for _, c := range contacts {
		res = append(res, gin.H{
			"id":       c.Contact.ID,
			"username": c.Contact.Username,
			"avatar":   c.Contact.Avatar,
			"lastmsg":  c.LastMsg,
		})
	}

	c.JSON(200, res)
}

func SearchUser(c *gin.Context) {
	var req model.SearchReq
	var user model.User
	var res model.SearchRes

	c.ShouldBindJSON(&req)

	if err := config.DB.Where("phone = ?", req.Phone).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"msg": "invalid credentials"})
		return
	}

	res.UserID = user.ID
	res.Username = user.Username
	res.Phone = user.Phone
	res.Avatar = user.Avatar

	c.JSON(200, gin.H{
		"Search": res,
	})
}
