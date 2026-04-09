-- Gerege Platform — SSO Client seed data
-- Safe to re-run (ON CONFLICT DO UPDATE)
-- Hashes verified with bcrypt.CompareHashAndPassword

INSERT INTO sso_clients (id, secret_hash, name, redirect_uris, scopes) VALUES
('gerege-developer-portal',
 '$2a$12$wExtdiDBLEWTINBXySE0Xef9nrHWI4QCPPDEn6o.YSZRtj56Mh9Me',
 'Gerege Developer Portal',
 '{https://developer.gerege.mn/api/auth/callback/gerege-sso,http://localhost:3002/api/auth/callback/gerege-sso}',
 '{openid,profile}'),
('cmnokjo8p0004v0ctdmssbr82',
 '$2a$12$7n0tMmVH3a08BhlkaQsFq.nRDGzrBgzIcTppBy4Ui0sHBc2nmijEu',
 'Gerege Test Sandbox',
 '{https://test.gerege.mn/api/auth/callback/gerege-sso,http://localhost:3003/api/auth/callback/gerege-sso}',
 '{openid,profile}'),
('cmnrijlr10004fap55twiocei',
 '$2a$12$WbzlkTlrMlWziXc19MSd9eNCAkk./jld1z1RDtDNCC56fz9XUUysC',
 'd-business.mn',
 '{https://d-business.mn/api/auth/callback/gerege-sso}',
 '{openid,profile}')
ON CONFLICT (id) DO UPDATE SET
  secret_hash = EXCLUDED.secret_hash,
  redirect_uris = EXCLUDED.redirect_uris;
