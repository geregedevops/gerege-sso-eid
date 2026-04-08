package handler

import "net/http"

func (h *Handler) Discovery(w http.ResponseWriter, r *http.Request) {
	issuer := h.cfg.Issuer
	w.Header().Set("Cache-Control", "public, max-age=3600")
	h.jsonOK(w, map[string]any{
		"issuer":                 issuer,
		"authorization_endpoint": issuer + "/oauth/authorize",
		"token_endpoint":         issuer + "/oauth/token",
		"userinfo_endpoint":      issuer + "/oauth/userinfo",
		"jwks_uri":               issuer + "/.well-known/jwks.json",
		"revocation_endpoint":    issuer + "/oauth/revoke",
		"introspection_endpoint": issuer + "/oauth/introspect",
		"scopes_supported":       []string{"openid", "profile", "pos", "social", "payment"},
		"response_types_supported": []string{"code"},
		"grant_types_supported":    []string{"authorization_code", "refresh_token"},
		"subject_types_supported":  []string{"public"},
		"id_token_signing_alg_values_supported": []string{"ES256"},
		"token_endpoint_auth_methods_supported": []string{"client_secret_basic", "client_secret_post"},
		"claims_supported": []string{
			"sub", "iss", "aud", "exp", "iat", "nonce",
			"name", "given_name", "family_name", "locale",
			"cert_serial", "reg_no", "identity_assurance_level", "amr",
			"tenant_id", "tenant_role", "plan",
		},
		"auth_methods_supported": []string{"eid", "dan"},
		"ui_locales_supported": []string{"mn", "en"},
	})
}
