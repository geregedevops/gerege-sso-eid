package handler

import (
	"net/http"
	"net/url"
	"time"

	"sso.gerege.mn/internal/model"
)

func (h *Handler) EIDCallback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	sessionID := q.Get("session")
	sub := q.Get("sub")
	name := q.Get("name")
	givenName := q.Get("given_name")
	familyName := q.Get("family_name")
	certSerial := q.Get("cert_serial")

	if sessionID == "" {
		h.jsonError(w, 400, "invalid_request", "missing session")
		return
	}

	// Get session from Redis
	var session model.AuthSession
	if err := h.cfg.Cache.Get(r.Context(), "sso:"+sessionID, &session); err != nil {
		logErr("eid_callback: session not found", err)
		h.jsonError(w, 400, "invalid_request", "session expired or not found")
		return
	}

	// Handle error from e-id.mn (user cancelled, timeout, etc.)
	if errParam := q.Get("error"); errParam != "" || sub == "" {
		h.cfg.Cache.Del(r.Context(), "sso:"+sessionID)
		errType := errParam
		if errType == "" {
			errType = "access_denied"
		}
		redirectWithError(w, r, session.RedirectURI, session.State, errType, "authentication failed")
		return
	}

	// OCSP check on cert_serial (fail closed — reject revoked certificates)
	if certSerial != "" && h.cfg.OCSP != nil {
		if err := h.cfg.OCSP.Check(r.Context(), certSerial); err != nil {
			logErr("eid_callback: OCSP check failed", err)
			h.cfg.Cache.Del(r.Context(), "sso:"+sessionID)
			redirectWithError(w, r, session.RedirectURI, session.State, "access_denied", "certificate validation failed")
			return
		}
	}

	// Generate auth code
	authCode := generateRandomString(32)

	// Store code in Redis (5 min, single use)
	codeData := model.AuthCode{
		Sub:         sub,
		Name:        name,
		GivenName:   givenName,
		FamilyName:  familyName,
		CertSerial:  certSerial,
		ClientID:    session.ClientID,
		RedirectURI: session.RedirectURI,
		Scope:       session.Scope,
		Nonce:       session.Nonce,
	}
	if err := h.cfg.Cache.Set(r.Context(), "code:"+authCode, codeData, 5*time.Minute); err != nil {
		logErr("eid_callback: redis set code", err)
		h.jsonError(w, 500, "server_error", "internal error")
		return
	}

	// Delete session
	h.cfg.Cache.Del(r.Context(), "sso:"+sessionID)

	// Redirect to client with code + state
	redirectWithCode(w, r, session.RedirectURI, authCode, session.State)
}

func redirectWithCode(w http.ResponseWriter, r *http.Request, redirectURI, code, state string) {
	u, _ := url.Parse(redirectURI)
	q := u.Query()
	q.Set("code", code)
	if state != "" {
		q.Set("state", state)
	}
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}
