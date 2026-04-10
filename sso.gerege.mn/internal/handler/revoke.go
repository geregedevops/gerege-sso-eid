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

	// RFC 7009: revocation endpoint always returns 200 to avoid leaking client info
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		clientID = r.FormValue("client_id")
		clientSecret = r.FormValue("client_secret")
	}

	if clientID == "" || clientSecret == "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	client, err := h.cfg.DB.GetClient(r.Context(), clientID)
	if err != nil || client == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(client.SecretHash), []byte(clientSecret)); err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	token := r.FormValue("token")
	if token != "" {
		h.cfg.Cache.Del(r.Context(), "at:"+token)
	}

	w.WriteHeader(http.StatusOK)
}
