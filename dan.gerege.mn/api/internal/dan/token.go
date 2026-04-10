package dan

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	ClientID     string
	ClientSecret string
	Scope        string
	CallbackURI  string
	TokenURL     string
	ServiceURL   string
}

var httpClient = &http.Client{Timeout: 15 * time.Second}

// GetAccessToken exchanges an authorization code for an access token from sso.gov.mn.
func GetAccessToken(cfg Config, code string) (string, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {cfg.CallbackURI},
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
	}

	resp, err := httpClient.Post(cfg.TokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("POST token: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	slog.Debug("dan token response", "status", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token returned %d: %s", resp.StatusCode, string(body))
	}

	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}
	if at, ok := raw["access_token"].(string); ok && at != "" {
		return at, nil
	}
	return "", fmt.Errorf("no access_token in response")
}
