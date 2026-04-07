# Gerege SSO Ecosystem — Claude Code Prompt
# sso.gerege.mn + developer.gerege.mn + test.gerege.mn

> gesign.mn ecosystem-ийг Gerege platform-д давтах.
> Гурван repo тус тусдаа Claude Code-д өгнө.

---

## Ерөнхий архитектур

```
e-id.mn                    ← SmartID auth backend (хөндөхгүй)
    │
    ▼
sso.gerege.mn              ← Gerege OIDC Authorization Server  [ШИНЭ]
    │
    ├──► developer.gerege.mn   ← Developer Portal              [ШИНЭ]
    ├──► gerege.mn             ← Gerege POS / main platform
    ├──► pos.gerege.mn         ← POS system
    ├──► social.gerege.mn      ← Social commerce
    └──► test.gerege.mn        ← API sandbox & testing         [ШИНЭ]
```

**sso.gesign.mn-аас ялгарах зүйл:**
- Issuer: `https://sso.gerege.mn`
- Gerege ecosystem-д зориулсан scopes: `pos`, `social`, `payment`
- Tenant-aware claims: `tenant_id`, `tenant_role`
- Gerege brand, өнгө, нэр

---

# PART 1 — sso.gerege.mn

## Зорилго

Gerege platform-ийн бүх бүрэлдэхүүн хэсгүүд (`gerege.mn`, `pos.gerege.mn`, `social.gerege.mn`, 3rd party plugin-ууд) энэ OIDC server-т нэвтрэнэ.

## sso.gesign.mn-тэй харьцуулбал

| | sso.gesign.mn | sso.gerege.mn |
|---|---|---|
| Issuer | https://sso.gesign.mn | https://sso.gerege.mn |
| Зориулалт | GeSign signing ecosystem | Gerege platform ecosystem |
| Scopes | openid, profile, sign | openid, profile, pos, social, payment |
| Extra claims | cert_serial, cert_type | tenant_id, tenant_role, plan |
| Brand | GeSign / dark gold | Gerege / platform theme |

## Stack

sso.gesign.mn-тэй **яг ижил stack** — зөвхөн config өөрчлөгдөнө:
- Go 1.22+, net/http, pgx/v5, valkey-go, golang-jwt ES256

## Directory Structure

```
sso.gerege.mn/
├── cmd/sso/main.go
├── internal/
│   ├── handler/
│   │   ├── discovery.go
│   │   ├── jwks.go
│   │   ├── authorize.go
│   │   ├── token.go
│   │   ├── userinfo.go
│   │   ├── revoke.go
│   │   ├── introspect.go
│   │   └── eid_callback.go
│   ├── model/
│   ├── store/
│   │   ├── postgres.go
│   │   └── redis.go
│   ├── token/
│   │   ├── jwt.go
│   │   └── jwks.go
│   ├── ocsp/
│   │   └── checker.go
│   └── middleware/
│       ├── cors.go
│       └── auth.go
├── migrations/
│   ├── 001_clients.sql
│   ├── 002_sessions.sql
│   └── 003_tenants.sql          # Gerege-specific
├── .env.example
├── CLAUDE.md
└── go.mod
```

## Gerege-specific нэмэлт — Tenant

```sql
-- migrations/003_tenants.sql

CREATE TABLE gerege_tenants (
    id          TEXT PRIMARY KEY,         -- tenant slug (e.g. "restaurant-govi")
    name        TEXT NOT NULL,
    plan        TEXT NOT NULL DEFAULT 'starter',  -- starter | pro | enterprise
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- sso_clients-д tenant холбох
ALTER TABLE sso_clients ADD COLUMN tenant_id TEXT REFERENCES gerege_tenants(id);

-- Developer-ийн tenant membership
CREATE TABLE tenant_members (
    tenant_id   TEXT REFERENCES gerege_tenants(id),
    sub         TEXT NOT NULL,        -- sso sub
    role        TEXT NOT NULL DEFAULT 'member',  -- owner | admin | member
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, sub)
);
```

## ID Token — Gerege Claims нэмэлт

```json
{
  "iss": "https://sso.gerege.mn",
  "sub": "MN-АА12345678-hash",
  "aud": "client_id",
  "exp": 1712349278,
  "iat": 1712345678,
  "nonce": "...",

  "name": "Батаа Дорж",
  "given_name": "Дорж",
  "family_name": "Батаа",
  "locale": "mn-MN",

  "cert_serial": "8d1d55d6...",
  "cert_type": "AUTH",
  "identity_assurance_level": "high",
  "amr": ["smartid", "pin1", "x509"],

  "tenant_id": "restaurant-govi",
  "tenant_role": "owner",
  "plan": "pro"
}
```

## OIDC Discovery

```json
{
  "issuer": "https://sso.gerege.mn",
  "authorization_endpoint": "https://sso.gerege.mn/oauth/authorize",
  "token_endpoint": "https://sso.gerege.mn/oauth/token",
  "userinfo_endpoint": "https://sso.gerege.mn/oauth/userinfo",
  "jwks_uri": "https://sso.gerege.mn/.well-known/jwks.json",
  "revocation_endpoint": "https://sso.gerege.mn/oauth/revoke",
  "introspection_endpoint": "https://sso.gerege.mn/oauth/introspect",
  "scopes_supported": ["openid", "profile", "pos", "social", "payment"],
  "response_types_supported": ["code"],
  "id_token_signing_alg_values_supported": ["ES256"],
  "claims_supported": [
    "sub", "iss", "aud", "exp", "iat", "nonce",
    "name", "given_name", "family_name", "locale",
    "cert_serial", "identity_assurance_level", "amr",
    "tenant_id", "tenant_role", "plan"
  ]
}
```

## Scope → Claims mapping

| Scope | Claims |
|---|---|
| `openid` | sub, iss, aud, exp, iat |
| `profile` | name, given_name, family_name, locale, cert_serial |
| `pos` | tenant_id, tenant_role, plan |
| `social` | tenant_id, social_page_ids |
| `payment` | tenant_id, payment_methods |

## Environment Variables

```env
SSO_ISSUER=https://sso.gerege.mn
SSO_PRIVATE_KEY_PATH=/etc/gerege/sso/ec-private.pem
SSO_PUBLIC_KEY_PATH=/etc/gerege/sso/ec-public.pem

DATABASE_URL=postgres://sso:pass@localhost:5432/gerege_sso_db
REDIS_URL=redis://localhost:6379/2

EID_BASE_URL=https://e-id.mn
OCSP_URL=https://ocsp.gesign.mn/ocsp

PORT=8443
TLS_CERT=/etc/gerege/sso/tls/cert.pem
TLS_KEY=/etc/gerege/sso/tls/key.pem
```

## Database Schema

```sql
-- migrations/001_clients.sql (sso.gesign.mn-тэй ижил)
CREATE TABLE sso_clients (
    id            TEXT PRIMARY KEY,
    secret_hash   TEXT NOT NULL,
    name          TEXT NOT NULL,
    redirect_uris TEXT[] NOT NULL,
    scopes        TEXT[] NOT NULL DEFAULT '{openid,profile}',
    tenant_id     TEXT,
    logo_url      TEXT,
    is_active     BOOLEAN NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- migrations/002_sessions.sql
CREATE TABLE sso_issued_tokens (
    id         BIGSERIAL PRIMARY KEY,
    client_id  TEXT NOT NULL REFERENCES sso_clients(id),
    sub        TEXT NOT NULL,
    scope      TEXT NOT NULL,
    issued_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked    BOOLEAN NOT NULL DEFAULT false
);
CREATE INDEX ON sso_issued_tokens(sub);
```

## Seed — Gerege platform clients

```sql
-- Gerege POS
INSERT INTO sso_clients (id, secret_hash, name, redirect_uris, scopes)
VALUES (
    'gerege-pos',
    '<bcrypt>',
    'Gerege POS',
    ARRAY['https://pos.gerege.mn/callback', 'https://pos.gerege.mn/api/auth/callback/gerege-sso'],
    ARRAY['openid', 'profile', 'pos', 'payment']
);

-- Gerege Social
INSERT INTO sso_clients (id, secret_hash, name, redirect_uris, scopes)
VALUES (
    'gerege-social',
    '<bcrypt>',
    'Gerege Social',
    ARRAY['https://social.gerege.mn/callback'],
    ARRAY['openid', 'profile', 'social']
);

-- Developer Portal
INSERT INTO sso_clients (id, secret_hash, name, redirect_uris, scopes)
VALUES (
    'gerege-developer-portal',
    '<bcrypt>',
    'Gerege Developer Portal',
    ARRAY['https://developer.gerege.mn/api/auth/callback/gerege-sso'],
    ARRAY['openid', 'profile']
);

-- Test sandbox
INSERT INTO sso_clients (id, secret_hash, name, redirect_uris, scopes)
VALUES (
    'dev-test-client',
    '<bcrypt of dev-secret-local>',
    'Local Dev Test',
    ARRAY['http://localhost:3000/callback', 'http://localhost:3000/api/auth/callback/gerege-sso'],
    ARRAY['openid', 'profile', 'pos', 'social', 'payment']
);
```

## CLAUDE.md

```markdown
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
```

## Implementation Order

```
sso.gesign.mn-ийн кодыг үндэс болгон:

1. go.mod — module: gesign.mn/gerege-sso
2. CLAUDE.md
3. .env.example
4. migrations/ (001, 002, 003)
5. seed.sql
6. internal/store/ (postgres + redis)
7. internal/ocsp/checker.go
8. internal/token/ (jwt + jwks)
9. internal/handler/ — бүх handler
   → discovery.go: issuer = https://sso.gerege.mn
   → scopes_supported нэмэх: pos, social, payment
   → token.go: tenant claims нэмэх
10. internal/middleware/
11. cmd/sso/main.go
```

---

# PART 2 — developer.gerege.mn

## Зорилго

Gerege platform-д app бүтээх хөгжүүлэгчид:
- POS plugin хийх
- Social commerce нэгтгэх
- Payment flow нэгтгэх
- Gerege API ашиглах

dev.gesign.mn-ийг **Gerege брэнд + нэмэлт Gerege-specific зүйлстэй** давтана.

## dev.gesign.mn-тэй харьцуулбал

| | dev.gesign.mn | developer.gerege.mn |
|---|---|---|
| SSO provider | sso.gesign.mn | sso.gerege.mn |
| Нэмэлт section | - | Tenant management |
| Scopes UI | openid, profile, sign | openid, profile, pos, social, payment |
| Docs | SSO / Sign API | Gerege POS API, Social API |
| Brand | GeSign dark gold | Gerege platform theme |

## Stack

```
Next.js 14 (App Router) + TypeScript
Tailwind CSS + shadcn/ui
NextAuth.js v5 — sso.gerege.mn OIDC
PostgreSQL + Prisma
```

## Directory Structure

```
developer.gerege.mn/
├── app/
│   ├── layout.tsx
│   ├── page.tsx                          # Landing
│   ├── auth/login/page.tsx
│   ├── dashboard/
│   │   ├── layout.tsx                   # Auth guard
│   │   ├── page.tsx                     # Overview
│   │   ├── apps/
│   │   │   ├── page.tsx
│   │   │   ├── new/page.tsx
│   │   │   └── [appId]/
│   │   │       ├── page.tsx
│   │   │       └── settings/page.tsx
│   │   ├── tenants/                     # Gerege-specific
│   │   │   ├── page.tsx                 # Tenant жагсаалт
│   │   │   ├── new/page.tsx             # Tenant үүсгэх
│   │   │   └── [tenantId]/page.tsx      # Tenant дэлгэрэнгүй
│   │   └── settings/page.tsx
│   └── docs/
│       ├── page.tsx
│       ├── quickstart/page.tsx
│       ├── api-reference/page.tsx       # Swagger UI
│       └── guides/
│           ├── pos-plugin/page.tsx      # POS plugin хийх
│           ├── social/page.tsx          # Social нэгтгэх
│           ├── payment/page.tsx         # Payment flow
│           ├── nextjs/page.tsx
│           └── go/page.tsx
├── components/
│   ├── layout/navbar.tsx
│   ├── layout/sidebar.tsx
│   ├── apps/
│   │   ├── app-card.tsx
│   │   ├── create-app-form.tsx
│   │   ├── credentials-display.tsx
│   │   └── scope-selector.tsx          # Gerege scopes
│   └── tenants/
│       ├── tenant-card.tsx
│       └── create-tenant-form.tsx
├── lib/
│   ├── auth.ts                         # sso.gerege.mn OIDC
│   └── db.ts
├── prisma/schema.prisma
├── .env.example
├── CLAUDE.md
└── package.json
```

## Auth Config

```typescript
// lib/auth.ts
export const { handlers, signIn, signOut, auth } = NextAuth({
  providers: [{
    id: "gerege-sso",
    name: "e-ID Mongolia",
    type: "oidc",
    issuer: "https://sso.gerege.mn",      // ← sso.gerege.mn
    clientId: process.env.EID_CLIENT_ID!,
    clientSecret: process.env.EID_CLIENT_SECRET!,
    authorization: {
      params: { scope: "openid profile" }
    },
  }],
  callbacks: {
    async jwt({ token, profile }) {
      if (profile) {
        token.sub        = profile.sub
        token.certSerial = profile.cert_serial
        token.tenantId   = profile.tenant_id
        token.tenantRole = profile.tenant_role
      }
      return token
    },
    async session({ session, token }) {
      session.user.sub        = token.sub
      session.user.tenantId   = token.tenantId
      session.user.tenantRole = token.tenantRole
      return session
    },
  },
  pages: { signIn: "/auth/login" },
})
```

## Prisma Schema

```prisma
model Developer {
  id         String   @id @default(cuid())
  sub        String   @unique
  name       String
  givenName  String
  familyName String
  certSerial String
  createdAt  DateTime @default(now())
  updatedAt  DateTime @updatedAt

  apps       App[]
  tenants    TenantMember[]
}

model App {
  id           String   @id @default(cuid())
  name         String
  description  String?
  clientId     String   @unique @default(cuid())
  clientSecret String
  redirectUris String[]
  scopes       String[] @default(["openid", "profile"])
  isActive     Boolean  @default(true)
  createdAt    DateTime @default(now())
  updatedAt    DateTime @updatedAt

  developerId  String
  developer    Developer @relation(fields: [developerId], references: [id])
  tenantId     String?
  tenant       Tenant?   @relation(fields: [tenantId], references: [id])
}

model Tenant {
  id        String   @id @default(cuid())
  name      String
  slug      String   @unique
  plan      String   @default("starter")
  isActive  Boolean  @default(true)
  createdAt DateTime @default(now())

  members   TenantMember[]
  apps      App[]
}

model TenantMember {
  tenantId    String
  developerId String
  role        String   @default("member")
  createdAt   DateTime @default(now())

  tenant      Tenant    @relation(fields: [tenantId], references: [id])
  developer   Developer @relation(fields: [developerId], references: [id])

  @@id([tenantId, developerId])
}
```

## Scope Selector UI

```tsx
// components/apps/scope-selector.tsx
const GEREGE_SCOPES = [
  {
    id: "openid",
    label: "OpenID",
    description: "Үндсэн нэвтрэлт",
    locked: true,
  },
  {
    id: "profile",
    label: "Profile",
    description: "Нэр, регистрийн мэдээлэл",
    default: true,
  },
  {
    id: "pos",
    label: "POS",
    description: "Борлуулалтын систем, захиалга, бараа",
    icon: "🏪",
  },
  {
    id: "social",
    label: "Social Commerce",
    description: "Нийгмийн сүлжээ, лайв худалдаа",
    icon: "📱",
  },
  {
    id: "payment",
    label: "Payment",
    description: "QPay, SocialPay, eBarimt",
    icon: "💳",
  },
]
```

## Landing Page

```
Gerege Developer Portal

Gerege platform-д app бүтээж эхлэ.
POS plugin, social commerce, payment — нэг API-аар.

[e-ID Mongolia-р нэвтрэх]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🏪 POS Plugin API
Борлуулалт, захиалга, бараа бүртгэл

📱 Social Commerce API
Лайв худалдаа, product feed, нийтлэл

💳 Payment API
QPay, SocialPay, eBarimt нэгтгэл

🔐 e-ID нэвтрэлт
SmartID + X.509, eIDAS High
━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## Docs — Gerege-specific guides

### `/docs/guides/pos-plugin`
```markdown
## Gerege POS Plugin хийх

1. App бүртгүүл → scope: pos
2. sso.gerege.mn-р нэвтрэх (scope: openid profile pos)
3. tenant_id claim авна
4. POS API дуудах: api.gerege.mn/pos/v1/...
```

### `/docs/guides/payment`
```markdown
## QPay / SocialPay нэгтгэх

QPay: api.gerege.mn/payment/v1/qpay
SocialPay: api.gerege.mn/payment/v1/socialpay
eBarimt: api.gerege.mn/payment/v1/ebarimt
```

## Environment Variables

```env
NEXT_PUBLIC_APP_URL=https://developer.gerege.mn
NEXT_PUBLIC_SSO_URL=https://sso.gerege.mn

NEXTAUTH_URL=https://developer.gerege.mn
NEXTAUTH_SECRET=<random>

EID_CLIENT_ID=gerege-developer-portal
EID_CLIENT_SECRET=<secret>

DATABASE_URL=postgres://dev:pass@localhost:5432/gerege_dev_portal
```

## Implementation Order

```
1. CLAUDE.md
2. package.json, next.config.ts, tailwind.config.ts
3. .env.example
4. prisma/schema.prisma + migrate
5. lib/auth.ts         ← sso.gerege.mn OIDC
6. app/layout.tsx
7. app/page.tsx        ← Gerege landing
8. app/auth/login/page.tsx
9. components/layout/
10. app/dashboard/layout.tsx
11. app/dashboard/page.tsx
12. app/dashboard/apps/  (CRUD)
13. app/dashboard/tenants/ (Gerege-specific)
14. app/docs/quickstart/
15. app/docs/api-reference/  ← Swagger UI
16. app/docs/guides/pos-plugin/
17. app/docs/guides/payment/
```

---

# PART 3 — test.gerege.mn

## Зорилго

Gerege API-г **sandbox орчинд** туршина. Хөгжүүлэгчид:
- POS API test хийх (бодит transaction биш)
- Payment flow simulate хийх
- Webhook test хийх

test.gesign.mn (document signing UI) биш — **Gerege API sandbox**.

## Stack

```
Next.js 14 + TypeScript
Tailwind CSS + shadcn/ui
NextAuth → sso.gerege.mn
```

## Directory Structure

```
test.gerege.mn/
├── app/
│   ├── layout.tsx
│   ├── page.tsx              # Sandbox нүүр
│   ├── auth/login/page.tsx
│   └── sandbox/
│       ├── layout.tsx        # Auth guard
│       ├── page.tsx          # API explorer нүүр
│       ├── pos/page.tsx      # POS API test
│       ├── payment/page.tsx  # Payment simulation
│       ├── social/page.tsx   # Social API test
│       └── webhook/page.tsx  # Webhook inspector
├── components/
│   ├── api-explorer/
│   │   ├── request-builder.tsx   # Method, URL, headers, body
│   │   ├── response-viewer.tsx   # JSON highlight
│   │   └── history-panel.tsx     # Сүүлийн requests
│   ├── payment/
│   │   ├── qpay-simulator.tsx
│   │   └── ebarimt-preview.tsx
│   └── webhook/
│       └── event-log.tsx
├── lib/
│   ├── auth.ts               # sso.gerege.mn
│   └── sandbox-client.ts    # Sandbox API wrapper
├── .env.example
├── CLAUDE.md
└── package.json
```

## Pages

### `/` — Landing

```
Gerege API Sandbox

Gerege platform-ийн API-г бодит transaction хийлгүй туршина.

[Sandbox нэвтрэх]  ← sso.gerege.mn

━━━━━━━━━━━━━━
🏪 POS Sandbox
📱 Social Sandbox
💳 Payment Simulator
🔔 Webhook Inspector
━━━━━━━━━━━━━━
```

### `/sandbox` — API Explorer (нүүр)

```
┌─────────────────────────────────────────────────────┐
│  [GET ▾]  https://sandbox.gerege.mn/pos/v1/______  │
│                                                     │
│  Headers:                                           │
│  Authorization: Bearer {token}  [Автоматаар]       │
│                                                     │
│  Body: { }                                          │
│                                          [Илгээх]  │
└─────────────────────────────────────────────────────┘

Response:
┌─────────────────────────────────────────────────────┐
│  200 OK  │  124ms                                   │
│  {                                                  │
│    "items": [...]                                   │
│  }                                                  │
└─────────────────────────────────────────────────────┘
```

### `/sandbox/pos` — POS API Test

Pre-built test scenarios:
- ✅ Бараа жагсаалт авах
- ✅ Захиалга үүсгэх (test mode)
- ✅ Гүйлгээ хийх (sandbox)
- ✅ Тайлан харах

### `/sandbox/payment` — Payment Simulator

```
QPay QR code → [Simulate амжилттай] [Simulate амжилтгүй]

eBarimt → НӨАТ баримт preview
```

### `/sandbox/webhook` — Webhook Inspector

```
Таны webhook endpoint: https://test.gerege.mn/webhook/{uuid}

Incoming events:
┌──────────────────────────────────────────────────┐
│ 16:55:23  payment.completed  { amount: 50000 }  │
│ 16:55:20  order.created      { id: "ord_123" }  │
│ 16:54:11  pos.sale           { total: 25000 }   │
└──────────────────────────────────────────────────┘
```

## Auth

```typescript
// lib/auth.ts — sso.gerege.mn
providers: [{
  id: "gerege-sso",
  name: "e-ID Mongolia",
  type: "oidc",
  issuer: "https://sso.gerege.mn",
  clientId: process.env.EID_CLIENT_ID!,
  clientSecret: process.env.EID_CLIENT_SECRET!,
  authorization: {
    params: { scope: "openid profile pos social payment" }
  },
}]
```

## Sandbox API wrapper

```typescript
// lib/sandbox-client.ts
export class SandboxClient {
  private baseURL = "https://sandbox.gerege.mn"
  private token: string

  constructor(token: string) {
    this.token = token
  }

  async request(method: string, path: string, body?: object) {
    const res = await fetch(`${this.baseURL}${path}`, {
      method,
      headers: {
        "Authorization": `Bearer ${this.token}`,
        "Content-Type": "application/json",
        "X-Sandbox": "true",
      },
      body: body ? JSON.stringify(body) : undefined,
    })

    return {
      status: res.status,
      headers: Object.fromEntries(res.headers),
      body: await res.json(),
      duration: 0,  // timing inject
    }
  }
}
```

## Environment Variables

```env
NEXT_PUBLIC_APP_URL=https://test.gerege.mn
NEXT_PUBLIC_SSO_URL=https://sso.gerege.mn
NEXT_PUBLIC_SANDBOX_URL=https://sandbox.gerege.mn

NEXTAUTH_URL=https://test.gerege.mn
NEXTAUTH_SECRET=<random>

EID_CLIENT_ID=gerege-test-sandbox
EID_CLIENT_SECRET=<secret>
```

## CLAUDE.md

```markdown
# test.gerege.mn — Gerege API Sandbox

## Зорилго
Gerege platform API-г sandbox орчинд туршина.
Бодит transaction хийгдэхгүй — X-Sandbox: true header.

## Auth
sso.gerege.mn-р нэвтэрнэ.
Access token → sandbox API Bearer болгон ашиглана.

## Sandbox endpoints
sandbox.gerege.mn/pos/v1/*      ← POS sandbox
sandbox.gerege.mn/payment/v1/*  ← Payment simulator
sandbox.gerege.mn/social/v1/*   ← Social sandbox

## Webhook test
test.gerege.mn/webhook/{uuid} → event log
```

## Implementation Order

```
1. CLAUDE.md
2. package.json, next.config.ts
3. .env.example
4. lib/auth.ts          ← sso.gerege.mn
5. app/layout.tsx
6. app/page.tsx         ← Landing
7. app/auth/login/
8. app/sandbox/layout.tsx  ← Auth guard
9. app/sandbox/page.tsx    ← API Explorer
10. components/api-explorer/request-builder.tsx
11. components/api-explorer/response-viewer.tsx
12. app/sandbox/pos/
13. app/sandbox/payment/
14. app/sandbox/webhook/
15. components/webhook/event-log.tsx
```

---

# Нэгтгэсэн дэс дараалал

```
Алхам 1: sso.gerege.mn     (Go)
  mkdir sso.gerege.mn && cd sso.gerege.mn && claude
  → PART 1 prompt өгөх

Алхам 2: sso.gerege.mn-ийн клиентуудыг seed хийх
  psql → gerege-developer-portal, gerege-test-sandbox client нэмэх

Алхам 3: developer.gerege.mn   (Next.js)
  mkdir developer.gerege.mn && cd developer.gerege.mn && claude
  → PART 2 prompt өгөх

Алхам 4: test.gerege.mn    (Next.js)
  mkdir test.gerege.mn && cd test.gerege.mn && claude
  → PART 3 prompt өгөх

Алхам 5: Шалгах
  curl https://sso.gerege.mn/.well-known/openid-configuration
  curl https://sso.gerege.mn/.well-known/jwks.json
  Browser → developer.gerege.mn → нэвтрэх
  Browser → test.gerege.mn → sandbox туршах
```

---

# Хурдан шалгах checklist

```
sso.gerege.mn:
□ /.well-known/openid-configuration → issuer: https://sso.gerege.mn
□ /.well-known/jwks.json → EC P-256 key
□ /oauth/authorize → e-id.mn руу redirect
□ /health → ok

developer.gerege.mn:
□ sso.gerege.mn-р нэвтрэх
□ App үүсгэх → client_id, secret
□ Gerege scopes харагдаж байна (pos, social, payment)
□ Tenant үүсгэх

test.gerege.mn:
□ sso.gerege.mn-р нэвтрэх
□ API explorer ажиллаж байна
□ Payment simulator
□ Webhook event log
```
