package handlers

import (
	"net/http"
	"talky-space-be/global"
	"talky-space-be/middleware"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ChatroomChannel(r chi.Router) {
	r.Post("/", h.CreateChatroom)
	r.Get("/", h.GetUserChatrooms)
	r.Get("/{chatroom_id}", h.GetChatroomDetails)
	r.Put("/{chatroom_id}", h.UpdateChatroom)
	r.Delete("/{chatroom_id}", h.DeleteChatroom)
	r.Get("/find-by-users/user1/{uid1}/user2/{uid2}", h.FindChatroomByUsers)
	r.Get("/user", h.FetchAllChatrooms)
}

func (h *Handler) CreateChatroom(w http.ResponseWriter, r *http.Request) {
	// Implementation for creating a chatroom
	writeJSON(w, http.StatusNotImplemented, map[string]string{"message": "Not implemented"})
}

func (h *Handler) GetUserChatrooms(w http.ResponseWriter, r *http.Request) {
	// Implementation for retrieving user's chatrooms
	writeJSON(w, http.StatusNotImplemented, map[string]string{"message": "Not implemented"})
}

func (h *Handler) GetChatroomDetails(w http.ResponseWriter, r *http.Request) {
	// Implementation for retrieving chatroom details
	chatroomID := chi.URLParam(r, "chatroom_id")

	res, err := h.service.GetChatroomByID(r.Context(), chatroomID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) UpdateChatroom(w http.ResponseWriter, r *http.Request) {
	// Implementation for updating a chatroom
	writeJSON(w, http.StatusNotImplemented, map[string]string{"message": "Not implemented"})
}

func (h *Handler) DeleteChatroom(w http.ResponseWriter, r *http.Request) {
	// Implementation for deleting a chatroom
	writeJSON(w, http.StatusNotImplemented, map[string]string{"message": "Not implemented"})
}

func (h *Handler) FindChatroomByUsers(w http.ResponseWriter, r *http.Request) {
	uid1 := chi.URLParam(r, "uid1")
	uid2 := chi.URLParam(r, "uid2")

	res, err := h.service.FindChatroomByUsers(r.Context(), uid1, uid2)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) FetchAllChatrooms(w http.ResponseWriter, r *http.Request) {
	// Implementation for fetching all chatrooms
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		writeError(w, http.StatusUnauthorized, global.CodeUnauthorized, "Unauthorized")
		return
	}
	chatrooms, err := h.service.FetchAllChatroomsForUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, global.CodeInternal, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, chatrooms)
}
