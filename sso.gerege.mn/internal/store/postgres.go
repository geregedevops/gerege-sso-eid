package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"sso.gerege.mn/internal/model"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, databaseURL string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("store.NewPostgres: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("store.NewPostgres ping: %w", err)
	}
	return &Postgres{pool: pool}, nil
}

func (p *Postgres) Close() {
	p.pool.Close()
}

func (p *Postgres) GetClient(ctx context.Context, clientID string) (*model.Client, error) {
	var c model.Client
	err := p.pool.QueryRow(ctx,
		`SELECT id, secret_hash, name, redirect_uris, scopes, COALESCE(tenant_id,''), COALESCE(logo_url,''), is_active, created_at, updated_at
		 FROM sso_clients WHERE id = $1`, clientID,
	).Scan(&c.ID, &c.SecretHash, &c.Name, &c.RedirectURIs, &c.Scopes, &c.TenantID, &c.LogoURL, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("store.GetClient: %w", err)
	}
	return &c, nil
}

func (p *Postgres) GetTenantMember(ctx context.Context, tenantID, sub string) (string, error) {
	var role string
	err := p.pool.QueryRow(ctx,
		`SELECT role FROM tenant_members WHERE tenant_id = $1 AND sub = $2`, tenantID, sub,
	).Scan(&role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("store.GetTenantMember: %w", err)
	}
	return role, nil
}

func (p *Postgres) GetTenantPlan(ctx context.Context, tenantID string) (string, error) {
	var plan string
	err := p.pool.QueryRow(ctx,
		`SELECT plan FROM gerege_tenants WHERE id = $1 AND is_active = true`, tenantID,
	).Scan(&plan)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("store.GetTenantPlan: %w", err)
	}
	return plan, nil
}

func (p *Postgres) RecordIssuedToken(ctx context.Context, clientID, sub, scope string, expiresAt time.Time) error {
	_, err := p.pool.Exec(ctx,
		`INSERT INTO sso_issued_tokens (client_id, sub, scope, expires_at)
		 VALUES ($1, $2, $3, $4)`, clientID, sub, scope, expiresAt)
	if err != nil {
		return fmt.Errorf("store.RecordIssuedToken: %w", err)
	}
	return nil
}

// --- DAN Clients ---

func (p *Postgres) GetDANClient(ctx context.Context, clientID string) (*model.DANClient, error) {
	var c model.DANClient
	err := p.pool.QueryRow(ctx,
		`SELECT id, secret_hash, hmac_key, name, callback_urls, active, created_at, updated_at
		 FROM dan_clients WHERE id = $1`, clientID,
	).Scan(&c.ID, &c.SecretHash, &c.HMACKey, &c.Name, &c.CallbackURLs, &c.Active, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("store.GetDANClient: %w", err)
	}
	return &c, nil
}

func (p *Postgres) ListDANClients(ctx context.Context) ([]model.DANClient, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT id, name, callback_urls, active, created_at, updated_at
		 FROM dan_clients ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("store.ListDANClients: %w", err)
	}
	defer rows.Close()

	var clients []model.DANClient
	for rows.Next() {
		var c model.DANClient
		if err := rows.Scan(&c.ID, &c.Name, &c.CallbackURLs, &c.Active, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("store.ListDANClients scan: %w", err)
		}
		clients = append(clients, c)
	}
	return clients, rows.Err()
}

func (p *Postgres) CreateDANClient(ctx context.Context, id, secretHash, hmacKey, name string, callbackURLs []string) error {
	_, err := p.pool.Exec(ctx,
		`INSERT INTO dan_clients (id, secret_hash, hmac_key, name, callback_urls)
		 VALUES ($1, $2, $3, $4, $5)`, id, secretHash, hmacKey, name, callbackURLs)
	if err != nil {
		return fmt.Errorf("store.CreateDANClient: %w", err)
	}
	return nil
}

func (p *Postgres) DeactivateDANClient(ctx context.Context, clientID string) error {
	_, err := p.pool.Exec(ctx,
		`UPDATE dan_clients SET active = false, updated_at = now() WHERE id = $1`, clientID)
	if err != nil {
		return fmt.Errorf("store.DeactivateDANClient: %w", err)
	}
	return nil
}
