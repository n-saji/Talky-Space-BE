package middleware

import (
	"log"
	"net/http"
	"strings"
	"talky-space-be/auth"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("access_token")
		if err != nil {
			cookie = ""
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" && cookie == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header or cookie"})
			c.Abort()
			return
		}

		if cookie != "" {
			authHeader = "Bearer " + cookie
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := auth.VerifyToken(tokenStr, auth.AccessSecret)
		if err != nil || !token.Valid {
			log.Println(tokenStr,authHeader,cookie)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, exists := token.Claims.(jwt.MapClaims)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		if len(claims) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No user ID in token"})
			c.Abort()
			return
		}

		// Assuming the user_id is stored as a string in the claims
		userID := claims["user_id"].(string)
		c.Set("user_id", userID)
		c.Next()
	}
}
