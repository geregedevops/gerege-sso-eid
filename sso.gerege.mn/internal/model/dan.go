package model

import "time"

type DANClient struct {
	ID           string    `json:"id"`
	SecretHash   string    `json:"-"`
	HMACKey      string    `json:"-"`
	Name         string    `json:"name"`
	CallbackURLs []string  `json:"callback_urls"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
