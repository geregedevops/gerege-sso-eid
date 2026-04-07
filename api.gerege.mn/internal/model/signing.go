package model

import "time"

type SigningSession struct {
	ID               string    `json:"id"`
	RequesterSub     string    `json:"requester_sub"`
	SignerSub        string    `json:"signer_sub,omitempty"`
	SignerName       string    `json:"signer_name,omitempty"`
	SignerReg        string    `json:"signer_reg,omitempty"`
	Status           string    `json:"status"`
	SmartIDSession   string    `json:"smartid_session,omitempty"`
	VerificationCode string    `json:"verification_code,omitempty"`
	DocumentName     string    `json:"document_name"`
	DocumentHash     string    `json:"document_hash"`
	DocumentSize     int       `json:"document_size"`
	DocumentPath     string    `json:"-"`
	SignedDocPath    string    `json:"-"`
	CertSerial       string    `json:"cert_serial,omitempty"`
	ErrorMessage     string    `json:"error_message,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	ExpiresAt        time.Time `json:"expires_at"`
}
