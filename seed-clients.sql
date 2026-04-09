-- Gerege Platform — SSO Client seed data
-- Safe to re-run (ON CONFLICT DO UPDATE to fix broken hashes)

INSERT INTO sso_clients (id, secret_hash, name, redirect_uris, scopes) VALUES
('gerege-developer-portal',
 '$2a$12$RRDNFloUkePn8q7HThclXOrlBY/mr/b9zuXgPjUj0zT1F4usJPPVy',
 'Gerege Developer Portal',
 '{https://developer.gerege.mn/api/auth/callback/gerege-sso,http://localhost:3002/api/auth/callback/gerege-sso}',
 '{openid,profile}'),
('cmnokjo8p0004v0ctdmssbr82',
 '$2a$12$84VQw1YZ7KfMQhfLS28yyuWANeM8IT4JSkY/eKwIuE3fNliAJdWx2',
 'Gerege Test Sandbox',
 '{https://test.gerege.mn/api/auth/callback/gerege-sso,http://localhost:3003/api/auth/callback/gerege-sso}',
 '{openid,profile}'),
('cmnrijlr10004fap55twiocei',
 '$2a$12$EEjiQMxykIToanlVWSh92ObdZb7ctq5FxlYNnW.DlGI0VsiENU73i',
 'd-business.mn',
 '{https://d-business.mn/api/auth/callback/gerege-sso}',
 '{openid,profile}')
ON CONFLICT (id) DO UPDATE SET
  secret_hash = EXCLUDED.secret_hash,
  redirect_uris = EXCLUDED.redirect_uris;
