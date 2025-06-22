package websocket

import (
	"fmt"
	"log"
	"sync"

	"simple-multiplayer-service/pkg/client"
	"simple-multiplayer-service/pkg/message"
)

// ConnectionManager manages all active WebSocket connections
type ConnectionManager struct {
	clients map[string]*client.Client
	mutex   sync.RWMutex
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		clients: make(map[string]*client.Client),
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
