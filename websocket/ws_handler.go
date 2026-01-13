package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "time"
	"aurora-im/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/olahol/melody"
)

var hub = NewHub()
var m = melody.New()

var jwtSecret = []byte("dev-secret")

func parseToken(tokenStr string) (int64, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	return int64(claims["uid"].(float64)), nil
}

type WSMessage struct {
	ReceiverID int64  	 `json:"receiver_id"`
	Content    string 	 `json:"content"`
	Type       string 	 `json:"type"`
}

func WSHandler() http.HandlerFunc {
	// Initialize handlers once
	m.HandleConnect(func(s *melody.Session) {
		tokenStr := s.Request.URL.Query().Get("token")
		if tokenStr == "" {
			s.Close()
			return
		}
		uid, err := parseToken(tokenStr)
		if err != nil {
			s.Close()
			return
		}
		hub.Register(uid, s)
		s.Set("uid", uid)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		if uidAny, ok := s.Get("uid"); ok {
			uid := uidAny.(int64)
			hub.Unregister(uid, s)
		}
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		if uidAny, ok := s.Get("uid"); ok {
			uid := uidAny.(int64)
			var incoming model.Message
			if err := json.Unmarshal(msg, &incoming); err != nil {
				fmt.Println("Invalid message format:", err)
				return
			}
			incoming.SenderID = uid
			// incoming.Timestamp = time.Now()
			// Here you can save the message to DB if needed
			hub.SendMessage(incoming.ReceiverID, incoming)
		}
	})

	return func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	}
}
