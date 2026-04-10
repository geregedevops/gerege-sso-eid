package handler

import (
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"sso.gerege.mn/internal/model"
)

func (h *Handler) Token(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.jsonError(w, 400, "invalid_request", "malformed form data")
		return
	}

	grantType := r.FormValue("grant_type")
	if grantType != "authorization_code" {
		h.jsonError(w, 400, "unsupported_grant_type", "only authorization_code supported")
		return
	}

	code := r.FormValue("code")
	redirectURI := r.FormValue("redirect_uri")

	// Client authentication (Basic or POST)
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		clientID = r.FormValue("client_id")
		clientSecret = r.FormValue("client_secret")
	}

	if clientID == "" || clientSecret == "" {
		h.jsonError(w, 401, "invalid_client", "client credentials required")
		return
	}

	// Verify client first, then rate limit (so invalid client_ids don't exhaust counters)
	client, err := h.cfg.DB.GetClient(r.Context(), clientID)
	if err != nil || client == nil || !client.IsActive {
		h.jsonError(w, 401, "invalid_client", "unknown client")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(client.SecretHash), []byte(clientSecret)); err != nil {
		h.jsonError(w, 401, "invalid_client", "invalid client credentials")
		return
	}

	// Rate limiting: 10 req/min per client (after validation)
	rlKey := "rl:token:" + clientID
	count, err := h.cfg.Cache.Incr(r.Context(), rlKey, time.Minute)
	if err == nil && count > 10 {
		h.jsonError(w, 429, "rate_limit", "too many requests")
		return
	}

	// Get and delete code (single use)
	var codeData model.AuthCode
	if err := h.cfg.Cache.GetAndDel(r.Context(), "code:"+code, &codeData); err != nil {
		h.jsonError(w, 400, "invalid_grant", "invalid or expired authorization code")
		return
	}

	// Verify code belongs to this client
	if codeData.ClientID != clientID {
		h.jsonError(w, 400, "invalid_grant", "code was issued to different client")
		return
	}

	// Verify redirect_uri
	if redirectURI != "" && redirectURI != codeData.RedirectURI {
		h.jsonError(w, 400, "invalid_grant", "redirect_uri mismatch")
		return
	}

	// Resolve tenant context if pos/social/payment scope requested
	var tenantID, tenantRole, plan string
	scopes := strings.Split(codeData.Scope, " ")
	needsTenant := false
	for _, s := range scopes {
		if s == "pos" || s == "social" || s == "payment" {
			needsTenant = true
			break
		}
	}
	if needsTenant && client.TenantID != "" {
		tenantID = client.TenantID
		role, _ := h.cfg.DB.GetTenantMember(r.Context(), tenantID, codeData.Sub)
		if role != "" {
			tenantRole = role
		}
		p, _ := h.cfg.DB.GetTenantPlan(r.Context(), tenantID)
		if p != "" {
			plan = p
		}
	}

	// Issue ID Token
	idToken, err := h.cfg.TokenIssuer.IssueIDToken(
		codeData.Sub, clientID, codeData.Nonce,
		codeData.Name, codeData.GivenName, codeData.FamilyName,
		codeData.CertSerial, codeData.RegNo,
		tenantID, tenantRole, plan,
	)
	if err != nil {
		logErr("token: issue id_token", err)
		h.jsonError(w, 500, "server_error", "failed to issue token")
		return
	}

	// Issue opaque access token
	accessToken := generateRandomString(32)
	now := time.Now()
	expiresIn := 3600

	atData := model.AccessTokenData{
		Sub:        codeData.Sub,
		ClientID:   clientID,
		Scope:      codeData.Scope,
		Name:       codeData.Name,
		GivenName:  codeData.GivenName,
		FamilyName: codeData.FamilyName,
		CertSerial: codeData.CertSerial,
		RegNo:      codeData.RegNo,
		TenantID:   tenantID,
		TenantRole: tenantRole,
		Plan:       plan,
		IssuedAt:   now.Unix(),
		ExpiresAt:  now.Add(time.Duration(expiresIn) * time.Second).Unix(),
	}
	if err := h.cfg.Cache.Set(r.Context(), "at:"+accessToken, atData, time.Duration(expiresIn)*time.Second); err != nil {
		logErr("token: redis set at", err)
		h.jsonError(w, 500, "server_error", "internal error")
		return
	}

	// Audit log
	h.cfg.DB.RecordIssuedToken(r.Context(), clientID, codeData.Sub, codeData.Scope, now.Add(time.Duration(expiresIn)*time.Second))

	// Response
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	h.jsonOK(w, map[string]any{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   expiresIn,
		"id_token":     idToken,
		"scope":        codeData.Scope,
	})
}
