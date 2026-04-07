DO $$ BEGIN
  CREATE TYPE signing_status AS ENUM ('PENDING','RUNNING','COMPLETE','ERROR','CANCELLED','EXPIRED');
EXCEPTION WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS signing_sessions (
    id                TEXT PRIMARY KEY,
    requester_sub     TEXT NOT NULL,
    signer_sub        TEXT,
    signer_name       TEXT,
    signer_reg        TEXT,
    status            signing_status NOT NULL DEFAULT 'PENDING',
    smartid_session   TEXT,
    verification_code TEXT,
    document_name     TEXT NOT NULL,
    document_hash     TEXT NOT NULL,
    document_size     INTEGER NOT NULL,
    document_path     TEXT,
    signed_doc_path   TEXT,
    cert_serial       TEXT,
    error_message     TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at        TIMESTAMPTZ NOT NULL DEFAULT now() + interval '10 minutes'
);

CREATE INDEX IF NOT EXISTS idx_signing_sessions_requester ON signing_sessions(requester_sub);
CREATE INDEX IF NOT EXISTS idx_signing_sessions_status ON signing_sessions(status);
