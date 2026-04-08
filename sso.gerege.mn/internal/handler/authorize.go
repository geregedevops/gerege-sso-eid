package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gesign.mn/gerege-sso/internal/model"
)

func (h *Handler) Authorize(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	clientID := q.Get("client_id")
	redirectURI := q.Get("redirect_uri")
	responseType := q.Get("response_type")
	scope := q.Get("scope")
	state := q.Get("state")
	nonce := q.Get("nonce")
	authMethod := q.Get("auth_method") // "eid" (default) or "dan"

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

	// Validate redirect_uri
	if !matchRedirectURI(client.RedirectURIs, redirectURI) {
		h.jsonError(w, 400, "invalid_request", "redirect_uri not registered")
		return
	}

	// Validate response_type
	if responseType != "code" {
		redirectWithError(w, r, redirectURI, state, "unsupported_response_type", "only code is supported")
		return
	}

	// Validate scope contains openid
	if !strings.Contains(scope, "openid") {
		redirectWithError(w, r, redirectURI, state, "invalid_scope", "openid scope required")
		return
	}

	// Generate session ID
	sessionID := generateRandomString(32)

	// Store session in Redis (10 min TTL)
	session := model.AuthSession{
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       scope,
		State:       state,
		Nonce:       nonce,
		AuthMethod:  authMethod,
	}
	if err := h.cfg.Cache.Set(r.Context(), "sso:"+sessionID, session, 10*time.Minute); err != nil {
		logErr("authorize: redis set", err)
		h.jsonError(w, 500, "server_error", "internal error")
		return
	}

	if authMethod == "dan" {
		// Redirect to sso.gov.mn for DAN verification
		h.redirectToDAN(w, r, sessionID)
		return
	}

	// Default: Redirect to e-id.mn
	eidURL := fmt.Sprintf("%s/auth?session=%s&callback_uri=%s/callback/eid&purpose=sso:%s",
		h.cfg.EIDBaseURL,
		url.QueryEscape(sessionID),
		url.QueryEscape(h.cfg.Issuer),
		url.QueryEscape(clientID),
	)
	http.Redirect(w, r, eidURL, http.StatusFound)
}

func (h *Handler) redirectToDAN(w http.ResponseWriter, r *http.Request, sessionID string) {
	// state = base64({"redirect_url":"...", "session":"..."})
	// redirect_url = dan.gerege.mn gateway exchanges code, then redirects here with citizen data
	stateJSON := fmt.Sprintf(`{"redirect_url":"%s/callback/dan","session":"%s"}`, h.cfg.Issuer, sessionID)
	stateB64 := base64.RawURLEncoding.EncodeToString([]byte(stateJSON))

	danURL := fmt.Sprintf("https://sso.gov.mn/login?state=%s&grant_type=authorization_code&response_type=code&client_id=%s&scope=%s&redirect_uri=%s",
		url.QueryEscape(stateB64),
		url.QueryEscape(h.cfg.DANClientID),
		url.QueryEscape(h.cfg.DANScope),
		url.QueryEscape(h.cfg.DANCallbackURI),
	)
	http.Redirect(w, r, danURL, http.StatusFound)
}

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
