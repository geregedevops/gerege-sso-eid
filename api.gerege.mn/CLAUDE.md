# api.gerege.mn — Gerege Sign API

## Architecture
api.gerege.mn   → REST Sign API (энэ repo)
sso.gerege.mn   → JWT auth provider
e-id.mn         → SmartID PIN2 signing
ocsp.gesign.mn  → Cert validity

## Sign Flow
1. Client → POST /v1/sign/request (PDF + signer_reg)
2. api → SmartID initiate (PIN2 push)
3. Client → GET /v1/sign/{id}/status (poll)
4. User → SmartID app PIN2 оруулна
5. api → SmartID status COMPLETE → PDF sign
6. Client → GET /v1/sign/{id}/result (signed PDF)

## Auth
sso.gerege.mn-ийн access_token Bearer болгон өгнө.
JWT verify: ES256, issuer: https://sso.gerege.mn
