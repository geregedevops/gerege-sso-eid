package dan

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// FieldMap maps sso.gov.mn field names to our field names.
var FieldMap = map[string]string{
	"regnum":               "reg_no",
	"surname":              "surname",
	"firstname":            "given_name",
	"lastname":             "family_name",
	"civilId":              "civil_id",
	"gender":               "gender",
	"birthDateAsText":      "birth_date",
	"birthPlace":           "birth_place",
	"nationality":          "nationality",
	"aimagCityName":        "aimag_name",
	"aimagCityCode":        "aimag_code",
	"soumDistrictName":     "sum_name",
	"soumDistrictCode":     "sum_code",
	"bagKhorooName":        "bag_name",
	"bagKhorooCode":        "bag_code",
	"addressDetail":        "address_detail",
	"passportAddress":      "passport_address",
	"passportExpireDate":   "passport_expire_date",
	"passportIssueDate":    "passport_issue_date",
	"addressApartmentName": "apartment_name",
	"addressStreetName":    "street_name",
}

// GetCitizenData fetches citizen data from sso.gov.mn using the access token.
func GetCitizenData(cfg Config, accessToken string) (map[string]string, error) {
	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"scope":         {cfg.Scope},
	}

	req, err := http.NewRequest("POST", cfg.ServiceURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("service returned %d: %s", resp.StatusCode, string(body))
	}

	var rawArr []any
	if err := json.Unmarshal(body, &rawArr); err != nil {
		return nil, err
	}

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
			if r, ok := svcObj["response"].(map[string]any); ok {
				citizen = r
				break
			}
		}
		if citizen != nil {
			break
		}
	}

	if citizen == nil {
		return nil, fmt.Errorf("no citizen data in response")
	}

	result := make(map[string]string)
	for ssoKey, ourKey := range FieldMap {
		if v, ok := citizen[ssoKey]; ok && v != nil {
			s := fmt.Sprintf("%v", v)
			if s != "" && s != "<nil>" {
				result[ourKey] = s
			}
		}
	}

	if img, ok := citizen["image"].(string); ok && img != "" {
		result["image"] = img
	}

	return result, nil
}
