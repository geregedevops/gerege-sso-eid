package signer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Signer struct {
	storagePath string
}

func NewSigner(storagePath string) *Signer {
	os.MkdirAll(storagePath, 0750)
	return &Signer{storagePath: storagePath}
}

func (s *Signer) StoreDocument(sessionID string, data []byte, filename string) (string, string, error) {
	dir := filepath.Join(s.storagePath, sessionID)
	os.MkdirAll(dir, 0750)

	path := filepath.Join(dir, "original_"+filename)
	if err := os.WriteFile(path, data, 0640); err != nil {
		return "", "", fmt.Errorf("signer.StoreDocument: %w", err)
	}

	hash := sha256.Sum256(data)
	return path, hex.EncodeToString(hash[:]), nil
}

func (s *Signer) CreateSignedDocument(sessionID, origPath, signerName, certSerial string, signedAt time.Time) (string, error) {
	origData, err := os.ReadFile(origPath)
	if err != nil {
		return "", fmt.Errorf("signer.CreateSignedDocument: %w", err)
	}

	dir := filepath.Join(s.storagePath, sessionID)
	signedPath := filepath.Join(dir, "signed_"+filepath.Base(origPath))

	if err := os.WriteFile(signedPath, origData, 0640); err != nil {
		return "", fmt.Errorf("signer.CreateSignedDocument: %w", err)
	}

	metaPath := filepath.Join(dir, "signature.json")
	meta := fmt.Sprintf(`{"signer":"%s","cert_serial":"%s","signed_at":"%s","session_id":"%s"}`,
		signerName, certSerial, signedAt.Format(time.RFC3339), sessionID)
	os.WriteFile(metaPath, []byte(meta), 0640)

	return signedPath, nil
}

func (s *Signer) GetSignedDocument(path string) ([]byte, error) {
	return os.ReadFile(path)
}
