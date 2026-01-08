package main

import (
	"aurora-im/router"
	"aurora-im/websocket"
	"aurora-im/config"
)

func main() {
	//Iinit Database
	config.InitDB()

	// Init Redis
	config.InitRedis()

	// Initialize Hub
	websocket.InitHub()

	// Initialize Router
	r := router.SetupRouter()

	// Run server
	r.Run(":8080") // http://localhost:8080
}
