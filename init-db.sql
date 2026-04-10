-- Gerege Platform — Full DB initialization script
-- Run on fresh DB or after recreate to ensure all tables exist
-- All statements use IF NOT EXISTS / ON CONFLICT — safe to re-run

-- =====================
-- SSO Server tables
-- =====================
CREATE TABLE IF NOT EXISTS sso_clients (
    id            TEXT PRIMARY KEY,
    secret_hash   TEXT        NOT NULL,
    name          TEXT        NOT NULL,
    redirect_uris TEXT[]      NOT NULL,
    scopes        TEXT[]      NOT NULL DEFAULT '{openid,profile}',
    is_active     BOOLEAN     NOT NULL DEFAULT true,
    tenant_id     TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sso_issued_tokens (
    id             TEXT PRIMARY KEY,
    client_id      TEXT        NOT NULL,
    subject        TEXT        NOT NULL,
    grant_type     TEXT        NOT NULL,
    scopes         TEXT[]      NOT NULL,
    access_token   TEXT        NOT NULL,
    refresh_token  TEXT,
    id_token       TEXT,
    expires_at     TIMESTAMPTZ NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_tokens_subject ON sso_issued_tokens(subject);
CREATE INDEX IF NOT EXISTS idx_tokens_access ON sso_issued_tokens(access_token);

CREATE TABLE IF NOT EXISTS gerege_tenants (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    slug       TEXT UNIQUE NOT NULL,
    plan       TEXT NOT NULL DEFAULT 'starter',
    is_active  BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tenant_members (
    tenant_id   TEXT NOT NULL,
    user_sub    TEXT NOT NULL,
    role        TEXT NOT NULL DEFAULT 'member',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, user_sub)
);

-- =====================
-- DAN Gateway tables
-- =====================
CREATE TABLE IF NOT EXISTS dan_clients (
    id            TEXT PRIMARY KEY,
    secret_hash   TEXT NOT NULL,
    hmac_key      TEXT NOT NULL DEFAULT '',
    name          TEXT NOT NULL,
    callback_urls TEXT[] NOT NULL,
    active        BOOLEAN NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_dan_clients_active ON dan_clients(active);

-- =====================
-- Developer Portal tables
-- =====================
CREATE TABLE IF NOT EXISTS dev_developers (
    id            TEXT PRIMARY KEY,
    sub           TEXT UNIQUE NOT NULL,
    name          TEXT NOT NULL,
    "givenName"   TEXT NOT NULL,
    "familyName"  TEXT NOT NULL,
    "certSerial"  TEXT NOT NULL DEFAULT '',
    "createdAt"   TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updatedAt"   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dev_tenants (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    slug       TEXT UNIQUE NOT NULL,
    plan       TEXT NOT NULL DEFAULT 'starter',
    "isActive" BOOLEAN NOT NULL DEFAULT true,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dev_apps (
    id            TEXT PRIMARY KEY,
    name          TEXT NOT NULL,
    description   TEXT,
    "logoUrl"     TEXT,
    "clientId"    TEXT UNIQUE NOT NULL,
    "secretHash"  TEXT NOT NULL,
    "redirectUris" TEXT[] NOT NULL,
    scopes        TEXT[] NOT NULL DEFAULT '{openid,profile}',
    "isActive"    BOOLEAN NOT NULL DEFAULT true,
    "createdAt"   TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updatedAt"   TIMESTAMPTZ NOT NULL DEFAULT now(),
    "developerId" TEXT NOT NULL REFERENCES dev_developers(id),
    "tenantId"    TEXT REFERENCES dev_tenants(id)
);
CREATE INDEX IF NOT EXISTS idx_dev_apps_dev ON dev_apps("developerId");

CREATE TABLE IF NOT EXISTS dev_tenant_members (
    "tenantId"    TEXT NOT NULL REFERENCES dev_tenants(id),
    "developerId" TEXT NOT NULL REFERENCES dev_developers(id),
    role          TEXT NOT NULL DEFAULT 'member',
    "createdAt"   TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY ("tenantId", "developerId")
);

-- =====================
-- d-business.mn tables
-- =====================
CREATE TABLE IF NOT EXISTS dbiz_users (
    id            TEXT PRIMARY KEY,
    sub           TEXT UNIQUE NOT NULL,
    name          TEXT NOT NULL,
    "givenName"   TEXT NOT NULL,
    "familyName"  TEXT NOT NULL,
    "certSerial"  TEXT NOT NULL DEFAULT '',
    "createdAt"   TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updatedAt"   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dbiz_organizations (
    id                   TEXT PRIMARY KEY,
    name                 TEXT NOT NULL,
    "registrationNumber" TEXT UNIQUE NOT NULL,
    type                 TEXT NOT NULL,
    address              TEXT,
    phone                TEXT,
    email                TEXT,
    "isActive"           BOOLEAN NOT NULL DEFAULT true,
    "isVerified"         BOOLEAN NOT NULL DEFAULT false,
    "createdAt"          TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updatedAt"          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dbiz_org_members (
    "organizationId" TEXT NOT NULL REFERENCES dbiz_organizations(id),
    "userId"         TEXT NOT NULL REFERENCES dbiz_users(id),
    role             TEXT NOT NULL DEFAULT 'viewer',
    "createdAt"      TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY ("organizationId", "userId")
);

CREATE TABLE IF NOT EXISTS dbiz_certificates (
    id               TEXT PRIMARY KEY,
    "organizationId" TEXT NOT NULL REFERENCES dbiz_organizations(id),
    "commonName"     TEXT NOT NULL,
    "serialNumber"   TEXT,
    status           TEXT NOT NULL DEFAULT 'pending',
    purpose          TEXT NOT NULL DEFAULT 'seal',
    "issuedAt"       TIMESTAMPTZ,
    "expiresAt"      TIMESTAMPTZ,
    "createdAt"      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_dbiz_certs_org ON dbiz_certificates("organizationId");

CREATE TABLE IF NOT EXISTS dbiz_documents (
    id               TEXT PRIMARY KEY,
    "organizationId" TEXT NOT NULL REFERENCES dbiz_organizations(id),
    "uploadedById"   TEXT NOT NULL REFERENCES dbiz_users(id),
    name             TEXT NOT NULL,
    "fileName"       TEXT NOT NULL,
    "fileSize"       INTEGER NOT NULL,
    "fileHash"       TEXT NOT NULL,
    status           TEXT NOT NULL DEFAULT 'uploaded',
    "createdAt"      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_dbiz_docs_org ON dbiz_documents("organizationId");
CREATE INDEX IF NOT EXISTS idx_dbiz_docs_user ON dbiz_documents("uploadedById");

CREATE TABLE IF NOT EXISTS dbiz_signatures (
    id                 TEXT PRIMARY KEY,
    "documentId"       TEXT NOT NULL REFERENCES dbiz_documents(id),
    "organizationId"   TEXT NOT NULL REFERENCES dbiz_organizations(id),
    "signedById"       TEXT NOT NULL REFERENCES dbiz_users(id),
    "signerName"       TEXT,
    "certSerial"       TEXT,
    "signedAt"         TIMESTAMPTZ,
    "sessionId"        TEXT,
    status             TEXT NOT NULL DEFAULT 'pending',
    "verificationCode" TEXT,
    "createdAt"        TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_dbiz_sigs_doc ON dbiz_signatures("documentId");
CREATE INDEX IF NOT EXISTS idx_dbiz_sigs_session ON dbiz_signatures("sessionId");

-- =====================
-- API Server tables
-- =====================
CREATE TYPE IF NOT EXISTS signing_status AS ENUM ('PENDING','RUNNING','COMPLETE','ERROR','CANCELLED','EXPIRED');

CREATE TABLE IF NOT EXISTS signing_sessions (
    id                TEXT PRIMARY KEY,
    requester_sub     TEXT NOT NULL,
    signer_sub        TEXT,
    signer_name       TEXT,
    signer_reg        TEXT,
    status            TEXT NOT NULL DEFAULT 'PENDING',
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
CREATE INDEX IF NOT EXISTS idx_signing_requester ON signing_sessions(requester_sub);
