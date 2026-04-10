package store

import (
	"context"
	"time"

	"sso.gerege.mn/internal/model"
)

// Cache defines the caching interface used by handlers.
type Cache interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string, dest any) error
	Del(ctx context.Context, key string) error
	GetAndDel(ctx context.Context, key string, dest any) error
	Incr(ctx context.Context, key string, ttl time.Duration) (int64, error)
	SetString(ctx context.Context, key, value string, ttl time.Duration) error
	GetString(ctx context.Context, key string) (string, error)
}

// DB defines the database interface used by handlers.
type DB interface {
	GetClient(ctx context.Context, clientID string) (*model.Client, error)
	GetTenantMember(ctx context.Context, tenantID, sub string) (string, error)
	GetTenantPlan(ctx context.Context, tenantID string) (string, error)
	RecordIssuedToken(ctx context.Context, clientID, sub, scope string, expiresAt time.Time) error

	// DAN clients
	GetDANClient(ctx context.Context, clientID string) (*model.DANClient, error)
	ListDANClients(ctx context.Context) ([]model.DANClient, error)
	CreateDANClient(ctx context.Context, id, secretHash, hmacKey, name string, callbackURLs []string) error
	DeactivateDANClient(ctx context.Context, clientID string) error
}
