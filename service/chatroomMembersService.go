package service

import (
	"context"
	"talky-space-be/dtos"
	"talky-space-be/models"
)

func (s *Service) CreateChatroomMember(ctx context.Context, req dtos.CreateChatroomMemberRequest) error {
	chatroomMember := models.CreateChatroomMemberRequestToChatroomMemberModel(&req)
	return s.daos.CreateChatroomMember(ctx, chatroomMember)
}

func (s *Service) GetChatroomMembersByChatroomID(ctx context.Context, chatroomID string) ([]dtos.ChatroomMemberResponse, error) {
	members, err := s.daos.GetChatroomMembersByChatroomID(ctx, chatroomID)
	if err != nil {
		return nil, err
	}

	var memberResponses []dtos.ChatroomMemberResponse
	for _, member := range members {
		memberResponses = append(memberResponses, *models.ChatroomMemberModelToChatroomMemberResponse(&member))
	}
	return memberResponses, nil
}

func (s *Service) DeleteChatroomMember(ctx context.Context, chatroomID string, userID string) error {
	return s.daos.DeleteChatroomMember(ctx, chatroomID, userID)
}

func (s *Service) IsUserInChatroom(ctx context.Context, chatroomID string, userID string) (bool, error) {
	return s.daos.IsUserInChatroom(ctx, chatroomID, userID)
}
