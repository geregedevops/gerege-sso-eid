package provider

import "context"

// CitizenInfo is our normalized output for the citizen lookup endpoint.
type CitizenInfo struct {
	RegNo           string `json:"reg_no"`
	LastName        string `json:"last_name"`
	FirstName       string `json:"first_name"`
	Surname         string `json:"surname,omitempty"`
	Gender          string `json:"gender,omitempty"`
	BirthDate       string `json:"birth_date,omitempty"`
	BirthPlace      string `json:"birth_place,omitempty"`
	Nationality     string `json:"nationality,omitempty"`
	CivilID         string `json:"civil_id,omitempty"`
	PassportNum     string `json:"passport_num,omitempty"`
	PassportAddress string `json:"passport_address,omitempty"`
	Image           string `json:"image,omitempty"`
}

// citizenResult maps the upstream /user/validate result object.
type citizenResult struct {
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	Surname         string `json:"surname"`
	Regnum          string `json:"regnum"`
	Gender          string `json:"gender"`
	BirthDate       string `json:"birthDateAsText"`
	BirthPlace      string `json:"birthPlace"`
	Nationality     string `json:"nationality"`
	CivilID         string `json:"civilId"`
	PassportNum     string `json:"passportNum"`
	PassportAddress string `json:"passportAddress"`
	Image           string `json:"image"`
	ResultCode      int    `json:"resultCode"`
}

type CitizenVerifyReq struct {
	RegNo     string `json:"reg_no"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// OrgInfo is our normalized output for the org lookup endpoint.
type OrgInfo struct {
	RegNo        string         `json:"reg_no"`
	Name         string         `json:"name"`
	Type         string         `json:"type,omitempty"`
	Capital      string         `json:"capital,omitempty"`
	CEO          string         `json:"ceo,omitempty"`
	CEORegNo     string         `json:"ceo_reg_no,omitempty"`
	CEOPosition  string         `json:"ceo_position,omitempty"`
	Phone        string         `json:"phone,omitempty"`
	Address      string         `json:"address,omitempty"`
	Industry     []string       `json:"industry,omitempty"`
	Founders     []OrgFounder   `json:"founders,omitempty"`
	StakeHolders []OrgStakeHolder `json:"stake_holders,omitempty"`
}

type OrgFounder struct {
	Name         string `json:"name"`
	RegNo        string `json:"reg_no"`
	Type         string `json:"type"`
	SharePercent string `json:"share_percent"`
}

type OrgStakeHolder struct {
	Name     string `json:"name"`
	RegNo    string `json:"reg_no"`
	Position string `json:"position"`
}

// orgResult maps the upstream /legalentity/info result object.
type orgResult struct {
	ResultCode int    `json:"resultCode"`
	GeneralR   orgCEO `json:"generalR"`
	Address    []struct {
		AddressStatus string `json:"addressStatus"`
		PhoneNumber   string `json:"phoneNumber"`
		StateCity     struct {
			Name string `json:"name"`
		} `json:"stateCity"`
		SoumDistrict struct {
			Name string `json:"name"`
		} `json:"soumDistrict"`
		BagKhoroo struct {
			Name string `json:"name"`
		} `json:"bagKhoroo"`
		Region struct {
			Name string `json:"name"`
		} `json:"region"`
		Door string `json:"door"`
	} `json:"address"`
	Capital []struct {
		RowStatusName string `json:"rowStatusName"`
		TotalAmount   string `json:"totalAmount"`
	} `json:"capital"`
	ChangeName []struct {
		RequestedName string `json:"requestedName"`
		CompanyType   string `json:"companyType"`
		CompanyRegnum string `json:"companyRegnum"`
	} `json:"changeName"`
	Induty []struct {
		IndustryName   string `json:"industryName"`
		IndustryStatus string `json:"industryStatus"`
	} `json:"induty"`
	Founder []struct {
		FirstName           string `json:"firstName"`
		LastName            string `json:"lastName"`
		StakeHolderRegnum   string `json:"stakeHolderRegnum"`
		StakeHolderTypeName string `json:"stakeHolderTypeName"`
		SharePercent        string `json:"sharePercent"`
		Status              string `json:"status"`
	} `json:"founder"`
	StakeHolders []struct {
		Firstname    string `json:"firstname"`
		Lastname     string `json:"lastname"`
		StateRegnum  string `json:"stateRegnum"`
		PositionName string `json:"positionName"`
		Status       string `json:"status"`
	} `json:"stakeHolders"`
}

type orgCEO struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Regnum       string `json:"regnum"`
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
