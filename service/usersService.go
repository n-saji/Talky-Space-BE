package service

import (
	"errors"
	"talky-space-be/dtos"
	"talky-space-be/models"
)

func ValidateUserCreation(req *dtos.CreateUserRequest) bool {
	// Add validation logic here (e.g., check email format, password strength)
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return false
	}
	return true
}

func (s *Service) CreateUser(req *dtos.CreateUserRequest) error {
	if !ValidateUserCreation(req) {
		return errors.New("invalid user creation request")
	}
	// Check if user with the same email already exists
	_, err := s.daos.GetUserByEmail(req.Email)
	if err == nil {
		return errors.New("user already exists")
	}
	_, err = s.daos.GetByPhoneNumber(req.PhoneNumber)
	if err == nil {
		return errors.New("user already exists")
	}
	err = s.daos.CreateUser(models.CreateUserRequestToUserModel(req))
	if err != nil {
		return errors.New("failed to create user: " + err.Error())
	}
	return nil
}

func (s *Service) GetUserByID(id string) (*dtos.UserResponse, error) {
	user, err := s.daos.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return models.UserModelToUserResponse(user), nil
}

func (s *Service) UpdateUser(id string, req *dtos.UpdateUserRequest) error {
	user, err := s.daos.GetUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	updatedUser := models.UpdateUserRequestToUserModel(req, user)
	err = s.daos.UpdateUser(updatedUser)
	if err != nil {
		return errors.New("failed to update user: " + err.Error())
	}
	return nil
}

func (s *Service) DeleteUser(id string) error {
	user, err := s.daos.GetUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	err = s.daos.DeleteUser(user)
	if err != nil {
		return errors.New("failed to delete user: " + err.Error())
	}
	
	return nil
}		