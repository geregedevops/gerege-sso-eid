package middleware

import (
	"log/slog"
	"net/http"

	"verify.gerege.mn/api/internal/store"
)

func RateLimit(rdb *store.Redis) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rdb == nil {
				next.ServeHTTP(w, r)
				return
			}

			clientID, _ := r.Context().Value(ClientIDKey).(string)
			rateLimit, _ := r.Context().Value(ClientRateLimitKey).(int)
			if clientID == "" || rateLimit <= 0 {
				next.ServeHTTP(w, r)
				return
			}

			count, err := rdb.IncrRateLimit(r.Context(), clientID)
			if err != nil {
				slog.Error("rate limit redis error", "error", err)
				next.ServeHTTP(w, r)
				return
			}

			if count > int64(rateLimit) {
				w.Header().Set("Retry-After", "60")
				jsonError(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
