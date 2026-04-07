package handler

import (
	"net/http"
	"time"
)

func (h *Handler) SignStatus(w http.ResponseWriter, r *http.Request) {
	sub := getSub(r)
	if sub == "" {
		h.jsonError(w, 401, "unauthorized")
		return
	}

	sessionID := r.PathValue("id")
	if sessionID == "" {
		h.jsonError(w, 400, "session_id required")
		return
	}

	session, err := h.db.GetSession(r.Context(), sessionID)
	if err != nil || session == nil {
		h.jsonError(w, 404, "session not found")
		return
	}

	if session.RequesterSub != sub {
		h.jsonError(w, 403, "not your session")
		return
	}

	if session.Status == "PENDING" || session.Status == "RUNNING" {
		if time.Now().After(session.ExpiresAt) {
			h.db.UpdateSessionStatus(r.Context(), sessionID, "EXPIRED")
			session.Status = "EXPIRED"
		} else if session.SmartIDSession != "" {
			status, err := h.smartid.Status(r.Context(), session.SmartIDSession)
			if err == nil {
				if status.State == "COMPLETE" && status.Result == "OK" {
					signedPath, signErr := h.signer.CreateSignedDocument(
						sessionID, session.DocumentPath,
						status.Name, status.CertSerial, time.Now(),
					)
					if signErr != nil {
						logErr("sign_status: create signed doc", signErr)
						h.db.UpdateSessionError(r.Context(), sessionID, signErr.Error())
						session.Status = "ERROR"
					} else {
						h.db.UpdateSessionComplete(r.Context(), sessionID, status.Sub, status.Name, status.CertSerial, signedPath)
						session.Status = "COMPLETE"
						session.SignerName = status.Name
						session.CertSerial = status.CertSerial
						session.SignedDocPath = signedPath
					}
				} else if status.State == "COMPLETE" && status.Result != "OK" {
					h.db.UpdateSessionError(r.Context(), sessionID, "user refused or timeout")
					session.Status = "ERROR"
					session.ErrorMessage = "user refused or timeout"
				} else {
					if session.Status == "PENDING" {
						h.db.UpdateSessionStatus(r.Context(), sessionID, "RUNNING")
						session.Status = "RUNNING"
					}
				}
			}
		}
	}

	resp := map[string]any{
		"session_id":        session.ID,
		"status":            session.Status,
		"verification_code": session.VerificationCode,
		"document_name":     session.DocumentName,
		"expires_at":        session.ExpiresAt.Format(time.RFC3339),
	}

	if session.Status == "COMPLETE" {
		resp["signer_name"] = session.SignerName
		resp["signed_at"] = session.UpdatedAt.Format(time.RFC3339)
		resp["result_url"] = "/v1/sign/" + session.ID + "/result"
	}
	if session.Status == "ERROR" {
		resp["error_message"] = session.ErrorMessage
	}

	h.jsonOK(w, resp)
}
