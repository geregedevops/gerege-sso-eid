package handler

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"gesign.mn/gerege-sso/internal/model"
)

// danState is the decoded state parameter from sso.gov.mn callback
type danState struct {
	RedirectURL string `json:"redirect_url"`
	Session     string `json:"session"`
}

// DANCallback handles the callback from sso.gov.mn after DAN verification.
// sso.gov.mn redirects to: {callback_uri}?reg_no={reg_no}&state={base64_state}
func (h *Handler) DANCallback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	regNo := q.Get("reg_no")
	stateB64 := q.Get("state")

	if stateB64 == "" {
		h.jsonError(w, 400, "invalid_request", "missing state")
		return
	}

	// Decode state to get session ID
	stateBytes, err := base64.RawURLEncoding.DecodeString(stateB64)
	if err != nil {
		// Try standard base64
		stateBytes, err = base64.StdEncoding.DecodeString(stateB64)
		if err != nil {
			h.jsonError(w, 400, "invalid_request", "invalid state encoding")
			return
		}
	}

	var state danState
	if err := json.Unmarshal(stateBytes, &state); err != nil {
		h.jsonError(w, 400, "invalid_request", "invalid state format")
		return
	}

	sessionID := state.Session
	if sessionID == "" {
		h.jsonError(w, 400, "invalid_request", "missing session in state")
		return
	}

	// Get session from Redis
	var session model.AuthSession
	if err := h.cfg.Cache.Get(r.Context(), "sso:"+sessionID, &session); err != nil {
		logErr("dan_callback: session not found", err)
		h.jsonError(w, 400, "invalid_request", "session expired or not found")
		return
	}

	// Handle error (user cancelled or reg_no empty)
	if regNo == "" {
		h.cfg.Cache.Del(r.Context(), "sso:"+sessionID)
		redirectWithError(w, r, session.RedirectURI, session.State, "access_denied", "DAN verification failed")
		return
	}

	// Generate auth code
	authCode := generateRandomString(32)

	// Store code in Redis (5 min, single use)
	// For DAN flow, sub = reg_no (registration number is the citizen identifier)
	codeData := model.AuthCode{
		Sub:         regNo,
		RegNo:       regNo,
		ClientID:    session.ClientID,
		RedirectURI: session.RedirectURI,
		Scope:       session.Scope,
		Nonce:       session.Nonce,
	}
	if err := h.cfg.Cache.Set(r.Context(), "code:"+authCode, codeData, 5*time.Minute); err != nil {
		logErr("dan_callback: redis set code", err)
		h.jsonError(w, 500, "server_error", "internal error")
		return
	}

	// Delete session
	h.cfg.Cache.Del(r.Context(), "sso:"+sessionID)

	// Redirect to client with code + state
	redirectWithCode(w, r, session.RedirectURI, authCode, session.State)
}
