package matchmaking

import (
	"log"

	"simple-multiplayer-service/internal/db"
	"simple-multiplayer-service/internal/message"
	"simple-multiplayer-service/internal/notification"

	"github.com/google/uuid"
)

type Service struct {
	SessionNumber       int
	SessionLimit        int
	SessionQueue        chan message.MatchmakingRequest
	SessionDB           db.Session
	NotificationService *notification.Service
}

func NewMatchmakingService(sessionLimit int, sessionDB db.Session, notificationService *notification.Service) *Service {
	sessionQueue := make(chan message.MatchmakingRequest, 100)
	return &Service{SessionLimit: sessionLimit, SessionQueue: sessionQueue, SessionDB: sessionDB, NotificationService: notificationService}
}

func (matchmakingService *Service) Start() {
	lookingForOpponent := ""
	for {
		select {
		case mmRequest := <-matchmakingService.SessionQueue:
			// if there is no opponent yet, put it into waiting for opponent variable
			if lookingForOpponent == "" {
				lookingForOpponent = mmRequest.ConnectionID
				continue
			}

			// if there is already player looking for opponent, match it
			newSession := Session{
				SessionID:           uuid.New().String(),
				Player1ConnectionID: mmRequest.ConnectionID,
				Player2ConnectionID: lookingForOpponent,
			}

			err := matchmakingService.SessionDB.CreateSession(newSession.SessionID, newSession.Player1ConnectionID, newSession.Player2ConnectionID)
			if err != nil {
				log.Println(err)
				continue
			}

			matchmakingService.SessionNumber++
			newSessionNotification := notification.SessionNotification{
				SessionID:           newSession.SessionID,
				Player1ConnectionID: mmRequest.ConnectionID,
				Player2ConnectionID: lookingForOpponent,
			}
			matchmakingService.NotificationService.Channel <- newSessionNotification
			lookingForOpponent = ""
		}
	}
}
