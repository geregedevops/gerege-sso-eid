package handler

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"api.gerege.mn/internal/model"
	"github.com/oklog/ulid/v2"
)

const maxDocSize = 10 * 1024 * 1024 // 10MB

type signRequest struct {
	SignerReg    string `json:"signer_reg"`
	DocumentName string `json:"document_name"`
	Document     string `json:"document"` // base64
}

func (h *Handler) SignRequest(w http.ResponseWriter, r *http.Request) {
	sub := getSub(r)
	if sub == "" {
		h.jsonError(w, 401, "unauthorized")
		return
	}

	var req signRequest
	if err := readJSON(r, &req); err != nil {
		h.jsonError(w, 400, "invalid request body")
		return
	}

	if req.DocumentName == "" || req.Document == "" {
		h.jsonError(w, 400, "document_name and document required")
		return
	}

	docData, err := base64.StdEncoding.DecodeString(req.Document)
	if err != nil {
		h.jsonError(w, 400, "invalid base64 document")
		return
	}

	if len(docData) > maxDocSize {
		h.jsonError(w, 400, "document too large (max 10MB)")
		return
	}

	sessionID := ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String()

	docPath, docHash, err := h.signer.StoreDocument(sessionID, docData, req.DocumentName)
	if err != nil {
		logErr("sign_request: store doc", err)
		h.jsonError(w, 500, "failed to store document")
		return
	}

	signerReg := req.SignerReg
	if signerReg == "" {
		h.jsonError(w, 400, "signer_reg required")
		return
	}

	displayText := fmt.Sprintf("Гарын үсэг: %s", req.DocumentName)
	if len(displayText) > 60 {
		displayText = displayText[:60]
	}

	smartidResp, err := h.smartid.Initiate(r.Context(), signerReg, displayText, "https://sso.gerege.mn/callback/eid")
	if err != nil {
		logErr("sign_request: smartid initiate", err)
		h.jsonError(w, 500, "SmartID initiate failed: "+err.Error())
		return
	}

	expiresAt := time.Now().Add(10 * time.Minute)
	session := &model.SigningSession{
		ID:               sessionID,
		RequesterSub:     sub,
		SignerReg:        signerReg,
		Status:           "PENDING",
		SmartIDSession:   smartidResp.SessionID,
		VerificationCode: smartidResp.VerificationCode,
		DocumentName:     req.DocumentName,
		DocumentHash:     docHash,
		DocumentSize:     len(docData),
		DocumentPath:     docPath,
		ExpiresAt:        expiresAt,
	}

	if err := h.db.CreateSession(r.Context(), session); err != nil {
		logErr("sign_request: db create", err)
		h.jsonError(w, 500, "internal error")
		return
	}

	w.WriteHeader(201)
	h.jsonOK(w, map[string]any{
		"session_id":        sessionID,
		"status":            "PENDING",
		"verification_code": smartidResp.VerificationCode,
		"expires_at":        expiresAt.Format(time.RFC3339),
	})
}

func readJSON(r *http.Request, dest any) error {
	return json.NewDecoder(io.LimitReader(r.Body, 15*1024*1024)).Decode(dest)
}

func generateNonce() *big.Int {
	n, _ := rand.Int(rand.Reader, big.NewInt(999999))
	return n
}
