package service

import (
	"talky-space-be/utils"
	"time"

	"github.com/google/uuid"
)

func (s *Service) SendMessage(chatroomID, senderID string, content string) error {

	chatroomUUID, err := uuid.Parse(chatroomID)
	if err != nil {
		return err
	}

	senderUUID, err := uuid.Parse(senderID)
	if err != nil {
		return err
	}

	// Then broadcast to WebSocket
	utils.BroadcastMessage(utils.MessagePayload{
		Type:       "message",
		ChatroomID: chatroomUUID,
		SenderID:   senderUUID,
		Content:    content,
		CreatedAt:  time.Now().Unix(),
		Source:     "server",
	})
	return nil
}
