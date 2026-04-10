package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"dan.gerege.mn/api/internal/dan"
)

func (h *Handler) Authorized(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	stateParam := r.URL.Query().Get("state")

	if code == "" {
		h.jsonError(w, 401, "sso.gov.mn-аас code ирсэнгүй")
		return
	}

	if stateParam == "" {
		h.jsonError(w, 400, "state параметр байхгүй")
		return
	}

	// Verify signed state
	state, err := verifyState(stateParam, h.cfg.StateSecret)
	if err != nil {
		slog.Error("authorized: invalid state", "error", err)
		h.jsonError(w, 400, "буруу эсвэл хуурамч state")
		return
	}

	// Standalone try mode — render citizen data directly
	if state["mode"] == "try" {
		h.HandleTryResult(w, r, code)
		return
	}

	cbURL := state["callback_url"]
	clientID := state["client_id"]

	if cbURL == "" || clientID == "" {
		h.jsonError(w, 400, "state-д callback_url эсвэл client_id байхгүй")
		return
	}

	// Re-validate client (may have been deactivated since verify)
	client, err := h.cfg.DB.GetDANClient(r.Context(), clientID)
	if err != nil || client == nil || !client.Active {
		h.jsonError(w, 400, "бүртгэлгүй эсвэл идэвхгүй client")
		return
	}

	// Re-validate callback URL
	if !matchCallbackURL(client.CallbackURLs, cbURL) {
		h.jsonError(w, 400, "callback_url бүртгэлгүй байна")
		return
	}

	// Exchange code for access token
	accessToken, err := dan.GetAccessToken(h.cfg.DAN, code)
	if err != nil {
		slog.Error("authorized: token exchange failed", "error", err)
		h.jsonError(w, 502, "access_token авахад алдаа гарлаа")
		return
	}

	// Fetch citizen data
	citizen, err := dan.GetCitizenData(h.cfg.DAN, accessToken)
	if err != nil {
		slog.Error("authorized: citizen data failed", "error", err)
		h.jsonError(w, 502, "иргэний мэдээлэл авахад алдаа гарлаа")
		return
	}

	slog.Info("authorized: success", "client_id", clientID)

	// Build POST payload
	postData := make(map[string]string)
	for k, v := range citizen {
		if v != "" {
			postData[k] = v
		}
	}
	postData["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())
	postData["client_id"] = clientID

	// HMAC signature using the client's hmac_key (not bcrypt hash)
	sig := dan.ComputeHMAC(postData, client.HMACKey)
	postData["signature"] = sig

	// POST full data (including image) to callback URL
	postJSON, err := json.Marshal(postData)
	if err != nil {
		slog.Error("authorized: marshal failed", "error", err)
		h.jsonError(w, 500, "internal error")
		return
	}

	postResp, err := http.Post(cbURL, "application/json", strings.NewReader(string(postJSON)))
	if err != nil {
		slog.Error("authorized: POST callback failed", "error", err)
		h.jsonError(w, 502, "callback URL руу POST хийхэд алдаа гарлаа")
		return
	}
	postResp.Body.Close()

	slog.Info("authorized: POSTed to callback", "client_id", clientID, "status", postResp.StatusCode)

	// Redirect browser to callback with ?status=ok
	redirectURL, err := url.Parse(cbURL)
	if err != nil {
		h.jsonError(w, 500, "invalid callback URL")
		return
	}

	params := redirectURL.Query()
	params.Set("status", "ok")
	params.Set("reg_no", citizen["reg_no"])
	params.Set("given_name", citizen["given_name"])
	params.Set("family_name", citizen["family_name"])
	params.Set("client_id", clientID)
	redirectURL.RawQuery = params.Encode()
	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}
