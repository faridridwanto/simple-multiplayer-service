package matchmaking

import (
	"testing"
	"time"

	"simple-multiplayer-service/internal/message"
	"simple-multiplayer-service/internal/notification"
)

// MockSessionDB is a mock implementation of the db.Session interface
type MockSessionDB struct {
	CreateSessionCalled bool
	SessionID           string
	Player1ID           string
	Player2ID           string
}

func (m *MockSessionDB) CreateSession(sessionID, player1ConnectionID, player2ConnectionID string) error {
	m.CreateSessionCalled = true
	m.SessionID = sessionID
	m.Player1ID = player1ConnectionID
	m.Player2ID = player2ConnectionID
	return nil
}

func TestNewMatchmakingService(t *testing.T) {
	// Setup
	sessionLimit := 10
	sessionDB := &MockSessionDB{}
	notificationService := notification.NewNotificationService()

	// Execute
	service := NewMatchmakingService(sessionLimit, sessionDB, notificationService)

	// Verify
	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}
	if service.SessionLimit != sessionLimit {
		t.Errorf("Expected SessionLimit to be %d, got %d", sessionLimit, service.SessionLimit)
	}
	if service.SessionDB != sessionDB {
		t.Errorf("Expected SessionDB to be set correctly")
	}
	if service.NotificationService != notificationService {
		t.Errorf("Expected NotificationService to be set correctly")
	}
	if service.SessionQueue == nil {
		t.Error("Expected SessionQueue to be initialized")
	}
}

func TestStart(t *testing.T) {
	// Setup
	sessionDB := &MockSessionDB{}
	notificationService := notification.NewNotificationService()
	service := NewMatchmakingService(10, sessionDB, notificationService)

	// Start the service in a goroutine since it has an infinite loop
	go func() {
		service.Start()
	}()

	// Send two matchmaking requests to create a session
	mmr1 := message.MatchmakingRequest{
		ConnectionID: "client1",
		Type:         message.MatchmakingRequestType,
	}
	mmr2 := message.MatchmakingRequest{
		ConnectionID: "client2",
		Type:         message.MatchmakingRequestType,
	}

	// Send the first request
	service.SessionQueue <- mmr1

	// Wait a bit to ensure the first request is processed
	time.Sleep(50 * time.Millisecond)

	// Send the second request
	service.SessionQueue <- mmr2

	// Wait for the session to be created
	time.Sleep(50 * time.Millisecond)

	// Verify
	if !sessionDB.CreateSessionCalled {
		t.Error("Expected CreateSession to be called")
	}
	if sessionDB.Player1ID != "client2" {
		t.Errorf("Expected Player1ID to be 'client2', got '%s'", sessionDB.Player1ID)
	}
	if sessionDB.Player2ID != "client1" {
		t.Errorf("Expected Player2ID to be 'client1', got '%s'", sessionDB.Player2ID)
	}

	// Verify notification was sent
	select {
	case notif := <-notificationService.Channel:
		if notif.Player1ConnectionID != "client2" {
			t.Errorf("Expected Player1ConnectionID to be 'client2', got '%s'", notif.Player1ConnectionID)
		}
		if notif.Player2ConnectionID != "client1" {
			t.Errorf("Expected Player2ConnectionID to be 'client1', got '%s'", notif.Player2ConnectionID)
		}
	default:
		t.Error("Expected notification to be sent")
	}
}
