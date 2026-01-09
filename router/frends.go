package router

import (
	"aurora-im/middleware"
	"aurora-im/controllers"
	"github.com/gin-gonic/gin"
)

func SetupFriendsRouter(r *gin.Engine) {
	friend := r.Group("/friends")
    friend.Use(middleware.Auth())
    friend.GET("/getFriends", controllers.GetContacts)
    friend.POST("/addFriend", controllers.AddContact)
	friend.POST("/searchFriend", controllers.SearchUser)
}