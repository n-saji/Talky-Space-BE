package service

import (
	"errors"
	"talky-space-be/dtos"
	"talky-space-be/models"
)

func (s *Service) CreateMessage(req *dtos.CreateMessageRequest) (*dtos.MessageResponse, error) {
	if (req.SenderID == "" || req.RecipientID == "") || req.Content == "" {
		return nil, errors.New("invalid request data")
	}
	if req.ChatroomID != "" {
		chatroom, err := s.daos.GetChatroomByID(req.ChatroomID)
		if err != nil {
			return nil, err
		}
		if chatroom == nil {
			return nil, errors.New("chatroom not found")
		}
	} else {
		chatroom, err := s.daos.CheckChatroomExistForSenderReceiver(req.SenderID, req.RecipientID)
		if err != nil {
			return nil, err
		}
		if chatroom == nil {
			chatroomRes, err := s.CreateChatroomForUsers(req.SenderID, req.RecipientID)
			if err != nil {
				return nil, err
			}

			req.ChatroomID = chatroomRes.Id
		} else {
			req.ChatroomID = chatroom.ID.String()
		}
	}
	messageModel := models.CreateMessageRequestToMessageModel(req)
	createdMessage, err := s.daos.CreateMessage(*messageModel)
	if err != nil {
		return nil, err
	}
	messageResponse := models.MessageModelToMessageResponse(createdMessage)
	err = s.SendMessage(req.ChatroomID, req.RecipientID, req.Content)
	if err != nil {
		return nil, err
	}
	return messageResponse, nil
}

func (s *Service) StoreMessage(req *dtos.CreateMessageRequest) error {
	if (req.SenderID == "" || req.RecipientID == "" || req.ChatroomID == "") || req.Content == "" {
		return errors.New("invalid request data")
	}

	chatroom, err := s.daos.GetChatroomByID(req.ChatroomID)
	if err != nil {
		return err
	}
	if chatroom == nil {
		return errors.New("chatroom not found")
	}

	messageModel := models.CreateMessageRequestToMessageModel(req)
	_, err = s.daos.CreateMessage(*messageModel)
	if err != nil {
		return err
	}
	return nil
}
