package main

import (
	"aurora-im/router"
	//"aurora-im/websocket"
	"aurora-im/config"
	"net/http"
	"fmt"

	//"github.com/gin-gonic/gin"
)

func main() {

	go func() {
        router.RegisterSocketRoutes()

	// Example REST endpoint
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"msg":"Hello server"}`))
	})

	fmt.Println("Server running on :8081")
	http.ListenAndServe(":8081", nil)
    }()

	//Iinit Database
	config.InitDB()

	// Initialize Hub
	//websocket.InitHub()

	// Initialize Router
	r := router.SetupRouter()

	// Run server
	r.Run(":8080") // http://localhost:8080

	
	
}
