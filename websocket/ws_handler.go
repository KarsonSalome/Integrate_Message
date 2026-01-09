package websocket

import (
	"aurora-im/config"
	"aurora-im/model"
	"aurora-im/dao"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"

	"fmt"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func parseToken(tokenStr string) (int64, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte("dev-secret"), nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)
	uid := int64(claims["uid"].(float64))
	return uid, nil
}

// WSHandler handles WebSocket connections
func WSHandler(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(401, gin.H{"msg": "missing token"})
		return
	}

	uid, err := parseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid token"})
		return
	}

	redisToken, err := config.RedisClient.Get(config.Ctx, "login:"+fmt.Sprint(uid)).Result()
	if err != nil || redisToken != tokenStr {
		c.JSON(401, gin.H{"msg": "token expired"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	fmt.Printf("New User Connected %d\n", uid)

	client := NewClient(uid, conn)
	RegisterClient(client)

	// Write loop: send messages from channel to client
	go func() {
		for msg := range client.Send {
			client.Write(msg)
		}
	}()

	// Read loop: receive messages from client
	for {
		var msg model.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("Error occured while reading msg %s\n", err)
			UnregisterClient(uid)
			conn.Close()
			break
		} else {
			fmt.Printf("New Message Arrived:%s\n", msg)
		}

		// Add timestamp
		msg.Timestamp = time.Now() // here!!!
		msg.SenderID = uid

		// Deliver message
		SendMessage(msg)

		if err := dao.SaveMessageToDB(msg); err != nil {
			fmt.Printf("Error saving message to DB: %s\n", err)
		}
	}
}