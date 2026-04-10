# api.gerege.mn — Gerege Sign API

## Architecture
api.gerege.mn   → REST Sign API (Go :8080)
sso.gerege.mn   → JWT auth provider (ES256 verify)
e-id.mn         → SmartID PIN2 signing
e-id.mn/ocsp    → Certificate validity check
Database        → gerege_sso_db (signing_sessions table)
Redis           → Session state caching

## Sign Flow
1. Client → POST /v1/sign/request (PDF + signer info, Bearer JWT)
2. API → SmartID initiate (PIN2 push to user)
3. Client → GET /v1/sign/{id}/status (poll PENDING→RUNNING→COMPLETE)
4. User → SmartID app PIN2 оруулна
5. API → SmartID status COMPLETE → PDF digitally signed
6. Client → GET /v1/sign/{id}/result (download signed PDF)

## Endpoints
- POST   /v1/sign/request — Initiate signing (JWT auth required)
- GET    /v1/sign/{id}/status — Poll signing status
- GET    /v1/sign/{id}/result — Retrieve signed PDF
- DELETE /v1/sign/{id} — Cancel signing session
- POST   /v1/verify — Verify document signature

## Auth
sso.gerege.mn-ийн access_token Bearer болгон өгнө.
JWT verify: ES256, issuer: https://sso.gerege.mn, JWKS: /.well-known/jwks.json

## Env vars
SSO_JWKS_URI, EID_API_URL, STORAGE_PATH
DATABASE_URL, REDIS_URL, PORT
