package handler

import (
	"net/http"

	"sso.gerege.mn/internal/model"
)

func (h *Handler) Introspect(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.jsonError(w, 400, "invalid_request", "malformed form data")
		return
	}

	if h.authenticateClient(w, r) == nil {
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
