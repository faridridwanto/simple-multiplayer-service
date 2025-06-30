package notification

type SessionNotification struct {
	SessionID           string `json:"session_id"`
	Player1ConnectionID string `json:"player_1_connection_id"`
	Player2ConnectionID string `json:"player_2_connection_id"`
}

type Service struct {
	Channel chan SessionNotification
}

func NewNotificationService() *Service {
	channel := make(chan SessionNotification, 100)
	return &Service{Channel: channel}
}
