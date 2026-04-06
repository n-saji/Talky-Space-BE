package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"talky-space-be/auth"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDContextKey contextKey = "user_id"

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookieValue := ""
			if cookie, err := r.Cookie("access_token"); err == nil {
				cookieValue = cookie.Value
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" && cookieValue == "" {
				writeUnauthorized(w, "Missing Authorization header or cookie")
				return
			}

			if cookieValue != "" {
				authHeader = "Bearer " + cookieValue
			}

			tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			token, err := auth.VerifyToken(tokenStr, auth.AccessSecret)
			if err != nil || !token.Valid {
				log.Println(tokenStr, authHeader, cookieValue)
				writeUnauthorized(w, "Invalid token")
				return
			}

			claims, exists := token.Claims.(jwt.MapClaims)
			if !exists {
				writeUnauthorized(w, "Invalid token claims")
				return
			}

			userIDRaw, exists := claims["user_id"]
			if !exists {
				writeUnauthorized(w, "No user ID in token")
				return
			}

			userID, ok := userIDRaw.(string)
			if !ok || userID == "" {
				writeUnauthorized(w, "No user ID in token")
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	return userID, ok
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
