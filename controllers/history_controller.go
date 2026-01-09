package controllers

import (
	"github.com/gin-gonic/gin"

	"aurora-im/dao"
	"aurora-im/model"

	"strconv"
)

func GetHistory(c *gin.Context) {
	var req model.HistoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	uid := c.GetString("uid") // From JWT middleware
	UserID, err := strconv.ParseInt(uid, 10, 64)

	if err != nil {
		c.JSON(400, gin.H{"error": "invalid uid"})
		return
	}

	PeerID := req.PeerID

	history, _ := dao.LoadHistory(UserID, PeerID)
	for _, msg := range history {
		c.JSON(200, gin.H{
			"history": msg,
		})
	}
}