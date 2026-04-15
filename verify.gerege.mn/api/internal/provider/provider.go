package provider

import "context"

// CitizenInfo maps the response from the /user/validate upstream API.
type CitizenInfo struct {
	RegNo       string `json:"reg_no"`
	LastName    string `json:"last_name,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	Gender      string `json:"gender,omitempty"`
}

type CitizenVerifyReq struct {
	RegNo     string `json:"reg_no"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// OrgInfo maps the response from the /legalentity/info upstream API.
type OrgInfo struct {
	RegNo           string `json:"reg_no"`
	Name            string `json:"name"`
	NameShort       string `json:"name_short,omitempty"`
	Type            string `json:"type,omitempty"`
	Status          string `json:"status,omitempty"`
	EstablishedDate string `json:"established_date,omitempty"`
}

type OrgVerifyReq struct {
	RegNo string `json:"reg_no"`
	Name  string `json:"name"`
}

type CitizenProvider interface {
	Lookup(ctx context.Context, regNo string) (*CitizenInfo, error)
	Verify(ctx context.Context, req CitizenVerifyReq) (bool, error)
}

type OrgProvider interface {
	Lookup(ctx context.Context, regNo string) (*OrgInfo, error)
	Verify(ctx context.Context, req OrgVerifyReq) (bool, error)
}
