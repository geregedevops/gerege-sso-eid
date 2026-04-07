package handler

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Revoke(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.jsonError(w, 400, "invalid_request", "malformed form data")
		return
	}

	// Client authentication
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		clientID = r.FormValue("client_id")
		clientSecret = r.FormValue("client_secret")
	}

	if clientID == "" || clientSecret == "" {
		h.jsonError(w, 401, "invalid_client", "client credentials required")
		return
	}

	client, err := h.cfg.DB.GetClient(r.Context(), clientID)
	if err != nil || client == nil {
		// RFC 7009: always return 200
		w.WriteHeader(http.StatusOK)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(client.SecretHash), []byte(clientSecret)); err != nil {
		h.jsonError(w, 401, "invalid_client", "invalid client credentials")
		return
	}

	token := r.FormValue("token")
	if token != "" {
		h.cfg.Cache.Del(r.Context(), "at:"+token)
	}

	// RFC 7009: always 200 OK
	w.WriteHeader(http.StatusOK)
}
