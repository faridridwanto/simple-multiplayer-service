package db

type Session interface {
	CreateSession(sessionID, player1ConnectionID, player2ConnectionID string) error
}
