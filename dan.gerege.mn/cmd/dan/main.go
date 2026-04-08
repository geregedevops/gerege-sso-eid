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

func authorizedHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			renderError(w, "Алдаа", "sso.gov.mn-аас code ирсэнгүй.")
			return
		}

		slog.Info("authorized: received callback", "has_code", true)

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

	var rows string
	for _, f := range fields {
		if v, ok := citizen[f.Key]; ok && v != "" {
			rows += fmt.Sprintf(`<tr><td style="padding:10px 16px;font-weight:600;color:#475569;white-space:nowrap;border-bottom:1px solid #f1f5f9">%s</td><td style="padding:10px 16px;color:#1e293b;border-bottom:1px solid #f1f5f9">%s</td></tr>`, f.Label, v)
		}
	}

	fmt.Fprintf(w, resultPage, rows, rawJSON)
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
<title>DAN - Иргэний мэдээлэл шалгах</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#f0f4ff;color:#1e293b;min-height:100vh;display:flex;align-items:center;justify-content:center}
.card{max-width:440px;width:100%%;margin:24px;background:#fff;border-radius:20px;box-shadow:0 4px 24px rgba(37,99,235,.08);padding:48px 40px;text-align:center}
.icon{width:64px;height:64px;background:linear-gradient(135deg,#2563eb,#1d4ed8);border-radius:16px;display:flex;align-items:center;justify-content:center;margin:0 auto 24px;color:#fff;font-weight:800;font-size:20px}
h1{font-size:24px;font-weight:700;margin-bottom:8px}
p{color:#64748b;font-size:14px;line-height:1.7;margin-bottom:28px}
.btn{display:inline-block;width:100%%;padding:16px;background:linear-gradient(135deg,#2563eb,#1d4ed8);color:#fff;font-weight:700;font-size:16px;border-radius:14px;text-decoration:none;transition:all .2s}
.btn:hover{background:linear-gradient(135deg,#3b82f6,#2563eb);transform:translateY(-1px);box-shadow:0 4px 12px rgba(37,99,235,.3)}
.hint{margin-top:20px;font-size:12px;color:#94a3b8}
.footer{margin-top:32px;font-size:12px;color:#94a3b8}
.footer a{color:#2563eb;text-decoration:none}
</style>
</head>
<body>
<div class="card">
  <div class="icon">DAN</div>
  <h1>DAN Нэвтрэлт</h1>
  <p>Монгол улсын иргэний цахим баримт бичгийн (DAN) систем ашиглан иргэний мэдээллийг баталгаажуулна.</p>
  <a href="%s" class="btn">sso.gov.mn-р нэвтрэх</a>
  <p class="hint">sso.gov.mn руу чиглүүлэн иргэний мэдээлэл авна</p>
  <div class="footer">
    <a href="https://gerege.mn">gerege.mn</a> &middot; <a href="https://dan.gov.mn">dan.gov.mn</a>
  </div>
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
