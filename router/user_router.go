package router

import (
    "aurora-im/controllers"
    "aurora-im/middleware"
    
    "github.com/gin-gonic/gin"
)

// Register user routes
func SetupUserRoutes(r *gin.Engine) {
    r.Static("/uploads", "./uploads")

    api := r.Group("/api")
    api.POST("/register", controllers.Register)
    api.POST("/login", controllers.Login)
    api.POST("/logout", controllers.Logout)

    user := r.Group("/user")
    user.Use(middleware.Auth())
    user.GET("/profile", controllers.Me)
    user.POST("/updateProfile", controllers.UpdateProfile)
    user.POST("/uploadAvatar", controllers.UploadAvatar)
}
