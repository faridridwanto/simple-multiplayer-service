package main

import (
	"log"
	"net/http"

	"simple-multiplayer-service/internal/config"
	"simple-multiplayer-service/internal/db/local"
	"simple-multiplayer-service/internal/matchmaking"
	"simple-multiplayer-service/internal/notification"
	"simple-multiplayer-service/internal/websocket"

	"github.com/caarlos0/env/v11"
)

func main() {
	// Parse env variables
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Create a Notification Service
	notificationService := notification.NewNotificationService()

	// Create a DB client
	localDB := local.DB{}

	// Create a Matchmaking Service
	matchmakingService := matchmaking.NewMatchmakingService(cfg.SessionLimit, localDB, notificationService)

	// Create a new connection manager
	manager := websocket.NewConnectionManager(matchmakingService, notificationService)

	// Start the matchmaking service
	go manager.StartMatchmakingService()

	// Set up the WebSocket endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleWebSocket(manager, w, r)
	})

	// Start the server
	port := ":8080"
	log.Printf("Starting WebSocket server on %s", port)
	log.Printf("Connect to ws://localhost%s/ws", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
