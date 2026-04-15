package provider

import "context"

// upstreamResponse is the common wrapper for both citizen and org APIs.
type upstreamResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  any    `json:"-"` // decoded separately per type
}

// CitizenInfo is our normalized output for the citizen lookup endpoint.
type CitizenInfo struct {
	RegNo       string `json:"reg_no"`
	LastName    string `json:"last_name"`
	FirstName   string `json:"first_name"`
	Surname     string `json:"surname,omitempty"`
	Gender      string `json:"gender,omitempty"`
	BirthDate   string `json:"birth_date,omitempty"`
	Nationality string `json:"nationality,omitempty"`
}

// citizenResult maps the upstream /user/validate result object.
type citizenResult struct {
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Surname      string `json:"surname"`
	Regnum       string `json:"regnum"`
	Gender       string `json:"gender"`
	BirthDate    string `json:"birthDateAsText"`
	Nationality  string `json:"nationality"`
	ResultCode   int    `json:"resultCode"`
}

type CitizenVerifyReq struct {
	RegNo     string `json:"reg_no"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// OrgInfo is our normalized output for the org lookup endpoint.
type OrgInfo struct {
	RegNo     string   `json:"reg_no"`
	Name      string   `json:"name"`
	Type      string   `json:"type,omitempty"`
	CEO       string   `json:"ceo,omitempty"`
	Industry  []string `json:"industry,omitempty"`
}

// orgResult maps the upstream /legalentity/info result object.
type orgResult struct {
	ResultCode int    `json:"resultCode"`
	GeneralR   orgCEO `json:"generalR"`
	ChangeName []struct {
		RequestedName string `json:"requestedName"`
		CompanyType   string `json:"companyType"`
		CompanyRegnum string `json:"companyRegnum"`
		CreatedDate   string `json:"createdDate"`
	} `json:"changeName"`
	Induty []struct {
		IndustryName   string `json:"industryName"`
		IndustryStatus string `json:"industryStatus"`
	} `json:"induty"`
}

type orgCEO struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	PositionName string `json:"positionName"`
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
