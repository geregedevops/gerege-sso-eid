# sso.gerege.mn — Gerege OIDC Authorization Server

## Architecture
e-id.mn         → SmartID auth (PIN1)
sso.gerege.mn   → Gerege OIDC server (энэ repo)
Database        → gerege_sso_db (PostgreSQL, shared with developer/test/dbusiness)
Cache           → Redis (session, auth code, access token)

## Scopes
- openid, profile (standard)
- pos, social, payment (Gerege-specific — tenant context)

## Claims
- Standard: sub, name, given_name, family_name, locale, cert_serial
- Gerege: tenant_id, tenant_role, plan, reg_no

## OIDC Flow
1. Client → GET /oauth/authorize (client_id, redirect_uri, scope, state, nonce)
2. → e-id.mn/auth?session=X&callback_uri=sso.gerege.mn/callback/eid
3. SmartID PIN1 баталгаажуулалт
4. → sso.gerege.mn/callback/eid (sub, name, cert_serial)
5. OCSP cert шалгалт (fail closed)
6. auth_code → Redis (5 min, single use)
7. Client → POST /oauth/token → ID token (ES256 JWT) + opaque access token

## Endpoints
- GET  /.well-known/openid-configuration — OIDC Discovery
- GET  /.well-known/jwks.json — Public key set
- GET  /oauth/authorize — Authorization code flow
- POST /oauth/token — Token exchange (client auth → rate limit)
- GET  /oauth/userinfo — User info (cert_serial included)
- POST /oauth/revoke — RFC 7009 compliant (always 200)
- POST /oauth/introspect — Token introspection

## Security
- JWT: ES256 (ECDSA P-256), KID from public key
- Client secrets: bcrypt cost 12
- Access tokens: opaque (stored in Redis, revocable)
- OCSP: fail closed (revoked cert = rejected)
- Rate limit: 10 req/min per client (after client validation)
- Revoke: RFC 7009 (never leaks client_id existence)

## EC Key
openssl ecparam -name prime256v1 -genkey -noout -out ec-private.pem
openssl ec -in ec-private.pem -pubout -out ec-public.pem

## Migrations
psql $DATABASE_URL -f migrations/001_clients.sql
psql $DATABASE_URL -f migrations/002_sessions.sql
psql $DATABASE_URL -f migrations/003_tenants.sql
psql $DATABASE_URL -f seed.sql

## Env vars
SSO_ISSUER, SSO_PRIVATE_KEY_PATH, EID_BASE_URL
DATABASE_URL, REDIS_URL, PORT
OCSP_URL, CA_ISSUING_URL (optional)
TLS_CERT, TLS_KEY, DEV_MODE
