package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (h *Handler) ListClients(w http.ResponseWriter, r *http.Request) {
	clients, err := h.cfg.DB.ListClients(r.Context())
	if err != nil {
		slog.Error("list clients failed", "error", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"clients": clients})
}

func (h *Handler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name         string `json:"name"`
		ContactEmail string `json:"contact_email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	client, secret, err := h.cfg.DB.CreateClient(r.Context(), req.Name, req.ContactEmail)
	if err != nil {
		slog.Error("create client failed", "error", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"client":        client,
		"client_secret": secret,
		"message":       "Save the client_secret now. It will not be shown again.",
	})
}

func (h *Handler) DeactivateClient(w http.ResponseWriter, r *http.Request) {
	clientID := r.PathValue("id")
	if clientID == "" {
		writeError(w, http.StatusBadRequest, "client id is required")
		return
	}

	if err := h.cfg.DB.DeactivateClient(r.Context(), clientID); err != nil {
		slog.Error("deactivate client failed", "error", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deactivated", "client_id": clientID})
}
