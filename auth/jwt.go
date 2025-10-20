package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	AccessSecret  = []byte(os.Getenv("ACCESS_SECRET_KEY"))  // store securely in env
	RefreshSecret = []byte(os.Getenv("REFRESH_SECRET_KEY")) // store securely in env
)

func GenerateTokens(userID uuid.UUID) (string, string, error) {
	// Access token (15 minutes)
	atClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString(AccessSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh token (7 days)
	rtClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString(RefreshSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func VerifyToken(tokenStr string, secret []byte) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return secret, nil
	})
}

func ParseUUID(idStr string) (uuid.UUID, error) {
	return uuid.Parse(idStr)
}
