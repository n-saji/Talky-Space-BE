package handlers

import (
	"encoding/json"
	"net/http"
	"talky-space-be/dtos"
	"talky-space-be/middleware"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RoutingMessage(r chi.Router) {
	r.Post("/new-chat", h.InitiateNewChat)
	r.Get("/chatroom/{cid}", h.FetchChatroomMessages)
}

func (h *Handler) InitiateNewChat(w http.ResponseWriter, r *http.Request) {
	req := &dtos.CreateMessageRequest{}
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}
	req.SenderID = userID

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		return
	}
	res, err := h.service.InitiateNewChat(r.Context(), req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create message with error: " + err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) FetchChatroomMessages(w http.ResponseWriter, r *http.Request) {
	chatroomID := chi.URLParam(r, "cid")
	res, err := h.service.FetchMessagesByChatroomID(r.Context(), chatroomID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch messages with error: " + err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, res)
}
