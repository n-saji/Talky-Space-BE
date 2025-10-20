package service

import (
	"talky-space-be/dtos"
	"talky-space-be/models"
)

func (s *Service) CreateChatroomMember(req dtos.CreateChatroomMemberRequest) error {
	chatroomMember := models.CreateChatroomMemberRequestToChatroomMemberModel(&req)
	return s.daos.CreateChatroomMember(chatroomMember)
}

func (s *Service) GetChatroomMembersByChatroomID(chatroomID string) ([]dtos.ChatroomMemberResponse, error) {
	members, err := s.daos.GetChatroomMembersByChatroomID(chatroomID)
	if err != nil {
		return nil, err
	}

	var memberResponses []dtos.ChatroomMemberResponse
	for _, member := range members {
		memberResponses = append(memberResponses, *models.ChatroomMemberModelToChatroomMemberResponse(&member))
	}
	return memberResponses, nil
}

func (s *Service) DeleteChatroomMember(chatroomID string, userID string) error {
	return s.daos.DeleteChatroomMember(chatroomID, userID)
}

func (s *Service) IsUserInChatroom(chatroomID string, userID string) (bool, error) {
	return s.daos.IsUserInChatroom(chatroomID, userID)
}
