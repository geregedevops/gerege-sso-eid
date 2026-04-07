CREATE TABLE IF NOT EXISTS sso_clients (
    id            TEXT PRIMARY KEY,
    secret_hash   TEXT        NOT NULL,
    name          TEXT        NOT NULL,
    redirect_uris TEXT[]      NOT NULL,
    scopes        TEXT[]      NOT NULL DEFAULT '{openid,profile}',
    tenant_id     TEXT,
    logo_url      TEXT,
    is_active     BOOLEAN     NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
