package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"xyp.gerege.mn/api/internal/provider"
)

// AuthenticateCitizen accepts reg_no + phone, looks up citizen from XYP,
// and returns limited verified info if found.
func (h *Handler) AuthenticateCitizen(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RegNo string `json:"reg_no"`
		Phone string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.RegNo == "" {
		writeError(w, http.StatusBadRequest, "reg_no is required")
		return
	}
	if req.Phone == "" {
		writeError(w, http.StatusBadRequest, "phone is required")
		return
	}

	info, err := h.cfg.Citizen.Lookup(r.Context(), req.RegNo)
	if err != nil {
		slog.Error("authenticate citizen lookup failed", "error", err, "reg_no", req.RegNo)
		writeError(w, http.StatusBadGateway, "upstream service error")
		return
	}
	if info == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"authenticated": false,
			"reason":        "citizen not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"authenticated": true,
		"citizen": map[string]any{
			"reg_no":     info.RegNo,
			"civil_id":   info.CivilID,
			"last_name":  info.LastName,
			"first_name": info.FirstName,
			"gender":     info.Gender,
			"birth_date": info.BirthDate,
			"image":      info.Image,
		},
	})
}

// AuthenticateOrg accepts reg_no + ceo_reg_no, looks up org from XYP.
// Validates: ceo_reg_no matches CEO OR largest shareholder (either one is valid).
// Returns org info + ultimate owner (largest shareholder) info.
func (h *Handler) AuthenticateOrg(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RegNo    string `json:"reg_no"`
		CEORegNo string `json:"ceo_reg_no"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.RegNo == "" {
		writeError(w, http.StatusBadRequest, "reg_no is required")
		return
	}
	if req.CEORegNo == "" {
		writeError(w, http.StatusBadRequest, "ceo_reg_no is required")
		return
	}

	info, err := h.cfg.Org.Lookup(r.Context(), req.RegNo)
	if err != nil {
		slog.Error("authenticate org lookup failed", "error", err, "reg_no", req.RegNo)
		writeError(w, http.StatusBadGateway, "upstream service error")
		return
	}
	if info == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"authenticated": false,
			"reason":        "organization not found",
		})
		return
	}

	inputRegNo := strings.ToLower(strings.TrimSpace(req.CEORegNo))

	// Check if matches CEO
	ceoMatch := strings.EqualFold(strings.TrimSpace(info.CEORegNo), inputRegNo)

	// Check if matches largest shareholder
	topOwner := findLargestFounder(info.Founders)
	ownerMatch := topOwner != nil && strings.EqualFold(strings.TrimSpace(topOwner.RegNo), inputRegNo)

	// Either one must match
	if !ceoMatch && !ownerMatch {
		writeJSON(w, http.StatusOK, map[string]any{
			"authenticated": false,
			"reason":        "ceo_reg_no does not match director or largest shareholder",
		})
		return
	}

	result := map[string]any{
		"authenticated": true,
		"organization": map[string]any{
			"reg_no":       info.RegNo,
			"name":         info.Name,
			"type":         info.Type,
			"ceo":          info.CEO,
			"ceo_reg_no":   info.CEORegNo,
			"ceo_position": info.CEOPosition,
		},
	}

	// Always include owner info if available
	if topOwner != nil {
		result["owner"] = map[string]any{
			"name":          topOwner.Name,
			"reg_no":        topOwner.RegNo,
			"type":          topOwner.Type,
			"share_percent": topOwner.SharePercent,
		}
	}

	writeJSON(w, http.StatusOK, result)
}

func findLargestFounder(founders []provider.OrgFounder) *provider.OrgFounder {
	if len(founders) == 0 {
		return nil
	}
	var top *provider.OrgFounder
	var topPct float64
	for i := range founders {
		pct, _ := strconv.ParseFloat(founders[i].SharePercent, 64)
		if top == nil || pct > topPct {
			top = &founders[i]
			topPct = pct
		}
	}
	return top
}
