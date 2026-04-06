package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"talky-space-be/auth"
	"talky-space-be/dtos"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

var authValidator = validator.New()

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req dtos.LoginRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	if err := authValidator.Var(req.Email, "required,email"); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}
	if err := authValidator.Var(req.Password, "required"); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	rememberMe := strings.EqualFold(r.URL.Query().Get("remember_me"), "true")
	accessToken, refreshToken, err := h.service.AuthenticateUser(r.Context(), &req)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   7 * 24 * 3600,
		SameSite: http.SameSiteNoneMode,
	})
	if !rememberMe {
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   15 * 60,
		SameSite: http.SameSiteNoneMode,
	})

	writeJSON(w, http.StatusOK, map[string]string{"access_token": accessToken})
}

func (h *Handler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil || refreshCookie.Value == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "No refresh token"})
		return
	}

	token, err := auth.VerifyToken(refreshCookie.Value, auth.RefreshSecret)
	if err != nil || !token.Valid {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid refresh token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid refresh token"})
		return
	}

	userIDValue, ok := claims["user_id"]
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid user ID in token"})
		return
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid user ID in token"})
		return
	}

	userIDParsed, err := auth.ParseUUID(strings.TrimSpace(userIDStr))
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid user ID in token"})
		return
	}

	accessToken, newRefresh, err := auth.GenerateTokens(userIDParsed)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to refresh"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefresh,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   7 * 24 * 3600,
		SameSite: http.SameSiteNoneMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   15 * 60,
		SameSite: http.SameSiteNoneMode,
	})

	writeJSON(w, http.StatusOK, map[string]string{"access_token": accessToken})
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, _ *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
		SameSite: http.SameSiteNoneMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
		SameSite: http.SameSiteNoneMode,
	})

	writeJSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

