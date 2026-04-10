package handler

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) requireDANAdmin(w http.ResponseWriter, r *http.Request) bool {
	if h.cfg.DANAdminKey == "" {
		h.jsonError(w, 403, "forbidden", "admin not configured")
		return false
	}
	auth := r.Header.Get("Authorization")
	expected := "Bearer " + h.cfg.DANAdminKey
	if subtle.ConstantTimeCompare([]byte(auth), []byte(expected)) != 1 {
		h.jsonError(w, 401, "unauthorized", "invalid admin key")
		return false
	}
	return true
}

func (h *Handler) ListDANClients(w http.ResponseWriter, r *http.Request) {
	if !h.requireDANAdmin(w, r) {
		return
	}

	clients, err := h.cfg.DB.ListDANClients(r.Context())
	if err != nil {
		logErr("dan_clients: list", err)
		h.jsonError(w, 500, "server_error", "internal error")
		return
	}
	if clients == nil {
		h.jsonOK(w, []any{})
		return
	}
	h.jsonOK(w, clients)
}

func (h *Handler) CreateDANClient(w http.ResponseWriter, r *http.Request) {
	if !h.requireDANAdmin(w, r) {
		return
	}

	var req struct {
		Name         string   `json:"name"`
		CallbackURLs []string `json:"callback_urls"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.jsonError(w, 400, "invalid_request", "invalid JSON body")
		return
	}
	if req.Name == "" {
		h.jsonError(w, 400, "invalid_request", "name is required")
		return
	}
	if len(req.CallbackURLs) == 0 {
		h.jsonError(w, 400, "invalid_request", "at least one callback_url is required")
		return
	}

	// Generate client credentials
	clientID := fmt.Sprintf("dan_%s", randomHex(16))
	clientSecret := randomBase64(32)
	hmacKey := randomBase64(32)

	hash, err := bcrypt.GenerateFromPassword([]byte(clientSecret), 12)
	if err != nil {
		logErr("dan_clients: bcrypt", err)
		h.jsonError(w, 500, "server_error", "internal error")
		return
	}

	if err := h.cfg.DB.CreateDANClient(r.Context(), clientID, string(hash), hmacKey, req.Name, req.CallbackURLs); err != nil {
		logErr("dan_clients: create", err)
		h.jsonError(w, 500, "server_error", "internal error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":            clientID,
		"secret":        clientSecret,
		"hmac_key":      hmacKey,
		"name":          req.Name,
		"callback_urls": req.CallbackURLs,
		"message":       "secret болон hmac_key-ийг хадгалж авна уу. Дахин харагдахгүй.",
	})
}

func (h *Handler) DeactivateDANClient(w http.ResponseWriter, r *http.Request) {
	if !h.requireDANAdmin(w, r) {
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.jsonError(w, 400, "invalid_request", "client id required")
		return
	}

	if err := h.cfg.DB.DeactivateDANClient(r.Context(), id); err != nil {
		logErr("dan_clients: deactivate", err)
		h.jsonError(w, 500, "server_error", "internal error")
		return
	}

	h.jsonOK(w, map[string]string{"status": "deactivated", "id": id})
}

func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand failed: " + err.Error())
	}
	return fmt.Sprintf("%x", b)
}

func randomBase64(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand failed: " + err.Error())
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
