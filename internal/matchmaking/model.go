package matchmaking

type Session struct {
	SessionID           string `json:"sessionId"`
	Player1ConnectionID string `json:"player1ConnectionId"`
	Player2ConnectionID string `json:"player2ConnectionId"`
}
