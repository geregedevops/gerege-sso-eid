package token

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IDTokenClaims struct {
	jwt.RegisteredClaims
	Nonce                  string   `json:"nonce,omitempty"`
	Name                   string   `json:"name,omitempty"`
	GivenName              string   `json:"given_name,omitempty"`
	FamilyName             string   `json:"family_name,omitempty"`
	Locale                 string   `json:"locale,omitempty"`
	CertSerial             string   `json:"cert_serial,omitempty"`
	CertType               string   `json:"cert_type,omitempty"`
	IdentityAssuranceLevel string   `json:"identity_assurance_level,omitempty"`
	AMR                    []string `json:"amr,omitempty"`
	// DAN (citizen registry) claim
	RegNo string `json:"reg_no,omitempty"`
	// Gerege-specific claims
	TenantID   string `json:"tenant_id,omitempty"`
	TenantRole string `json:"tenant_role,omitempty"`
	Plan       string `json:"plan,omitempty"`
}

type Issuer struct {
	privateKey *ecdsa.PrivateKey
	kid        string
	issuer     string
}

func NewIssuer(privKey *ecdsa.PrivateKey, kid, issuer string) *Issuer {
	return &Issuer{
		privateKey: privKey,
		kid:        kid,
		issuer:     issuer,
	}
}

func (i *Issuer) IssueIDToken(sub, aud, nonce, name, givenName, familyName, certSerial, regNo, tenantID, tenantRole, plan string) (string, error) {
	now := time.Now()

	amr := []string{"smartid", "pin1", "x509"}
	if regNo != "" {
		amr = []string{"dan", "sso_gov_mn"}
	}

	claims := IDTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    i.issuer,
			Subject:   sub,
			Audience:  jwt.ClaimStrings{aud},
			ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Nonce:                  nonce,
		Name:                   name,
		GivenName:              givenName,
		FamilyName:             familyName,
		Locale:                 "mn-MN",
		CertSerial:             certSerial,
		CertType:               "AUTH",
		IdentityAssuranceLevel: "high",
		AMR:                    amr,
		RegNo:                  regNo,
		TenantID:               tenantID,
		TenantRole:             tenantRole,
		Plan:                   plan,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = i.kid

	signed, err := token.SignedString(i.privateKey)
	if err != nil {
		return "", fmt.Errorf("token.IssueIDToken: %w", err)
	}
	return signed, nil
}

func (i *Issuer) VerifyIDToken(tokenString string, pubKey *ecdsa.PublicKey) (*IDTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &IDTokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return pubKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("token.Verify: %w", err)
	}
	claims, ok := token.Claims.(*IDTokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("token.Verify: invalid claims")
	}
	return claims, nil
}
