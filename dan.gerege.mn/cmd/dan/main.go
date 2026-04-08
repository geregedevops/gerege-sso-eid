package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type config struct {
	ClientID     string
	ClientSecret string
	Scope        string
	CallbackURI  string
	TokenURL     string
	ServiceURL   string
	Port         string
}

func main() {
	slog.Info("starting dan.gerege.mn")

	cfg := config{
		ClientID:     envOrDefault("DAN_CLIENT_ID", ""),
		ClientSecret: envOrDefault("DAN_CLIENT_SECRET", ""),
		Scope:        envOrDefault("DAN_SCOPE", ""),
		CallbackURI:  envOrDefault("DAN_CALLBACK_URI", "http://dan.gerege.mn/authorized"),
		TokenURL:     envOrDefault("DAN_TOKEN_URL", "https://sso.gov.mn/oauth2/token"),
		ServiceURL:   envOrDefault("DAN_SERVICE_URL", "https://sso.gov.mn/oauth2/api/v1/service"),
		Port:         envOrDefault("PORT", "8444"),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", indexHandler(cfg))
	mux.HandleFunc("GET /docs", docsHandler)
	mux.HandleFunc("GET /verify", verifyHandler(cfg))
	mux.HandleFunc("GET /authorized", authorizedHandler(cfg))
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"dan.gerege.mn"}`))
	})
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 40 40"><rect width="40" height="40" rx="10" fill="#2563eb"/><text x="50%" y="54%" dominant-baseline="middle" text-anchor="middle" fill="#fff" font-family="sans-serif" font-weight="700" font-size="18">DAN</text></svg>`))
	})

	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      logMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		slog.Info("listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	slog.Info("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

// --- Handlers ---

func indexHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		nonce := randomString(16)
		loginURL := fmt.Sprintf("https://sso.gov.mn/login?state=%s&grant_type=authorization_code&response_type=code&client_id=%s&scope=%s&redirect_uri=%s",
			url.QueryEscape(nonce),
			url.QueryEscape(cfg.ClientID),
			url.QueryEscape(cfg.Scope),
			url.QueryEscape(cfg.CallbackURI),
		)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, indexPage, loginURL)
	}
}

// verifyHandler starts DAN verification for 3rd party clients.
// GET /verify?callback_url=https://test.gerege.mn/api/dan/callback
// Encodes callback_url in state, redirects to sso.gov.mn.
func verifyHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		callbackURL := r.URL.Query().Get("callback_url")
		if callbackURL == "" {
			renderError(w, "Алдаа", "callback_url параметр шаардлагатай.")
			return
		}

		// Encode callback_url in state as base64 JSON
		stateJSON, _ := json.Marshal(map[string]string{
			"callback_url": callbackURL,
		})
		stateB64 := base64.RawURLEncoding.EncodeToString(stateJSON)

		loginURL := fmt.Sprintf("https://sso.gov.mn/login?state=%s&grant_type=authorization_code&response_type=code&client_id=%s&scope=%s&redirect_uri=%s",
			url.QueryEscape(stateB64),
			url.QueryEscape(cfg.ClientID),
			url.QueryEscape(cfg.Scope),
			url.QueryEscape(cfg.CallbackURI),
		)

		slog.Info("verify: redirecting to sso.gov.mn", "callback_url", callbackURL)
		http.Redirect(w, r, loginURL, http.StatusFound)
	}
}

func authorizedHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		stateB64 := r.URL.Query().Get("state")
		if code == "" {
			renderError(w, "Алдаа", "sso.gov.mn-аас code ирсэнгүй.")
			return
		}

		slog.Info("authorized: received callback", "has_code", true, "has_state", stateB64 != "")

		// Step 1: Exchange code for access_token
		accessToken, err := getAccessToken(cfg, code)
		if err != nil {
			slog.Error("authorized: token exchange failed", "error", err)
			renderError(w, "Token алдаа", fmt.Sprintf("sso.gov.mn-аас access_token авахад алдаа гарлаа: %v", err))
			return
		}

		// Step 2: Get citizen data
		citizen, rawJSON, err := getCitizenData(cfg, accessToken)
		if err != nil {
			slog.Error("authorized: citizen data failed", "error", err)
			renderError(w, "Мэдээлэл авах алдаа", fmt.Sprintf("Иргэний мэдээлэл авахад алдаа гарлаа: %v", err))
			return
		}

		slog.Info("authorized: success", "reg_no", citizen["reg_no"], "given_name", citizen["given_name"])

		// Check if state contains a callback_url (3rd party flow)
		if stateB64 != "" {
			stateBytes, err := base64.RawURLEncoding.DecodeString(stateB64)
			if err != nil {
				stateBytes, _ = base64.StdEncoding.DecodeString(stateB64)
			}
			var state map[string]string
			if json.Unmarshal(stateBytes, &state) == nil {
				if cbURL := state["callback_url"]; cbURL != "" {
					// Redirect to callback_url with citizen data as query params
					redirectURL, err := url.Parse(cbURL)
					if err == nil {
						params := redirectURL.Query()
						for k, v := range citizen {
							if k != "image" && v != "" {
								params.Set(k, v)
							}
						}
						redirectURL.RawQuery = params.Encode()
						slog.Info("authorized: redirecting to callback", "url", redirectURL.Host)
						http.Redirect(w, r, redirectURL.String(), http.StatusFound)
						return
					}
				}
			}
		}

		// No callback_url — render result page (standalone mode)
		renderResult(w, citizen, rawJSON)
	}
}

// --- sso.gov.mn API ---

func getAccessToken(cfg config, code string) (string, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {cfg.CallbackURI},
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Post(cfg.TokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("POST token: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	slog.Info("token response", "status", resp.StatusCode, "body", string(body))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return "", fmt.Errorf("parse token response: %w", err)
	}

	if at, ok := raw["access_token"].(string); ok && at != "" {
		return at, nil
	}
	return "", fmt.Errorf("no access_token in response")
}

func getCitizenData(cfg config, accessToken string) (map[string]string, string, error) {
	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"scope":         {cfg.Scope},
	}

	req, err := http.NewRequest("POST", cfg.ServiceURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("POST service: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("service returned %d: %s", resp.StatusCode, string(body))
	}

	// sso.gov.mn returns: [{citizen_loginType:7}, {services: {WS100101_...: {response: {...}}}}]
	var rawArr []any
	if err := json.Unmarshal(body, &rawArr); err != nil {
		return nil, "", fmt.Errorf("parse response: %w", err)
	}

	var citizen map[string]any
	for _, item := range rawArr {
		obj, ok := item.(map[string]any)
		if !ok {
			continue
		}
		services, ok := obj["services"].(map[string]any)
		if !ok {
			continue
		}
		for _, svc := range services {
			svcObj, ok := svc.(map[string]any)
			if !ok {
				continue
			}
			if resp, ok := svcObj["response"].(map[string]any); ok {
				citizen = resp
				break
			}
		}
		if citizen != nil {
			break
		}
	}

	if citizen == nil {
		return nil, "", fmt.Errorf("no citizen data in response")
	}

	// Map sso.gov.mn fields to display names
	fieldMap := map[string]string{
		"regnum":           "reg_no",
		"surname":          "surname",
		"firstname":        "given_name",
		"lastname":         "family_name",
		"civilId":          "civil_id",
		"gender":           "gender",
		"birthDateAsText":  "birth_date",
		"birthPlace":       "birth_place",
		"nationality":      "nationality",
		"aimagCityName":    "aimag_name",
		"aimagCityCode":    "aimag_code",
		"soumDistrictName": "sum_name",
		"soumDistrictCode": "sum_code",
		"bagKhorooName":    "bag_name",
		"bagKhorooCode":    "bag_code",
		"addressDetail":    "address_detail",
		"passportAddress":  "passport_address",
		"passportExpireDate":  "passport_expire_date",
		"passportIssueDate":  "passport_issue_date",
		"addressApartmentName": "apartment_name",
		"addressStreetName":    "street_name",
	}

	result := make(map[string]string)
	for ssoKey, ourKey := range fieldMap {
		if v, ok := citizen[ssoKey]; ok && v != nil {
			s := fmt.Sprintf("%v", v)
			if s != "" && s != "<nil>" {
				result[ourKey] = s
			}
		}
	}

	// Extract base64 image if present
	if img, ok := citizen["image"].(string); ok && img != "" {
		result["image"] = img
	}

	// Build raw JSON without image for display
	citizenClean := make(map[string]any)
	for k, v := range citizen {
		if k != "image" {
			citizenClean[k] = v
		}
	}
	rawJSONBytes, _ := json.MarshalIndent(citizenClean, "", "  ")

	return result, string(rawJSONBytes), nil
}

// --- Docs handler ---

func docsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(docsPage))
}

// --- HTML rendering ---

func renderResult(w http.ResponseWriter, citizen map[string]string, rawJSON string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Ordered display fields
	fields := []struct{ Key, Label string }{
		{"reg_no", "Регистрийн дугаар"},
		{"family_name", "Овог"},
		{"given_name", "Нэр"},
		{"surname", "Ургийн овог"},
		{"civil_id", "Иргэний ID"},
		{"gender", "Хүйс"},
		{"birth_date", "Төрсөн огноо"},
		{"birth_place", "Төрсөн газар"},
		{"nationality", "Үндэс"},
		{"aimag_name", "Аймаг/Хот"},
		{"sum_name", "Сум/Дүүрэг"},
		{"bag_name", "Баг/Хороо"},
		{"address_detail", "Хаягийн дэлгэрэнгүй"},
		{"passport_address", "Паспортын хаяг"},
		{"apartment_name", "Байр"},
		{"street_name", "Гудамж"},
		{"passport_issue_date", "Паспорт олгосон"},
		{"passport_expire_date", "Паспорт дуусах"},
		{"aimag_code", "Аймаг код"},
		{"sum_code", "Сум код"},
		{"bag_code", "Баг код"},
	}

	// Photo
	photoHTML := ""
	if img, ok := citizen["image"]; ok && img != "" {
		photoHTML = fmt.Sprintf(`<div style="text-align:center;margin-bottom:24px"><img src="data:image/jpeg;base64,%s" style="width:160px;height:200px;object-fit:cover;border-radius:12px;border:3px solid #e2e8f0;box-shadow:0 2px 8px rgba(0,0,0,.1)" alt="Иргэний зураг"></div>`, img)
	}

	var rows string
	for _, f := range fields {
		if v, ok := citizen[f.Key]; ok && v != "" {
			rows += fmt.Sprintf(`<tr><td style="padding:10px 16px;font-weight:600;color:#475569;white-space:nowrap;border-bottom:1px solid #f1f5f9">%s</td><td style="padding:10px 16px;color:#1e293b;border-bottom:1px solid #f1f5f9">%s</td></tr>`, f.Label, v)
		}
	}

	fmt.Fprintf(w, resultPage, photoHTML, rows, rawJSON)
}

func renderError(w http.ResponseWriter, title, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(502)
	fmt.Fprintf(w, errorPage, title, msg)
}

// --- Helpers ---

func randomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request", "method", r.Method, "path", r.URL.Path, "ip", r.RemoteAddr, "latency_ms", time.Since(start).Milliseconds())
	})
}

// --- HTML Templates ---

const indexPage = `<!DOCTYPE html>
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
.hero{text-align:center;padding:80px 24px 48px}
.badge{display:inline-flex;align-items:center;gap:6px;padding:6px 16px;background:rgba(37,99,235,.1);border:1px solid rgba(37,99,235,.2);border-radius:24px;font-size:12px;color:#60a5fa;font-weight:500;margin-bottom:32px}
.hero h1{font-size:48px;font-weight:800;line-height:1.1;margin-bottom:20px;color:#fff}
.hero h1 span{background:linear-gradient(135deg,#3b82f6,#2563eb);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
.hero p{max-width:600px;margin:0 auto 16px;color:#94a3b8;font-size:16px;line-height:1.7}
.hero .sub{font-size:13px;color:#64748b;margin-bottom:40px}
.verify-btn{display:inline-flex;align-items:center;gap:10px;padding:18px 40px;background:linear-gradient(135deg,#2563eb,#1d4ed8);color:#fff;font-weight:700;font-size:17px;border-radius:14px;text-decoration:none;transition:all .2s;box-shadow:0 4px 20px rgba(37,99,235,.3)}
.verify-btn:hover{transform:translateY(-2px);box-shadow:0 8px 30px rgba(37,99,235,.4)}
.verify-btn svg{width:22px;height:22px}
.verify-hint{margin-top:16px;font-size:12px;color:#475569}
.features{display:grid;grid-template-columns:repeat(auto-fit,minmax(260px,1fr));gap:20px;max-width:900px;margin:60px auto 0;padding:0 24px}
.feature{background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06);border-radius:16px;padding:28px}
.feature-icon{width:44px;height:44px;border-radius:12px;display:flex;align-items:center;justify-content:center;margin-bottom:16px;font-size:20px}
.feature-icon.blue{background:rgba(37,99,235,.1);color:#60a5fa}
.feature-icon.green{background:rgba(22,163,74,.1);color:#4ade80}
.feature-icon.amber{background:rgba(245,158,11,.1);color:#fbbf24}
.feature h3{font-size:15px;font-weight:700;color:#fff;margin-bottom:8px}
.feature p{font-size:13px;color:#94a3b8;line-height:1.6;margin:0}
.footer{text-align:center;padding:48px 24px 32px;font-size:12px;color:#475569}
.footer a{color:#60a5fa;text-decoration:none}
</style>
</head>
<body>
<nav>
  <div class="nav-logo">
    <svg viewBox="0 0 32 32" fill="none"><rect width="32" height="32" rx="8" fill="#2563eb"/><text x="50%%" y="55%%" dominant-baseline="middle" text-anchor="middle" fill="#fff" font-family="sans-serif" font-weight="800" font-size="11">DAN</text></svg>
    DAN Gateway
  </div>
  <div class="nav-links">
    <a href="/">Тойм</a>
    <a href="/docs">Холболтын заавар</a>
    <a href="/health">Health</a>
  </div>
</nav>

<div class="hero">
  <div class="badge">Gerege Systems LLC</div>
  <h1>DAN <span>Verify</span></h1>
  <p>Монгол Улсын ДАН (Цахим Гарын Үсэг) системтэй холбогдох OAuth2 gateway. sso.gov.mn-р дамжуулан иргэний мэдээллийг баталгаажуулна.</p>
  <p class="sub">Зөвхөн Gerege Systems LLC-ийн дотоод систем, бүтээгдэхүүнүүдэд зориулагдсан.</p>

  <a href="%s" class="verify-btn">
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 12l2 2 4-4"/><path d="M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z"/></svg>
    DAN Verify
  </a>
  <p class="verify-hint">sso.gov.mn руу чиглүүлэн иргэний мэдээлэл авна</p>
</div>

<div class="features">
  <div class="feature">
    <div class="feature-icon blue">&#128196;</div>
    <h3>Иргэний мэдээлэл</h3>
    <p>Регистрийн дугаар, овог нэр, хаяг, иргэний ID зэрэг бүрэн мэдээллийг sso.gov.mn-аас авна.</p>
  </div>
  <div class="feature">
    <div class="feature-icon green">&#9989;</div>
    <h3>Баталгаажуулалт</h3>
    <p>sso.gov.mn OAuth2 протоколоор дамжуулан баталгаажсан, итгэмжлэгдсэн мэдээлэл авна.</p>
  </div>
  <div class="feature">
    <div class="feature-icon amber">&#9889;</div>
    <h3>Хурдан холболт</h3>
    <p>Нэг товч дарахад sso.gov.mn-р нэвтэрч, иргэний мэдээллийг шууд авч харуулна.</p>
  </div>
</div>

<div class="footer">
  <a href="https://gerege.mn">gerege.mn</a> &middot; <a href="https://dan.gov.mn">dan.gov.mn</a>
</div>
</body>
</html>`

const resultPage = `<!DOCTYPE html>
<html lang="mn">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>DAN - Иргэний мэдээлэл</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#f0f4ff;color:#1e293b;min-height:100vh;padding:32px 16px}
.container{max-width:720px;margin:0 auto}
.header{text-align:center;margin-bottom:32px}
.icon{width:56px;height:56px;background:linear-gradient(135deg,#16a34a,#15803d);border-radius:14px;display:flex;align-items:center;justify-content:center;margin:0 auto 16px;color:#fff;font-size:24px}
h1{font-size:22px;font-weight:700;margin-bottom:4px}
.subtitle{color:#64748b;font-size:14px}
.badge{display:inline-block;padding:4px 12px;border-radius:20px;font-size:12px;font-weight:600;background:#dcfce7;color:#166534;margin-top:8px}
.card{background:#fff;border-radius:16px;box-shadow:0 1px 4px rgba(0,0,0,.06);overflow:hidden;margin-bottom:24px}
.card-title{padding:16px 20px;font-size:14px;font-weight:700;background:#f8fafc;border-bottom:1px solid #f1f5f9;color:#475569}
table{width:100%%;border-collapse:collapse}
.raw{padding:16px 20px}
pre{background:#f8fafc;border-radius:10px;padding:16px;font-size:12px;overflow-x:auto;color:#334155;line-height:1.6;max-height:400px;overflow-y:auto}
.actions{text-align:center;margin-top:24px}
.btn{display:inline-block;padding:12px 32px;background:#2563eb;color:#fff;font-weight:600;font-size:14px;border-radius:12px;text-decoration:none;transition:all .2s}
.btn:hover{background:#1d4ed8}
.footer{text-align:center;margin-top:32px;font-size:12px;color:#94a3b8}
.footer a{color:#2563eb;text-decoration:none}
</style>
</head>
<body>
<div class="container">
  <div class="header">
    <div class="icon">&#10003;</div>
    <h1>Иргэний мэдээлэл амжилттай авлаа</h1>
    <p class="subtitle">sso.gov.mn-р баталгаажуулсан</p>
    <span class="badge">DAN Verified</span>
  </div>

  %s

  <div class="card">
    <div class="card-title">Иргэний мэдээлэл</div>
    <table>%s</table>
  </div>

  <div class="card">
    <div class="card-title">sso.gov.mn-ийн бүтэн хариу (JSON)</div>
    <div class="raw"><pre>%s</pre></div>
  </div>

  <div class="actions">
    <a href="/" class="btn">Дахин шалгах</a>
  </div>

  <div class="footer">
    <a href="https://gerege.mn">gerege.mn</a> &middot; <a href="https://dan.gov.mn">dan.gov.mn</a>
  </div>
</div>
</body>
</html>`

const errorPage = `<!DOCTYPE html>
<html lang="mn">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>DAN - Алдаа</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#f0f4ff;color:#1e293b;min-height:100vh;display:flex;align-items:center;justify-content:center}
.card{max-width:440px;width:100%%;margin:24px;background:#fff;border-radius:20px;box-shadow:0 4px 24px rgba(239,68,68,.08);padding:48px 40px;text-align:center}
.icon{width:64px;height:64px;background:linear-gradient(135deg,#ef4444,#dc2626);border-radius:16px;display:flex;align-items:center;justify-content:center;margin:0 auto 24px;color:#fff;font-size:28px}
h1{font-size:22px;font-weight:700;margin-bottom:12px;color:#dc2626}
p{color:#64748b;font-size:14px;line-height:1.7;margin-bottom:28px;word-break:break-word}
.btn{display:inline-block;padding:14px 32px;background:#2563eb;color:#fff;font-weight:600;font-size:14px;border-radius:12px;text-decoration:none}
</style>
</head>
<body>
<div class="card">
  <div class="icon">&#10007;</div>
  <h1>%s</h1>
  <p>%s</p>
  <a href="/" class="btn">Буцах</a>
</div>
</body>
</html>`

const docsPage = `<!DOCTYPE html>
<html lang="mn">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>DAN Gateway — Холболтын заавар</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#0b1120;color:#e2e8f0;min-height:100vh}
nav{display:flex;align-items:center;justify-content:space-between;padding:16px 32px;border-bottom:1px solid rgba(255,255,255,.06)}
.nav-logo{display:flex;align-items:center;gap:10px;font-weight:700;font-size:16px;color:#fff}
.nav-logo svg{width:32px;height:32px}
.nav-links{display:flex;gap:24px}
.nav-links a{color:#94a3b8;font-size:13px;text-decoration:none;font-weight:500}
.nav-links a:hover{color:#fff}
.nav-links a.active{color:#60a5fa}
.container{max-width:800px;margin:0 auto;padding:48px 24px}
h1{font-size:36px;font-weight:800;color:#fff;margin-bottom:8px}
.subtitle{color:#94a3b8;font-size:16px;margin-bottom:40px}
h2{font-size:22px;font-weight:700;color:#fff;margin:40px 0 16px;padding-top:24px;border-top:1px solid rgba(255,255,255,.06)}
h2:first-of-type{border-top:none;margin-top:0}
h3{font-size:16px;font-weight:600;color:#e2e8f0;margin:24px 0 8px}
p{color:#94a3b8;font-size:14px;line-height:1.8;margin-bottom:12px}
code{background:rgba(255,255,255,.06);padding:2px 6px;border-radius:4px;font-size:13px;color:#60a5fa;font-family:'SF Mono',Monaco,monospace}
pre{background:rgba(255,255,255,.04);border:1px solid rgba(255,255,255,.08);border-radius:12px;padding:20px;margin:16px 0;overflow-x:auto;font-size:13px;line-height:1.7;color:#e2e8f0;font-family:'SF Mono',Monaco,monospace}
.step{display:flex;gap:16px;margin:20px 0;padding:20px;background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06);border-radius:12px}
.step-num{width:36px;height:36px;min-width:36px;background:linear-gradient(135deg,#2563eb,#1d4ed8);border-radius:10px;display:flex;align-items:center;justify-content:center;font-weight:700;font-size:15px;color:#fff}
.step-content h3{margin:0 0 6px;font-size:15px}
.step-content p{margin:0;font-size:13px}
table{width:100%;border-collapse:collapse;margin:16px 0;font-size:13px}
th{text-align:left;padding:10px 14px;background:rgba(255,255,255,.04);color:#94a3b8;font-weight:600;border-bottom:1px solid rgba(255,255,255,.08)}
td{padding:10px 14px;border-bottom:1px solid rgba(255,255,255,.04);color:#e2e8f0}
td code{color:#60a5fa}
.badge{display:inline-block;padding:2px 8px;border-radius:4px;font-size:11px;font-weight:600}
.badge.get{background:rgba(34,197,94,.1);color:#4ade80}
.badge.required{background:rgba(239,68,68,.1);color:#f87171}
.badge.optional{background:rgba(250,204,21,.1);color:#fbbf24}
.note{padding:16px 20px;background:rgba(37,99,235,.08);border:1px solid rgba(37,99,235,.15);border-radius:10px;margin:16px 0;font-size:13px;color:#93c5fd;line-height:1.7}
.footer{text-align:center;padding:48px 24px 32px;font-size:12px;color:#475569}
.footer a{color:#60a5fa;text-decoration:none}
</style>
</head>
<body>
<nav>
  <div class="nav-logo">
    <svg viewBox="0 0 32 32" fill="none"><rect width="32" height="32" rx="8" fill="#2563eb"/><text x="50%" y="55%" dominant-baseline="middle" text-anchor="middle" fill="#fff" font-family="sans-serif" font-weight="800" font-size="11">DAN</text></svg>
    DAN Gateway
  </div>
  <div class="nav-links">
    <a href="/">Тойм</a>
    <a href="/docs" class="active">Холболтын заавар</a>
    <a href="/health">Health</a>
  </div>
</nav>

<div class="container">
  <h1>3-р талын систем холбох заавар</h1>
  <p class="subtitle">DAN Gateway-р дамжуулан sso.gov.mn-аас иргэний мэдээлэл авах</p>

  <h2>Ерөнхий бүтэц</h2>
  <p>Таны систем DAN Gateway-г ашиглан sso.gov.mn-тай шууд холбогдохгүйгээр иргэний мэдээллийг баталгаажуулж авна. OAuth2 бүртгэл, client_id, client_secret хэрэггүй — зөвхөн нэг URL дуудахад хангалттай.</p>

  <div class="step">
    <div class="step-num">1</div>
    <div class="step-content">
      <h3>Хэрэглэгчийг DAN Gateway руу чиглүүлэх</h3>
      <p>Таны системээс хэрэглэгчийг <code>dan.gerege.mn/verify</code> руу redirect хийнэ. <code>callback_url</code> параметрт өөрийн callback endpoint-г зааж өгнө.</p>
    </div>
  </div>

  <div class="step">
    <div class="step-num">2</div>
    <div class="step-content">
      <h3>Хэрэглэгч sso.gov.mn-р нэвтрэх</h3>
      <p>DAN Gateway автоматаар sso.gov.mn руу чиглүүлнэ. Хэрэглэгч ДАН-аар нэвтэрнэ.</p>
    </div>
  </div>

  <div class="step">
    <div class="step-num">3</div>
    <div class="step-content">
      <h3>Callback URL руу иргэний мэдээлэл буцна</h3>
      <p>Амжилттай нэвтэрсний дараа DAN Gateway таны <code>callback_url</code> руу иргэний мэдээллийг query parameter-ээр дамжуулж redirect хийнэ.</p>
    </div>
  </div>

  <h2>API Reference</h2>

  <h3><span class="badge get">GET</span> /verify</h3>
  <p>DAN баталгаажуулалт эхлүүлэх endpoint. Хэрэглэгчийг энэ URL руу redirect хийнэ.</p>

  <table>
    <thead><tr><th>Параметр</th><th>Тайлбар</th><th>Заавал</th></tr></thead>
    <tbody>
      <tr><td><code>callback_url</code></td><td>Иргэний мэдээлэл буцаах URL. Баталгаажуулалт амжилттай болсны дараа энэ URL руу redirect хийнэ.</td><td><span class="badge required">заавал</span></td></tr>
    </tbody>
  </table>

  <h3>Жишээ URL</h3>
  <pre>https://dan.gerege.mn/verify?callback_url=https%3A%2F%2Fmyapp.mn%2Fapi%2Fdan%2Fcallback</pre>

  <h3>HTML товчны жишээ</h3>
  <pre>&lt;a href="https://dan.gerege.mn/verify?callback_url=https%3A%2F%2Fmyapp.mn%2Fapi%2Fdan%2Fcallback"&gt;
  DAN Verify
&lt;/a&gt;</pre>

  <h2>Callback Response</h2>
  <p>Амжилттай баталгаажуулалтын дараа таны <code>callback_url</code> руу дараах query parameter-ууд дамжина:</p>

  <table>
    <thead><tr><th>Параметр</th><th>Тайлбар</th><th>Жишээ</th></tr></thead>
    <tbody>
      <tr><td><code>reg_no</code></td><td>Регистрийн дугаар</td><td>АБ12345678</td></tr>
      <tr><td><code>given_name</code></td><td>Нэр</td><td>БАТБОЛД</td></tr>
      <tr><td><code>family_name</code></td><td>Овог</td><td>Дорж</td></tr>
      <tr><td><code>surname</code></td><td>Ургийн овог</td><td>Боржигон</td></tr>
      <tr><td><code>civil_id</code></td><td>Иргэний ID</td><td>110012345678</td></tr>
      <tr><td><code>gender</code></td><td>Хүйс</td><td>Эрэгтэй</td></tr>
      <tr><td><code>birth_date</code></td><td>Төрсөн огноо</td><td>1990-01-15 00:00</td></tr>
      <tr><td><code>birth_place</code></td><td>Төрсөн газар</td><td>Улаанбаатар,Баянгол</td></tr>
      <tr><td><code>nationality</code></td><td>Үндэс</td><td>Халх</td></tr>
      <tr><td><code>aimag_name</code></td><td>Аймаг/Хот</td><td>Улаанбаатар</td></tr>
      <tr><td><code>aimag_code</code></td><td>Аймаг код</td><td>11</td></tr>
      <tr><td><code>sum_name</code></td><td>Сум/Дүүрэг</td><td>Хан-Уул</td></tr>
      <tr><td><code>sum_code</code></td><td>Сум код</td><td>22</td></tr>
      <tr><td><code>bag_name</code></td><td>Баг/Хороо</td><td>1-р хороо</td></tr>
      <tr><td><code>bag_code</code></td><td>Баг код</td><td>01</td></tr>
      <tr><td><code>address_detail</code></td><td>Хаягийн дэлгэрэнгүй</td><td>59/17</td></tr>
      <tr><td><code>passport_address</code></td><td>Паспортын хаяг</td><td>УБ, Хан-Уул, 1-р хороо ...</td></tr>
    </tbody>
  </table>

  <h3>Callback URL жишээ</h3>
  <pre>https://myapp.mn/api/dan/callback?reg_no=АБ12345678&amp;given_name=БАТБОЛД&amp;family_name=Дорж&amp;civil_id=110012345678&amp;gender=Эрэгтэй&amp;birth_date=1990-01-15+00%3A00&amp;nationality=Халх&amp;aimag_name=Улаанбаатар&amp;sum_name=Хан-Уул&amp;bag_name=1-р+хороо&amp;address_detail=59%2F17</pre>

  <h2>Callback Handler жишээ</h2>

  <h3>Next.js (TypeScript)</h3>
  <pre>// app/api/dan/callback/route.ts
import { NextRequest, NextResponse } from "next/server";

export async function GET(req: NextRequest) {
  const params = req.nextUrl.searchParams;
  const regNo = params.get("reg_no");
  const givenName = params.get("given_name");
  const familyName = params.get("family_name");
  const civilId = params.get("civil_id");

  if (!regNo) {
    return NextResponse.redirect("/login?error=dan_failed");
  }

  // Иргэний мэдээллийг DB-д хадгалах, session үүсгэх гэх мэт
  // await db.upsertCitizen({ regNo, givenName, familyName, civilId });

  return NextResponse.redirect("/dashboard");
}</pre>

  <h3>Go (net/http)</h3>
  <pre>func danCallback(w http.ResponseWriter, r *http.Request) {
    regNo := r.URL.Query().Get("reg_no")
    givenName := r.URL.Query().Get("given_name")
    familyName := r.URL.Query().Get("family_name")
    civilId := r.URL.Query().Get("civil_id")

    if regNo == "" {
        http.Redirect(w, r, "/login?error=dan_failed", 302)
        return
    }

    // Иргэний мэдээллийг боловсруулах
    // ...

    http.Redirect(w, r, "/dashboard", 302)
}</pre>

  <h3>Python (Flask)</h3>
  <pre>@app.route("/api/dan/callback")
def dan_callback():
    reg_no = request.args.get("reg_no")
    given_name = request.args.get("given_name")
    family_name = request.args.get("family_name")
    civil_id = request.args.get("civil_id")

    if not reg_no:
        return redirect("/login?error=dan_failed")

    # Иргэний мэдээллийг боловсруулах
    # ...

    return redirect("/dashboard")</pre>

  <div class="note">
    <strong>Анхааруулга:</strong> callback_url нь HTTPS протоколтой, нийтэд нээлттэй URL байх ёстой. localhost дээр ажиллахгүй. Туршилтын зорилгоор <code>test.gerege.mn</code> ашиглаж болно.
  </div>

  <h2>Алдааны тохиолдол</h2>
  <p>Хэрэв хэрэглэгч sso.gov.mn дээр нэвтрэлтээ цуцалсан, эсвэл алдаа гарсан бол таны callback_url руу redirect хийгдэхгүй. Хэрэглэгч DAN Gateway-ийн алдааны хуудас дээр үлдэнэ.</p>

  <h2>Туршилт</h2>
  <p>Доорх линкээр DAN Verify-г шууд туршиж болно:</p>
  <pre>https://dan.gerege.mn/verify?callback_url=https%3A%2F%2Ftest.gerege.mn%2Fapi%2Fdan%2Fcallback</pre>
</div>

<div class="footer">
  <a href="https://gerege.mn">gerege.mn</a> &middot; <a href="https://dan.gov.mn">dan.gov.mn</a>
</div>
</body>
</html>`
