package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api.gerege.mn/internal/handler"
	"api.gerege.mn/internal/middleware"
	"api.gerege.mn/internal/signer"
	"api.gerege.mn/internal/smartid"
	"api.gerege.mn/internal/store"
)

func main() {
	slog.Info("starting api.gerege.mn")

	port := env("PORT", "8080")
	jwksURI := env("SSO_JWKS_URI", "https://sso.gerege.mn/.well-known/jwks.json")
	eidURL := env("EID_API_URL", "https://ca.gerege.mn")
	storagePath := env("STORAGE_PATH", "/var/api/signed-docs")
	databaseURL := env("DATABASE_URL", "postgres://sso:pass@localhost:5432/gerege_sso_db")
	redisURL := env("REDIS_URL", "redis://localhost:6379/3")

	ctx := context.Background()

	db, err := store.NewPostgres(ctx, databaseURL)
	if err != nil {
		slog.Error("db connect failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	cache, err := store.NewRedis(redisURL)
	if err != nil {
		slog.Error("redis connect failed", "error", err)
		os.Exit(1)
	}
	defer cache.Close()

	sid := smartid.NewClient(eidURL)
	sig := signer.NewSigner(storagePath)

	h := handler.New(db, cache, sid, sig)

	jwtAuth := middleware.JWTAuth(jwksURI)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.Handle("POST /v1/sign/request", jwtAuth(http.HandlerFunc(h.SignRequest)))
	mux.Handle("GET /v1/sign/{id}/status", jwtAuth(http.HandlerFunc(h.SignStatus)))
	mux.Handle("GET /v1/sign/{id}/result", jwtAuth(http.HandlerFunc(h.SignResult)))
	mux.Handle("DELETE /v1/sign/{id}", jwtAuth(http.HandlerFunc(h.SignCancel)))
	mux.Handle("POST /v1/verify", jwtAuth(http.HandlerFunc(h.Verify)))

	var root http.Handler = mux
	root = handler.Logger(root)
	root = middleware.CORS(root)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      root,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		slog.Info("listening", "addr", ":"+port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	slog.Info("shutting down")
	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutCtx)
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
