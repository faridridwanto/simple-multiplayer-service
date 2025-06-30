package local

import (
	"simple-multiplayer-service/internal/matchmaking"
)

var LocalDB = make(map[string]interface{})

type DB struct {
}

func (l DB) CreateSession(sessionID, player1ConnectionID, player2ConnectionID string) error {
	session := matchmaking.Session{
		SessionID:           sessionID,
		Player1ConnectionID: player2ConnectionID,
		Player2ConnectionID: player1ConnectionID,
	}
	LocalDB[sessionID] = session
	return nil
}
