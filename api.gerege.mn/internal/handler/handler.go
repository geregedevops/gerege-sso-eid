package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"gesign.mn/gerege-api/internal/middleware"
	"gesign.mn/gerege-api/internal/signer"
	"gesign.mn/gerege-api/internal/smartid"
	"gesign.mn/gerege-api/internal/store"
)

type Handler struct {
	db      *store.Postgres
	cache   *store.Redis
	smartid *smartid.Client
	signer  *signer.Signer
}

func New(db *store.Postgres, cache *store.Redis, sid *smartid.Client, sig *signer.Signer) *Handler {
	return &Handler{db: db, cache: cache, smartid: sid, signer: sig}
}

func (h *Handler) jsonOK(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) jsonError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	h.jsonOK(w, map[string]string{"status": "ok", "service": "api.gerege.mn"})
}

func getSub(r *http.Request) string {
	sub, _ := r.Context().Value(middleware.SubKey).(string)
	return sub
}

func logErr(msg string, err error) {
	slog.Error(msg, "error", err)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request", "method", r.Method, "path", r.URL.Path, "latency_ms", time.Since(start).Milliseconds())
	})
}
