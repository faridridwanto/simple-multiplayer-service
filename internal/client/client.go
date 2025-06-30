package client

import (
	"encoding/json"
	"log"

	"simple-multiplayer-service/internal/matchmaking"
	"simple-multiplayer-service/internal/message"
	"simple-multiplayer-service/internal/notification"

	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection
type Client struct {
	ID                  string
	IPAddress           string
	Connection          *websocket.Conn
	MatchmakingService  *matchmaking.Service
	NotificationService *notification.Service
	SendMessageFunc     func(message message.Message) error
	UnregisterFunc      func(clientID string)
}

// HandleMessage processes incoming messages from clients
func (c *Client) HandleMessage(message message.Message) {
	// Set the sender ID
	message.From = c.ID

	// Log the message
	log.Printf("Message from %s to %s: %s", message.From, message.To, message.Content)

	// Send the message to the target client
	err := c.SendMessageFunc(message)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// HandleMatchmakingRequest processes incoming matchmaking requests from clients
func (c *Client) HandleMatchmakingRequest(mmr message.MatchmakingRequest) {
	// put matchmaking request into the queue
	c.MatchmakingService.SessionQueue <- mmr
}

// ReadMessages continuously reads messages from the client
func (c *Client) ReadMessages() {
	defer func() {
		c.Connection.Close()
		c.UnregisterFunc(c.ID)
	}()

	for {
		var msg message.Message
		err := c.Connection.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		log.Printf("Message from %s: %s", msg.From, msg.Content)

		// Handling Matchmaking Request
		var mmr message.MatchmakingRequest
		err = json.Unmarshal([]byte(msg.Content), &mmr)
		if err == nil && mmr.Type == message.MatchmakingRequestType {
			c.HandleMatchmakingRequest(mmr)
			continue
		}

		c.HandleMessage(msg)
	}
}

func (c *Client) HandleSessionCreatedNotification(session notification.SessionNotification, targetUserID1, targetUserID2 string) error {
	sessionByte, err := json.Marshal(session)
	if err != nil {
		log.Printf("Error marshalling session: %v", err)
		return err
	}
	notificationMessage1 := message.Message{
		From:    "Server",
		To:      targetUserID1,
		Content: string(sessionByte),
	}

	log.Printf("Sending message to %s", notificationMessage1.To)
	err = c.SendMessageFunc(notificationMessage1)
	if err != nil {
		log.Printf("Error sending notificationMessage: %v", err)
		return err
	}

	notificationMessage2 := message.Message{
		From:    "Server",
		To:      targetUserID2,
		Content: string(sessionByte),
	}

	log.Printf("Sending message to %s", notificationMessage2.To)
	err = c.SendMessageFunc(notificationMessage2)
	if err != nil {
		log.Printf("Error sending notificationMessage: %v", err)
		return err
	}

	return nil
}

// CheckNotifications check notifications from the service's notification channel
func (c *Client) CheckNotifications() {
	for {
		select {
		case notif := <-c.NotificationService.Channel:
			if notif.Player1ConnectionID == c.ID || notif.Player2ConnectionID == c.ID {
				log.Println("handling session created notification")
				err := c.HandleSessionCreatedNotification(notif, notif.Player1ConnectionID, notif.Player2ConnectionID)
				if err != nil {
					log.Printf("Error handling session created notification: %v", err)
					continue
				}
			}
		}
	}
}
