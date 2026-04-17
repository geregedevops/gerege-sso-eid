-- DAN Gateway (dan.gerege.mn) — dan_clients seed data
-- Database: gerege_sso_db, User: sso
-- Generated: 2026-04-10

CREATE TABLE IF NOT EXISTS dan_clients (
    id            TEXT PRIMARY KEY,
    secret_hash   TEXT NOT NULL,
    name          TEXT NOT NULL,
    callback_urls TEXT[] NOT NULL DEFAULT '{}',
    active        BOOLEAN NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_dan_clients_active ON dan_clients(active);

INSERT INTO dan_clients (id, secret_hash, name, callback_urls, active, created_at) VALUES
(
    'dan_c203a108a3241f5ada7227c4a591f289',
    '$2a$10$IM8AbWyGBMR543SSM91rHuJ/vBvja8n6jM5gNCGWhiuh55ClB9O2W',
    'GeregeID',
    ARRAY['https://ca.gerege.mn/mobile/v1/registration/kyc/dan/callback'],
    true,
    '2026-04-09 15:05:15.376751+00'
),
(
    'dan_24be55348c9643b0f9857ab49e4e94a5',
    '$2a$10$P2iPMkzlarS2WFVhDOvIdO1EI8lii.xQ2LPIsarQPjj9BBCT8IAsy',
    'screening',
    ARRAY['https://screening.gov.mn/api/auth/dan/gateway/callback/ncc', 'https://screening.gov.mn/api/auth/dan/gateway/callback/civil'],
    true,
    '2026-04-10 00:41:11.7604+00'
),
(
    'dan_ddcf12e133b9a5632f1ec6616160628c',
    '$2a$10$HOCEEjBrlj3lUOFrtNn1MOnGyJ3xDAOnhM5Uo.2lcrMEgLk1VggDq',
    'canreg',
    ARRAY['https://canreg.gov.mn/api/auth/dan/gateway/callback'],
    true,
    '2026-04-10 00:42:16.229201+00'
),
(
    'dan_ef8e0c4ef288a0a571c31dda7abf0f11',
    '$2a$10$nK6CC8LeUHiSvWYl2h6EmeAbi.0jp/Alt/1mqVk2lzEtknMmzPH7y',
    'EC',
    ARRAY['https://ec.transport.ub.gov.mn/api/auth/dan/callback'],
    true,
    '2026-04-10 03:04:36.388805+00'
),
(
    'dan_0f4e3d993b6c0547d569cf25b69da0b3',
    '$2a$10$D2Lc/XonW3DG6u7ZlLeK9e/nw.JzZjvhwIPdHjiRARS30f357LFXe',
    'insure',
    ARRAY['https://insure.gerege.mn/api/auth/dan/callback'],
    true,
    '2026-04-10 09:22:53.331252+00'
)
ON CONFLICT (id) DO NOTHING;
