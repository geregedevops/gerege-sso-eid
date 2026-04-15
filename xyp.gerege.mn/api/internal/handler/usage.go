package handler

import (
	"log/slog"
	"net/http"
)

func (h *Handler) Usage(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	stats, err := h.cfg.DB.GetUsage(r.Context(), clientID, from, to)
	if err != nil {
		slog.Error("get usage failed", "error", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"usage": stats})
}
