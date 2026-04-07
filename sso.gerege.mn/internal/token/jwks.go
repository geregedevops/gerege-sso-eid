package token

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"math/big"
)

type JWK struct {
	KTY string `json:"kty"`
	CRV string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
	KID string `json:"kid"`
	Use string `json:"use"`
	ALG string `json:"alg"`
}

type JWKSet struct {
	Keys []JWK `json:"keys"`
}

func BuildJWKSet(pubKey *ecdsa.PublicKey, kid string) *JWKSet {
	return &JWKSet{
		Keys: []JWK{
			{
				KTY: "EC",
				CRV: "P-256",
				X:   base64URLEncodeBigInt(pubKey.X),
				Y:   base64URLEncodeBigInt(pubKey.Y),
				KID: kid,
				Use: "sig",
				ALG: "ES256",
			},
		},
	}
}

func ComputeKID(pubKey *ecdsa.PublicKey) string {
	thumbInput := map[string]string{
		"crv": "P-256",
		"kty": "EC",
		"x":   base64URLEncodeBigInt(pubKey.X),
		"y":   base64URLEncodeBigInt(pubKey.Y),
	}
	canonical, _ := json.Marshal(thumbInput)
	hash := sha256.Sum256(canonical)
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func base64URLEncodeBigInt(n *big.Int) string {
	b := n.Bytes()
	if len(b) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(b):], b)
		b = padded
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
