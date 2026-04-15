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

type OrgHTTP struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewOrgHTTP(baseURL, apiKey string) *OrgHTTP {
	return &OrgHTTP{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 15 * time.Second},
	}
}

// Lookup calls the upstream legalentity/info API.
// Upstream: POST {baseURL}/legalentity/info  body: {"reg_no":"6235972"}
func (o *OrgHTTP) Lookup(ctx context.Context, regNo string) (*OrgInfo, error) {
	if o.baseURL == "" {
		return nil, fmt.Errorf("org API not configured")
	}

	body, _ := json.Marshal(map[string]string{"reg_no": regNo})
	req, err := http.NewRequestWithContext(ctx, "POST", o.baseURL+"/legalentity/info", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("org lookup request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if o.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+o.apiKey)
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("org lookup: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("org lookup status %d: %s", resp.StatusCode, string(respBody))
	}

	var info OrgInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("org lookup decode: %w", err)
	}
	return &info, nil
}

// Verify calls the upstream legalentity/info API and compares the name.
func (o *OrgHTTP) Verify(ctx context.Context, req OrgVerifyReq) (bool, error) {
	info, err := o.Lookup(ctx, req.RegNo)
	if err != nil {
		return false, err
	}
	if info == nil {
		return false, nil
	}

	// Compare name case-insensitively
	return strings.EqualFold(strings.TrimSpace(info.Name), strings.TrimSpace(req.Name)), nil
}
