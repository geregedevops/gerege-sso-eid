package handler

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gesign.mn/gerege-sso/internal/model"
)

func (h *Handler) Introspect(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.jsonError(w, 400, "invalid_request", "malformed form data")
		return
	}

	// Client authentication (Basic auth)
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
		h.jsonError(w, 401, "invalid_client", "unknown client")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(client.SecretHash), []byte(clientSecret)); err != nil {
		h.jsonError(w, 401, "invalid_client", "invalid client credentials")
		return
	}

	tokenStr := r.FormValue("token")
	if tokenStr == "" {
		h.jsonOK(w, map[string]any{"active": false})
		return
	}

	var atData model.AccessTokenData
	if err := h.cfg.Cache.Get(r.Context(), "at:"+tokenStr, &atData); err != nil {
		h.jsonOK(w, map[string]any{"active": false})
		return
	}

	h.jsonOK(w, map[string]any{
		"active":      true,
		"sub":         atData.Sub,
		"scope":       atData.Scope,
		"client_id":   atData.ClientID,
		"exp":         atData.ExpiresAt,
		"iat":         atData.IssuedAt,
		"token_type":  "Bearer",
		"iss":         h.cfg.Issuer,
		"tenant_id":   atData.TenantID,
		"tenant_role": atData.TenantRole,
		"plan":        atData.Plan,
	})
}
