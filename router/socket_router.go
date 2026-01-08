package router

import (
	"aurora-im/websocket"
	"github.com/gin-gonic/gin"
)

// InitRouter initializes Gin routes
func SetupWebSocketRoutes(r *gin.Engine) {
	r.GET("/ws", websocket.WSHandler)
}
