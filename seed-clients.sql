-- Gerege Platform — SSO Client seed data
-- Safe to re-run (ON CONFLICT DO NOTHING)

INSERT INTO sso_clients (id, secret_hash, name, redirect_uris, scopes) VALUES
('gerege-sso-client', '$2a$12$LjRHqhLBbQr.Nwl0CUhXl.AMWBVI3FkJR.I/CgKbmOsIGWJFaRWFa', 'Gerege SSO Client',
 '{https://sso.gerege.mn/callback,https://localhost:8443/callback,http://localhost:8443/callback,https://sso.gerege.mn/oauth/callback}',
 '{openid,profile,pos,social,payment}'),
('gerege-developer-portal', '$2a$12$LjRHqhLBbQr.Nwl0CUhXl.AMWBVI3FkJR.I/CgKbmOsIGWJFaRWFa', 'Gerege Developer Portal',
 '{https://developer.gerege.mn/api/auth/callback/gerege-sso,http://localhost:3002/api/auth/callback/gerege-sso}',
 '{openid,profile}'),
('cmnokjo8p0004v0ctdmssbr82', '$2a$12$hhxSGRHWFkRHwNrAcx6oteaoiH1pFojy/a6G/MhJfKKWhpKpDn7Oe', 'Gerege Test Sandbox',
 '{https://test.gerege.mn/api/auth/callback/gerege-sso,http://localhost:3003/api/auth/callback/gerege-sso}',
 '{openid,profile}'),
('cmnrijlr10004fap55twiocei', '$2a$12$EEjiQMxykIToanlVWSh92ObdZb7ctq5FxlYNnW.DlGI0VsiENU73i', 'd-business.mn',
 '{https://d-business.mn/api/auth/callback/gerege-sso}',
 '{openid,profile}')
ON CONFLICT (id) DO NOTHING;
