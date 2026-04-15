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

// Lookup calls POST {baseURL}/legalentity/info with {"reg_no":"..."}
// Upstream wraps data in {"code":200,"status":"success","result":{...}}
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

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 128*1024))

	// Upstream returns 400 with {"status":"error"} when not found
	if resp.StatusCode != http.StatusOK {
		var env struct {
			Status string `json:"status"`
		}
		if json.Unmarshal(respBody, &env) == nil && env.Status == "error" {
			return nil, nil
		}
		return nil, fmt.Errorf("org lookup status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse envelope: {"code":200,"status":"success","result":{...}}
	var env struct {
		Code   int       `json:"code"`
		Status string    `json:"status"`
		Result orgResult `json:"result"`
	}
	if err := json.Unmarshal(respBody, &env); err != nil {
		return nil, fmt.Errorf("org lookup decode: %w", err)
	}
	if env.Status != "success" || env.Result.ResultCode != 200 {
		return nil, nil
	}

	r := env.Result

	// Get current company name from changeName (first entry is most recent)
	var name, companyType string
	if len(r.ChangeName) > 0 {
		name = r.ChangeName[0].RequestedName
		companyType = r.ChangeName[0].CompanyType
		regNo = r.ChangeName[0].CompanyRegnum
	}

	// CEO info
	var ceo string
	if r.GeneralR.FirstName != "" {
		ceo = r.GeneralR.LastName + " " + r.GeneralR.FirstName
	}

	// Phone and address from active address entry
	var phone, address string
	for _, addr := range r.Address {
		if addr.AddressStatus == "Тийм" {
			phone = addr.PhoneNumber
			parts := []string{}
			if addr.StateCity.Name != "" {
				parts = append(parts, addr.StateCity.Name)
			}
			if addr.SoumDistrict.Name != "" {
				parts = append(parts, addr.SoumDistrict.Name)
			}
			if addr.BagKhoroo.Name != "" {
				parts = append(parts, addr.BagKhoroo.Name)
			}
			if addr.Region.Name != "" {
				parts = append(parts, addr.Region.Name)
			}
			if addr.Door != "" {
				parts = append(parts, addr.Door)
			}
			address = strings.Join(parts, ", ")
			break
		}
	}

	// Capital (active)
	var capital string
	for _, c := range r.Capital {
		if c.RowStatusName == "Тийм" {
			capital = c.TotalAmount
			break
		}
	}

	// Active industries
	var industries []string
	for _, ind := range r.Induty {
		if ind.IndustryStatus == "Тийм" {
			industries = append(industries, ind.IndustryName)
		}
	}

	// Active founders
	var founders []OrgFounder
	for _, f := range r.Founder {
		if f.Status == "Тийм" {
			n := strings.TrimSpace(f.LastName + " " + f.FirstName)
			founders = append(founders, OrgFounder{
				Name:         n,
				RegNo:        f.StakeHolderRegnum,
				Type:         f.StakeHolderTypeName,
				SharePercent: f.SharePercent,
			})
		}
	}

	// Active stake holders (board members)
	var stakeHolders []OrgStakeHolder
	for _, s := range r.StakeHolders {
		if s.Status == "Тийм" {
			n := strings.TrimSpace(s.Lastname + " " + s.Firstname)
			stakeHolders = append(stakeHolders, OrgStakeHolder{
				Name:     n,
				RegNo:    s.StateRegnum,
				Position: s.PositionName,
			})
		}
	}

	return &OrgInfo{
		RegNo:        regNo,
		Name:         name,
		Type:         companyType,
		Capital:      capital,
		CEO:          ceo,
		CEORegNo:     r.GeneralR.Regnum,
		CEOPosition:  r.GeneralR.PositionName,
		Phone:        phone,
		Address:      address,
		Industry:     industries,
		Founders:     founders,
		StakeHolders: stakeHolders,
	}, nil
}

func (o *OrgHTTP) Verify(ctx context.Context, req OrgVerifyReq) (bool, error) {
	info, err := o.Lookup(ctx, req.RegNo)
	if err != nil {
		return false, err
	}
	if info == nil {
		return false, nil
	}

	return strings.EqualFold(strings.TrimSpace(info.Name), strings.TrimSpace(req.Name)), nil
}
