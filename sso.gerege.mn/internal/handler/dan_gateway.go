package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// danGatewayState is the state we encode before redirecting to sso.gov.mn
type danGatewayState struct {
	RedirectURL string `json:"redirect_url"`
}

// DANGatewayAuthorized handles the callback from sso.gov.mn on dan.gerege.mn/authorized.
// It exchanges the authorization code for citizen data via sso.gov.mn token endpoint,
// then redirects to redirect_url (sso.gerege.mn) with all citizen data as query params.
func (h *Handler) DANGatewayAuthorized(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	code := q.Get("code")
	stateB64 := q.Get("state")

	slog.Info("dan_gateway: received callback", "has_code", code != "", "has_state", stateB64 != "")

	if code == "" || stateB64 == "" {
		h.jsonError(w, 400, "invalid_request", "missing code or state")
		return
	}

	// Decode state to get redirect_url
	stateBytes, err := base64.RawURLEncoding.DecodeString(stateB64)
	if err != nil {
		stateBytes, err = base64.StdEncoding.DecodeString(stateB64)
		if err != nil {
			h.jsonError(w, 400, "invalid_request", "invalid state encoding")
			return
		}
	}

	var state danGatewayState
	if err := json.Unmarshal(stateBytes, &state); err != nil {
		h.jsonError(w, 400, "invalid_request", "invalid state format")
		return
	}

	if state.RedirectURL == "" {
		h.jsonError(w, 400, "invalid_request", "missing redirect_url in state")
		return
	}

	// Exchange code for citizen data at sso.gov.mn token endpoint
	citizenData, err := h.exchangeDANCode(code)
	if err != nil {
		slog.Error("dan_gateway: code exchange failed", "error", err)
		h.jsonError(w, 502, "upstream_error", "failed to exchange code with sso.gov.mn")
		return
	}

	slog.Info("dan_gateway: citizen data received",
		"reg_no", citizenData["reg_no"],
		"surname", citizenData["surname"],
		"given_name", citizenData["given_name"],
	)

	// Redirect to redirect_url with all citizen data as query params
	redirectURL, err := url.Parse(state.RedirectURL)
	if err != nil {
		h.jsonError(w, 400, "invalid_request", "invalid redirect_url")
		return
	}

	params := redirectURL.Query()
	for k, v := range citizenData {
		if v != "" {
			params.Set(k, v)
		}
	}
	redirectURL.RawQuery = params.Encode()

	slog.Info("dan_gateway: redirecting", "url", redirectURL.String())
	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}

// exchangeDANCode exchanges the authorization code with sso.gov.mn for citizen data
func (h *Handler) exchangeDANCode(code string) (map[string]string, error) {
	tokenURL := "https://api.sso.gov.mn/oauth/api/token"

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", h.cfg.DANCallbackURI)
	form.Set("client_id", h.cfg.DANClientID)
	form.Set("client_secret", h.cfg.DANGatewaySecret)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("POST token endpoint: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	slog.Info("dan_gateway: token response", "status", resp.StatusCode, "body", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	// Parse response — sso.gov.mn may return citizen data in various formats
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	// Extract citizen fields from response
	result := make(map[string]string)
	citizenFields := []string{
		"reg_no", "surname", "given_name", "family_name",
		"civil_id", "gender", "birth_date",
		"phone_no", "email", "nationality",
		"aimag_name", "sum_name", "bag_name", "address_detail",
		"aimag_id", "aimag_code", "sum_id", "sum_code",
		"bag_id", "bag_code",
		"residential_aimag_name", "residential_sum_name",
		"residential_bag_name", "residential_address_detail",
	}

	// Try top-level fields
	for _, field := range citizenFields {
		if v, ok := raw[field]; ok {
			result[field] = fmt.Sprintf("%v", v)
		}
	}

	// Try nested "result" or "data" object
	for _, key := range []string{"result", "data", "citizen"} {
		if nested, ok := raw[key]; ok {
			if m, ok := nested.(map[string]any); ok {
				for _, field := range citizenFields {
					if v, ok := m[field]; ok {
						result[field] = fmt.Sprintf("%v", v)
					}
				}
			}
		}
	}

	if result["reg_no"] == "" {
		return nil, fmt.Errorf("no reg_no in response: %s", string(body))
	}

	return result, nil
}
