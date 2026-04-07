package handler

import "net/http"

func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	h.jsonOK(w, map[string]any{
		"valid":      false,
		"message":    "PDF signature verification not yet implemented",
		"signatures": []any{},
	})
}
