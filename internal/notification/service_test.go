package notification

import (
	"testing"
)

func TestNewNotificationService(t *testing.T) {
	// Execute
	service := NewNotificationService()

	// Verify
	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}
	if service.Channel == nil {
		t.Error("Expected Channel to be initialized")
	}

	// Test that the channel is buffered with capacity 100
	// We can send 100 messages without blocking
	for i := 0; i < 100; i++ {
		service.Channel <- SessionNotification{
			SessionID:           "test",
			Player1ConnectionID: "player1",
			Player2ConnectionID: "player2",
		}
	}

	// Verify we can receive all 100 messages
	for i := 0; i < 100; i++ {
		select {
		case <-service.Channel:
			// Successfully received message
		default:
			t.Errorf("Expected to receive message %d, but channel was empty", i)
		}
	}
}
