package service

import (
	"errors"
	"talky-space-be/dtos"
	"talky-space-be/models"
)

func ValidateChatroomCreation(req *dtos.CreateChatroomRequest) bool {
	// Add validation logic here (e.g., check name length)
	if req.IsGroup && (req.Name == "" || req.Description == "" || req.CreatedBy == "") {
		return false
	}
	return true
}

func (s *Service) CreateChatroom(req *dtos.CreateChatroomRequest) (*models.Chatroom, error) {
	if !ValidateChatroomCreation(req) {
		return nil, errors.New("invalid chatroom creation request")
	}
	chatroomModel := models.CreateChatroomRequestToChatroomModel(req)
	err := s.daos.CreateChatroom(chatroomModel)
	if err != nil {
		return nil, err
	}
	return chatroomModel, nil
}

func (s *Service) GetChatroomByID(id string) (*dtos.ChatroomResponse, error) {
	chatroom, err := s.daos.GetChatroomByID(id)
	if err != nil {
		return nil, err
	}
	return models.ChatroomModelToChatroomResponse(chatroom), nil
}

func (s *Service) UpdateChatroom(id string, req *dtos.UpdateChatroomRequest) error {
	chatroom, err := s.daos.GetChatroomByID(id)
	if err != nil {
		return err
	}

	updatedChatroom := models.UpdateChatroomRequestToChatroomModel(req, chatroom)
	err = s.daos.UpdateChatroom(updatedChatroom)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteChatroom(id string) error {
	err := s.daos.DeleteChatroom(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) CreateChatroomForUsers(userID1, userID2 string) (*dtos.ChatroomResponse, error) {
	chatroomReq := &dtos.CreateChatroomRequest{
		Name:        "Private Chat",
		Description: "Chat between users",
		IsGroup:     false,
		CreatedBy:   userID1,
	}
	chatroom, err := s.CreateChatroom(chatroomReq)
	if err != nil {
		return nil, err
	}

	err = s.CreateChatroomMember(dtos.CreateChatroomMemberRequest{
		ChatroomID: chatroom.ID.String(),
		UserID:     userID1,
	})
	if err != nil {
		return nil, err
	}

	err = s.CreateChatroomMember(dtos.CreateChatroomMemberRequest{
		ChatroomID: chatroom.ID.String(),
		UserID:     userID2,
	})
	if err != nil {
		return nil, err
	}

	chatroom, err = s.daos.CheckChatroomExistForSenderReceiver(userID1, userID2)
	if err != nil {
		return nil, err
	}
	return models.ChatroomModelToChatroomResponse(chatroom), nil
}
