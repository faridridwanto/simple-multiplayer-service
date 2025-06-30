package websocket

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"simple-multiplayer-service/internal/db/local"
	"simple-multiplayer-service/internal/matchmaking"
	"simple-multiplayer-service/internal/message"
	"simple-multiplayer-service/internal/notification"

	"github.com/gorilla/websocket"
)

// TestSendMessageToClient tests the SendMessageToClient functionality
func TestSendMessageToClient(t *testing.T) {
	notifSvc := notification.NewNotificationService()
	sessionDB := local.DB{}
	mmSvc := matchmaking.NewMatchmakingService(10, sessionDB, notifSvc)
	manager := NewConnectionManager(mmSvc, notifSvc)

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleWebSocket(manager, w, r)
	}))
	defer server.Close()

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// Connect a client
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Could not connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	// Read the welcome message
	var welcomeMsg message.Message
	err = conn.ReadJSON(&welcomeMsg)
	if err != nil {
		t.Fatalf("Error reading welcome message: %v", err)
	}

	// Extract the client ID from the welcome message
	clientID := welcomeMsg.To

	// Verify the client was registered
	time.Sleep(100 * time.Millisecond) // Give some time for the client to be registered
	_, exists := manager.GetClient(clientID)
	if !exists {
		t.Fatalf("Client was not registered properly")
	}

	// Test sending a message to self
	testMsg := message.Message{
		To:      clientID,
		Content: "Test message",
	}

	// Find the client in the manager and send the message
	client, _ := manager.GetClient(clientID)
	client.HandleMessage(testMsg)

	// Read the message back
	var receivedMsg message.Message
	err = conn.ReadJSON(&receivedMsg)
	if err != nil {
		t.Fatalf("Error reading message: %v", err)
	}

	// Verify the message content
	if receivedMsg.Content != "Test message" {
		t.Errorf("Expected message content 'Test message', got '%s'", receivedMsg.Content)
	}
	if receivedMsg.From != clientID {
		t.Errorf("Expected message from '%s', got '%s'", clientID, receivedMsg.From)
	}
}

// TestMultipleClients tests communication between multiple clients
func TestMultipleClients(t *testing.T) {
	notifSvc := notification.NewNotificationService()
	sessionDB := local.DB{}
	mmSvc := matchmaking.NewMatchmakingService(10, sessionDB, notifSvc)
	manager := NewConnectionManager(mmSvc, notifSvc)

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleWebSocket(manager, w, r)
	}))
	defer server.Close()

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// Connect first client
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Could not connect first client: %v", err)
	}
	defer conn1.Close()

	// Read welcome message for first client
	var welcomeMsg1 message.Message
	err = conn1.ReadJSON(&welcomeMsg1)
	if err != nil {
		t.Fatalf("Error reading welcome message for first client: %v", err)
	}
	clientID1 := welcomeMsg1.To

	// Connect second client
	conn2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Could not connect second client: %v", err)
	}
	defer conn2.Close()

	// Read welcome message for second client
	var welcomeMsg2 message.Message
	err = conn2.ReadJSON(&welcomeMsg2)
	if err != nil {
		t.Fatalf("Error reading welcome message for second client: %v", err)
	}
	clientID2 := welcomeMsg2.To

	// Give some time for both clients to be registered
	time.Sleep(100 * time.Millisecond)

	// Send message from client 1 to client 2
	testMsg := message.Message{
		To:      clientID2,
		Content: "Hello from client 1",
	}
	err = conn1.WriteJSON(testMsg)
	if err != nil {
		t.Fatalf("Error sending message from client 1: %v", err)
	}

	// Read the message on client 2
	var receivedMsg message.Message
	err = conn2.ReadJSON(&receivedMsg)
	if err != nil {
		t.Fatalf("Error reading message on client 2: %v", err)
	}

	// Verify the message
	if receivedMsg.Content != "Hello from client 1" {
		t.Errorf("Expected message content 'Hello from client 1', got '%s'", receivedMsg.Content)
	}
	if receivedMsg.From != clientID1 {
		t.Errorf("Expected message from '%s', got '%s'", clientID1, receivedMsg.From)
	}
	if receivedMsg.To != clientID2 {
		t.Errorf("Expected message to '%s', got '%s'", clientID2, receivedMsg.To)
	}

	// Test client disconnection
	conn1.Close()
	time.Sleep(100 * time.Millisecond) // Give some time for the client to be unregistered

	// Verify client 1 was unregistered
	_, exists := manager.GetClient(clientID1)
	if exists {
		t.Error("Client 1 was not properly unregistered after disconnection")
	}
}
