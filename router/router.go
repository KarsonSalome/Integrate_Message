package router

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
	"time"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    //Use Cors
    SetupCors(r)

    // User APIs
    SetupUserRoutes(r)
    
    // Friend APIs
    SetupFriendsRouter(r)

    SetupChatRoutes(r)

    // WebSocket routes (already exist)
    //SetupWebSocketRoutes(r)

    return r
}

func SetupCors(r *gin.Engine) {
    r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://172.20.250.37:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}
