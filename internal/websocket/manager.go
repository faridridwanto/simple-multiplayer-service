package websocket

import (
	"fmt"
	"log"
	"sync"

	"simple-multiplayer-service/internal/client"
	"simple-multiplayer-service/internal/matchmaking"
	"simple-multiplayer-service/internal/message"
	"simple-multiplayer-service/internal/notification"
)

// ConnectionManager manages all active WebSocket connections
type ConnectionManager struct {
	clients             map[string]*client.Client
	mutex               sync.RWMutex
	matchmakingService  *matchmaking.Service
	notificationService *notification.Service
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager(mmSvc *matchmaking.Service, notifSvc *notification.Service) *ConnectionManager {
	return &ConnectionManager{
		clients:             make(map[string]*client.Client),
		matchmakingService:  mmSvc,
		notificationService: notifSvc,
	}
}

// RegisterClient adds a new client to the manager
func (cm *ConnectionManager) RegisterClient(client *client.Client) {
	// Set the client's SendMessageFunc and UnregisterFunc
	client.SendMessageFunc = cm.SendMessageToClient
	client.UnregisterFunc = cm.UnregisterClient

	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.clients[client.ID] = client
	log.Printf("Client registered: %s (IP: %s)", client.ID, client.IPAddress)
}

// UnregisterClient removes a client from the manager
func (cm *ConnectionManager) UnregisterClient(clientID string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	if _, exists := cm.clients[clientID]; exists {
		delete(cm.clients, clientID)
		log.Printf("Client unregistered: %s", clientID)
		cm.matchmakingService.ClientDisconnects <- clientID
	}
}

// GetClient retrieves a client by ID
func (cm *ConnectionManager) GetClient(clientID string) (*client.Client, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	client, exists := cm.clients[clientID]
	return client, exists
}

// SendMessageToClient sends a message to a specific client
func (cm *ConnectionManager) SendMessageToClient(message message.Message) error {
	targetClient, exists := cm.GetClient(message.To)
	if !exists {
		return fmt.Errorf("client with ID %s not found", message.To)
	}

	return targetClient.Connection.WriteJSON(message)
}

// StartMatchmakingService start the matchmaking service
func (cm *ConnectionManager) StartMatchmakingService() {
	cm.matchmakingService.Start()
}
