-- bcrypt hash of "dev-secret-local" (cost 12)

-- Gerege POS
INSERT INTO sso_clients (id, secret_hash, name, redirect_uris, scopes)
VALUES (
    'gerege-pos',
    '$2a$12$05qLIc1GfAKtDAuf99O67uDMs9I7Yqdb16rnes85xfh427OtzzBQS',
    'Gerege POS',
    ARRAY['https://pos.gerege.mn/callback', 'https://pos.gerege.mn/api/auth/callback/gerege-sso'],
    ARRAY['openid', 'profile', 'pos', 'payment']
) ON CONFLICT (id) DO NOTHING;

-- Gerege Social
INSERT INTO sso_clients (id, secret_hash, name, redirect_uris, scopes)
VALUES (
    'gerege-social',
    '$2a$12$05qLIc1GfAKtDAuf99O67uDMs9I7Yqdb16rnes85xfh427OtzzBQS',
    'Gerege Social',
    ARRAY['https://social.gerege.mn/callback'],
    ARRAY['openid', 'profile', 'social']
) ON CONFLICT (id) DO NOTHING;

-- Developer Portal
INSERT INTO sso_clients (id, secret_hash, name, redirect_uris, scopes)
VALUES (
    'gerege-developer-portal',
    '$2a$12$05qLIc1GfAKtDAuf99O67uDMs9I7Yqdb16rnes85xfh427OtzzBQS',
    'Gerege Developer Portal',
    ARRAY['https://developer.gerege.mn/api/auth/callback/gerege-sso'],
    ARRAY['openid', 'profile']
) ON CONFLICT (id) DO NOTHING;

-- Test sandbox
INSERT INTO sso_clients (id, secret_hash, name, redirect_uris, scopes)
VALUES (
    'dev-test-client',
    '$2a$12$05qLIc1GfAKtDAuf99O67uDMs9I7Yqdb16rnes85xfh427OtzzBQS',
    'Local Dev Test',
    ARRAY['http://localhost:3000/callback', 'http://localhost:3000/api/auth/callback/gerege-sso'],
    ARRAY['openid', 'profile', 'pos', 'social', 'payment']
) ON CONFLICT (id) DO NOTHING;
