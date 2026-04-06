package handlers

import (
	"encoding/json"
	"net/http"
	"talky-space-be/dtos"
	"talky-space-be/middleware"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RoutingUser(r chi.Router) {
	r.Post("/register", h.CreateUser)
	r.Get("/me", h.GetUser)
	r.Put("/update", h.UpdateUser)
	r.Delete("/delete", h.DeleteUser)
	r.Get("/look-up", h.LookUpUser)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateUserRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}
	if err := h.service.CreateUser(r.Context(), &req); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, exists := middleware.UserIDFromContext(r.Context())
	if !exists {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}
	user, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, exists := middleware.UserIDFromContext(r.Context())
	if !exists {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}
	var req dtos.UpdateUserRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}
	if err := h.service.UpdateUser(r.Context(), id, &req); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, exists := middleware.UserIDFromContext(r.Context())
	if !exists {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}
	if err := h.service.DeleteUser(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: "", Path: "/", Secure: true, HttpOnly: true, MaxAge: -1})
	http.SetCookie(w, &http.Cookie{Name: "access_token", Value: "", Path: "/", Secure: true, HttpOnly: true, MaxAge: -1})
	writeJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func (h *Handler) LookUpUser(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	id, exists := middleware.UserIDFromContext(r.Context())
	if !exists {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}
	if query == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Query parameter 'q' is required"})
		return
	}
	users, err := h.service.LookUpUser(r.Context(), query, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, users)
}
