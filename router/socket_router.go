package router

import (
	"aurora-im/websocket"
	"aurora-im/middleware"
	"aurora-im/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

// InitRouter initializes Gin routes
func RegisterSocketRoutes() {
	http.HandleFunc("/ws", websocket.WSHandler())
}

func SetupChatRoutes(r *gin.Engine) {
	ws := r.Group("/chat")
    ws.Use(middleware.Auth())
    ws.POST("/history", controllers.GetHistory)
}