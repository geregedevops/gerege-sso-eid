package handler

import (
	"net/http"

	"gesign.mn/gerege-sso/internal/token"
)

func (h *Handler) JWKS(w http.ResponseWriter, r *http.Request) {
	jwkSet := token.BuildJWKSet(h.cfg.PubKey, h.cfg.KID)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	h.jsonOK(w, jwkSet)
}
