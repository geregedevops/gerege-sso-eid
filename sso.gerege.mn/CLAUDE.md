# sso.gerege.mn — Gerege OIDC Authorization Server

## Architecture
e-id.mn         → SmartID auth (хөндөхгүй)
sso.gerege.mn   → Gerege OIDC server (энэ repo)
ocsp.gesign.mn  → Cert validity (gesign CA ашиглана)

## sso.gesign.mn-тэй ялгаа
- Issuer: https://sso.gerege.mn
- Нэмэлт scopes: pos, social, payment
- Нэмэлт claims: tenant_id, tenant_role, plan
- Tenant middleware: token-д tenant context inject хийнэ

## Flow (sso.gesign.mn-тэй ижил)
1. Client → GET /oauth/authorize
2. → e-id.mn/auth?session=X&callback_uri=sso.gerege.mn/callback/eid
3. SmartID PIN1
4. → sso.gerege.mn/callback/eid
5. OCSP verify → auth_code → JWT

## EC Key (тусдаа — gesign-ийнхтэй хутгахгүй)
openssl ecparam -name prime256v1 -genkey -noout -out ec-private.pem
openssl ec -in ec-private.pem -pubout -out ec-public.pem

## Migrations
psql $DATABASE_URL -f migrations/001_clients.sql
psql $DATABASE_URL -f migrations/002_sessions.sql
psql $DATABASE_URL -f migrations/003_tenants.sql
psql $DATABASE_URL -f seed.sql
