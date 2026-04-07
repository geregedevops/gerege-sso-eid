CREATE TABLE IF NOT EXISTS gerege_tenants (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    plan        TEXT NOT NULL DEFAULT 'starter',
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE sso_clients ADD COLUMN IF NOT EXISTS tenant_id TEXT REFERENCES gerege_tenants(id);

CREATE TABLE IF NOT EXISTS tenant_members (
    tenant_id   TEXT REFERENCES gerege_tenants(id),
    sub         TEXT NOT NULL,
    role        TEXT NOT NULL DEFAULT 'member',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, sub)
);
