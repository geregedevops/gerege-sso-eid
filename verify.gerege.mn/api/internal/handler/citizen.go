package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"xyp.gerege.mn/api/internal/provider"
)

func (h *Handler) CitizenLookup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RegNo string `json:"reg_no"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.RegNo == "" {
		writeError(w, http.StatusBadRequest, "reg_no is required")
		return
	}

	info, err := h.cfg.Citizen.Lookup(r.Context(), req.RegNo)
	if err != nil {
		slog.Error("citizen lookup failed", "error", err, "reg_no", req.RegNo)
		writeError(w, http.StatusBadGateway, "upstream service error")
		return
	}

	if info == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"found":   false,
			"citizen": nil,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"found":   true,
		"citizen": info,
	})
}

func (h *Handler) CitizenVerify(w http.ResponseWriter, r *http.Request) {
	var req provider.CitizenVerifyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.RegNo == "" {
		writeError(w, http.StatusBadRequest, "reg_no is required")
		return
	}
	if req.FirstName == "" && req.LastName == "" {
		writeError(w, http.StatusBadRequest, "first_name or last_name is required")
		return
	}

	match, err := h.cfg.Citizen.Verify(r.Context(), req)
	if err != nil {
		slog.Error("citizen verify failed", "error", err, "reg_no", req.RegNo)
		writeError(w, http.StatusBadGateway, "upstream service error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"match":  match,
		"reg_no": req.RegNo,
	})
}
