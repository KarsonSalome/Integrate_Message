package router

import (
    "aurora-im/controllers"
    "aurora-im/middleware"
    
    "github.com/gin-gonic/gin"
)

// Register user routes
func SetupUserRoutes(r *gin.Engine) {
    r.POST("/register", controllers.Register)
    r.POST("/login", controllers.Login)
    r.POST("/logout", controllers.Logout)

    auth := r.Group("/")
    auth.Use(middleware.Auth())
    auth.GET("/profile", controllers.Me)
    auth.POST("/updateProfile", controllers.UpdateProfile)
}
