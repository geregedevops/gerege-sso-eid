# gsign.gerege.mn — G-Sign Digital Signature Gateway

## Architecture
gsign.gerege.mn → Go G-Sign gateway (:8445)
MSSP            → Mobile Signature Service Provider (SmartID PIN2)
d-business.mn   → Primary client (e-Seal platform)

## Sign Flow
1. Client → POST /sign (document hash + signer info)
2. G-Sign → MSSP initiate (SmartID PIN2 push)
3. Client → poll status
4. User → SmartID PIN2 оруулна
5. G-Sign → signed hash буцаана

## Env vars
PORT, ESIGN_TOKEN, MSSP_URL, APP_URL
