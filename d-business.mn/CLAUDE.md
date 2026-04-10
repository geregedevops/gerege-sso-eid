# d-business.mn — Байгууллагын e-Seal Platform

## Architecture
sso.gerege.mn   → OIDC provider (нэвтрэлт)
gsign.gerege.mn → G-Sign gateway (цахим гарын үсэг)
d-business.mn   → Next.js 14 Platform
Database        → gerege_sso_db (raw SQL, pg driver)

## Features
- Байгууллагын бүртгэл, баталгаажуулалт
- Гишүүд удирдлага (owner, admin, member)
- Цахим тамга (e-Seal) сертификат удирдлага
- Баримт бичиг upload, цахим гарын үсэг зурах
- Гарын үсэг шалгах (verification code)

## Auth
NextAuth.js v5 → sso.gerege.mn OIDC
Env: GEREGE_SSO_CLIENT, GEREGE_SSO_SECRET (CLAUDE memory дагуу)

## Database
Raw SQL (pg driver) — gerege_sso_db
- dbiz_users, dbiz_organizations, dbiz_org_members
- dbiz_certificates, dbiz_documents, dbiz_signatures

## Key pages
/auth/login                    → SSO login
/dashboard                     → Overview
/dashboard/org                 → Байгууллагын жагсаалт
/dashboard/org/[id]/members    → Гишүүд удирдлага
/dashboard/org/[id]/certificates → Сертификат удирдлага
/dashboard/documents           → Баримт бичиг
/verify                        → Гарын үсэг шалгах

## Env vars
GEREGE_SSO_CLIENT, GEREGE_SSO_SECRET, NEXT_PUBLIC_SSO_URL
DATABASE_URL, NEXTAUTH_SECRET, NEXTAUTH_URL, API_URL
