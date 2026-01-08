package router

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
    r := gin.Default()

    // User APIs
    SetupUserRoutes(r)

    // WebSocket routes (already exist)
    SetupWebSocketRoutes(r)

    return r
}
