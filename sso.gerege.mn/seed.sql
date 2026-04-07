-- bcrypt hash of "gerege-sso-secret-2026" (cost 12)

-- Gerege SSO Client (main integration client)
INSERT INTO sso_clients (id, secret_hash, name, redirect_uris, scopes)
VALUES (
    'gerege-sso-client',
    '$2b$12$u0pRrLyt/9Csd/Y6nrPuUOr8ZIG8MKRDOMu40.s7fWBIWsHi3ui26',
    'Gerege SSO Client',
    ARRAY['https://gerege.mn/callback', 'https://gerege.mn/api/auth/callback/gerege-sso', 'http://localhost:3000/callback', 'http://localhost:3000/api/auth/callback/gerege-sso'],
    ARRAY['openid', 'profile', 'pos', 'social', 'payment']
) ON CONFLICT (id) DO UPDATE SET secret_hash='$2b$12$u0pRrLyt/9Csd/Y6nrPuUOr8ZIG8MKRDOMu40.s7fWBIWsHi3ui26';

-- Gerege POS
INSERT INTO sso_clients (id, secret_hash, name, redirect_uris, scopes)
VALUES (
    'gerege-pos',
    '$2b$12$u0pRrLyt/9Csd/Y6nrPuUOr8ZIG8MKRDOMu40.s7fWBIWsHi3ui26',
    'Gerege POS',
    ARRAY['https://pos.gerege.mn/callback', 'https://pos.gerege.mn/api/auth/callback/gerege-sso'],
    ARRAY['openid', 'profile', 'pos', 'payment']
) ON CONFLICT (id) DO UPDATE SET secret_hash='$2b$12$u0pRrLyt/9Csd/Y6nrPuUOr8ZIG8MKRDOMu40.s7fWBIWsHi3ui26';

-- Gerege Social
INSERT INTO sso_clients (id, secret_hash, name, redirect_uris, scopes)
VALUES (
    'gerege-social',
    '$2b$12$u0pRrLyt/9Csd/Y6nrPuUOr8ZIG8MKRDOMu40.s7fWBIWsHi3ui26',
    'Gerege Social',
    ARRAY['https://social.gerege.mn/callback'],
    ARRAY['openid', 'profile', 'social']
) ON CONFLICT (id) DO UPDATE SET secret_hash='$2b$12$u0pRrLyt/9Csd/Y6nrPuUOr8ZIG8MKRDOMu40.s7fWBIWsHi3ui26';

-- Developer Portal
INSERT INTO sso_clients (id, secret_hash, name, redirect_uris, scopes)
VALUES (
    'gerege-developer-portal',
    '$2b$12$u0pRrLyt/9Csd/Y6nrPuUOr8ZIG8MKRDOMu40.s7fWBIWsHi3ui26',
    'Gerege Developer Portal',
    ARRAY['https://developer.gerege.mn/api/auth/callback/gerege-sso'],
    ARRAY['openid', 'profile']
) ON CONFLICT (id) DO UPDATE SET secret_hash='$2b$12$u0pRrLyt/9Csd/Y6nrPuUOr8ZIG8MKRDOMu40.s7fWBIWsHi3ui26';

-- Test sandbox
INSERT INTO sso_clients (id, secret_hash, name, redirect_uris, scopes)
VALUES (
    'dev-test-client',
    '$2b$12$u0pRrLyt/9Csd/Y6nrPuUOr8ZIG8MKRDOMu40.s7fWBIWsHi3ui26',
    'Local Dev Test',
    ARRAY['http://localhost:3000/callback', 'http://localhost:3000/api/auth/callback/gerege-sso', 'https://test.gerege.mn/api/auth/callback/gerege-sso'],
    ARRAY['openid', 'profile', 'pos', 'social', 'payment']
) ON CONFLICT (id) DO UPDATE SET secret_hash='$2b$12$u0pRrLyt/9Csd/Y6nrPuUOr8ZIG8MKRDOMu40.s7fWBIWsHi3ui26';
