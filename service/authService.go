package service

import (
	"context"
	"errors"
	"talky-space-be/auth"
	"talky-space-be/dtos"
	"talky-space-be/models"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) AuthenticateUser(ctx context.Context, req *dtos.LoginRequest) (string, string, error) {
	user, err := s.daos.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return "", "", errors.New("user not found")
		}
		return "", "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, refreshToken, err := auth.GenerateTokens(user.Id)
	if err != nil {
		return "", "", err
	}

	// Save refresh token in DB
	// session := models.Sessions{
	// 	UserId:       user.Id,
	// 	RefreshToken: refreshToken,
	// 	ExpiresAt:    time.Now().Add(7 * 24 * time.Hour).Unix(),
	// 	CreatedAt:    time.Now().Unix(),
	// }
	// if err := s.daos.SaveSession(ctx, &session); err != nil {
	// 	return "", "", err
	// }

	return accessToken, refreshToken, nil
}
