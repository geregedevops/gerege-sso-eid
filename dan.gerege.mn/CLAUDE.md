# dan.gerege.mn — DAN Verify Gateway

## Architecture
```
dan.gerege.mn/
├── api/             → Go DAN gateway backend (:8444)
│   ├── cmd/dan/     → Entry point, config, server
│   └── internal/
│       ├── handler/ → verify, authorized, try, clients, index
│       ├── dan/     → token.go, citizen.go, hmac.go (sso.gov.mn API)
│       ├── middleware/ → logger (status code), cors (restricted)
│       └── store/   → postgres.go (dan_clients CRUD + auto-migrate)
├── web/             → Next.js admin dashboard (:3000) (deploy хийгдээгүй)
│   ├── app/         → Pages (dashboard, docs, auth)
│   └── lib/         → auth.ts (NextAuth), api.ts (DAN API client)
└── CLAUDE.md
```

## Database
- Тусдаа DB: `gerege_dan_db` (SSO-с бүрэн салсан)
- Auto-migration: DAN API эхлэхэд `dan_clients` table автоматаар үүснэ
- Client CRUD: DAN API `/api/clients` endpoint (DAN_ADMIN_KEY auth)

## DAN Verify Flow (client бүртгэлтэй)
1. Client → GET /verify?client_id=X&callback_url=Y
2. Validate client + callback URL (domain match, HTTPS required)
3. → sso.gov.mn/login (HMAC signed state)
4. User ДАН нэвтрэлт
5. → GET /authorized?code=Z&state=S
6. Verify signed state, re-validate client
7. sso.gov.mn token exchange → citizen data
8. POST callback_url (citizen JSON + image + HMAC signature)
9. Redirect browser → callback_url?status=ok&reg_no=...

## DAN Try Flow (бие даасан, client шаардахгүй)
1. GET /try → sso.gov.mn/login (state mode=try)
2. User ДАН нэвтрэлт
3. → GET /authorized (mode=try detected)
4. Citizen data → HTML хуудсаар харуулна (зураг + бүх мэдээлэл)

## Admin API
- GET    /api/clients — жагсаалт (DAN_ADMIN_KEY Bearer auth)
- POST   /api/clients — шинэ client үүсгэх (returns secret + hmac_key once)
- DELETE /api/clients/{id} — идэвхгүйжүүлэх

## Security
- Callback URL: scheme + host exact match, path prefix, HTTPS заавал
- HMAC: тусдаа hmac_key column (bcrypt hash биш)
- State: HMAC-SHA256 signed (CSRF/tampering хамгаалалт)
- CORS: зөвхөн https://dan.gerege.mn
- Admin key: subtle.ConstantTimeCompare (timing attack хамгаалалт)
- Error codes: 400/401/502 зөв ялгана
- rand.Read: error шалгадаг

## Env vars (api)
DAN_CLIENT_ID, DAN_CLIENT_SECRET, DAN_SCOPE — sso.gov.mn credentials
DAN_CALLBACK_URI — http://dan.gerege.mn/authorized (sso.gov.mn-д бүртгэлтэй)
DAN_TOKEN_URL, DAN_SERVICE_URL — sso.gov.mn endpoints
DAN_STATE_SECRET — state signing key
DAN_ADMIN_KEY — admin API key
DAN_DATABASE_URL — postgres://...gerege_dan_db
CORS_ORIGIN, PORT

## Env vars (web)
GEREGE_SSO_CLIENT, GEREGE_SSO_SECRET — SSO OIDC credentials
NEXT_PUBLIC_SSO_URL, SSO_API_URL — SSO endpoints
DAN_ADMIN_KEY, AUTH_SECRET, AUTH_URL

## Production
- sso.gov.mn redirect_uri: http:// (https биш!)
- DAN_CALLBACK_URI=http://dan.gerege.mn/authorized
- nginx: /verify, /authorized, /try, /api/clients → dan-api:8444
- nginx: / → dan-api:8444 (Next.js deploy болтол)
