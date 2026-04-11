package handler

import (
	"net/http"
	"strings"

	"sso.gerege.mn/internal/model"
)

// ProfilePage — /user/profile хуудас (access_token cookie эсвэл Bearer header-ээр)
func (h *Handler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	// Extract token from cookie or Authorization header
	var accessToken string
	if cookie, err := r.Cookie("access_token"); err == nil {
		accessToken = cookie.Value
	}
	if accessToken == "" {
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, "Bearer ") {
			accessToken = strings.TrimPrefix(auth, "Bearer ")
		}
	}
	if accessToken == "" {
		// Query param fallback
		accessToken = r.URL.Query().Get("token")
	}

	if accessToken == "" {
		h.jsonError(w, 401, "unauthorized", "not logged in")
		return
	}

	var atData model.AccessTokenData
	if err := h.cfg.Cache.Get(r.Context(), "at:"+accessToken, &atData); err != nil {
		h.jsonError(w, 401, "unauthorized", "session expired")
		return
	}

	// Return profile as JSON (same as userinfo)
	resp := map[string]any{
		"sub":                      atData.Sub,
		"name":                     atData.Name,
		"given_name":               atData.GivenName,
		"family_name":              atData.FamilyName,
		"locale":                   "mn-MN",
		"identity_assurance_level": "high",
		"iss":                      h.cfg.Issuer,
		"tenant_id":               atData.TenantID,
		"tenant_role":             atData.TenantRole,
		"plan":                    atData.Plan,
	}

	h.jsonOK(w, resp)
}
