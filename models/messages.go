package models

import (
	"talky-space-be/dtos"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Messages struct {
	Id         uuid.UUID `gorm:"primary_key;type:uuid;unique"`
	ChatroomID string
	UserId     uuid.UUID
	Content    string
	CreatedAt  int64
}

func (m *Messages) TableName() string {
	return "messages"
}

func (m *Messages) BeforeCreate(tx *gorm.DB) (err error) {
	m.Id = uuid.New()
	m.CreatedAt = time.Now().Unix()
	return
}

func (m *Messages) BeforeSave(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now().Unix()
	return
}

func (m *Messages) BeforeUpdate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now().Unix()
	return
}

func CreateMessageRequestToMessageModel(req *dtos.CreateMessageRequest) *Messages {
	userID, _ := uuid.Parse(req.SenderID)

	return &Messages{
		Id:         uuid.New(),
		ChatroomID: req.ChatroomID,
		UserId:     userID,
		Content:    req.Content,
		CreatedAt:  time.Now().Unix(),
	}
}

func MessageModelToMessageResponse(m *Messages) *dtos.MessageResponse {
	return &dtos.MessageResponse{
		Id:         m.Id.String(),
		ChatroomID: m.ChatroomID,
		UserId:     m.UserId.String(),
		Content:    m.Content,
		CreatedAt:  m.CreatedAt,
	}
}
