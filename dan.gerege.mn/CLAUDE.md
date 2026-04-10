# dan.gerege.mn — DAN Verify Gateway

## Architecture

```
dan.gerege.mn/
├── api/             → Go DAN gateway backend (:8444)
│   ├── cmd/dan/     → Entry point
│   └── internal/    → handler, dan, middleware, store
├── web/             → Next.js admin dashboard (:3000)
│   ├── app/         → Pages (dashboard, docs, auth)
│   └── lib/         → auth.ts (NextAuth), api.ts (SSO API client)
└── CLAUDE.md
```

## DAN Flow
1. Client → GET /verify?client_id=X&callback_url=Y
2. → sso.gov.mn/login (signed state)
3. User DAN authentication
4. → GET /authorized?code=Z&state=S
5. DAN API → sso.gov.mn token exchange
6. DAN API → sso.gov.mn citizen data
7. DAN API → POST callback_url (citizen JSON + HMAC signature)
8. DAN API → redirect browser → callback_url?status=ok

## Client Management
- dan_clients table lives in sso.gerege.mn's PostgreSQL
- CRUD via SSO API: /api/dan/clients (DAN_ADMIN_KEY auth)
- DAN Go backend reads dan_clients directly from DB (read-only)
- Next.js dashboard calls SSO API for CRUD

## Security
- Callback URL: domain match (scheme + host), HTTPS required
- HMAC: uses dedicated hmac_key column (not bcrypt hash)
- State: HMAC-SHA256 signed (prevents CSRF/tampering)
- CORS: restricted to https://dan.gerege.mn
- Errors: proper HTTP status codes (400/401/502)

## Env vars (api)
DAN_CLIENT_ID, DAN_CLIENT_SECRET, DAN_SCOPE, DAN_CALLBACK_URI
DAN_TOKEN_URL, DAN_SERVICE_URL, DAN_STATE_SECRET
DATABASE_URL, CORS_ORIGIN, PORT

## Env vars (web)
GEREGE_SSO_CLIENT, GEREGE_SSO_SECRET, NEXT_PUBLIC_SSO_URL
SSO_API_URL, DAN_ADMIN_KEY, AUTH_SECRET, AUTH_URL
