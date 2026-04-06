package service

import (
	"context"
	"errors"
	"talky-space-be/dtos"
	"talky-space-be/models"
)

func ValidateChatroomCreation(ctx context.Context, req *dtos.CreateChatroomRequest) bool {
	// Add validation logic here (e.g., check name length)
	if req.IsGroup && (req.Name == "" || req.Description == "" || req.CreatedBy == "") {
		return false
	}
	return true
}

func (s *Service) CreateChatroom(ctx context.Context, req *dtos.CreateChatroomRequest) (*models.Chatroom, error) {
	if !ValidateChatroomCreation(ctx, req) {
		return nil, errors.New("invalid chatroom creation request")
	}
	chatroomModel := models.CreateChatroomRequestToChatroomModel(req)
	err := s.daos.CreateChatroom(ctx, chatroomModel)
	if err != nil {
		return nil, err
	}
	return chatroomModel, nil
}

func (s *Service) GetChatroomByID(ctx context.Context, id string) (*dtos.ChatroomResponse, error) {
	chatroom, err := s.daos.GetChatroomByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return models.ChatroomModelToChatroomResponse(chatroom), nil
}

func (s *Service) UpdateChatroom(ctx context.Context, id string, req *dtos.UpdateChatroomRequest) error {
	chatroom, err := s.daos.GetChatroomByID(ctx, id)
	if err != nil {
		return err
	}

	updatedChatroom := models.UpdateChatroomRequestToChatroomModel(req, chatroom)
	err = s.daos.UpdateChatroom(ctx, updatedChatroom)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteChatroom(ctx context.Context, id string) error {
	err := s.daos.DeleteChatroom(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) CreateChatroomForUsers(ctx context.Context, userID1, userID2 string) (*dtos.ChatroomResponse, error) {
	chatroomReq := &dtos.CreateChatroomRequest{
		Name:        "Private Chat",
		Description: "Chat between users",
		IsGroup:     false,
		CreatedBy:   userID1,
	}
	chatroom, err := s.CreateChatroom(ctx, chatroomReq)
	if err != nil {
		return nil, err
	}

	err = s.CreateChatroomMember(ctx, dtos.CreateChatroomMemberRequest{
		ChatroomID: chatroom.ID.String(),
		UserID:     userID1,
	})
	if err != nil {
		return nil, err
	}

	err = s.CreateChatroomMember(ctx, dtos.CreateChatroomMemberRequest{
		ChatroomID: chatroom.ID.String(),
		UserID:     userID2,
	})
	if err != nil {
		return nil, err
	}

	chatroom, err = s.daos.CheckChatroomExistForSenderReceiver(ctx, userID1, userID2)
	if err != nil {
		return nil, err
	}
	return models.ChatroomModelToChatroomResponse(chatroom), nil
}

func (s *Service) FindChatroomByUsers(ctx context.Context, userID1, userID2 string) (*dtos.ChatroomResponse, error) {
	chatroom, err := s.daos.CheckChatroomExistForSenderReceiver(ctx, userID1, userID2)
	if err != nil {
		return nil, err
	}
	if chatroom == nil {
		return nil, errors.New("chatroom not found for the given users")
	}
	return models.ChatroomModelToChatroomResponse(chatroom), nil
}

func (s *Service) FetchAllChatroomsForUser(ctx context.Context, userID string) ([]dtos.ChatroomResponse, error) {
	chatrooms, err := s.daos.FetchChatroomsForUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	var chatroomResponses []dtos.ChatroomResponse
	for _, chatroom := range chatrooms {
		if !chatroom.IsGroup {
			member, err := s.daos.GetPrivateChatroomOtherMember(ctx, chatroom.ID.String(), userID)
			if err != nil {
				return nil, err
			}
			if member != nil {
				chatroom.Name = member.Username
			}
		}
		chatroomResponses = append(chatroomResponses, *models.ChatroomModelToChatroomResponse(chatroom))
	}
	return chatroomResponses, nil
}
