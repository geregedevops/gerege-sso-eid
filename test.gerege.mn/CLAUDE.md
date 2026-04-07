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
