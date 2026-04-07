package smartid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

type InitiateResponse struct {
	SessionID        string `json:"session_id"`
	VerificationCode string `json:"verification_code"`
	ExpiresAt        string `json:"expires_at"`
}

type StatusResponse struct {
	State            string `json:"state"`
	Result           string `json:"result"`
	Sub              string `json:"sub"`
	Name             string `json:"name"`
	GivenName        string `json:"given_name"`
	FamilyName       string `json:"family_name"`
	CertSerial       string `json:"cert_serial"`
	VerificationCode string `json:"verification_code"`
}

func (c *Client) Initiate(ctx context.Context, nationalID, displayText, callbackURI string) (*InitiateResponse, error) {
	body := map[string]string{
		"national_id":  nationalID,
		"display_text": displayText,
		"callback_uri": callbackURI,
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/web/auth/sso/initiate", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("smartid.Initiate: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("smartid.Initiate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("smartid.Initiate: %d %s", resp.StatusCode, string(b))
	}

	var result InitiateResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return &result, nil
}

func (c *Client) Status(ctx context.Context, sessionID string) (*StatusResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/web/auth/sso/status/"+sessionID, nil)
	if err != nil {
		return nil, fmt.Errorf("smartid.Status: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("smartid.Status: %w", err)
	}
	defer resp.Body.Close()

	var result StatusResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return &result, nil
}
