package handler

import (
	"crypto/ecdsa"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	ocspChecker "gesign.mn/gerege-sso/internal/ocsp"
	"gesign.mn/gerege-sso/internal/store"
	"gesign.mn/gerege-sso/internal/token"
)

type Config struct {
	Issuer      string
	EIDBaseURL  string
	PrivKey     *ecdsa.PrivateKey
	PubKey      *ecdsa.PublicKey
	KID         string
	DB          *store.Postgres
	Cache       *store.Redis
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
<title>sso.gerege.mn — OpenID Connect</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#f8fafc;color:#1e293b;min-height:100vh;display:flex;align-items:center;justify-content:center}
.card{max-width:520px;width:100%%;margin:24px;background:#fff;border-radius:16px;box-shadow:0 1px 3px rgba(0,0,0,.1);padding:40px}
.logo{display:flex;align-items:center;gap:10px;margin-bottom:24px}
.logo-icon{width:40px;height:40px;background:#16a34a;border-radius:10px;display:flex;align-items:center;justify-content:center;color:#fff;font-weight:700;font-size:18px}
.logo span{font-size:18px;font-weight:700;color:#1e293b}
h1{font-size:22px;font-weight:700;margin-bottom:8px}
p{color:#64748b;font-size:14px;line-height:1.6;margin-bottom:16px}
.badge{display:inline-block;padding:3px 10px;border-radius:20px;font-size:11px;font-weight:600;background:#dcfce7;color:#166534;margin-bottom:20px}
.endpoints{background:#f1f5f9;border-radius:10px;padding:16px;margin:20px 0}
.endpoints h3{font-size:13px;font-weight:600;color:#475569;margin-bottom:10px}
.ep{display:flex;align-items:center;gap:8px;padding:6px 0;font-size:13px}
.ep .method{font-family:monospace;font-size:11px;font-weight:600;padding:2px 6px;border-radius:4px;min-width:40px;text-align:center}
.ep .method.get{background:#dcfce7;color:#166534}
.ep .method.post{background:#fef3c7;color:#92400e}
.ep .path{font-family:monospace;color:#334155;font-size:12px}
a{color:#16a34a;text-decoration:none}
a:hover{text-decoration:underline}
.footer{text-align:center;margin-top:24px;font-size:12px;color:#94a3b8}
</style>
</head>
<body>
<div class="card">
  <div class="logo"><div class="logo-icon">G</div><span>sso.gerege.mn</span></div>
  <span class="badge">OpenID Connect Provider</span>
  <h1>Gerege OIDC Authorization Server</h1>
  <p>Gerege platform ecosystem-ийн OpenID Connect сервер. e-ID Mongolia (SmartID.mn) ашиглан POS, Social Commerce, Payment зэрэг бүрэлдэхүүн хэсгүүдэд стандарт OAuth2/OIDC баталгаажуулалт хийнэ.</p>

  <div class="endpoints">
    <h3>Endpoints</h3>
    <div class="ep"><span class="method get">GET</span><a class="path" href="/.well-known/openid-configuration">/.well-known/openid-configuration</a></div>
    <div class="ep"><span class="method get">GET</span><a class="path" href="/.well-known/jwks.json">/.well-known/jwks.json</a></div>
    <div class="ep"><span class="method get">GET</span><span class="path">/oauth/authorize</span></div>
    <div class="ep"><span class="method post">POST</span><span class="path">/oauth/token</span></div>
    <div class="ep"><span class="method get">GET</span><span class="path">/oauth/userinfo</span></div>
    <div class="ep"><span class="method post">POST</span><span class="path">/oauth/revoke</span></div>
    <div class="ep"><span class="method post">POST</span><span class="path">/oauth/introspect</span></div>
    <div class="ep"><span class="method get">GET</span><a class="path" href="/health">/health</a></div>
  </div>

  <p>Issuer: <code>{{ISSUER}}</code></p>

  <div class="footer">
    Powered by <a href="https://e-id.mn">e-ID Mongolia</a> &middot; <a href="https://gerege.mn">gerege.mn</a>
  </div>
</div>
</body>
</html>`

func logErr(msg string, err error) {
	slog.Error(msg, "error", err)
}
