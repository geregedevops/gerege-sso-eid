# developer.gerege.mn — Gerege Developer Portal

## Architecture
sso.gerege.mn   -> OIDC provider (энэ portal нэвтрэхэд ашиглана)
developer.gerege.mn -> энэ repo (Developer Portal)

## Features
- Нэмэлт scopes: pos, social, payment
- Tenant management section
- Gerege-specific docs: POS plugin, Social commerce, Payment

## Auth
NextAuth.js v5 -> sso.gerege.mn OIDC
Provider ID: "gerege-sso"

## Database
Prisma + PostgreSQL
Shared DB with sso.gerege.mn (sso_clients table sync)

## Key pages
/dashboard/apps     -> App CRUD
/dashboard/tenants  -> Tenant management (Gerege-specific)
/docs              -> API docs, guides
