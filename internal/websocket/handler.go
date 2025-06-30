package websocket

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"simple-multiplayer-service/internal/client"
	"simple-multiplayer-service/internal/message"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Upgrader for WebSocket connections
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for now (in production, this should be restricted)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleWebSocket handles WebSocket connection requests
func HandleWebSocket(manager *ConnectionManager, w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}

	// Get the wsClient's IP address
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	// Create a unique ID for the wsClient using UUID
	clientID := uuid.New().String()

	// Create a new wsClient
	wsClient := &client.Client{
		ID:                  clientID,
		IPAddress:           ip,
		Connection:          conn,
		MatchmakingService:  manager.matchmakingService,
		NotificationService: manager.notificationService,
	}

	// Register the wsClient
	manager.RegisterClient(wsClient)

	// Send the wsClient their ID
	welcomeMsg := message.Message{
		From:    "server",
		To:      clientID,
		Content: fmt.Sprintf("Welcome! Your connection ID is: %s", clientID),
	}
	err = conn.WriteJSON(welcomeMsg)
	if err != nil {
		log.Printf("Error sending welcome message: %v", err)
		conn.Close()
		manager.UnregisterClient(clientID)
		return
	}

	// Start reading messages from the wsClient
	go wsClient.ReadMessages()

	// Start checking notifications from notification service
	go wsClient.CheckNotifications()
}
