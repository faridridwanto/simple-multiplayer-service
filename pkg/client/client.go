package client

import (
	"log"

	"github.com/gorilla/websocket"
	"simple-multiplayer-service/pkg/message"
)

// Client represents a single WebSocket connection
type Client struct {
	ID              string
	IPAddress       string
	Connection      *websocket.Conn
	SendMessageFunc func(message message.Message) error
	UnregisterFunc  func(clientID string)
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

		c.HandleMessage(msg)
	}
}
