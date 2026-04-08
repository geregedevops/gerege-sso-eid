package handler

import (
	"log/slog"
	"net/http"
	"time"

	"sso.gerege.mn/internal/model"
)

// DANCallback handles the redirect from dan.gerege.mn gateway.
// dan.gerege.mn exchanges the sso.gov.mn code for citizen data,
// then redirects here with all citizen fields as query params:
//
//	/callback/dan?reg_no=XX&surname=...&given_name=...&family_name=...&state=...
func (h *Handler) DANCallback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	regNo := q.Get("reg_no")
	sessionID := r.PathValue("session")

	slog.Info("dan_callback", "reg_no", regNo, "session", sessionID, "params", r.URL.RawQuery)

	if sessionID == "" {
		h.jsonError(w, 400, "invalid_request", "missing session")
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
