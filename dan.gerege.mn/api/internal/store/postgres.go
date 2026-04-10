package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DANClient struct {
	ID           string
	SecretHash   string
	HMACKey      string
	Name         string
	CallbackURLs []string
	Active       bool
}

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

func (p *Postgres) GetDANClient(ctx context.Context, clientID string) (*DANClient, error) {
	var c DANClient
	err := p.pool.QueryRow(ctx,
		`SELECT id, secret_hash, hmac_key, name, callback_urls, active
		 FROM dan_clients WHERE id = $1`, clientID,
	).Scan(&c.ID, &c.SecretHash, &c.HMACKey, &c.Name, &c.CallbackURLs, &c.Active)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("store.GetDANClient: %w", err)
	}
	return &c, nil
}
