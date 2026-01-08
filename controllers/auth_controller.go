package controllers

import (
	"aurora-im/utils"
	"aurora-im/model"
	"aurora-im/config"

	"time"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
    var req struct {
        Username string `json:"username" binding:"required"`
        Phone    string `json:"phone" binding:"required"`
        Password string `json:"password" binding:"required,min=6"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"msg": err.Error()})
        return
    }

    hash, _ := utils.HashPassword(req.Password)

    user := model.User{
        Username: req.Username,
        Phone:    req.Phone,
        Password: hash,
        Avatar: "uploads/avatars/default.png",
    }

    if err := config.DB.Create(&user).Error; err != nil {
        c.JSON(400, gin.H{"msg": "user exists"})
        return
    }

    c.JSON(200, gin.H{"msg": "register success"})
}

func Login(c *gin.Context) {
    var req struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password"`
    }

    c.ShouldBindJSON(&req)

    var user model.User
    if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
        c.JSON(401, gin.H{"msg": "invalid credentials"})
        return
    }

    if !utils.CheckPassword(user.Password, req.Password) {
        c.JSON(401, gin.H{"msg": "invalid credentials"})
        return
    }

    token, _ := utils.GenerateToken(user.ID)

    config.RedisClient.Set(
        config.Ctx,
        "login:"+fmt.Sprint(user.ID),
        token,
        24*time.Hour,
    )

    c.JSON(200, gin.H{"token": token})
}

func Logout(c *gin.Context) {
    uid := c.GetString("uid")
    config.RedisClient.Del(config.Ctx, "login:"+uid)
    c.JSON(200, gin.H{"msg": "logout success"})
}

func Me(c *gin.Context) {
    uid := c.GetString("uid")

    var user model.User
    if err := config.DB.First(&user, uid).Error; err != nil {
        c.JSON(404, gin.H{"msg": "not found"})
        return
    }

    c.JSON(200, user)
}
