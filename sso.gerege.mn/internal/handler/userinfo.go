package handler

import (
	"net/http"
	"strings"

	"gesign.mn/gerege-sso/internal/model"
)

func (h *Handler) UserInfo(w http.ResponseWriter, r *http.Request) {
	// Extract Bearer token
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		w.Header().Set("WWW-Authenticate", "Bearer")
		h.jsonError(w, 401, "invalid_token", "missing bearer token")
		return
	}
	accessToken := strings.TrimPrefix(auth, "Bearer ")

	// Lookup in Redis
	var atData model.AccessTokenData
	if err := h.cfg.Cache.Get(r.Context(), "at:"+accessToken, &atData); err != nil {
		w.Header().Set("WWW-Authenticate", "Bearer error=\"invalid_token\"")
		h.jsonError(w, 401, "invalid_token", "token expired or invalid")
		return
	}

	// Build response based on scope
	resp := map[string]any{
		"sub": atData.Sub,
		"iss": h.cfg.Issuer,
	}

	scopes := strings.Split(atData.Scope, " ")
	for _, s := range scopes {
		switch s {
		case "profile":
			resp["name"] = atData.Name
			resp["given_name"] = atData.GivenName
			resp["family_name"] = atData.FamilyName
			resp["locale"] = "mn-MN"
			resp["cert_serial"] = ""
			resp["cert_type"] = "AUTH"
			resp["identity_assurance_level"] = "high"
		case "pos":
			resp["tenant_id"] = atData.TenantID
			resp["tenant_role"] = atData.TenantRole
			resp["plan"] = atData.Plan
		case "social":
			resp["tenant_id"] = atData.TenantID
		case "payment":
			resp["tenant_id"] = atData.TenantID
		}
	}

	h.jsonOK(w, resp)
}
