package handlers

import (
	"net/http"
	"talky-space-be/auth"
	"talky-space-be/dtos"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) AuthenticationChannel(rg *gin.RouterGroup) {
	authGrp := rg.Group("/auth")
	{
		authGrp.POST("/login", func(c *gin.Context) {
			var req dtos.LoginRequest

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}
			rememberMe := c.Query("remember_me") == "true"
			accessToken, refreshToken, err := h.service.AuthenticateUser(&req)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}

			// Set refresh token in cookie (secure, HTTP-only)
			if rememberMe {
				http.SetCookie(c.Writer, &http.Cookie{
					Name:     "refresh_token",
					Value:    refreshToken,
					Path:     "/",
					HttpOnly: true,
					Secure:   true,
					MaxAge:   7 * 24 * 3600,
					SameSite: http.SameSiteNoneMode,
				})
			} else {
				http.SetCookie(c.Writer, &http.Cookie{
					Name:     "refresh_token",
					Value:    refreshToken,
					Path:     "/",
					HttpOnly: true,
					Secure:   true,
					MaxAge:   0,
					SameSite: http.SameSiteNoneMode,
				})
			}
			// Set access token in cookie (secure, HTTP-only)
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     "access_token",
				Value:    accessToken,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				MaxAge:   15*60,
				SameSite: http.SameSiteNoneMode,
			})

			c.JSON(http.StatusOK, gin.H{
				"access_token": accessToken,
			})
		})

		authGrp.POST("/refresh", func(c *gin.Context) {
			refreshToken, err := c.Cookie("refresh_token")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token"})
				return
			}

			token, err := auth.VerifyToken(refreshToken, auth.RefreshSecret)
			if err != nil || !token.Valid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			userID := claims["user_id"].(string)
			userIDParsed, err := auth.ParseUUID(userID)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
				return
			}

			accessToken, newRefresh, err := auth.GenerateTokens(userIDParsed)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh"})
				return
			}

			// Update refresh token cookie
			c.SetCookie("refresh_token", newRefresh, 7*24*3600, "/", "", true, true)
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     "access_token",
				Value:    accessToken,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				MaxAge:   15*60,
				SameSite: http.SameSiteNoneMode,
			})
			c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
		})

		authGrp.POST("/logout", func(c *gin.Context) {
			c.SetCookie("refresh_token", "", -1, "/", "", true, true)
			c.SetCookie("access_token", "", -1, "/", "", true, true)
			c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
		})

	}
}
