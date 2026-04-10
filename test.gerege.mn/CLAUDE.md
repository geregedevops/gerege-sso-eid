# test.gerege.mn — Gerege API Sandbox

## Зорилго
Gerege platform API-г sandbox орчинд туршина.
Бодит transaction хийгдэхгүй — X-Sandbox: true header.

## Auth
- SSO: sso.gerege.mn-р нэвтэрнэ (e-ID SmartID PIN1)
- DAN: dan.gerege.mn/verify-р ДАН баталгаажуулалт хийнэ
  - DAN client_id: dan_8b2a7b8e256e812ef98fda52062c2046 (env: DAN_CLIENT_ID)
  - Callback: https://test.gerege.mn/api/dan/callback

## Sandbox endpoints
sandbox.gerege.mn/pos/v1/*      ← POS sandbox
sandbox.gerege.mn/payment/v1/*  ← Payment simulator
sandbox.gerege.mn/social/v1/*   ← Social sandbox

## Webhook test
test.gerege.mn/webhook/{uuid} → event log

## Key pages
/auth/login      → SSO + DAN login
/sandbox         → Main sandbox
/sandbox/pos     → POS API testing
/sandbox/payment → Payment simulator
/sandbox/social  → Social API testing

## Env vars
NEXT_PUBLIC_SSO_URL, EID_CLIENT_ID, EID_CLIENT_SECRET — SSO auth
DAN_URL — DAN API internal URL (http://dan-api:8444)
DAN_CLIENT_ID — DAN verify client ID
NEXTAUTH_SECRET, NEXTAUTH_URL
