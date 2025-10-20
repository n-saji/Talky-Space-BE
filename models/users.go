package models

import (
	"talky-space-be/dtos"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

type User struct {
	Id           uuid.UUID `gorm:"primary_key;type:uuid;unique"`
	Username     string
	Email        string
	PhoneNumber  string
	PasswordHash string
	AvatarURL    string
	CreatedAt    int64
	UpdatedAt    int64
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id = uuid.New()
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now().Unix()
	return
}
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now().Unix()
	return
}
func (u *User) TableName() string {
	return "users"
}

func CreateUserRequestToUserModel(req *dtos.CreateUserRequest) *User {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	now := time.Now().Unix()

	return &User{
		Id:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hashed),
		Username:     req.Name,
		AvatarURL:    req.Avatar,
		PhoneNumber:  req.PhoneNumber,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func UserModelToUserResponse(user *User) *dtos.UserResponse {
	return &dtos.UserResponse{
		Id:          user.Id.String(),
		Email:       user.Email,
		Name:        user.Username,
		Avatar:      user.AvatarURL,
		PhoneNumber: user.PhoneNumber,
	}
}

func UpdateUserRequestToUserModel(req *dtos.UpdateUserRequest, existingUser *User) *User {
	if req.Email != "" {
		existingUser.Email = req.Email
	}
	if req.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		existingUser.PasswordHash = string(hashed)
	}
	if req.Name != "" {
		existingUser.Username = req.Name
	}
	if req.Avatar != "" {
		existingUser.AvatarURL = req.Avatar
	}
	if req.PhoneNumber != "" {
		existingUser.PhoneNumber = req.PhoneNumber
	}
	return existingUser
}