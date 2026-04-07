package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gesign.mn/gerege-api/internal/model"
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

func (p *Postgres) Close() { p.pool.Close() }

func (p *Postgres) CreateSession(ctx context.Context, s *model.SigningSession) error {
	_, err := p.pool.Exec(ctx,
		`INSERT INTO signing_sessions (id, requester_sub, signer_reg, status, smartid_session, verification_code,
		 document_name, document_hash, document_size, document_path, expires_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		s.ID, s.RequesterSub, s.SignerReg, s.Status, s.SmartIDSession, s.VerificationCode,
		s.DocumentName, s.DocumentHash, s.DocumentSize, s.DocumentPath, s.ExpiresAt)
	return err
}

func (p *Postgres) GetSession(ctx context.Context, id string) (*model.SigningSession, error) {
	var s model.SigningSession
	var signerSub, signerName, signerReg, smartidSess, vc, docPath, signedPath, certSerial, errMsg *string
	err := p.pool.QueryRow(ctx,
		`SELECT id, requester_sub, signer_sub, signer_name, signer_reg, status::text,
		 smartid_session, verification_code, document_name, document_hash, document_size,
		 document_path, signed_doc_path, cert_serial, error_message, created_at, updated_at, expires_at
		 FROM signing_sessions WHERE id = $1`, id,
	).Scan(&s.ID, &s.RequesterSub, &signerSub, &signerName, &signerReg, &s.Status,
		&smartidSess, &vc, &s.DocumentName, &s.DocumentHash, &s.DocumentSize,
		&docPath, &signedPath, &certSerial, &errMsg, &s.CreatedAt, &s.UpdatedAt, &s.ExpiresAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if signerSub != nil { s.SignerSub = *signerSub }
	if signerName != nil { s.SignerName = *signerName }
	if signerReg != nil { s.SignerReg = *signerReg }
	if smartidSess != nil { s.SmartIDSession = *smartidSess }
	if vc != nil { s.VerificationCode = *vc }
	if docPath != nil { s.DocumentPath = *docPath }
	if signedPath != nil { s.SignedDocPath = *signedPath }
	if certSerial != nil { s.CertSerial = *certSerial }
	if errMsg != nil { s.ErrorMessage = *errMsg }
	return &s, nil
}

func (p *Postgres) UpdateSessionStatus(ctx context.Context, id, status string) error {
	_, err := p.pool.Exec(ctx,
		`UPDATE signing_sessions SET status=$2, updated_at=now() WHERE id=$1`, id, status)
	return err
}

func (p *Postgres) UpdateSessionComplete(ctx context.Context, id, signerSub, signerName, certSerial, signedDocPath string) error {
	_, err := p.pool.Exec(ctx,
		`UPDATE signing_sessions SET status='COMPLETE', signer_sub=$2, signer_name=$3, cert_serial=$4,
		 signed_doc_path=$5, updated_at=now() WHERE id=$1`,
		id, signerSub, signerName, certSerial, signedDocPath)
	return err
}

func (p *Postgres) UpdateSessionError(ctx context.Context, id, errMsg string) error {
	_, err := p.pool.Exec(ctx,
		`UPDATE signing_sessions SET status='ERROR', error_message=$2, updated_at=now() WHERE id=$1`, id, errMsg)
	return err
}
