package ocsp

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	ocspLib "golang.org/x/crypto/ocsp"
	"gesign.mn/gerege-sso/internal/store"
)

type Checker struct {
	ocspURL      string
	caIssuingURL string
	cache        *store.Redis
	httpClient   *http.Client
}

func NewChecker(ocspURL, caIssuingURL string, cache *store.Redis) *Checker {
	return &Checker{
		ocspURL:      ocspURL,
		caIssuingURL: caIssuingURL,
		cache:        cache,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Checker) Check(ctx context.Context, certSerialHex string) error {
	cacheKey := "ocsp:" + certSerialHex
	cached, err := c.cache.GetString(ctx, cacheKey)
	if err == nil && cached == "GOOD" {
		return nil
	}

	issuerCert, err := c.getIssuingCA(ctx)
	if err != nil {
		return fmt.Errorf("ocsp: failed to get issuing CA: %w", err)
	}

	serial := new(big.Int)
	serial.SetString(certSerialHex, 16)

	ocspReq, err := ocspLib.CreateRequest(&x509.Certificate{
		SerialNumber: serial,
		Issuer:       issuerCert.Subject,
	}, issuerCert, nil)
	if err != nil {
		return fmt.Errorf("ocsp: create request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.ocspURL, bytes.NewReader(ocspReq))
	if err != nil {
		return fmt.Errorf("ocsp: new request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/ocsp-request")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("ocsp: request failed (fail closed): %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ocsp: read response: %w", err)
	}

	ocspResp, err := ocspLib.ParseResponse(respBytes, issuerCert)
	if err != nil {
		return fmt.Errorf("ocsp: parse response: %w", err)
	}

	switch ocspResp.Status {
	case ocspLib.Good:
		c.cache.SetString(ctx, cacheKey, "GOOD", 4*time.Hour)
		return nil
	case ocspLib.Revoked:
		return fmt.Errorf("ocsp: certificate revoked")
	default:
		return fmt.Errorf("ocsp: certificate status unknown")
	}
}

func (c *Checker) getIssuingCA(ctx context.Context) (*x509.Certificate, error) {
	cacheKey := "ca:issuing"
	cached, err := c.cache.GetString(ctx, cacheKey)
	if err == nil && cached != "" {
		return parsePEMCert([]byte(cached))
	}

	req, err := http.NewRequestWithContext(ctx, "GET", c.caIssuingURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ocsp: fetch issuing CA: %w", err)
	}
	defer resp.Body.Close()

	pemBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	c.cache.SetString(ctx, cacheKey, string(pemBytes), 24*time.Hour)

	return parsePEMCert(pemBytes)
}

func parsePEMCert(pemData []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("ocsp: no PEM block found")
	}
	return x509.ParseCertificate(block.Bytes)
}
