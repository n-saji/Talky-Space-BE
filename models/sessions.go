package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Sessions struct {
	Id           uuid.UUID `gorm:"primary_key;type:uuid;unique"`
	UserId       uuid.UUID `gorm:"type:uuid;not null"`
	RefreshToken string    `gorm:"type:text;not null"`
	ExpiresAt    int64     `gorm:"not null"`
	CreatedAt    int64     `gorm:"not null"`
}

func (s *Sessions) BeforeCreate(tx *gorm.DB) (err error) {
	s.Id = uuid.New()
	return
}

func (s *Sessions) TableName() string {
	return "sessions"
}
