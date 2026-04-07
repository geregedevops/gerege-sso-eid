package handler

import (
	"net/http"
	"path/filepath"
)

func (h *Handler) SignResult(w http.ResponseWriter, r *http.Request) {
	sub := getSub(r)
	if sub == "" {
		h.jsonError(w, 401, "unauthorized")
		return
	}

	sessionID := r.PathValue("id")
	session, err := h.db.GetSession(r.Context(), sessionID)
	if err != nil || session == nil {
		h.jsonError(w, 404, "session not found")
		return
	}

	if session.RequesterSub != sub {
		h.jsonError(w, 403, "not your session")
		return
	}

	if session.Status != "COMPLETE" || session.SignedDocPath == "" {
		h.jsonError(w, 400, "document not yet signed")
		return
	}

	data, err := h.signer.GetSignedDocument(session.SignedDocPath)
	if err != nil {
		logErr("sign_result: read file", err)
		h.jsonError(w, 500, "failed to read signed document")
		return
	}

	filename := session.DocumentName
	if ext := filepath.Ext(filename); ext != "" {
		filename = filename[:len(filename)-len(ext)] + "_signed" + ext
	} else {
		filename += "_signed"
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.Write(data)
}
