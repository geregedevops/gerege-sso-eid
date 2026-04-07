CREATE TABLE IF NOT EXISTS sso_issued_tokens (
    id         BIGSERIAL    PRIMARY KEY,
    client_id  TEXT         NOT NULL REFERENCES sso_clients(id),
    sub        TEXT         NOT NULL,
    scope      TEXT         NOT NULL,
    issued_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ  NOT NULL,
    revoked    BOOLEAN      NOT NULL DEFAULT false
);
CREATE INDEX IF NOT EXISTS idx_sso_issued_tokens_sub ON sso_issued_tokens(sub);
CREATE INDEX IF NOT EXISTS idx_sso_issued_tokens_client ON sso_issued_tokens(client_id);
