package message

const MatchmakingRequestType = "matchmakingRequest"

// Message represents a message sent between clients
type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

type MatchmakingRequest struct {
	ConnectionID string `json:"connection_id"`
	Type         string `json:"type"`
}
