package handler

import (
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"gesign.mn/gerege-sso/internal/model"
)

// danState is the decoded state parameter from dan.gerege.mn gateway
type danState struct {
	RedirectURL string `json:"redirect_url"`
	Session     string `json:"session"`
}

// DANCallback handles the redirect from dan.gerege.mn gateway.
// dan.gerege.mn exchanges the sso.gov.mn code for citizen data,
// then redirects here with all citizen fields as query params:
//
//	/callback/dan?reg_no=XX&surname=...&given_name=...&family_name=...&state=...
func (h *Handler) DANCallback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	stateB64 := q.Get("state")
	regNo := q.Get("reg_no")

	slog.Info("dan_callback", "reg_no", regNo, "has_state", stateB64 != "", "params", r.URL.RawQuery)

	if stateB64 == "" {
		h.jsonError(w, 400, "invalid_request", "missing state")
		return
	}

	// Decode state — try RawURL first, then StdEncoding
	stateBytes, err := base64.RawURLEncoding.DecodeString(stateB64)
	if err != nil {
		stateBytes, err = base64.StdEncoding.DecodeString(stateB64)
		if err != nil {
			// Try RawStdEncoding (no padding, standard alphabet)
			stateBytes, err = base64.RawStdEncoding.DecodeString(stateB64)
			if err != nil {
				h.jsonError(w, 400, "invalid_request", "invalid state encoding")
				return
			}
		}
	}

	var state danState
	if err := json.Unmarshal(stateBytes, &state); err != nil {
		h.jsonError(w, 400, "invalid_request", "invalid state format")
		return
	}

	sessionID := state.Session
	if sessionID == "" {
		h.jsonError(w, 400, "invalid_request", "missing session in state")
		return
	}

	// Get session from Redis
	var session model.AuthSession
	if err := h.cfg.Cache.Get(r.Context(), "sso:"+sessionID, &session); err != nil {
		logErr("dan_callback: session not found", err)
		h.jsonError(w, 400, "invalid_request", "session expired or not found")
		return
	}

	// Handle error (user cancelled or no reg_no)
	if regNo == "" {
		h.cfg.Cache.Del(r.Context(), "sso:"+sessionID)
		redirectWithError(w, r, session.RedirectURI, session.State, "access_denied", "DAN verification failed or cancelled")
		return
	}

	// Collect all citizen data from query params
	givenName := q.Get("given_name")
	familyName := q.Get("family_name")
	surname := q.Get("surname")

	// Build display name
	name := givenName
	if familyName != "" {
		name = givenName + " " + familyName
	}
	if name == "" && surname != "" {
		name = surname
	}

	// Generate auth code
	authCode := generateRandomString(32)

	// Store code in Redis (5 min, single use)
	codeData := model.AuthCode{
		Sub:           regNo,
		Name:          name,
		GivenName:     givenName,
		FamilyName:    familyName,
		RegNo:         regNo,
		ClientID:      session.ClientID,
		RedirectURI:   session.RedirectURI,
		Scope:         session.Scope,
		Nonce:         session.Nonce,
		Surname:       surname,
		CivilID:       q.Get("civil_id"),
		Gender:        q.Get("gender"),
		BirthDate:     q.Get("birth_date"),
		Nationality:   q.Get("nationality"),
		PhoneNo:       q.Get("phone_no"),
		Email:         q.Get("email"),
		AimagName:     q.Get("aimag_name"),
		SumName:       q.Get("sum_name"),
		BagName:       q.Get("bag_name"),
		AddressDetail: q.Get("address_detail"),
	}

	slog.Info("dan_callback: citizen data",
		"reg_no", regNo,
		"name", name,
		"civil_id", codeData.CivilID,
		"gender", codeData.Gender,
		"client_id", session.ClientID,
	)

	if err := h.cfg.Cache.Set(r.Context(), "code:"+authCode, codeData, 5*time.Minute); err != nil {
		logErr("dan_callback: redis set code", err)
		h.jsonError(w, 500, "server_error", "internal error")
		return
	}

	// Delete session
	h.cfg.Cache.Del(r.Context(), "sso:"+sessionID)

	// Redirect to client with code + state
	redirectWithCode(w, r, session.RedirectURI, authCode, session.State)
}
