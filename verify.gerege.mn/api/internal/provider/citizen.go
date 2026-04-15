package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type CitizenHTTP struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewCitizenHTTP(baseURL, apiKey string) *CitizenHTTP {
	return &CitizenHTTP{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 15 * time.Second},
	}
}

// Lookup calls the upstream user/validate API.
// Upstream: POST {baseURL}/user/validate  body: {"reg_no":"АА12345678"}
func (c *CitizenHTTP) Lookup(ctx context.Context, regNo string) (*CitizenInfo, error) {
	if c.baseURL == "" {
		return nil, fmt.Errorf("citizen API not configured")
	}

	body, _ := json.Marshal(map[string]string{"reg_no": regNo})
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/user/validate", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("citizen lookup request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("citizen lookup: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("citizen lookup status %d: %s", resp.StatusCode, string(respBody))
	}

	var info CitizenInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("citizen lookup decode: %w", err)
	}
	return &info, nil
}

// Verify calls the upstream user/validate API and compares the name fields.
func (c *CitizenHTTP) Verify(ctx context.Context, req CitizenVerifyReq) (bool, error) {
	info, err := c.Lookup(ctx, req.RegNo)
	if err != nil {
		return false, err
	}
	if info == nil {
		return false, nil
	}

	firstMatch := req.FirstName == "" || strings.EqualFold(strings.TrimSpace(info.FirstName), strings.TrimSpace(req.FirstName))
	lastMatch := req.LastName == "" || strings.EqualFold(strings.TrimSpace(info.LastName), strings.TrimSpace(req.LastName))

	return firstMatch && lastMatch, nil
}
