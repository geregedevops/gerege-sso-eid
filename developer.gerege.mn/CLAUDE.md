# developer.gerege.mn — Gerege Developer Portal

## Architecture
sso.gerege.mn        → OIDC provider (portal нэвтрэлт)
developer.gerege.mn  → Next.js 14 Developer Portal
Database             → gerege_sso_db (Prisma ORM)

## Features
- App CRUD: OAuth2 client бүртгэл, secret удирдлага
- Tenant management: pos, social, payment scope-ийн tenant
- API docs: SSO integration, POS plugin, Social, Payment, DAN заавар
- Gerege-specific scopes: pos, social, payment

## Auth
NextAuth.js v5 → sso.gerege.mn OIDC
Provider ID: "gerege-sso"
Env: EID_CLIENT_ID, EID_CLIENT_SECRET
Scope: openid profile

## Database
Prisma + PostgreSQL (gerege_sso_db)
- dev_developers — Developer profiles
- dev_tenants — Organization tenants
- dev_apps — OAuth2 applications
- dev_tenant_members — Team membership

## Key pages
/auth/login         → SSO login
/dashboard          → Overview
/dashboard/apps     → App CRUD (client_id, secret, redirect_uris)
/dashboard/apps/new → Create app
/dashboard/tenants  → Tenant management
/docs               → API documentation + guides

## Env vars
NEXT_PUBLIC_SSO_URL, EID_CLIENT_ID, EID_CLIENT_SECRET
DATABASE_URL, NEXTAUTH_SECRET, NEXTAUTH_URL
