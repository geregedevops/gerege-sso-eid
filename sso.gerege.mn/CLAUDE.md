# sso.gerege.mn — Gerege OIDC Authorization Server

## Architecture
e-id.mn         → SmartID auth (хөндөхгүй)
sso.gerege.mn   → Gerege OIDC server (энэ repo)

## Scopes
- openid, profile (standard)
- pos, social, payment (Gerege-specific)

## Claims
- Standard: sub, name, given_name, family_name, locale
- Gerege: tenant_id, tenant_role, plan, reg_no

## Flow
1. Client → GET /oauth/authorize
2. → e-id.mn/auth?session=X&callback_uri=sso.gerege.mn/callback/eid
3. SmartID PIN1
4. → sso.gerege.mn/callback/eid
5. auth_code → JWT

## EC Key
openssl ecparam -name prime256v1 -genkey -noout -out ec-private.pem
openssl ec -in ec-private.pem -pubout -out ec-public.pem

## Migrations
psql $DATABASE_URL -f migrations/001_clients.sql
psql $DATABASE_URL -f migrations/002_sessions.sql
psql $DATABASE_URL -f migrations/003_tenants.sql
psql $DATABASE_URL -f seed.sql
