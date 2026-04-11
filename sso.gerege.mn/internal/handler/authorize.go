package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"sso.gerege.mn/internal/model"
)

func (h *Handler) Authorize(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	clientID := q.Get("client_id")
	redirectURI := q.Get("redirect_uri")
	responseType := q.Get("response_type")
	scope := q.Get("scope")
	state := q.Get("state")
	nonce := q.Get("nonce")

	// Validate client
	client, err := h.cfg.DB.GetClient(r.Context(), clientID)
	if err != nil {
		logErr("authorize: db error", err)
		h.jsonError(w, 500, "server_error", "internal error")
		return
	}
	if client == nil || !client.IsActive {
		h.jsonError(w, 400, "invalid_request", "unknown client_id")
		return
	}
	if !matchRedirectURI(client.RedirectURIs, redirectURI) {
		h.jsonError(w, 400, "invalid_request", "redirect_uri not registered")
		return
	}
	if responseType != "code" {
		redirectWithError(w, r, redirectURI, state, "unsupported_response_type", "only code is supported")
		return
	}
	if !strings.Contains(scope, "openid") {
		redirectWithError(w, r, redirectURI, state, "invalid_scope", "openid scope required")
		return
	}

	// Store auth session
	sessionID := generateRandomString(32)
	session := model.AuthSession{
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       scope,
		State:       state,
		Nonce:       nonce,
	}
	if err := h.cfg.Cache.Set(r.Context(), "sso:"+sessionID, session, 10*time.Minute); err != nil {
		logErr("authorize: redis set", err)
		h.jsonError(w, 500, "server_error", "internal error")
		return
	}

	// Render login page (SSO owns the auth UI)
	clientName := client.Name
	if clientName == "" {
		clientName = clientID
	}
	h.renderLoginPage(w, sessionID, clientName)
}

// AuthInitiateAPI — login хуудаснаас national_id оруулахад дуудагдана
func (h *Handler) AuthInitiateAPI(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SessionID  string `json:"session_id"`
		NationalID string `json:"national_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.SessionID == "" || body.NationalID == "" {
		h.jsonError(w, 400, "invalid_request", "session_id and national_id required")
		return
	}

	// Verify session exists
	var session model.AuthSession
	if err := h.cfg.Cache.Get(r.Context(), "sso:"+body.SessionID, &session); err != nil {
		h.jsonError(w, 400, "invalid_request", "session not found")
		return
	}

	// Resolve client name for display_text (e-id.mn activity log)
	displayText := "SSO нэвтрэлт"
	if c, _ := h.cfg.DB.GetClient(r.Context(), session.ClientID); c != nil && c.Name != "" {
		displayText = c.Name
	}

	// Call e-ID RP API to initiate auth
	eidResp, err := h.eidAuthInitiate(body.NationalID, displayText)
	if err != nil {
		logErr("auth_initiate: eid api failed", err)
		h.jsonError(w, 502, "server_error", "e-ID холболт амжилтгүй")
		return
	}

	// Store e-ID session ID in our session
	h.cfg.Cache.Set(r.Context(), "sso-eid:"+body.SessionID, eidResp.SessionID, 10*time.Minute)

	h.jsonOK(w, map[string]interface{}{
		"eid_session_id":    eidResp.SessionID,
		"verification_code": eidResp.VerificationCode,
	})
}

// AuthPollAPI — login хуудаснаас e-ID session status poll хийнэ
func (h *Handler) AuthPollAPI(w http.ResponseWriter, r *http.Request) {
	ssoSessionID := r.URL.Query().Get("session_id")
	if ssoSessionID == "" {
		h.jsonError(w, 400, "invalid_request", "session_id required")
		return
	}

	// Get e-ID session ID
	var eidSessionID string
	if err := h.cfg.Cache.Get(r.Context(), "sso-eid:"+ssoSessionID, &eidSessionID); err != nil {
		h.jsonError(w, 400, "invalid_request", "e-ID session not found")
		return
	}

	// Poll e-ID RP API
	eidStatus, err := h.eidAuthStatus(eidSessionID)
	if err != nil {
		logErr("auth_poll: eid status failed", err)
		h.jsonError(w, 502, "server_error", "e-ID холболт амжилтгүй")
		return
	}

	if eidStatus.State == "COMPLETE" && eidStatus.Result == "OK" {
		// Auth successful — generate code
		var session model.AuthSession
		if err := h.cfg.Cache.Get(r.Context(), "sso:"+ssoSessionID, &session); err != nil {
			h.jsonError(w, 400, "invalid_request", "session expired")
			return
		}

		authCode := generateRandomString(32)
		codeData := model.AuthCode{
			Sub:         eidStatus.Identity.NationalID,
			Name:        eidStatus.Identity.FullName,
			GivenName:   "",
			FamilyName:  "",
			CertSerial:  eidStatus.Certificate.SerialNumber,
			ClientID:    session.ClientID,
			RedirectURI: session.RedirectURI,
			Scope:       session.Scope,
			Nonce:       session.Nonce,
		}
		h.cfg.Cache.Set(r.Context(), "code:"+authCode, codeData, 5*time.Minute)
		h.cfg.Cache.Del(r.Context(), "sso:"+ssoSessionID)
		h.cfg.Cache.Del(r.Context(), "sso-eid:"+ssoSessionID)

		redirectURL := session.RedirectURI + "?code=" + authCode
		if session.State != "" {
			redirectURL += "&state=" + url.QueryEscape(session.State)
		}

		h.jsonOK(w, map[string]interface{}{
			"status":       "complete",
			"redirect_url": redirectURL,
		})
		return
	}

	h.jsonOK(w, map[string]interface{}{
		"status": eidStatus.State,
	})
}

// ── e-ID RP API calls ──────────────────────────────────────────────

type eidAuthResp struct {
	SessionID        string `json:"session_id"`
	VerificationCode string `json:"verification_code"`
}

type eidStatusResp struct {
	State  string `json:"state"`
	Result string `json:"result"`
	Identity struct {
		NationalID string `json:"national_id"`
		FullName   string `json:"full_name"`
		KYCLevel   string `json:"kyc_level"`
	} `json:"identity"`
	Certificate struct {
		PEM          string `json:"pem"`
		Subject      string `json:"subject"`
		SerialNumber string `json:"serial_number"`
		ValidUntil   string `json:"valid_until"`
	} `json:"certificate"`
}

func (h *Handler) eidAuthInitiate(nationalID, displayText string) (*eidAuthResp, error) {
	body, _ := json.Marshal(map[string]string{
		"national_id":  nationalID,
		"display_text": displayText,
	})
	req, _ := http.NewRequest("POST", h.cfg.EIDBaseURL+"/rp/v1/auth/initiate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.cfg.EIDRPApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("eid request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("eid status %d: %s", resp.StatusCode, string(respBody))
	}

	var result eidAuthResp
	json.Unmarshal(respBody, &result)
	return &result, nil
}

func (h *Handler) eidAuthStatus(sessionID string) (*eidStatusResp, error) {
	req, _ := http.NewRequest("GET", h.cfg.EIDBaseURL+"/rp/v1/auth/session/"+sessionID+"?timeout_ms=5000", nil)
	req.Header.Set("Authorization", "Bearer "+h.cfg.EIDRPApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("eid request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("eid status %d: %s", resp.StatusCode, string(respBody))
	}

	var result eidStatusResp
	json.Unmarshal(respBody, &result)
	return &result, nil
}

// ── Helpers ─────────────────────────────────────────────────────────

func matchRedirectURI(registered []string, uri string) bool {
	for _, u := range registered {
		if u == uri {
			return true
		}
	}
	return false
}

func redirectWithError(w http.ResponseWriter, r *http.Request, redirectURI, state, errType, desc string) {
	u, _ := url.Parse(redirectURI)
	q := u.Query()
	q.Set("error", errType)
	q.Set("error_description", desc)
	if state != "" {
		q.Set("state", state)
	}
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func generateRandomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
