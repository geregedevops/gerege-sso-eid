package main

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// --- Config ---

type config struct {
	// sso.gov.mn credentials
	DANClientID     string
	DANClientSecret string
	DANScope        string
	DANCallbackURI  string
	DANTokenURL     string
	DANServiceURL   string
	// server
	Port        string
	DatabaseURL string
	AdminKey    string
}

// --- Client model ---

type Client struct {
	ID           string   `json:"id"`
	Secret       string   `json:"secret,omitempty"` // only returned on create
	SecretHash   string   `json:"-"`
	Name         string   `json:"name"`
	CallbackURLs []string `json:"callback_urls"`
	Active       bool     `json:"active"`
	CreatedAt    string   `json:"created_at"`
}

// --- Main ---

func main() {
	slog.Info("starting dan.gerege.mn")

	cfg := config{
		DANClientID:     envOrDefault("DAN_CLIENT_ID", ""),
		DANClientSecret: envOrDefault("DAN_CLIENT_SECRET", ""),
		DANScope:        envOrDefault("DAN_SCOPE", ""),
		DANCallbackURI:  envOrDefault("DAN_CALLBACK_URI", "http://dan.gerege.mn/authorized"),
		DANTokenURL:     envOrDefault("DAN_TOKEN_URL", "https://sso.gov.mn/oauth2/token"),
		DANServiceURL:   envOrDefault("DAN_SERVICE_URL", "https://sso.gov.mn/oauth2/api/v1/service"),
		Port:            envOrDefault("PORT", "8444"),
		DatabaseURL:     envOrDefault("DATABASE_URL", ""),
		AdminKey:        envOrDefault("DAN_ADMIN_KEY", ""),
	}

	// Database (optional - if not configured, run without client registration)
	var db *pgxpool.Pool
	if cfg.DatabaseURL != "" {
		var err error
		db, err = pgxpool.New(context.Background(), cfg.DatabaseURL)
		if err != nil {
			slog.Error("failed to connect to database", "error", err)
			os.Exit(1)
		}
		defer db.Close()
		slog.Info("connected to database")

		// Run migration
		migration, _ := os.ReadFile("migrations/001_clients.sql")
		if len(migration) > 0 {
			if _, err := db.Exec(context.Background(), string(migration)); err != nil {
				slog.Warn("migration error (may already exist)", "error", err)
			}
		}
	}

	mux := http.NewServeMux()

	// Public pages
	mux.HandleFunc("GET /", indexHandler(cfg))
	mux.HandleFunc("GET /docs", docsHandler)
	mux.HandleFunc("GET /verify", verifyHandler(cfg, db))
	mux.HandleFunc("GET /authorized", authorizedHandler(cfg, db))
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"dan.gerege.mn"}`))
	})
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 40 40"><rect width="40" height="40" rx="10" fill="#2563eb"/><text x="50%" y="54%" dominant-baseline="middle" text-anchor="middle" fill="#fff" font-family="sans-serif" font-weight="700" font-size="18">DAN</text></svg>`))
	})

	// Admin API (requires DAN_ADMIN_KEY)
	mux.HandleFunc("GET /api/clients", adminListClients(cfg, db))
	mux.HandleFunc("POST /api/clients", adminCreateClient(cfg, db))
	mux.HandleFunc("DELETE /api/clients/{id}", adminDeleteClient(cfg, db))

	// Admin dashboard static
	mux.HandleFunc("GET /admin", adminDashboard)

	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      corsMiddleware(logMiddleware(mux)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
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

// =====================
// PUBLIC HANDLERS
// =====================

func indexHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, indexPage)
	}
}

func verifyHandler(cfg config, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		callbackURL := r.URL.Query().Get("callback_url")
		clientID := r.URL.Query().Get("client_id")

		if callbackURL == "" {
			renderError(w, "Алдаа", "callback_url параметр шаардлагатай.")
			return
		}

		// If DB is configured, validate client
		if db != nil {
			if clientID == "" {
				renderError(w, "Алдаа", "client_id параметр шаардлагатай.")
				return
			}
			client, err := getClient(db, clientID)
			if err != nil || client == nil || !client.Active {
				renderError(w, "Алдаа", "Бүртгэлгүй эсвэл идэвхгүй client.")
				return
			}
			if !matchCallbackURL(client.CallbackURLs, callbackURL) {
				renderError(w, "Алдаа", "callback_url бүртгэлгүй байна.")
				return
			}
		}

		stateJSON, _ := json.Marshal(map[string]string{
			"callback_url": callbackURL,
			"client_id":    clientID,
		})
		stateB64 := base64.RawURLEncoding.EncodeToString(stateJSON)

		loginURL := fmt.Sprintf("https://sso.gov.mn/login?state=%s&grant_type=authorization_code&response_type=code&client_id=%s&scope=%s&redirect_uri=%s",
			url.QueryEscape(stateB64),
			url.QueryEscape(cfg.DANClientID),
			url.QueryEscape(cfg.DANScope),
			url.QueryEscape(cfg.DANCallbackURI),
		)

		slog.Info("verify: redirecting", "client_id", clientID, "callback_url", callbackURL)
		http.Redirect(w, r, loginURL, http.StatusFound)
	}
}

func authorizedHandler(cfg config, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		stateB64 := r.URL.Query().Get("state")
		if code == "" {
			renderError(w, "Алдаа", "sso.gov.mn-аас code ирсэнгүй.")
			return
		}

		accessToken, err := getAccessToken(cfg, code)
		if err != nil {
			slog.Error("authorized: token exchange failed", "error", err)
			renderError(w, "Token алдаа", fmt.Sprintf("access_token авахад алдаа: %v", err))
			return
		}

		citizen, _, err := getCitizenData(cfg, accessToken)
		if err != nil {
			slog.Error("authorized: citizen data failed", "error", err)
			renderError(w, "Мэдээлэл авах алдаа", fmt.Sprintf("Иргэний мэдээлэл авахад алдаа: %v", err))
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
					clientID := state["client_id"]

					// POST full citizen data (including image) to callback URL server-to-server
					postData := make(map[string]string)
					for k, v := range citizen {
						if v != "" {
							postData[k] = v
						}
					}
					postData["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())
					if clientID != "" {
						postData["client_id"] = clientID
					}

					// Compute HMAC on non-image fields for signature
					if clientID != "" && db != nil {
						hmacParams := url.Values{}
						for k, v := range postData {
							if k != "image" {
								hmacParams.Set(k, v)
							}
						}
						client, _ := getClient(db, clientID)
						if client != nil {
							sig := computeHMAC(hmacParams, client.SecretHash)
							postData["signature"] = sig
						}
					}

					postJSON, _ := json.Marshal(postData)
					postResp, postErr := http.Post(cbURL, "application/json", strings.NewReader(string(postJSON)))
					if postErr != nil {
						slog.Error("authorized: POST to callback failed", "error", postErr)
					} else {
						postResp.Body.Close()
						slog.Info("authorized: POSTed to callback", "client_id", clientID, "status", postResp.StatusCode)
					}

					// Redirect browser with non-image params
					redirectURL, err := url.Parse(cbURL)
					if err == nil {
						params := redirectURL.Query()
						for k, v := range citizen {
							if k != "image" && v != "" {
								params.Set(k, v)
							}
						}
						params.Set("timestamp", postData["timestamp"])
						if sig, ok := postData["signature"]; ok {
							params.Set("signature", sig)
						}
						if clientID != "" {
							params.Set("client_id", clientID)
						}
						redirectURL.RawQuery = params.Encode()
						slog.Info("authorized: redirecting to callback", "client_id", clientID, "host", redirectURL.Host)
						http.Redirect(w, r, redirectURL.String(), http.StatusFound)
						return
					}
				}
			}
		}

		renderError(w, "Алдаа", "callback_url байхгүй байна. Зөвхөн бүртгэлтэй client-ээр дамжуулан ашиглана.")
	}
}

// =====================
// ADMIN API
// =====================

func requireAdmin(cfg config, w http.ResponseWriter, r *http.Request) bool {
	if cfg.AdminKey == "" {
		jsonErr(w, 403, "admin not configured")
		return false
	}
	auth := r.Header.Get("Authorization")
	if auth != "Bearer "+cfg.AdminKey {
		jsonErr(w, 401, "unauthorized")
		return false
	}
	return true
}

func adminListClients(cfg config, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !requireAdmin(cfg, w, r) {
			return
		}
		if db == nil {
			jsonErr(w, 500, "database not configured")
			return
		}

		rows, err := db.Query(r.Context(), `SELECT id, name, callback_urls, active, created_at FROM dan_clients ORDER BY created_at DESC`)
		if err != nil {
			jsonErr(w, 500, err.Error())
			return
		}
		defer rows.Close()

		clients := []Client{}
		for rows.Next() {
			var c Client
			var createdAt time.Time
			if err := rows.Scan(&c.ID, &c.Name, &c.CallbackURLs, &c.Active, &createdAt); err != nil {
				continue
			}
			c.CreatedAt = createdAt.Format(time.RFC3339)
			clients = append(clients, c)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clients)
	}
}

func adminCreateClient(cfg config, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !requireAdmin(cfg, w, r) {
			return
		}
		if db == nil {
			jsonErr(w, 500, "database not configured")
			return
		}

		var req struct {
			Name         string   `json:"name"`
			CallbackURLs []string `json:"callback_urls"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonErr(w, 400, "invalid JSON")
			return
		}
		if req.Name == "" || len(req.CallbackURLs) == 0 {
			jsonErr(w, 400, "name and callback_urls required")
			return
		}

		clientID := generateClientID()
		clientSecret := generateClientSecret()
		hash, _ := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)

		_, err := db.Exec(r.Context(),
			`INSERT INTO dan_clients (id, secret_hash, name, callback_urls) VALUES ($1, $2, $3, $4)`,
			clientID, string(hash), req.Name, req.CallbackURLs,
		)
		if err != nil {
			jsonErr(w, 500, err.Error())
			return
		}

		slog.Info("admin: client created", "id", clientID, "name", req.Name)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{
			"id":            clientID,
			"secret":        clientSecret,
			"name":          req.Name,
			"callback_urls": req.CallbackURLs,
			"message":       "Secret зөвхөн нэг удаа харагдана. Хадгалаарай!",
		})
	}
}

func adminDeleteClient(cfg config, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !requireAdmin(cfg, w, r) {
			return
		}
		if db == nil {
			jsonErr(w, 500, "database not configured")
			return
		}

		id := r.PathValue("id")
		_, err := db.Exec(r.Context(), `UPDATE dan_clients SET active = false, updated_at = now() WHERE id = $1`, id)
		if err != nil {
			jsonErr(w, 500, err.Error())
			return
		}

		slog.Info("admin: client deactivated", "id", id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deactivated", "id": id})
	}
}

// =====================
// DB HELPERS
// =====================

func getClient(db *pgxpool.Pool, clientID string) (*Client, error) {
	var c Client
	var createdAt time.Time
	err := db.QueryRow(context.Background(),
		`SELECT id, secret_hash, name, callback_urls, active, created_at FROM dan_clients WHERE id = $1`,
		clientID,
	).Scan(&c.ID, &c.SecretHash, &c.Name, &c.CallbackURLs, &c.Active, &createdAt)
	if err != nil {
		return nil, err
	}
	c.CreatedAt = createdAt.Format(time.RFC3339)
	return &c, nil
}

func matchCallbackURL(registered []string, target string) bool {
	for _, u := range registered {
		if u == target {
			return true
		}
		// Allow prefix match (e.g., "https://myapp.mn" matches "https://myapp.mn/api/dan/callback")
		if strings.HasPrefix(target, u) {
			return true
		}
	}
	return false
}

// =====================
// HMAC
// =====================

func computeHMAC(params url.Values, secret string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "signature" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(url.QueryEscape(k))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(params.Get(k)))
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(buf.String()))
	return hex.EncodeToString(mac.Sum(nil))
}

// =====================
// sso.gov.mn API
// =====================

func getAccessToken(cfg config, code string) (string, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {cfg.DANCallbackURI},
		"client_id":     {cfg.DANClientID},
		"client_secret": {cfg.DANClientSecret},
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Post(cfg.DANTokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("POST token: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	slog.Info("token response", "status", resp.StatusCode, "body", string(body))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token returned %d: %s", resp.StatusCode, string(body))
	}

	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}
	if at, ok := raw["access_token"].(string); ok && at != "" {
		return at, nil
	}
	return "", fmt.Errorf("no access_token")
}

func getCitizenData(cfg config, accessToken string) (map[string]string, string, error) {
	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {cfg.DANClientID},
		"client_secret": {cfg.DANClientSecret},
		"scope":         {cfg.DANScope},
	}

	req, err := http.NewRequest("POST", cfg.DANServiceURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("service returned %d: %s", resp.StatusCode, string(body))
	}

	var rawArr []any
	if err := json.Unmarshal(body, &rawArr); err != nil {
		return nil, "", err
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
			if r, ok := svcObj["response"].(map[string]any); ok {
				citizen = r
				break
			}
		}
		if citizen != nil {
			break
		}
	}

	if citizen == nil {
		return nil, "", fmt.Errorf("no citizen data")
	}

	fieldMap := map[string]string{
		"regnum": "reg_no", "surname": "surname", "firstname": "given_name",
		"lastname": "family_name", "civilId": "civil_id", "gender": "gender",
		"birthDateAsText": "birth_date", "birthPlace": "birth_place",
		"nationality": "nationality", "aimagCityName": "aimag_name",
		"aimagCityCode": "aimag_code", "soumDistrictName": "sum_name",
		"soumDistrictCode": "sum_code", "bagKhorooName": "bag_name",
		"bagKhorooCode": "bag_code", "addressDetail": "address_detail",
		"passportAddress": "passport_address", "passportExpireDate": "passport_expire_date",
		"passportIssueDate": "passport_issue_date", "addressApartmentName": "apartment_name",
		"addressStreetName": "street_name",
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

	if img, ok := citizen["image"].(string); ok && img != "" {
		result["image"] = img
	}

	citizenClean := make(map[string]any)
	for k, v := range citizen {
		if k != "image" {
			citizenClean[k] = v
		}
	}
	rawJSONBytes, _ := json.MarshalIndent(citizenClean, "", "  ")

	return result, string(rawJSONBytes), nil
}

// =====================
// HTML RENDERING
// =====================

func renderError(w http.ResponseWriter, title, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(502)
	fmt.Fprintf(w, errorPage, title, msg)
}

// =====================
// HELPERS
// =====================

func generateClientID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("dan_%x", b)
}

func generateClientSecret() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

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

func jsonErr(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request", "method", r.Method, "path", r.URL.Path, "ip", r.RemoteAddr, "latency_ms", time.Since(start).Milliseconds())
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// =====================
// HTML TEMPLATES
// =====================

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
.hero{text-align:center;padding:64px 24px 40px;max-width:720px;margin:0 auto}
.internal-badge{display:inline-flex;align-items:center;gap:6px;padding:6px 16px;background:rgba(245,158,11,.1);border:1px solid rgba(245,158,11,.25);border-radius:24px;font-size:12px;color:#fbbf24;font-weight:600;margin-bottom:28px}
.hero h1{font-size:42px;font-weight:800;line-height:1.15;margin-bottom:16px;color:#fff}
.hero h1 span{background:linear-gradient(135deg,#3b82f6,#2563eb);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
.hero p{max-width:560px;margin:0 auto 12px;color:#94a3b8;font-size:15px;line-height:1.7}
.hero .sub{font-size:13px;color:#64748b;margin-bottom:32px}
.cta-row{display:flex;flex-wrap:wrap;justify-content:center;gap:12px;margin-bottom:8px}
.btn{padding:14px 32px;font-weight:700;font-size:15px;border-radius:12px;text-decoration:none;transition:all .2s;display:inline-flex;align-items:center;gap:8px}
.btn-primary{background:linear-gradient(135deg,#2563eb,#1d4ed8);color:#fff;box-shadow:0 4px 16px rgba(37,99,235,.3)}
.btn-primary:hover{transform:translateY(-2px);box-shadow:0 8px 24px rgba(37,99,235,.4)}
.btn-outline{border:1px solid rgba(255,255,255,.15);color:#fff;background:transparent}
.btn-outline:hover{background:rgba(255,255,255,.05)}
.verify-hint{margin-top:12px;font-size:12px;color:#475569}

.sections{max-width:960px;margin:0 auto;padding:0 24px}
.section-title{font-size:13px;font-weight:700;color:#64748b;text-transform:uppercase;letter-spacing:1px;margin:48px 0 16px;text-align:center}

.warning-box{max-width:640px;margin:0 auto;padding:16px 20px;background:rgba(239,68,68,.06);border:1px solid rgba(239,68,68,.2);border-radius:12px;display:flex;align-items:flex-start;gap:12px}
.warning-box .icon{font-size:18px;color:#f87171;margin-top:2px}
.warning-box p{font-size:13px;color:#fca5a5;line-height:1.6;margin:0}
.warning-box a{color:#60a5fa;text-decoration:none}
.warning-box a:hover{text-decoration:underline}
.warning-box strong{color:#f87171}

.modes{display:grid;grid-template-columns:1fr 1fr;gap:16px;max-width:640px;margin:0 auto}
.mode{background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06);border-radius:14px;padding:24px;position:relative}
.mode-badge{font-size:10px;font-weight:700;padding:3px 8px;border-radius:6px;position:absolute;top:16px;right:16px}
.mode-badge.std{background:rgba(37,99,235,.15);color:#60a5fa}
.mode-badge.full{background:rgba(22,163,74,.15);color:#4ade80}
.mode h3{font-size:15px;font-weight:700;color:#fff;margin-bottom:4px}
.mode code{font-size:12px;color:#60a5fa;background:rgba(37,99,235,.1);padding:2px 8px;border-radius:4px}
.mode p{font-size:12px;color:#94a3b8;line-height:1.6;margin-top:10px}
.mode ul{font-size:12px;color:#94a3b8;padding-left:16px;margin-top:8px;line-height:1.8}

.endpoints{background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06);border-radius:16px;padding:24px;max-width:640px;margin:0 auto}
.ep{display:flex;align-items:center;gap:10px;padding:8px 0;border-bottom:1px solid rgba(255,255,255,.04);font-size:13px}
.ep:last-child{border-bottom:none}
.ep .method{font-family:monospace;font-size:11px;font-weight:700;padding:3px 8px;border-radius:6px;min-width:44px;text-align:center}
.ep .method.get{background:rgba(37,99,235,.15);color:#60a5fa}
.ep .method.post{background:rgba(245,158,11,.12);color:#fbbf24}
.ep .method.del{background:rgba(239,68,68,.12);color:#f87171}
.ep .path{font-family:monospace;color:#e2e8f0;font-size:12px}
.ep .desc{color:#64748b;font-size:11px;margin-left:auto}

.data-table{max-width:640px;margin:0 auto;background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06);border-radius:16px;overflow:hidden}
.data-table table{width:100%%;border-collapse:collapse;font-size:12px}
.data-table th{text-align:left;padding:10px 16px;background:rgba(255,255,255,.04);color:#94a3b8;font-weight:600}
.data-table td{padding:8px 16px;border-bottom:1px solid rgba(255,255,255,.03);color:#e2e8f0}
.data-table td:first-child{font-family:monospace;color:#60a5fa;font-size:11px}

.features{display:grid;grid-template-columns:repeat(auto-fit,minmax(200px,1fr));gap:16px;max-width:700px;margin:0 auto}
.feature{background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06);border-radius:12px;padding:20px}
.feature h4{font-size:13px;font-weight:700;color:#fff;margin-bottom:6px}
.feature p{font-size:12px;color:#94a3b8;line-height:1.6;margin:0}

.footer{text-align:center;padding:48px 24px 32px;font-size:12px;color:#475569}
.footer a{color:#60a5fa;text-decoration:none}
@media(max-width:640px){.hero h1{font-size:28px}.modes{grid-template-columns:1fr}}
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
    <a href="/docs">Заавар</a>
    <a href="/admin">Admin</a>
    <a href="https://docs.gerege.mn/dan/overview">Wiki</a>
  </div>
</nav>

<div class="hero">
  <div class="internal-badge">&#9888; Зөвхөн Gerege дотоод platform</div>
  <h1>DAN <span>Verify</span></h1>
  <p>sso.gov.mn-ийн ДАН системээр иргэний бүртгэлийн мэдээлэл баталгаажуулах OAuth2 gateway. Регистрийн дугаар, нэр, хаяг, зураг зэргийг авна.</p>
  <p class="sub">sso.gov.mn OAuth2 &middot; HMAC-SHA256 &middot; One-time Token</p>
  <div class="cta-row">
    <a href="%s" class="btn btn-primary">DAN Verify</a>
    <a href="/docs" class="btn btn-outline">Холболтын заавар</a>
    <a href="/admin" class="btn btn-outline">Admin</a>
  </div>
  <p class="verify-hint">sso.gov.mn руу чиглүүлэн иргэний мэдээлэл авна</p>
</div>

<div class="sections">

  <div class="warning-box">
    <span class="icon">&#9888;</span>
    <p><strong>DAN Verify-ийн мэдээллийг 3-р талд дамжуулах хориотой.</strong>
    sso.gov.mn-ийн зөвшөөрлийн дагуу зөвхөн Gerege Systems-ийн дотоод platform-ууд ашиглах эрхтэй.
    3-р тал нэвтрэлт нэгтгэхдээ <a href="https://sso.gerege.mn">sso.gerege.mn SSO</a> ашиглана.</p>
  </div>

  <div class="section-title">Verify Flow</div>
  <div class="modes">
    <div class="mode" style="max-width:420px;margin:0 auto">
      <h3>GET /verify</h3>
      <code style="display:block;margin-top:8px">/verify?client_id=XXX&amp;callback_url=XXX</code>
      <p>Иргэний мэдээлэл callback URL-д query param-р дамжина.</p>
      <ul>
        <li>РД, нэр, овог, хүйс, огноо</li>
        <li>Хаяг (аймаг, сум, баг)</li>
        <li>HMAC-SHA256 signature</li>
        <li>Timestamp (replay хамгаалалт)</li>
      </ul>
    </div>
  </div>

  <div class="section-title">API Endpoints</div>
  <div class="endpoints">
    <div class="ep"><span class="method get">GET</span><span class="path">/verify</span><span class="desc">DAN verify (зургүй)</span></div>
    <div class="ep"><span class="method get">GET</span><span class="path">/authorized</span><span class="desc">sso.gov.mn callback</span></div>
    <div class="ep"><span class="method get">GET</span><span class="path">/api/clients</span><span class="desc">Client жагсаалт (admin)</span></div>
    <div class="ep"><span class="method post">POST</span><span class="path">/api/clients</span><span class="desc">Client бүртгэх (admin)</span></div>
    <div class="ep"><span class="method del">DEL</span><span class="path">/api/clients/{id}</span><span class="desc">Client устгах (admin)</span></div>
  </div>

  <div class="section-title">Иргэний мэдээлэл (Callback параметрүүд)</div>
  <div class="data-table">
    <table>
      <thead><tr><th>Параметр</th><th>Тайлбар</th></tr></thead>
      <tbody>
        <tr><td>reg_no</td><td>Регистрийн дугаар</td></tr>
        <tr><td>given_name</td><td>Нэр</td></tr>
        <tr><td>family_name</td><td>Овог</td></tr>
        <tr><td>civil_id</td><td>Иргэний ID</td></tr>
        <tr><td>gender</td><td>Хүйс</td></tr>
        <tr><td>birth_date</td><td>Төрсөн огноо</td></tr>
        <tr><td>aimag_name</td><td>Аймаг/Хот</td></tr>
        <tr><td>sum_name</td><td>Сум/Дүүрэг</td></tr>
        <tr><td>bag_name</td><td>Баг/Хороо</td></tr>
        <tr><td>address_detail</td><td>Дэлгэрэнгүй хаяг</td></tr>
        <tr><td>signature</td><td>HMAC-SHA256</td></tr>
        <tr><td>timestamp</td><td>Unix timestamp</td></tr>
      </tbody>
    </table>
  </div>

  <div class="section-title">Онцлог</div>
  <div class="features">
    <div class="feature"><h4>Client бүртгэл</h4><p>Бүртгэлтэй client_id + callback URL шалгалт. Admin dashboard-аас удирдана.</p></div>
    <div class="feature"><h4>HMAC баталгаажуулалт</h4><p>HMAC-SHA256 signature-р мэдээллийн бүрэн бүтэн байдлыг шалгана.</p></div>
    <div class="feature"><h4>Хурдан холболт</h4><p>Нэг URL дуудахад хангалттай. sso.gov.mn credential шаардлагагүй.</p></div>
    <div class="feature"><h4>Replay хамгаалалт</h4><p>Timestamp 5 мин + token нэг удаа. Replay attack-аас хамгаална.</p></div>
  </div>

</div>

<div class="footer">
  <a href="https://docs.gerege.mn/dan/overview">Docs</a> &middot; <a href="/docs">Холболтын заавар</a> &middot; <a href="https://gerege.mn">gerege.mn</a>
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

// Admin dashboard - single page app
func adminDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(adminPage))
}

const adminPage = `<!DOCTYPE html>
<html lang="mn">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>DAN Gateway — Admin</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#0b1120;color:#e2e8f0;min-height:100vh}
nav{display:flex;align-items:center;justify-content:space-between;padding:16px 32px;border-bottom:1px solid rgba(255,255,255,.06)}
.nav-logo{display:flex;align-items:center;gap:10px;font-weight:700;font-size:16px;color:#fff}
.nav-links{display:flex;gap:24px}
.nav-links a{color:#94a3b8;font-size:13px;text-decoration:none;font-weight:500}
.nav-links a:hover{color:#fff}
.nav-links a.active{color:#60a5fa}
.container{max-width:900px;margin:0 auto;padding:32px 24px}
h1{font-size:28px;font-weight:800;color:#fff;margin-bottom:24px}
.auth-section{margin-bottom:32px;padding:20px;background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.08);border-radius:12px}
.auth-section label{display:block;font-size:13px;color:#94a3b8;margin-bottom:6px}
.auth-section input{width:100%;padding:10px 14px;background:rgba(255,255,255,.06);border:1px solid rgba(255,255,255,.1);border-radius:8px;color:#fff;font-size:14px;outline:none}
.auth-section input:focus{border-color:#2563eb}
.btn{padding:10px 20px;background:#2563eb;color:#fff;border:none;border-radius:8px;font-size:14px;font-weight:600;cursor:pointer}
.btn:hover{background:#1d4ed8}
.btn-red{background:#dc2626}
.btn-red:hover{background:#b91c1c}
.card{background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06);border-radius:12px;padding:20px;margin-bottom:16px}
.card h3{font-size:16px;font-weight:700;color:#fff;margin-bottom:4px}
.card .meta{font-size:12px;color:#64748b;margin-bottom:8px}
.card .urls{font-size:13px;color:#60a5fa;font-family:monospace}
.card .actions{margin-top:12px;display:flex;gap:8px;align-items:center}
.card .id{font-family:monospace;font-size:12px;color:#94a3b8}
.inactive{opacity:.5}
.form-group{margin-bottom:16px}
.form-group label{display:block;font-size:13px;color:#94a3b8;margin-bottom:6px}
.form-group input{width:100%;padding:10px 14px;background:rgba(255,255,255,.06);border:1px solid rgba(255,255,255,.1);border-radius:8px;color:#fff;font-size:14px;outline:none}
.secret-box{margin-top:16px;padding:16px;background:rgba(250,204,21,.08);border:1px solid rgba(250,204,21,.2);border-radius:10px}
.secret-box p{color:#fbbf24;font-size:13px;margin-bottom:8px}
.secret-box code{color:#fff;font-family:monospace;font-size:13px;word-break:break-all}
#clients-list{min-height:100px}
.empty{text-align:center;padding:40px;color:#475569}
</style>
</head>
<body>
<nav>
  <div class="nav-logo">DAN Gateway</div>
  <div class="nav-links">
    <a href="/">Тойм</a>
    <a href="/docs">Заавар</a>
    <a href="/admin" class="active">Admin</a>
  </div>
</nav>
<div class="container">
  <h1>Client бүртгэл</h1>

  <div class="auth-section">
    <label>Admin Key</label>
    <input type="password" id="admin-key" placeholder="DAN_ADMIN_KEY оруулна уу">
  </div>

  <div style="display:flex;gap:16px;margin-bottom:32px">
    <div style="flex:1">
      <h2 style="font-size:18px;font-weight:700;color:#fff;margin-bottom:16px">Шинэ client бүртгэх</h2>
      <div class="form-group"><label>Нэр</label><input type="text" id="client-name" placeholder="My App"></div>
      <div class="form-group"><label>Callback URL (олон бол таслалаар)</label><input type="text" id="callback-urls" placeholder="https://myapp.mn/api/dan/callback"></div>
      <button class="btn" onclick="createClient()">Бүртгэх</button>
      <div id="create-result"></div>
    </div>
  </div>

  <h2 style="font-size:18px;font-weight:700;color:#fff;margin-bottom:16px">Бүртгэлтэй client-ууд</h2>
  <button class="btn" onclick="loadClients()" style="margin-bottom:16px">Шинэчлэх</button>
  <div id="clients-list"><div class="empty">Admin key оруулж "Шинэчлэх" дарна уу</div></div>
</div>

<script>
function getKey() { return document.getElementById('admin-key').value; }
function headers() { return { 'Authorization': 'Bearer ' + getKey(), 'Content-Type': 'application/json' }; }

async function loadClients() {
  try {
    const res = await fetch('/api/clients', { headers: headers() });
    if (!res.ok) { document.getElementById('clients-list').innerHTML = '<div class="empty">Алдаа: ' + res.status + '</div>'; return; }
    const clients = await res.json();
    if (!clients.length) { document.getElementById('clients-list').innerHTML = '<div class="empty">Client бүртгэл байхгүй</div>'; return; }
    document.getElementById('clients-list').innerHTML = clients.map(c => '<div class="card' + (c.active ? '' : ' inactive') + '">' +
      '<h3>' + c.name + (c.active ? '' : ' <span style="color:#f87171">(идэвхгүй)</span>') + '</h3>' +
      '<div class="meta">Үүсгэсэн: ' + c.created_at + '</div>' +
      '<div class="id">ID: ' + c.id + '</div>' +
      '<div class="urls">' + (c.callback_urls||[]).join(', ') + '</div>' +
      (c.active ? '<div class="actions"><button class="btn btn-red" onclick="deleteClient(\'' + c.id + '\')">Идэвхгүй болгох</button></div>' : '') +
    '</div>').join('');
  } catch(e) { document.getElementById('clients-list').innerHTML = '<div class="empty">Холболтын алдаа</div>'; }
}

async function createClient() {
  const name = document.getElementById('client-name').value;
  const urls = document.getElementById('callback-urls').value.split(',').map(s => s.trim()).filter(Boolean);
  if (!name || !urls.length) { alert('Нэр болон callback URL оруулна уу'); return; }
  try {
    const res = await fetch('/api/clients', { method: 'POST', headers: headers(), body: JSON.stringify({ name, callback_urls: urls }) });
    const data = await res.json();
    if (!res.ok) { alert('Алдаа: ' + (data.error || res.status)); return; }
    document.getElementById('create-result').innerHTML = '<div class="secret-box"><p>&#9888; Secret зөвхөн нэг удаа харагдана!</p><code>client_id: ' + data.id + '<br>client_secret: ' + data.secret + '</code></div>';
    document.getElementById('client-name').value = '';
    document.getElementById('callback-urls').value = '';
    loadClients();
  } catch(e) { alert('Холболтын алдаа'); }
}

async function deleteClient(id) {
  if (!confirm('Энэ client-г идэвхгүй болгох уу?')) return;
  await fetch('/api/clients/' + id, { method: 'DELETE', headers: headers() });
  loadClients();
}
</script>
</body>
</html>`

// Docs page handler uses the same docsHandler but simplified reference
func docsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(docsPage))
}

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
.nav-logo{font-weight:700;font-size:16px;color:#fff}
.nav-links{display:flex;gap:24px}
.nav-links a{color:#94a3b8;font-size:13px;text-decoration:none}
.nav-links a:hover{color:#fff}
.nav-links a.active{color:#60a5fa}
.container{max-width:800px;margin:0 auto;padding:48px 24px}
h1{font-size:32px;font-weight:800;color:#fff;margin-bottom:8px}
.subtitle{color:#94a3b8;font-size:16px;margin-bottom:40px}
h2{font-size:20px;font-weight:700;color:#fff;margin:36px 0 16px;padding-top:20px;border-top:1px solid rgba(255,255,255,.06)}
h3{font-size:15px;font-weight:600;color:#e2e8f0;margin:20px 0 8px}
p{color:#94a3b8;font-size:14px;line-height:1.8;margin-bottom:12px}
code{background:rgba(255,255,255,.06);padding:2px 6px;border-radius:4px;font-size:13px;color:#60a5fa}
pre{background:rgba(255,255,255,.04);border:1px solid rgba(255,255,255,.08);border-radius:12px;padding:20px;margin:16px 0;overflow-x:auto;font-size:13px;line-height:1.7;color:#e2e8f0}
.step{display:flex;gap:16px;margin:20px 0;padding:20px;background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06);border-radius:12px}
.step-num{width:36px;height:36px;min-width:36px;background:#2563eb;border-radius:10px;display:flex;align-items:center;justify-content:center;font-weight:700;color:#fff}
.step-content h3{margin:0 0 6px}
.step-content p{margin:0;font-size:13px}
table{width:100%;border-collapse:collapse;margin:16px 0;font-size:13px}
th{text-align:left;padding:10px 14px;background:rgba(255,255,255,.04);color:#94a3b8;font-weight:600}
td{padding:10px 14px;border-bottom:1px solid rgba(255,255,255,.04)}
.note{padding:16px;background:rgba(37,99,235,.08);border:1px solid rgba(37,99,235,.15);border-radius:10px;margin:16px 0;font-size:13px;color:#93c5fd;line-height:1.7}
</style>
</head>
<body>
<nav>
  <div class="nav-logo">DAN Gateway</div>
  <div class="nav-links">
    <a href="/">Тойм</a>
    <a href="/docs" class="active">Заавар</a>
    <a href="/admin">Admin</a>
  </div>
</nav>
<div class="container">
<h1>3-р талын систем холбох</h1>
<p class="subtitle">DAN Gateway-р дамжуулан sso.gov.mn-аас иргэний мэдээлэл авах</p>

<h2>1. Client бүртгүүлэх</h2>
<p>Admin-аас <code>client_id</code> болон <code>client_secret</code> авна. <a href="/admin" style="color:#60a5fa">/admin</a> хуудаснаас бүртгүүлнэ.</p>

<h2>2. Flow</h2>
<div class="step"><div class="step-num">1</div><div class="step-content"><h3>Хэрэглэгчийг DAN Gateway руу чиглүүлэх</h3><p><code>dan.gerege.mn/verify?client_id=YOUR_ID&callback_url=YOUR_CALLBACK</code></p></div></div>
<div class="step"><div class="step-num">2</div><div class="step-content"><h3>sso.gov.mn-р нэвтрэх</h3><p>Gateway автоматаар sso.gov.mn руу redirect хийнэ.</p></div></div>
<div class="step"><div class="step-num">3</div><div class="step-content"><h3>Callback URL руу мэдээлэл буцна</h3><p>Citizen data + <code>timestamp</code> + HMAC <code>signature</code> query param-аар дамжина.</p></div></div>

<h2>3. API</h2>
<h3>GET /verify</h3>
<table>
<tr><th>Параметр</th><th>Тайлбар</th></tr>
<tr><td><code>client_id</code></td><td>Бүртгэлтэй client ID</td></tr>
<tr><td><code>callback_url</code></td><td>Бүртгэлтэй callback URL</td></tr>
</table>

<pre>https://dan.gerege.mn/verify?client_id=dan_abc123&callback_url=https://myapp.mn/api/dan/callback</pre>

<h2>4. Callback параметрүүд</h2>
<table>
<tr><th>Параметр</th><th>Тайлбар</th></tr>
<tr><td><code>reg_no</code></td><td>Регистрийн дугаар</td></tr>
<tr><td><code>given_name</code></td><td>Нэр</td></tr>
<tr><td><code>family_name</code></td><td>Овог</td></tr>
<tr><td><code>civil_id</code></td><td>Иргэний ID</td></tr>
<tr><td><code>gender</code></td><td>Хүйс</td></tr>
<tr><td><code>birth_date</code></td><td>Төрсөн огноо</td></tr>
<tr><td><code>aimag_name</code></td><td>Аймаг/Хот</td></tr>
<tr><td><code>sum_name</code></td><td>Сум/Дүүрэг</td></tr>
<tr><td><code>timestamp</code></td><td>Unix timestamp</td></tr>
<tr><td><code>client_id</code></td><td>Таны client ID</td></tr>
<tr><td><code>signature</code></td><td>HMAC-SHA256 signature</td></tr>
</table>

<h2>5. HMAC шалгалт</h2>
<p>Signature-г <code>client_secret</code> ашиглан шалгана:</p>
<pre>// Go жишээ
func verifySignature(params url.Values, secret string) bool {
    expected := params.Get("signature")
    keys := []string{}
    for k := range params {
        if k != "signature" { keys = append(keys, k) }
    }
    sort.Strings(keys)
    var buf strings.Builder
    for i, k := range keys {
        if i > 0 { buf.WriteByte('&') }
        buf.WriteString(url.QueryEscape(k) + "=" + url.QueryEscape(params.Get(k)))
    }
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(buf.String()))
    return hex.EncodeToString(mac.Sum(nil)) == expected
}</pre>

<pre>// Python жишээ
import hmac, hashlib, urllib.parse
def verify(params, secret):
    sig = params.pop('signature', '')
    canonical = '&'.join(f'{urllib.parse.quote(k)}={urllib.parse.quote(params[k])}'
                         for k in sorted(params))
    expected = hmac.new(secret.encode(), canonical.encode(), hashlib.sha256).hexdigest()
    return sig == expected</pre>

<div class="note"><strong>Анхааруулга:</strong> <code>timestamp</code> 5 минутаас хэтэрсэн бол хүлээж авахгүй байхыг зөвлөж байна (replay attack-аас хамгаалах).</div>
</div>
</body>
</html>`
