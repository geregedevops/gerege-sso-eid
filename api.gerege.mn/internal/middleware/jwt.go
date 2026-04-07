package middleware

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const SubKey contextKey = "sub"
const NameKey contextKey = "name"

type jwksCache struct {
	mu      sync.RWMutex
	keys    map[string]*ecdsa.PublicKey
	fetched time.Time
	uri     string
}

func JWTAuth(jwksURI string) func(http.Handler) http.Handler {
	cache := &jwksCache{uri: jwksURI, keys: make(map[string]*ecdsa.PublicKey)}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, `{"error":"missing bearer token"}`, http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")

			if err := cache.refresh(); err != nil {
				slog.Error("jwks fetch failed", "error", err)
				http.Error(w, `{"error":"auth service unavailable"}`, http.StatusServiceUnavailable)
				return
			}

			var sub, name string

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodECDSA); !ok {
					return nil, fmt.Errorf("unexpected alg: %v", t.Header["alg"])
				}
				kid, _ := t.Header["kid"].(string)
				cache.mu.RLock()
				key, ok := cache.keys[kid]
				cache.mu.RUnlock()
				if !ok {
					return nil, fmt.Errorf("unknown kid: %s", kid)
				}
				return key, nil
			})
			if err == nil {
				claims, ok := token.Claims.(jwt.MapClaims)
				if ok && token.Valid {
					sub, _ = claims["sub"].(string)
					name, _ = claims["name"].(string)
				}
			}

			// Fallback: opaque token → userinfo endpoint
			if sub == "" {
				ssoBase := strings.TrimSuffix(strings.TrimSuffix(jwksURI, "/.well-known/jwks.json"), "/")
				uiReq, _ := http.NewRequestWithContext(r.Context(), "GET", ssoBase+"/oauth/userinfo", nil)
				uiReq.Header.Set("Authorization", "Bearer "+tokenStr)
				uiResp, uiErr := http.DefaultClient.Do(uiReq)
				if uiErr != nil || uiResp.StatusCode != 200 {
					http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
					return
				}
				defer uiResp.Body.Close()
				var ui map[string]any
				json.NewDecoder(uiResp.Body).Decode(&ui)
				sub, _ = ui["sub"].(string)
				name, _ = ui["name"].(string)
				if sub == "" {
					http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
					return
				}
			}

			ctx := context.WithValue(r.Context(), SubKey, sub)
			ctx = context.WithValue(ctx, NameKey, name)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (c *jwksCache) refresh() error {
	c.mu.RLock()
	if time.Since(c.fetched) < time.Hour && len(c.keys) > 0 {
		c.mu.RUnlock()
		return nil
	}
	c.mu.RUnlock()

	resp, err := http.Get(c.uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var jwks struct {
		Keys []struct {
			KID string `json:"kid"`
			KTY string `json:"kty"`
			CRV string `json:"crv"`
			X   string `json:"x"`
			Y   string `json:"y"`
		} `json:"keys"`
	}
	if err := json.Unmarshal(body, &jwks); err != nil {
		return err
	}

	keys := make(map[string]*ecdsa.PublicKey)
	for _, k := range jwks.Keys {
		if k.KTY != "EC" || k.CRV != "P-256" {
			continue
		}
		xBytes, _ := base64.RawURLEncoding.DecodeString(k.X)
		yBytes, _ := base64.RawURLEncoding.DecodeString(k.Y)
		keys[k.KID] = &ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     new(big.Int).SetBytes(xBytes),
			Y:     new(big.Int).SetBytes(yBytes),
		}
	}

	c.mu.Lock()
	c.keys = keys
	c.fetched = time.Now()
	c.mu.Unlock()

	return nil
}
