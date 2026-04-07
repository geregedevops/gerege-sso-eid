package handler

import "net/http"

func (h *Handler) SignCancel(w http.ResponseWriter, r *http.Request) {
	sub := getSub(r)
	if sub == "" {
		h.jsonError(w, 401, "unauthorized")
		return
	}

	sessionID := r.PathValue("id")
	session, err := h.db.GetSession(r.Context(), sessionID)
	if err != nil || session == nil {
		h.jsonError(w, 404, "session not found")
		return
	}

	if session.RequesterSub != sub {
		h.jsonError(w, 403, "not your session")
		return
	}

	if session.Status != "PENDING" && session.Status != "RUNNING" {
		h.jsonError(w, 400, "session cannot be cancelled")
		return
	}

	h.db.UpdateSessionStatus(r.Context(), sessionID, "CANCELLED")
	h.jsonOK(w, map[string]string{"status": "CANCELLED"})
}
