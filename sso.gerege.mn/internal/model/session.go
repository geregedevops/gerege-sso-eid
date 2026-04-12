package model

// AuthSession is stored in Redis during the authorize flow
type AuthSession struct {
	ClientID    string `json:"client_id"`
	RedirectURI string `json:"redirect_uri"`
	Scope       string `json:"scope"`
	State       string `json:"state"`
	Nonce       string `json:"nonce"`
}

// AuthCode is stored in Redis after gerege.mn or DAN callback
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
}

// AccessTokenData is stored in Redis for opaque access tokens
type AccessTokenData struct {
	Sub        string `json:"sub"`
	ClientID   string `json:"client_id"`
	Scope      string `json:"scope"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	CertSerial string `json:"cert_serial,omitempty"`
	RegNo      string `json:"reg_no,omitempty"`
	TenantID   string `json:"tenant_id,omitempty"`
	TenantRole string `json:"tenant_role,omitempty"`
	Plan       string `json:"plan,omitempty"`
	IssuedAt   int64  `json:"iat"`
	ExpiresAt  int64  `json:"exp"`
}
