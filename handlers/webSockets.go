package handlers

import (
	"net/http"
	"talky-space-be/middleware"
	"talky-space-be/utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Handler) RoutingWebSockets(r chi.Router) {
	r.Get("/connect", h.HandleWebSocket)
}

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	userID, exists := middleware.UserIDFromContext(r.Context())
	if !exists {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &utils.Client{
		UserID: uuid.MustParse(userID),
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	utils.HubInstance.Register <- client

	go client.ReadPump()
	go client.WritePump()
}
