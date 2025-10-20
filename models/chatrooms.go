package models

import (
	"talky-space-be/dtos"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Chatroom struct {
	ID          uuid.UUID `gorm:"primary_key;type:uuid;unique"`
	Name        string
	Description string
	IsGroup     bool
	CreatedBy   string
	CreatedAt   int64
	UpdatedAt   int64
}

func (c *Chatroom) TableName() string {
	return "chatrooms"
}

func (c *Chatroom) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	c.CreatedAt = time.Now().Unix()
	c.UpdatedAt = time.Now().Unix()
	return
}
func (c *Chatroom) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now().Unix()
	return
}

func (c *Chatroom) BeforeSave(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now().Unix()
	return
}

func CreateChatroomRequestToChatroomModel(req *dtos.CreateChatroomRequest) *Chatroom {
	now := time.Now().Unix()
	return &Chatroom{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		IsGroup:     req.IsGroup,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func ChatroomModelToChatroomResponse(chatroom *Chatroom) *dtos.ChatroomResponse {
	return &dtos.ChatroomResponse{
		Id:          chatroom.ID.String(),
		Name:        chatroom.Name,
		Description: chatroom.Description,
		IsGroup:     chatroom.IsGroup,
		CreatedBy:   chatroom.CreatedBy,
		CreatedAt:   chatroom.CreatedAt,
		UpdatedAt:   chatroom.UpdatedAt,
	}
}

func UpdateChatroomRequestToChatroomModel(req *dtos.UpdateChatroomRequest, chatroom *Chatroom) *Chatroom {
	if req.Name != "" {
		chatroom.Name = req.Name
	}
	if req.Description != "" {
		chatroom.Description = req.Description
	}
	return chatroom
}
