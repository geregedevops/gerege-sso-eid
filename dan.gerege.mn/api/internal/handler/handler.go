package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"dan.gerege.mn/api/internal/dan"
	"dan.gerege.mn/api/internal/store"
)

type Config struct {
	DAN           dan.Config
	DB            *store.Postgres
	StateSecret   string // HMAC key for signing state parameter
	AllowedOrigin string
	AdminKey      string // DAN admin API key
}

type Handler struct {
	cfg Config
}

func New(cfg Config) *Handler {
	return &Handler{cfg: cfg}
}

func (h *Handler) jsonError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok","service":"dan.gerege.mn"}`))
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}

func (h *Handler) Favicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 40 40"><rect width="40" height="40" rx="10" fill="#2563eb"/><text x="50%" y="54%" dominant-baseline="middle" text-anchor="middle" fill="#fff" font-family="sans-serif" font-weight="700" font-size="18">DAN</text></svg>`))
}

const indexHTML = `<!DOCTYPE html>
<html lang="mn">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>DAN Gateway — dan.gerege.mn</title>
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
.badge{display:inline-flex;align-items:center;gap:6px;padding:6px 16px;background:rgba(245,158,11,.1);border:1px solid rgba(245,158,11,.25);border-radius:24px;font-size:12px;color:#fbbf24;font-weight:600;margin-bottom:28px}
.hero h1{font-size:42px;font-weight:800;line-height:1.15;margin-bottom:16px;color:#fff}
.hero h1 span{background:linear-gradient(135deg,#3b82f6,#2563eb);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
.hero p{max-width:560px;margin:0 auto 12px;color:#94a3b8;font-size:15px;line-height:1.7}
.hero .sub{font-size:13px;color:#64748b;margin-bottom:32px}
.cta-row{display:flex;flex-wrap:wrap;justify-content:center;gap:12px;margin-bottom:8px}
.btn{padding:14px 32px;font-weight:700;font-size:15px;border-radius:12px;text-decoration:none;transition:all .2s;display:inline-flex;align-items:center;gap:8px}
.btn-dan{background:linear-gradient(135deg,#16a34a,#15803d);color:#fff;box-shadow:0 4px 16px rgba(22,163,74,.4);font-size:17px;padding:16px 40px}
.btn-dan:hover{transform:translateY(-2px);box-shadow:0 8px 24px rgba(22,163,74,.5)}
.btn-primary{background:linear-gradient(135deg,#2563eb,#1d4ed8);color:#fff;box-shadow:0 4px 16px rgba(37,99,235,.3)}
.btn-primary:hover{transform:translateY(-2px);box-shadow:0 8px 24px rgba(37,99,235,.4)}
.btn-outline{border:1px solid rgba(255,255,255,.15);color:#fff;background:transparent}
.btn-outline:hover{background:rgba(255,255,255,.05)}
.hint{margin-top:10px;font-size:12px;color:#475569}
.sections{max-width:960px;margin:0 auto;padding:0 24px}
.section-title{font-size:13px;font-weight:700;color:#64748b;text-transform:uppercase;letter-spacing:1px;margin:48px 0 16px;text-align:center}
.features{display:grid;grid-template-columns:repeat(auto-fit,minmax(200px,1fr));gap:16px;max-width:700px;margin:0 auto}
.feature{background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06);border-radius:12px;padding:20px}
.feature h4{font-size:13px;font-weight:700;color:#fff;margin-bottom:6px}
.feature p{font-size:12px;color:#94a3b8;line-height:1.6;margin:0}
.footer{text-align:center;padding:48px 24px 32px;font-size:12px;color:#475569}
.footer a{color:#60a5fa;text-decoration:none}
@media(max-width:640px){.hero h1{font-size:28px}}
</style>
</head>
<body>
<nav>
  <div class="nav-logo">
    <svg viewBox="0 0 32 32" fill="none"><rect width="32" height="32" rx="8" fill="#2563eb"/><text x="50%%" y="55%%" dominant-baseline="middle" text-anchor="middle" fill="#fff" font-family="sans-serif" font-weight="800" font-size="11">DAN</text></svg>
    DAN Gateway
  </div>
  <div class="nav-links">
    <a href="/">Нүүр</a>
    <a href="https://docs.gerege.mn/dan/overview">Заавар</a>
  </div>
</nav>

<div class="hero">
  <div class="badge">&#9888; Зөвхөн Gerege дотоод platform</div>
  <h1>DAN <span>Verify</span></h1>
  <p>sso.gov.mn-ийн ДАН системээр иргэний бүртгэлийн мэдээлэл баталгаажуулах OAuth2 gateway. Регистрийн дугаар, нэр, хаяг, зураг зэргийг авна.</p>
  <p class="sub">sso.gov.mn OAuth2 &middot; POST callback (зураг бүхий) &middot; HMAC-SHA256</p>
  <div class="cta-row">
    <a href="/try" class="btn btn-dan">&#9889; DAN Verify</a>
  </div>
  <p class="hint">Шууд sso.gov.mn-р нэвтэрч иргэний мэдээллийг харна. Client бүртгэл шаардахгүй.</p>
  <div class="cta-row" style="margin-top:20px">
    <a href="https://docs.gerege.mn/dan/overview" class="btn btn-primary">Холболтын заавар</a>
    <a href="https://developer.gerege.mn" class="btn btn-outline">Developer Portal</a>
  </div>
</div>

<div class="sections">
  <div class="section-title">Онцлог</div>
  <div class="features">
    <div class="feature"><h4>ДАН баталгаажуулалт</h4><p>sso.gov.mn OAuth2 flow-р иргэний бүрэн мэдээлэл авна</p></div>
    <div class="feature"><h4>HMAC-SHA256</h4><p>Callback дата бүрэн бүтэн, өөрчлөгдөөгүйг signature-р баталгаажуулна</p></div>
    <div class="feature"><h4>Зураг + Мэдээлэл</h4><p>Иргэний цээж зураг, хаяг, бүртгэл зэрэг бүрэн мэдээлэл POST-р ирнэ</p></div>
  </div>
</div>

<div class="footer">
  <a href="https://sso.gerege.mn">sso.gerege.mn</a> &middot; <a href="https://developer.gerege.mn">Developer Portal</a> &middot; <a href="https://docs.gerege.mn">Docs</a> &middot; <a href="https://gerege.mn">gerege.mn</a>
</div>
</body>
</html>`

// matchCallbackURL validates that target URL matches a registered callback URL
// by comparing scheme and host exactly, and path by prefix.
func matchCallbackURL(registered []string, target string) bool {
	targetURL, err := url.Parse(target)
	if err != nil {
		return false
	}
	// Require HTTPS for callbacks
	if targetURL.Scheme != "https" {
		return false
	}
	for _, u := range registered {
		regURL, err := url.Parse(u)
		if err != nil {
			continue
		}
		if targetURL.Scheme == regURL.Scheme &&
			targetURL.Host == regURL.Host &&
			strings.HasPrefix(targetURL.Path, regURL.Path) {
			return true
		}
	}
	return false
}
