package router

import (
	"aurora-im/websocket"
	"aurora-im/middleware"
	"aurora-im/controllers"

	"github.com/gin-gonic/gin"
)

// InitRouter initializes Gin routes
func SetupWebSocketRoutes(r *gin.Engine) {
	r.GET("/ws", websocket.WSHandler)
	ws := r.Group("/chat")
    ws.Use(middleware.Auth())
    ws.POST("/history", controllers.GetHistory)
}
