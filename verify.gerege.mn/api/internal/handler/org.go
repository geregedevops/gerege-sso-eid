package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"xyp.gerege.mn/api/internal/provider"
)

func (h *Handler) OrgLookup(w http.ResponseWriter, r *http.Request) {
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

	info, err := h.cfg.Org.Lookup(r.Context(), req.RegNo)
	if err != nil {
		slog.Error("org lookup failed", "error", err, "reg_no", req.RegNo)
		writeError(w, http.StatusBadGateway, "upstream service error")
		return
	}

	if info == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"found":        false,
			"organization": nil,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"found":        true,
		"organization": info,
	})
}

func (h *Handler) OrgVerify(w http.ResponseWriter, r *http.Request) {
	var req provider.OrgVerifyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.RegNo == "" {
		writeError(w, http.StatusBadRequest, "reg_no is required")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	match, err := h.cfg.Org.Verify(r.Context(), req)
	if err != nil {
		slog.Error("org verify failed", "error", err, "reg_no", req.RegNo)
		writeError(w, http.StatusBadGateway, "upstream service error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"match":  match,
		"reg_no": req.RegNo,
	})
}
