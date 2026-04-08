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
// Flow: code → access_token (oauth2/token) → citizen data (oauth2/api/v1/service)
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

	// Step 1: Exchange code for access_token
	accessToken, err := h.danGetAccessToken(code)
	if err != nil {
		slog.Error("dan_gateway: token exchange failed", "error", err)
		h.jsonError(w, 502, "upstream_error", "failed to get access token from sso.gov.mn")
		return
	}

	// Step 2: Call service endpoint with access_token to get citizen data
	citizenData, err := h.danGetCitizenData(accessToken)
	if err != nil {
		slog.Error("dan_gateway: service call failed", "error", err)
		h.jsonError(w, 502, "upstream_error", "failed to get citizen data from sso.gov.mn")
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

// danGetAccessToken exchanges authorization code for access_token at sso.gov.mn
func (h *Handler) danGetAccessToken(code string) (string, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", h.cfg.DANCallbackURI)
	form.Set("client_id", h.cfg.DANClientID)
	form.Set("client_secret", h.cfg.DANClientSecret)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Post(h.cfg.DANTokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("POST token endpoint: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	slog.Info("dan_gateway: token response", "status", resp.StatusCode, "body", string(body))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return "", fmt.Errorf("parse token response: %w", err)
	}

	// Try access_token field
	if at, ok := raw["access_token"].(string); ok && at != "" {
		return at, nil
	}

	return "", fmt.Errorf("no access_token in response: %s", string(body))
}

// danGetCitizenData calls the service endpoint with access_token to retrieve citizen data
func (h *Handler) danGetCitizenData(accessToken string) (map[string]string, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	req, err := http.NewRequest("POST", h.cfg.DANServiceURL, strings.NewReader(url.Values{
		"grant_type":   {"client_credentials"},
		"client_id":    {h.cfg.DANClientID},
		"client_secret": {h.cfg.DANClientSecret},
		"scope":        {h.cfg.DANScope},
	}.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create service request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("POST service endpoint: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read service response: %w", err)
	}

	slog.Info("dan_gateway: service response", "status", resp.StatusCode, "body_len", len(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("service endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	// sso.gov.mn returns a JSON array:
	// [{citizen_loginType:7}, {services: {WS100101_getCitizenIDCardInfo: {response: {...}}}}]
	var rawArr []any
	if err := json.Unmarshal(body, &rawArr); err != nil {
		return nil, fmt.Errorf("parse service response: %w", err)
	}

	// Navigate to services -> WS100101_getCitizenIDCardInfo -> response
	var citizen map[string]any
	for _, item := range rawArr {
		obj, ok := item.(map[string]any)
		if !ok {
			continue
		}
		services, ok := obj["services"].(map[string]any)
		if !ok {
			continue
		}
		for _, svc := range services {
			svcObj, ok := svc.(map[string]any)
			if !ok {
				continue
			}
			if resp, ok := svcObj["response"].(map[string]any); ok {
				citizen = resp
				break
			}
		}
		if citizen != nil {
			break
		}
	}

	if citizen == nil {
		return nil, fmt.Errorf("no citizen data in response: %s", string(body))
	}

	// Map sso.gov.mn field names to our field names
	fieldMap := map[string]string{
		"regnum":              "reg_no",
		"surname":             "surname",
		"firstname":           "given_name",
		"lastname":            "family_name",
		"civilId":             "civil_id",
		"gender":              "gender",
		"birthDateAsText":     "birth_date",
		"nationality":         "nationality",
		"aimagCityName":       "aimag_name",
		"aimagCityCode":       "aimag_code",
		"soumDistrictName":    "sum_name",
		"soumDistrictCode":    "sum_code",
		"bagKhorooName":       "bag_name",
		"bagKhorooCode":       "bag_code",
		"addressDetail":       "address_detail",
		"passportAddress":     "passport_address",
		"birthPlace":          "birth_place",
	}

	result := make(map[string]string)
	for ssoKey, ourKey := range fieldMap {
		if v, ok := citizen[ssoKey]; ok && v != nil {
			s := fmt.Sprintf("%v", v)
			if s != "" && s != "<nil>" {
				result[ourKey] = s
			}
		}
	}

	if result["reg_no"] == "" {
		return nil, fmt.Errorf("no reg_no in response: %s", string(body))
	}

	slog.Info("dan_gateway: parsed citizen", "reg_no", result["reg_no"], "given_name", result["given_name"])

	return result, nil
}
