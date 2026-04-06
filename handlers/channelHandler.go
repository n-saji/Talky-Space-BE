package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RoutingChannel(r chi.Router) {
	r.Get("/{id}", h.GetChannelByID)
	r.Post("/", h.CreateChannel)
}

func (h *Handler) GetChannelByID(w http.ResponseWriter, r *http.Request) {
	_ = chi.URLParam(r, "id")
	writeJSON(w, http.StatusNotImplemented, map[string]string{"message": "Not implemented"})
}

func (h *Handler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		return
	}

	writeJSON(w, http.StatusNotImplemented, map[string]string{"message": "Not implemented"})
}
