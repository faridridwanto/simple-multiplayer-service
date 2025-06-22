package main

import (
	"log"
	"net/http"

	"simple-multiplayer-service/internal/websocket"
)

func main() {
	// Create a new connection manager
	manager := websocket.NewConnectionManager()

	// Set up the WebSocket endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleWebSocket(manager, w, r)
	})

	// Start the server
	port := ":8080"
	log.Printf("Starting WebSocket server on %s", port)
	log.Printf("Connect to ws://localhost%s/ws", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}