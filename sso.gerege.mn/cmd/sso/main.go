package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gesign.mn/gerege-sso/internal/handler"
	"gesign.mn/gerege-sso/internal/middleware"
	ocspChecker "gesign.mn/gerege-sso/internal/ocsp"
	"gesign.mn/gerege-sso/internal/store"
	"gesign.mn/gerege-sso/internal/token"
)

func main() {
	slog.Info("starting sso.gerege.mn")

	// Config
	issuer := envOrDefault("SSO_ISSUER", "https://sso.gerege.mn")
	privKeyPath := envOrDefault("SSO_PRIVATE_KEY_PATH", "ec-private.pem")
	eidBaseURL := envOrDefault("EID_BASE_URL", "https://e-id.mn")
	ocspURL := envOrDefault("OCSP_URL", "https://ocsp.gesign.mn/ocsp")
	caIssuingURL := envOrDefault("CA_ISSUING_URL", "https://ca.gesign.mn/api/ca/download/issuing.pem")
	danClientID := envOrDefault("DAN_CLIENT_ID", "f3f14ab1af2cf74fd7ade8a0-964f9d4992277df04d43aef0c80a1152")
	danScope := envOrDefault("DAN_SCOPE", "W3sic2VydmljZXMiOlsiV1MxMDAxMDFfZ2V0Q2l0aXplbklEQ2FyZEluZm8iXSwid3NkbCI6Imh0dHBzOlwvXC94eXAuZ292Lm1uXC9jaXRpemVuLTEuMy4wXC93cz9XU0RMIn1d")
	danCallbackURI := envOrDefault("DAN_CALLBACK_URI", "https://sso.gerege.mn/callback/dan")
	databaseURL := envOrDefault("DATABASE_URL", "postgres://sso:pass@localhost:5432/gerege_sso_db")
	redisURL := envOrDefault("REDIS_URL", "redis://localhost:6379/2")
	port := envOrDefault("PORT", "8443")
	tlsCert := os.Getenv("TLS_CERT")
	tlsKey := os.Getenv("TLS_KEY")
	devMode := os.Getenv("DEV_MODE") == "true"

	// Load EC key pair
	privKey, pubKey, err := loadECKeys(privKeyPath)
	if err != nil {
		slog.Error("failed to load EC keys", "error", err)
		os.Exit(1)
	}

	kid := token.ComputeKID(pubKey)
	slog.Info("loaded EC key", "kid", kid)

	// Database
	ctx := context.Background()
	db, err := store.NewPostgres(ctx, databaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Redis
	cache, err := store.NewRedis(redisURL)
	if err != nil {
		slog.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}
	defer cache.Close()

	// OCSP checker
	ocsp := ocspChecker.NewChecker(ocspURL, caIssuingURL, cache)

	// Token issuer
	tokenIssuer := token.NewIssuer(privKey, kid, issuer)

	// Handler
	h := handler.New(handler.Config{
		Issuer:         issuer,
		EIDBaseURL:     eidBaseURL,
		PrivKey:        privKey,
		PubKey:         pubKey,
		KID:            kid,
		DB:             db,
		Cache:          cache,
		OCSP:           ocsp,
		TokenIssuer:    tokenIssuer,
		DANClientID:    danClientID,
		DANScope:       danScope,
		DANCallbackURI: danCallbackURI,
	})

	// Router
	mux := http.NewServeMux()
	mux.HandleFunc("GET /.well-known/openid-configuration", h.Discovery)
	mux.HandleFunc("GET /.well-known/jwks.json", h.JWKS)
	mux.HandleFunc("GET /oauth/authorize", h.Authorize)
	mux.HandleFunc("POST /oauth/token", h.Token)
	mux.HandleFunc("GET /oauth/userinfo", h.UserInfo)
	mux.HandleFunc("POST /oauth/revoke", h.Revoke)
	mux.HandleFunc("POST /oauth/introspect", h.Introspect)
	mux.HandleFunc("GET /callback/eid", h.EIDCallback)
	mux.HandleFunc("GET /callback/dan", h.DANCallback)
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /", h.Index)

	// Middleware
	var root http.Handler = mux
	root = middleware.Logger(root)
	root = middleware.CORS(nil)(root)

	// Server
	addr := ":" + port
	srv := &http.Server{
		Addr:         addr,
		Handler:      root,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start
	go func() {
		if !devMode && tlsCert != "" && tlsKey != "" {
			slog.Info("starting HTTPS server", "addr", addr)
			if err := srv.ListenAndServeTLS(tlsCert, tlsKey); err != nil && err != http.ErrServerClosed {
				slog.Error("server error", "error", err)
				os.Exit(1)
			}
		} else {
			slog.Info("starting HTTP server (dev mode)", "addr", addr)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				slog.Error("server error", "error", err)
				os.Exit(1)
			}
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	slog.Info("shutting down")
	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutCtx)
}

func loadECKeys(privKeyPath string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privPEM, err := os.ReadFile(privKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read private key: %w", err)
	}
	block, _ := pem.Decode(privPEM)
	if block == nil {
		return nil, nil, fmt.Errorf("no PEM block in private key file")
	}

	privKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8
		key, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err2 != nil {
			return nil, nil, fmt.Errorf("parse private key: %w (pkcs8: %w)", err, err2)
		}
		var ok bool
		privKey, ok = key.(*ecdsa.PrivateKey)
		if !ok {
			return nil, nil, fmt.Errorf("private key is not ECDSA")
		}
	}

	return privKey, &privKey.PublicKey, nil
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
