package client

import (
	"encoding/json"
	"testing"

	"simple-multiplayer-service/internal/matchmaking"
	"simple-multiplayer-service/internal/message"
	"simple-multiplayer-service/internal/notification"
)

func TestHandleMessage(t *testing.T) {
	// Setup
	sendMessageCalled := false
	var sentMessage message.Message

	client := &Client{
		ID: "client1",
		SendMessageFunc: func(msg message.Message) error {
			sendMessageCalled = true
			sentMessage = msg
			return nil
		},
	}

	// Test message
	msg := message.Message{
		To:      "client2",
		Content: "Hello",
	}

	// Execute
	client.HandleMessage(msg)

	// Verify
	if !sendMessageCalled {
		t.Error("SendMessageFunc was not called")
	}
	if sentMessage.From != "client1" {
		t.Errorf("Expected From to be 'client1', got '%s'", sentMessage.From)
	}
	if sentMessage.To != "client2" {
		t.Errorf("Expected To to be 'client2', got '%s'", sentMessage.To)
	}
	if sentMessage.Content != "Hello" {
		t.Errorf("Expected Content to be 'Hello', got '%s'", sentMessage.Content)
	}
}

func TestHandleMatchmakingRequest(t *testing.T) {
	// Setup
	sessionQueue := make(chan message.MatchmakingRequest, 1)
	matchmakingService := &matchmaking.Service{
		SessionQueue: sessionQueue,
	}

	client := &Client{
		ID:                 "client1",
		MatchmakingService: matchmakingService,
	}

	// Test request
	mmr := message.MatchmakingRequest{
		ConnectionID: "client1",
		Type:         message.MatchmakingRequestType,
	}

	// Execute
	client.HandleMatchmakingRequest(mmr)

	// Verify
	select {
	case receivedMMR := <-sessionQueue:
		if receivedMMR.ConnectionID != "client1" {
			t.Errorf("Expected ConnectionID to be 'client1', got '%s'", receivedMMR.ConnectionID)
		}
		if receivedMMR.Type != message.MatchmakingRequestType {
			t.Errorf("Expected Type to be '%s', got '%s'", message.MatchmakingRequestType, receivedMMR.Type)
		}
	default:
		t.Error("No matchmaking request was sent to the queue")
	}
}

func TestHandleSessionCreatedNotification(t *testing.T) {
	// Setup
	sentMessages := make([]message.Message, 0)
	client := &Client{
		ID: "server",
		SendMessageFunc: func(msg message.Message) error {
			sentMessages = append(sentMessages, msg)
			return nil
		},
	}

	// Test notification
	session := notification.SessionNotification{
		SessionID:           "session1",
		Player1ConnectionID: "client1",
		Player2ConnectionID: "client2",
	}

	// Execute
	err := client.HandleSessionCreatedNotification(session, "client1", "client2")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(sentMessages) != 2 {
		t.Errorf("Expected 2 messages to be sent, got %d", len(sentMessages))
	}

	// Check first message
	if sentMessages[0].From != "Server" {
		t.Errorf("Expected From to be 'Server', got '%s'", sentMessages[0].From)
	}
	if sentMessages[0].To != "client1" {
		t.Errorf("Expected To to be 'client1', got '%s'", sentMessages[0].To)
	}

	// Check second message
	if sentMessages[1].From != "Server" {
		t.Errorf("Expected From to be 'Server', got '%s'", sentMessages[1].From)
	}
	if sentMessages[1].To != "client2" {
		t.Errorf("Expected To to be 'client2', got '%s'", sentMessages[1].To)
	}

	// Verify content contains session data
	var sessionData1 notification.SessionNotification
	err = json.Unmarshal([]byte(sentMessages[0].Content), &sessionData1)
	if err != nil {
		t.Errorf("Failed to unmarshal message content: %v", err)
	}
	if sessionData1.SessionID != "session1" {
		t.Errorf("Expected SessionID to be 'session1', got '%s'", sessionData1.SessionID)
	}
}

// Note: ReadMessages is not tested here because it relies heavily on the websocket.Conn interface,
// which is difficult to mock effectively. In a real-world scenario, you might use a library like
// github.com/stretchr/testify/mock to create a proper mock for websocket.Conn.

// Note: CheckNotifications is not tested here because it contains an infinite loop that
// continuously reads from a channel, making it difficult to test in a unit test context.
// In a real-world scenario, you might refactor the function to accept a context for cancellation
// or use a testing framework that supports testing goroutines with infinite loops.
