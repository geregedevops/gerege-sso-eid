package handler

import (
	"crypto/ecdsa"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	ocspChecker "sso.gerege.mn/internal/ocsp"
	"sso.gerege.mn/internal/store"
	"sso.gerege.mn/internal/token"
)

type Config struct {
	Issuer      string
	EIDBaseURL  string
	PrivKey     *ecdsa.PrivateKey
	PubKey      *ecdsa.PublicKey
	KID         string
	DB          store.DB
	Cache       store.Cache
	OCSP        *ocspChecker.Checker
	TokenIssuer *token.Issuer
}

type Handler struct {
	cfg Config
}

func New(cfg Config) *Handler {
	return &Handler{cfg: cfg}
}

func (h *Handler) jsonError(w http.ResponseWriter, code int, errType, desc string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error":             errType,
		"error_description": desc,
	})
}

func (h *Handler) jsonOK(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	h.jsonOK(w, map[string]string{
		"status": "ok",
		"issuer": h.cfg.Issuer,
	})
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	out := strings.ReplaceAll(indexHTML, "{{ISSUER}}", h.cfg.Issuer)
	w.Write([]byte(out))
}

const indexHTML = `<!DOCTYPE html>
<html lang="mn">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>sso.gerege.mn — e-ID SSO Provider</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#0b1120;color:#e2e8f0;min-height:100vh}
nav{display:flex;align-items:center;justify-content:space-between;padding:16px 32px;border-bottom:1px solid rgba(255,255,255,.06)}
.nav-logo{display:flex;align-items:center;gap:10px;font-weight:700;font-size:16px;color:#fff}
.nav-logo svg{width:32px;height:32px}
.nav-links{display:flex;gap:24px}
.nav-links a{color:#94a3b8;font-size:13px;text-decoration:none;font-weight:500}
.nav-links a:hover{color:#fff}
.hero{text-align:center;padding:64px 24px 40px;max-width:720px;margin:0 auto}
.open-badge{display:inline-flex;align-items:center;gap:6px;padding:6px 16px;background:rgba(22,163,74,.1);border:1px solid rgba(22,163,74,.25);border-radius:24px;font-size:12px;color:#4ade80;font-weight:600;margin-bottom:28px}
.hero h1{font-size:42px;font-weight:800;line-height:1.15;margin-bottom:16px;color:#fff}
.hero h1 span{color:#16a34a}
.hero p{max-width:560px;margin:0 auto 12px;color:#94a3b8;font-size:15px;line-height:1.7}
.hero .sub{font-size:13px;color:#64748b;margin-bottom:32px}
.cta-row{display:flex;flex-wrap:wrap;justify-content:center;gap:12px;margin-bottom:12px}
.btn{padding:14px 32px;font-weight:700;font-size:15px;border-radius:12px;text-decoration:none;transition:all .2s;display:inline-flex;align-items:center;gap:8px}
.btn-primary{background:linear-gradient(135deg,#16a34a,#15803d);color:#fff;box-shadow:0 4px 16px rgba(22,163,74,.3)}
.btn-primary:hover{transform:translateY(-2px);box-shadow:0 8px 24px rgba(22,163,74,.4)}
.btn-outline{border:1px solid rgba(255,255,255,.15);color:#fff;background:transparent}
.btn-outline:hover{background:rgba(255,255,255,.05)}

.sections{max-width:960px;margin:0 auto;padding:0 24px}
.section-title{font-size:13px;font-weight:700;color:#64748b;text-transform:uppercase;letter-spacing:1px;margin:48px 0 16px;text-align:center}

.flow{display:flex;gap:12px;justify-content:center;flex-wrap:wrap;margin:0 auto 8px}
.flow-step{background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06);border-radius:12px;padding:16px;flex:1;min-width:140px;max-width:180px;text-align:center}
.flow-num{width:28px;height:28px;background:#16a34a;border-radius:8px;display:flex;align-items:center;justify-content:center;font-weight:700;color:#fff;font-size:13px;margin:0 auto 10px}
.flow-step h4{font-size:12px;font-weight:600;color:#fff;margin-bottom:4px}
.flow-step p{font-size:11px;color:#94a3b8;margin:0;line-height:1.4}

.endpoints{background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06);border-radius:16px;padding:24px;max-width:600px;margin:0 auto}
.ep{display:flex;align-items:center;gap:10px;padding:8px 0;border-bottom:1px solid rgba(255,255,255,.04);font-size:13px}
.ep:last-child{border-bottom:none}
.ep .method{font-family:monospace;font-size:11px;font-weight:700;padding:3px 8px;border-radius:6px;min-width:44px;text-align:center}
.ep .method.get{background:rgba(22,163,74,.15);color:#4ade80}
.ep .method.post{background:rgba(245,158,11,.12);color:#fbbf24}
.ep .path{font-family:monospace;color:#e2e8f0;font-size:12px}
.ep a{color:#e2e8f0;text-decoration:none}
.ep a:hover{color:#16a34a}

.scopes{display:grid;grid-template-columns:repeat(auto-fit,minmax(160px,1fr));gap:12px;max-width:600px;margin:0 auto}
.scope{background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06);border-radius:12px;padding:16px}
.scope h4{font-size:13px;font-weight:700;color:#fff;margin-bottom:4px}
.scope code{font-size:11px;color:#16a34a;background:rgba(22,163,74,.1);padding:2px 6px;border-radius:4px}
.scope p{font-size:11px;color:#94a3b8;margin-top:6px;line-height:1.5}

.claims-box{background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06);border-radius:16px;padding:24px;max-width:600px;margin:0 auto}
.claims-box pre{background:rgba(0,0,0,.3);border-radius:10px;padding:16px;font-size:12px;color:#e2e8f0;line-height:1.6;overflow-x:auto;margin:0}

.features{display:grid;grid-template-columns:repeat(auto-fit,minmax(200px,1fr));gap:16px;max-width:700px;margin:0 auto}
.feature{background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06);border-radius:12px;padding:20px}
.feature h4{font-size:13px;font-weight:700;color:#fff;margin-bottom:6px}
.feature p{font-size:12px;color:#94a3b8;line-height:1.6;margin:0}

.footer{text-align:center;padding:48px 24px 32px;font-size:12px;color:#475569}
.footer a{color:#16a34a;text-decoration:none}
@media(max-width:640px){.hero h1{font-size:28px}.flow{flex-direction:column;align-items:center}}
</style>
</head>
<body>
<nav>
  <div class="nav-logo">
    <svg viewBox="0 0 32 32" fill="none"><rect width="32" height="32" rx="8" fill="#16a34a"/><text x="50%%" y="55%%" dominant-baseline="middle" text-anchor="middle" fill="#fff" font-family="sans-serif" font-weight="800" font-size="14">G</text></svg>
    sso.gerege.mn
  </div>
  <div class="nav-links">
    <a href="https://developer.gerege.mn">Developer Portal</a>
    <a href="https://docs.gerege.mn">Docs</a>
    <a href="/.well-known/openid-configuration">Discovery</a>
  </div>
</nav>

<div class="hero">
  <div class="open-badge">&#10003; Бүх platform-д нээлттэй</div>
  <h1>e-ID <span>SSO</span> Provider</h1>
  <p>OpenID Connect сервер — e-ID Mongolia смарт картаар баталгаажуулна. Аливаа 3-р талын систем чөлөөтэй холбогдож, иргэний мэдээллийг стандарт OIDC-р авна.</p>
  <p class="sub">ES256 JWT &middot; OAuth 2.0 Authorization Code Flow &middot; OIDC Discovery</p>
  <div class="cta-row">
    <a href="https://developer.gerege.mn/dashboard/apps/new" class="btn btn-primary">App бүртгүүлэх</a>
    <a href="https://developer.gerege.mn/docs/guides/sso-integration" class="btn btn-outline">Нэгтгэх заавар</a>
    <a href="/.well-known/openid-configuration" class="btn btn-outline">OIDC Discovery</a>
  </div>
</div>

<div class="sections">

  <div class="section-title">Integration Flow</div>
  <div class="flow">
    <div class="flow-step"><div class="flow-num">1</div><h4>App бүртгэл</h4><p>developer.gerege.mn дээр client_id авах</p></div>
    <div class="flow-step"><div class="flow-num">2</div><h4>Authorize</h4><p>/oauth/authorize руу redirect</p></div>
    <div class="flow-step"><div class="flow-num">3</div><h4>e-ID нэвтрэх</h4><p>SmartID PIN1 оруулах</p></div>
    <div class="flow-step"><div class="flow-num">4</div><h4>Token</h4><p>code &#8594; access_token + id_token</p></div>
  </div>

  <div class="section-title">OIDC Endpoints</div>
  <div class="endpoints">
    <div class="ep"><span class="method get">GET</span><a class="path" href="/.well-known/openid-configuration">/.well-known/openid-configuration</a></div>
    <div class="ep"><span class="method get">GET</span><a class="path" href="/.well-known/jwks.json">/.well-known/jwks.json</a></div>
    <div class="ep"><span class="method get">GET</span><span class="path">/oauth/authorize</span></div>
    <div class="ep"><span class="method post">POST</span><span class="path">/oauth/token</span></div>
    <div class="ep"><span class="method get">GET</span><span class="path">/oauth/userinfo</span></div>
    <div class="ep"><span class="method post">POST</span><span class="path">/oauth/revoke</span></div>
    <div class="ep"><span class="method post">POST</span><span class="path">/oauth/introspect</span></div>
  </div>

  <div class="section-title">Scopes</div>
  <div class="scopes">
    <div class="scope"><h4>Нэвтрэлт</h4><code>openid</code> <code>profile</code><p>sub, name, given_name, family_name, cert_serial</p></div>
    <div class="scope"><h4>POS Plugin</h4><code>pos</code><p>tenant_id, tenant_role, plan</p></div>
    <div class="scope"><h4>Social</h4><code>social</code><p>tenant_id, tenant_role, plan</p></div>
    <div class="scope"><h4>Payment</h4><code>payment</code><p>tenant_id, tenant_role, plan</p></div>
  </div>

  <div class="section-title">ID Token (ES256 JWT)</div>
  <div class="claims-box">
    <pre>{
  "sub": "eid-12345678",
  "name": "БАТБОЛД Ганбаатар",
  "given_name": "Ганбаатар",
  "family_name": "БАТБОЛД",
  "cert_serial": "ABC123DEF456",
  "identity_assurance_level": "high",
  "amr": ["eid"],
  "tenant_id": "t_abc123",
  "tenant_role": "owner",
  "iss": "{{ISSUER}}",
  "aud": "your-client-id"
}</pre>
  </div>

  <div class="section-title">Яагаад sso.gerege.mn?</div>
  <div class="features">
    <div class="feature"><h4>e-ID Mongolia</h4><p>SmartID.mn смарт картаар баталгаажуулсан, X.509 сертификат суурилсан</p></div>
    <div class="feature"><h4>Стандарт OIDC</h4><p>OpenID Connect Discovery, JWKS, token introspection — ямар ч хэл, framework дэмждэг</p></div>
    <div class="feature"><h4>Нэг API</h4><p>POS, Social, Payment бүгд нэг access token-р — tenant бүр тусдаа scope</p></div>
    <div class="feature"><h4>Үнэгүй</h4><p>developer.gerege.mn дээр бүртгүүлэхэд үнэгүй. Rate limit: 1000 req/min</p></div>
  </div>

</div>

<div class="footer">
  Issuer: <code style="background:rgba(255,255,255,.06);padding:2px 8px;border-radius:4px;font-size:12px;color:#16a34a">{{ISSUER}}</code><br><br>
  <a href="https://e-id.mn">e-ID Mongolia</a> &middot; <a href="https://developer.gerege.mn">Developer Portal</a> &middot; <a href="https://docs.gerege.mn">Docs</a> &middot; <a href="https://gerege.mn">gerege.mn</a>
</div>
</body>
</html>`

func (h *Handler) Favicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 40 40"><rect width="40" height="40" rx="10" fill="#16a34a"/><text x="50%" y="54%" dominant-baseline="middle" text-anchor="middle" fill="#fff" font-family="sans-serif" font-weight="700" font-size="22">G</text></svg>`))
}

func logErr(msg string, err error) {
	slog.Error(msg, "error", err)
}
