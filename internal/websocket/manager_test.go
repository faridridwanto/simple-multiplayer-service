package websocket

import (
	"testing"

	"simple-multiplayer-service/pkg/client"
)

// TestConnectionManager tests the ConnectionManager functionality
func TestConnectionManager(t *testing.T) {
	manager := NewConnectionManager()

	// Create a mock client
	client := &client.Client{
		ID:        "test-client",
		IPAddress: "127.0.0.1",
	}

	// Test RegisterClient
	manager.RegisterClient(client)
	if len(manager.clients) != 1 {
		t.Errorf("Expected 1 client, got %d", len(manager.clients))
	}

	// Test GetClient
	retrievedClient, exists := manager.GetClient("test-client")
	if !exists {
		t.Error("Expected client to exist, but it doesn't")
	}
	if retrievedClient != client {
		t.Error("Retrieved client is not the same as the registered client")
	}

	// Test GetClient for non-existent client
	_, exists = manager.GetClient("non-existent")
	if exists {
		t.Error("Expected client to not exist, but it does")
	}

	// Test UnregisterClient
	manager.UnregisterClient("test-client")
	if len(manager.clients) != 0 {
		t.Errorf("Expected 0 clients, got %d", len(manager.clients))
	}
}
