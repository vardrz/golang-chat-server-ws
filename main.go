package main

import (
	"log"
	"os"

	"chat-server/config"
	"chat-server/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create log file
	f, err := os.OpenFile("chat-server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	// Set log output to file
	log.SetOutput(f)

    db := config.SetupDB()
	config.InitMigration(db)

	r := gin.Default()
	wsHandler := handlers.NewWebSocketHandler(db)

	r.GET("/ws", wsHandler.HandleWebSocket)

	log.Fatal(r.Run(":6000"))
}