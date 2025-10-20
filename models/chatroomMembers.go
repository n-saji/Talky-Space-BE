package models

import (
	"talky-space-be/dtos"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatroomMember struct {
	Id         uuid.UUID `gorm:"primary_key;type:uuid;unique"`
	ChatroomID uuid.UUID
	UserID     uuid.UUID
	JoinedAt   int64
}

func (cm *ChatroomMember) TableName() string {
	return "chatroom_members"
}

func (cm *ChatroomMember) BeforeCreate(tx *gorm.DB) (err error) {
	cm.Id = uuid.New()
	cm.JoinedAt = time.Now().Unix()
	return
}

func CreateChatroomMemberRequestToChatroomMemberModel(req *dtos.CreateChatroomMemberRequest) *ChatroomMember {
	chatroomID, _ := uuid.Parse(req.ChatroomID)
	userID, _ := uuid.Parse(req.UserID)

	return &ChatroomMember{
		Id:         uuid.New(),
		ChatroomID: chatroomID,
		UserID:     userID,
		JoinedAt:   time.Now().Unix(),
	}
}

func ChatroomMemberModelToChatroomMemberResponse(cm *ChatroomMember) *dtos.ChatroomMemberResponse {
	return &dtos.ChatroomMemberResponse{
		Id:         cm.Id.String(),
		ChatroomID: cm.ChatroomID.String(),
		UserID:     cm.UserID.String(),
		JoinedAt:   cm.JoinedAt,
	}
}
