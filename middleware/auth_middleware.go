package middleware

import (
    "aurora-im/config"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"

	"fmt"
)

func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenStr := c.GetHeader("Authorization")

        token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
            return []byte("dev-secret"), nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(401, gin.H{"msg": "unauthorized"})
            return
        }

        claims := token.Claims.(jwt.MapClaims)
        uid := fmt.Sprint(claims["uid"])

        val, _ := config.RedisClient.Get(config.Ctx, "login:"+uid).Result()
        if val != tokenStr {
            c.AbortWithStatusJSON(401, gin.H{"msg": "expired"})
            return
        }

        c.Set("uid", uid)
        c.Next()
    }
}
