package model

import "time"

type Client struct {
	ID           string    `json:"id"`
	SecretHash   string    `json:"-"`
	Name         string    `json:"name"`
	RedirectURIs []string  `json:"redirect_uris"`
	Scopes       []string  `json:"scopes"`
	TenantID     string    `json:"tenant_id,omitempty"`
	LogoURL      string    `json:"logo_url,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
