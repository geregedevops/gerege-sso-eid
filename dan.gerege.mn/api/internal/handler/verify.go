package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
)

func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	callbackURL := r.URL.Query().Get("callback_url")
	clientID := r.URL.Query().Get("client_id")

	if callbackURL == "" {
		h.jsonError(w, 400, "callback_url параметр шаардлагатай")
		return
	}

	if clientID == "" {
		h.jsonError(w, 400, "client_id параметр шаардлагатай")
		return
	}

	// Validate client
	client, err := h.cfg.DB.GetDANClient(r.Context(), clientID)
	if err != nil {
		slog.Error("verify: db error", "error", err)
		h.jsonError(w, 500, "internal error")
		return
	}
	if client == nil || !client.Active {
		h.jsonError(w, 400, "бүртгэлгүй эсвэл идэвхгүй client")
		return
	}

	// Validate callback URL (domain match, HTTPS required)
	if !matchCallbackURL(client.CallbackURLs, callbackURL) {
		h.jsonError(w, 400, "callback_url бүртгэлгүй байна")
		return
	}

	// Build signed state
	statePayload, _ := json.Marshal(map[string]string{
		"callback_url": callbackURL,
		"client_id":    clientID,
	})
	signedState := signState(statePayload, h.cfg.StateSecret)

	loginURL := fmt.Sprintf("https://sso.gov.mn/login?state=%s&grant_type=authorization_code&response_type=code&client_id=%s&scope=%s&redirect_uri=%s",
		url.QueryEscape(signedState),
		url.QueryEscape(h.cfg.DAN.ClientID),
		url.QueryEscape(h.cfg.DAN.Scope),
		url.QueryEscape(h.cfg.DAN.CallbackURI),
	)

	slog.Info("verify: redirecting", "client_id", clientID)
	http.Redirect(w, r, loginURL, http.StatusFound)
}

// signState creates a signed state: base64(payload).base64(hmac)
func signState(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	sig := mac.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(payload) + "." +
		base64.RawURLEncoding.EncodeToString(sig)
}

// verifyState validates and extracts payload from signed state.
func verifyState(signed, secret string) (map[string]string, error) {
	parts := splitState(signed)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid state format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}

	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decode signature: %w", err)
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	if !hmac.Equal(mac.Sum(nil), sig) {
		return nil, fmt.Errorf("signature mismatch")
	}

	var state map[string]string
	if err := json.Unmarshal(payload, &state); err != nil {
		return nil, fmt.Errorf("unmarshal state: %w", err)
	}
	return state, nil
}

func splitState(s string) []string {
	idx := len(s) - 1
	for idx >= 0 && s[idx] != '.' {
		idx--
	}
	if idx <= 0 {
		return nil
	}
	return []string{s[:idx], s[idx+1:]}
}
