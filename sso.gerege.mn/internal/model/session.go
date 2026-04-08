package model

// AuthSession is stored in Redis during the authorize flow
type AuthSession struct {
	ClientID    string `json:"client_id"`
	RedirectURI string `json:"redirect_uri"`
	Scope       string `json:"scope"`
	State       string `json:"state"`
	Nonce       string `json:"nonce"`
	AuthMethod  string `json:"auth_method,omitempty"` // "eid" (default) or "dan"
}

// AuthCode is stored in Redis after e-id.mn or DAN callback
type AuthCode struct {
	Sub         string `json:"sub"`
	Name        string `json:"name"`
	GivenName   string `json:"given_name"`
	FamilyName  string `json:"family_name"`
	CertSerial  string `json:"cert_serial"`
	RegNo       string `json:"reg_no,omitempty"`
	ClientID    string `json:"client_id"`
	RedirectURI string `json:"redirect_uri"`
	Scope       string `json:"scope"`
	Nonce       string `json:"nonce"`
	// DAN citizen data
	Surname       string `json:"surname,omitempty"`
	CivilID       string `json:"civil_id,omitempty"`
	Gender        string `json:"gender,omitempty"`
	BirthDate     string `json:"birth_date,omitempty"`
	Nationality   string `json:"nationality,omitempty"`
	PhoneNo       string `json:"phone_no,omitempty"`
	Email         string `json:"email,omitempty"`
	AimagName     string `json:"aimag_name,omitempty"`
	SumName       string `json:"sum_name,omitempty"`
	BagName       string `json:"bag_name,omitempty"`
	AddressDetail string `json:"address_detail,omitempty"`
}

// AccessTokenData is stored in Redis for opaque access tokens
type AccessTokenData struct {
	Sub        string `json:"sub"`
	ClientID   string `json:"client_id"`
	Scope      string `json:"scope"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	RegNo      string `json:"reg_no,omitempty"`
	TenantID   string `json:"tenant_id,omitempty"`
	TenantRole string `json:"tenant_role,omitempty"`
	Plan       string `json:"plan,omitempty"`
	IssuedAt   int64  `json:"iat"`
	ExpiresAt  int64  `json:"exp"`
}
