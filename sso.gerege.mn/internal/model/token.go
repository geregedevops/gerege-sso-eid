package model

import "time"

type IssuedToken struct {
	ID        int64     `json:"id"`
	ClientID  string    `json:"client_id"`
	Sub       string    `json:"sub"`
	Scope     string    `json:"scope"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
}
