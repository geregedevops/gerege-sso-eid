export type DocPage = {
  slug: string;
  title: string;
  description?: string;
  content: string;
};

export type DocSection = {
  title: string;
  icon?: string;
  pages: DocPage[];
};

export const sections: DocSection[] = [
  {
    title: "Эхлэх",
    icon: "rocket",
    pages: [
      {
        slug: "introduction",
        title: "Танилцуулга",
        description: "Gerege platform-ийн тойм",
        content: `# Gerege Platform

Gerege нь Монгол Улсын GeregeID дэд бүтцэд суурилсан нэгдсэн platform юм.

## Platform-ийн бүрэлдэхүүн

| Service | URL | Зориулалт |
|---------|-----|-----------|
| **SSO Server** | sso.gerege.mn | OpenID Connect provider, GeregeID нэвтрэлт |
| **Developer Portal** | developer.gerege.mn | App бүртгэл, API docs, dashboard |
| **API Server** | api.gerege.mn | PDF signing, баримт бичиг API |
| **DAN Gateway** | dan.gerege.mn | ДАН иргэний мэдээлэл (дотоод) |
| **G-Sign Gateway** | gsign.gerege.mn | Тоон гарын үсэг (MSSP) |
| **Test Sandbox** | test.gerege.mn | API sandbox орчин |
| **Docs** | docs.gerege.mn | Энэ wiki site |

## Хэнд зориулсан?

- **3-р талын хөгжүүлэгчид** — sso.gerege.mn OIDC нэвтрэлт нэгтгэх
- **Gerege дотоод баг** — DAN, G-Sign, API серверийн тохиргоо
- **Системийн администратор** — Deploy, тохиргоо, мониторинг`,
      },
      {
        slug: "architecture",
        title: "Архитектур",
        description: "Системийн бүтэц, технологи",
        content: `# Архитектур

## Технологийн stack

| Давхарга | Технологи |
|----------|-----------|
| SSO Server | Go, PostgreSQL, Redis, EC JWT (ES256) |
| DAN Gateway | Go, PostgreSQL, sso.gov.mn OAuth2 |
| G-Sign Gateway | Go, MSSP ETSI TS 102 204, PKCS7 |
| API Server | Go, PostgreSQL, Redis, PDF signing |
| Developer Portal | Next.js 14, Prisma, NextAuth v5 |
| Test Sandbox | Next.js 14, NextAuth v5 |
| Docs Wiki | Next.js 14 (static) |
| Reverse Proxy | Nginx, Let's Encrypt |
| Infrastructure | Docker Compose, PostgreSQL 15, Redis 7 |

## Сүлжээний бүтэц

\`\`\`
                    ┌──────────┐
    Internet ──────▶│  Nginx   │ :80/:443
                    └────┬─────┘
          ┌──────────────┼──────────────┐
          │              │              │
    ┌─────▼─────┐ ┌─────▼─────┐ ┌─────▼─────┐
    │    SSO    │ │    DAN    │ │   G-Sign  │
    │  :8443   │ │  :8444   │ │  :8445   │
    └─────┬─────┘ └─────┬─────┘ └───────────┘
          │              │
    ┌─────▼──────────────▼─────┐
    │     PostgreSQL :5432     │
    │       Redis :6379        │
    └──────────────────────────┘
\`\`\`

## Docker Compose

Бүх service Docker Compose-р удирддаг. \`.env\` файлд бүх secret хадгалагдана.

\`\`\`bash
# Бүх service эхлүүлэх
docker compose up -d

# Тодорхой service rebuild
docker compose up -d --build sso

# Log харах
docker compose logs -f sso
\`\`\``,
      },
    ],
  },
  {
    title: "SSO Server",
    icon: "shield",
    pages: [
      {
        slug: "sso/overview",
        title: "SSO тойм",
        description: "OpenID Connect provider",
        content: `# SSO Server (sso.gerege.mn)

OpenID Connect 1.0 provider — GeregeID смарт картаар баталгаажуулна.

## Нээлттэй — бүх 3-р талын platform-д

sso.gerege.mn нь аливаа систем, platform, апп-д нээлттэй.
developer.gerege.mn дээр app бүртгүүлж client_id авахад хангалттай.

## OIDC Endpoints

| Endpoint | URL |
|----------|-----|
| Discovery | \`/.well-known/openid-configuration\` |
| Authorization | \`/oauth/authorize\` |
| Token | \`/oauth/token\` |
| UserInfo | \`/oauth/userinfo\` |
| JWKS | \`/.well-known/jwks.json\` |
| Introspect | \`/oauth/introspect\` |
| Revoke | \`/oauth/revoke\` |

## Дэмждэг scopes

| Scope | Тайлбар |
|-------|---------|
| \`openid\` | Заавал — sub, iss, aud |
| \`profile\` | name, given_name, family_name, cert_serial |
| \`pos\` | POS Plugin API — tenant_id, tenant_role, plan |
| \`social\` | Social Commerce API |
| \`payment\` | Payment API |

## ID Token Claims

\`\`\`json
{
  "sub": "eid-12345678",
  "name": "БАТБОЛД Ганбаатар",
  "given_name": "Ганбаатар",
  "family_name": "БАТБОЛД",
  "cert_serial": "ABC123DEF456",
  "identity_assurance_level": "high",
  "amr": ["eid"],
  "tenant_id": "t_abc123",
  "tenant_role": "owner",
  "plan": "pro"
}
\`\`\`

## Тохиргоо (Environment)

| Variable | Тайлбар |
|----------|---------|
| \`SSO_ISSUER\` | Issuer URL (https://sso.gerege.mn) |
| \`SSO_PRIVATE_KEY_PATH\` | EC private key path |
| \`DATABASE_URL\` | PostgreSQL connection string |
| \`REDIS_URL\` | Redis connection string |
| \`EID_BASE_URL\` | GeregeID API URL |`,
      },
      {
        slug: "sso/integration",
        title: "SSO нэгтгэх",
        description: "3-р талын platform холбох",
        content: `# SSO нэгтгэх заавар

Дэлгэрэнгүй зааврыг [developer.gerege.mn/docs/guides/sso-integration](https://developer.gerege.mn/docs/guides/sso-integration) дээрээс үзнэ.

## Товч алхам

1. **App бүртгүүлэх** — [developer.gerege.mn/dashboard/apps/new](https://developer.gerege.mn/dashboard/apps/new)
2. **client_id, client_secret авах** — зөвхөн нэг удаа харагдана
3. **OIDC холбох** — Authorization Code Flow
4. **ID Token задлах** — иргэний мэдээлэл авах

## Жишээ (Next.js)

\`\`\`typescript
import NextAuth from "next-auth"

export const { handlers, signIn, signOut, auth } = NextAuth({
  providers: [{
    id: "gerege-sso",
    name: "GeregeID",
    type: "oidc",
    issuer: "https://sso.gerege.mn",
    clientId: process.env.EID_CLIENT_ID!,
    clientSecret: process.env.EID_CLIENT_SECRET!,
  }],
})
\`\`\`

Бусад хэлний жишээ (Go, Python, PHP): [developer.gerege.mn/docs/guides/sso-integration](https://developer.gerege.mn/docs/guides/sso-integration)`,
      },
    ],
  },
  {
    title: "DAN Gateway",
    icon: "id-card",
    pages: [
      {
        slug: "dan/overview",
        title: "DAN тойм",
        description: "ДАН иргэний мэдээлэл",
        content: `# DAN Gateway (dan.gerege.mn)

sso.gov.mn-ийн ДАН системээс иргэний бүртгэлийн мэдээлэл авах gateway.

> **Анхааруулга:** DAN Verify нь зөвхөн Gerege-ийн дотоод platform-уудад зориулагдсан.
> 3-р талын системд DAN мэдээлэл дамжуулах хориотой.
> 3-р тал нэвтрэлт нэгтгэхдээ SSO (sso.gerege.mn) ашиглана.

## Хоёр горим

| Горим | Endpoint | Тайлбар |
|-------|----------|---------|
| Зургүй | \`/verify\` | Callback URL-д query param-р дамжина |
| Зураг бүхий | \`/verify-full\` | Token-р дамжуулж API-р бүтэн data авна |

## API Endpoints

| Endpoint | Тайлбар |
|----------|---------|
| \`GET /verify\` | DAN verify эхлүүлэх (зургүй) |
| \`GET /verify-full\` | DAN verify эхлүүлэх (зураг бүхий) |
| \`GET /authorized\` | sso.gov.mn callback |
| \`GET /api/citizen?token=xxx\` | Зураг бүхий бүтэн data авах |
| \`GET /api/clients\` | Client жагсаалт (admin) |
| \`POST /api/clients\` | Client бүртгэх (admin) |
| \`GET /docs\` | Холболтын заавар |
| \`GET /admin\` | Admin dashboard |

## Тохиргоо

| Variable | Тайлбар |
|----------|---------|
| \`DAN_CLIENT_ID\` | sso.gov.mn client ID |
| \`DAN_CLIENT_SECRET\` | sso.gov.mn client secret |
| \`DAN_SCOPE\` | ДАН service scope (base64) |
| \`DAN_CALLBACK_URI\` | OAuth2 redirect URI |
| \`DAN_ADMIN_KEY\` | Admin API key |
| \`DATABASE_URL\` | PostgreSQL connection |`,
      },
    ],
  },
  {
    title: "G-Sign Gateway",
    icon: "pen",
    pages: [
      {
        slug: "gsign/overview",
        title: "G-Sign тойм",
        description: "Тоон гарын үсэг",
        content: `# G-Sign Gateway (gsign.gerege.mn)

УБЕГ-ийн GSign клауд тоон гарын үсгийн gateway.
MSSP ETSI TS 102 204 протоколоор ажиллана.

## Flow

1. Хэрэглэгч утасны дугаараа оруулна
2. MSSP руу MSS_SignatureReq илгээгдэнэ
3. GSign апп дээр PIN оруулна (synch, 120 сек)
4. Base64Signature буцна → PKCS7/CMS decode
5. Сертификатын SubjectDN-ээс иргэний мэдээлэл задлана

## API Endpoints

| Endpoint | Method | Тайлбар |
|----------|--------|---------|
| \`/\` | GET | Landing page + утасны дугаар form |
| \`/sign\` | POST | Signature request (\`{phoneNo, callbackUrl?}\`) |
| \`/verify?callback_url=...\` | GET | 3-р тал flow |

## Техникийн stack

- **MSSP**: Methics Kiuru MSSP (ETSI TS 102 204)
- **CMS Decode**: Go \`go.mozilla.org/pkcs7\` (Java dependency-гүй)
- **Сертификат**: SubjectDN → SERIALNUMBER (РД), CN, GivenName, Surname

## Тохиргоо

| Variable | Тайлбар |
|----------|---------|
| \`ESIGN_TOKEN\` | MSSP Basic Auth token |
| \`MSSP_URL\` | MSSP REST endpoint |
| \`APP_URL\` | Public URL |`,
      },
    ],
  },
  {
    title: "API Server",
    icon: "server",
    pages: [
      {
        slug: "api/overview",
        title: "API тойм",
        description: "REST API server",
        content: `# API Server (api.gerege.mn)

PDF signing, баримт бичиг боловсруулах REST API.

## Баталгаажуулалт

SSO access token ашиглан Bearer authentication хийнэ:

\`\`\`
Authorization: Bearer <access_token>
\`\`\`

JWKS URI: \`https://sso.gerege.mn/.well-known/jwks.json\`

## Тохиргоо

| Variable | Тайлбар |
|----------|---------|
| \`SSO_JWKS_URI\` | JWT шалгах JWKS endpoint |
| \`EID_API_URL\` | GeregeID API (https://ca.gerege.mn) |
| \`STORAGE_PATH\` | Гарын үсэг зурсан файлын хадгалах зам |
| \`DATABASE_URL\` | PostgreSQL connection |
| \`REDIS_URL\` | Redis connection |`,
      },
    ],
  },
  {
    title: "Deploy & Ops",
    icon: "terminal",
    pages: [
      {
        slug: "ops/deploy",
        title: "Deploy заавар",
        description: "Server тохиргоо, deploy",
        content: `# Deploy заавар

## Шаардлага

- Docker + Docker Compose
- Domain DNS тохиргоо (\`*.gerege.mn\` → server IP)
- SSH хандалт

## Анхны тохиргоо

\`\`\`bash
# 1. Repo clone
git clone git@github.com:erdenebatt/gerege-sso-eid.git
cd gerege-sso-eid

# 2. .env файл үүсгэх
cp .env.example .env
# .env файл засах — бүх secret-ээ оруулна

# 3. SSL сертификат авах
docker compose up -d nginx
docker compose run --rm certbot certonly --webroot -w /var/www/certbot \\
  -d sso.gerege.mn -d developer.gerege.mn -d test.gerege.mn \\
  -d api.gerege.mn -d dan.gerege.mn -d gsign.gerege.mn -d docs.gerege.mn

# 4. Бүх service эхлүүлэх
docker compose up -d
\`\`\`

## Шинэчлэх

\`\`\`bash
git pull origin main
docker compose up -d --build <service_name>
docker compose exec nginx nginx -s reload  # nginx config өөрчилсөн бол
\`\`\`

## SSL сертификат шинэчлэх

\`\`\`bash
docker compose run --rm certbot renew
docker compose exec nginx nginx -s reload
\`\`\`

## Мониторинг

\`\`\`bash
# Бүх service статус
docker compose ps

# Log
docker compose logs -f <service_name>

# Health check
curl https://sso.gerege.mn/.well-known/openid-configuration
curl https://dan.gerege.mn/health
curl https://gsign.gerege.mn/health
curl https://api.gerege.mn/health
\`\`\``,
      },
      {
        slug: "ops/env",
        title: "Environment variables",
        description: "Бүх тохиргооны хувьсагч",
        content: `# Environment Variables

\`.env.example\` файлд бүх тохиргоо байна.

## Postgres
| Variable | Тайлбар |
|----------|---------|
| \`POSTGRES_USER\` | DB хэрэглэгчийн нэр |
| \`POSTGRES_PASSWORD\` | DB нууц үг |
| \`POSTGRES_DB\` | Database нэр |

## Redis
| Variable | Тайлбар |
|----------|---------|
| \`REDIS_PASSWORD\` | Redis нууц үг |

## DAN
| Variable | Тайлбар |
|----------|---------|
| \`DAN_CLIENT_ID\` | sso.gov.mn client ID |
| \`DAN_CLIENT_SECRET\` | sso.gov.mn client secret |
| \`DAN_SCOPE\` | Service scope (base64 JSON) |
| \`DAN_CALLBACK_URI\` | OAuth2 redirect URI |
| \`DAN_ADMIN_KEY\` | Admin API bearer token |

## Developer Portal
| Variable | Тайлбар |
|----------|---------|
| \`DEV_NEXTAUTH_SECRET\` | NextAuth session encryption |
| \`DEV_EID_CLIENT_ID\` | SSO app client ID |
| \`DEV_EID_CLIENT_SECRET\` | SSO app client secret |

## Test Sandbox
| Variable | Тайлбар |
|----------|---------|
| \`TEST_NEXTAUTH_SECRET\` | NextAuth session encryption |
| \`TEST_EID_CLIENT_ID\` | SSO app client ID |
| \`TEST_EID_CLIENT_SECRET\` | SSO app client secret |

## G-Sign
| Variable | Тайлбар |
|----------|---------|
| \`ESIGN_TOKEN\` | MSSP Basic Auth token (base64) |
| \`MSSP_URL\` | MSSP REST endpoint URL |`,
      },
    ],
  },
];

export function getPage(slug: string[]): DocPage | undefined {
  const path = slug.join("/");
  for (const section of sections) {
    for (const page of section.pages) {
      if (page.slug === path) return page;
    }
  }
  return undefined;
}

export function getAllSlugs(): string[][] {
  const slugs: string[][] = [];
  for (const section of sections) {
    for (const page of section.pages) {
      slugs.push(page.slug.split("/"));
    }
  }
  return slugs;
}
