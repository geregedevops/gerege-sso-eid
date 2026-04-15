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

// Lookup calls POST {baseURL}/user/validate with {"reg_no":"..."}
// Upstream wraps data in {"code":200,"status":"success","result":{...}}
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

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))

	// Upstream returns 400 with {"status":"error"} when not found
	if resp.StatusCode != http.StatusOK {
		var env struct {
			Status string `json:"status"`
		}
		if json.Unmarshal(respBody, &env) == nil && env.Status == "error" {
			return nil, nil
		}
		return nil, fmt.Errorf("citizen lookup status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse envelope: {"code":200,"status":"success","result":{...}}
	var env struct {
		Code   int           `json:"code"`
		Status string        `json:"status"`
		Result citizenResult `json:"result"`
	}
	if err := json.Unmarshal(respBody, &env); err != nil {
		return nil, fmt.Errorf("citizen lookup decode: %w", err)
	}
	if env.Status != "success" || env.Result.ResultCode != 200 {
		return nil, nil
	}

	r := env.Result
	return &CitizenInfo{
		RegNo:           r.Regnum,
		LastName:        r.Lastname,
		FirstName:       r.Firstname,
		Surname:         r.Surname,
		Gender:          r.Gender,
		BirthDate:       r.BirthDate,
		BirthPlace:      r.BirthPlace,
		Nationality:     r.Nationality,
		CivilID:         r.CivilID,
		PassportNum:     r.PassportNum,
		PassportAddress: r.PassportAddress,
		Image:           r.Image,
	}, nil
}

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
