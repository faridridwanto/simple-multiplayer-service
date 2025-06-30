package local

import (
	"testing"

	"simple-multiplayer-service/internal/matchmaking"
)

func TestCreateSession(t *testing.T) {
	// Clear the LocalDB before the test
	for k := range LocalDB {
		delete(LocalDB, k)
	}

	// Setup
	db := DB{}
	sessionID := "test-session"
	player1ID := "player1"
	player2ID := "player2"

	// Execute
	err := db.CreateSession(sessionID, player1ID, player2ID)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check that the session was stored in the LocalDB
	session, ok := LocalDB[sessionID]
	if !ok {
		t.Errorf("Expected session with ID %s to be stored in LocalDB", sessionID)
	}

	// Check that the session has the correct values
	sessionObj, ok := session.(matchmaking.Session)
	if !ok {
		t.Errorf("Expected session to be of type matchmaking.Session, got %T", session)
	}

	// Note: The implementation in matchmaking.go swaps player1 and player2
	if sessionObj.SessionID != sessionID {
		t.Errorf("Expected SessionID to be %s, got %s", sessionID, sessionObj.SessionID)
	}
	if sessionObj.Player1ConnectionID != player2ID {
		t.Errorf("Expected Player1ConnectionID to be %s, got %s", player2ID, sessionObj.Player1ConnectionID)
	}
	if sessionObj.Player2ConnectionID != player1ID {
		t.Errorf("Expected Player2ConnectionID to be %s, got %s", player1ID, sessionObj.Player2ConnectionID)
	}
}
