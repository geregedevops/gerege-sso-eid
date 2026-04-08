package token

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

func generateTestKey(t *testing.T) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	t.Helper()
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return privKey, &privKey.PublicKey
}

func TestIssueAndVerifyIDToken(t *testing.T) {
	privKey, pubKey := generateTestKey(t)
	kid := ComputeKID(pubKey)
	issuer := NewIssuer(privKey, kid, "https://sso.gerege.mn")

	tokenStr, err := issuer.IssueIDToken(
		"user1", "client1", "nonce123",
		"Test User", "Test", "User",
		"CERT123", "",
		"", "", "",
	)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("empty token")
	}

	claims, err := issuer.VerifyIDToken(tokenStr, pubKey)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}

	if claims.Subject != "user1" {
		t.Fatalf("expected sub user1, got %s", claims.Subject)
	}
	if claims.Issuer != "https://sso.gerege.mn" {
		t.Fatalf("expected issuer https://sso.gerege.mn, got %s", claims.Issuer)
	}
	if claims.Name != "Test User" {
		t.Fatalf("expected name Test User, got %s", claims.Name)
	}
	if claims.GivenName != "Test" {
		t.Fatalf("expected given_name Test, got %s", claims.GivenName)
	}
	if claims.Nonce != "nonce123" {
		t.Fatalf("expected nonce nonce123, got %s", claims.Nonce)
	}
	if claims.Locale != "mn-MN" {
		t.Fatalf("expected locale mn-MN, got %s", claims.Locale)
	}
	if claims.CertSerial != "CERT123" {
		t.Fatalf("expected cert_serial CERT123, got %s", claims.CertSerial)
	}
	// EID auth → smartid AMR
	if len(claims.AMR) != 3 || claims.AMR[0] != "smartid" {
		t.Fatalf("expected AMR [smartid pin1 x509], got %v", claims.AMR)
	}
}

func TestIssueIDToken_DAN_AMR(t *testing.T) {
	privKey, pubKey := generateTestKey(t)
	kid := ComputeKID(pubKey)
	issuer := NewIssuer(privKey, kid, "https://sso.gerege.mn")

	tokenStr, err := issuer.IssueIDToken(
		"AA12345678", "client1", "",
		"Бат", "Бат", "",
		"", "AA12345678",
		"tenant1", "admin", "pro",
	)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}

	claims, err := issuer.VerifyIDToken(tokenStr, pubKey)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}

	if claims.RegNo != "AA12345678" {
		t.Fatalf("expected reg_no AA12345678, got %s", claims.RegNo)
	}
	// DAN auth → dan AMR
	if len(claims.AMR) != 2 || claims.AMR[0] != "dan" {
		t.Fatalf("expected AMR [dan sso_gov_mn], got %v", claims.AMR)
	}
	if claims.TenantID != "tenant1" {
		t.Fatalf("expected tenant_id tenant1, got %s", claims.TenantID)
	}
	if claims.TenantRole != "admin" {
		t.Fatalf("expected tenant_role admin, got %s", claims.TenantRole)
	}
	if claims.Plan != "pro" {
		t.Fatalf("expected plan pro, got %s", claims.Plan)
	}
}

func TestVerifyIDToken_WrongKey(t *testing.T) {
	privKey, _ := generateTestKey(t)
	_, otherPub := generateTestKey(t)
	kid := ComputeKID(&privKey.PublicKey)
	issuer := NewIssuer(privKey, kid, "https://sso.gerege.mn")

	tokenStr, _ := issuer.IssueIDToken("user1", "c1", "", "N", "N", "", "", "", "", "", "")

	_, err := issuer.VerifyIDToken(tokenStr, otherPub)
	if err == nil {
		t.Fatal("expected verification failure with wrong key")
	}
}

func TestVerifyIDToken_InvalidString(t *testing.T) {
	privKey, pubKey := generateTestKey(t)
	kid := ComputeKID(pubKey)
	issuer := NewIssuer(privKey, kid, "https://sso.gerege.mn")

	_, err := issuer.VerifyIDToken("not.a.jwt", pubKey)
	if err == nil {
		t.Fatal("expected error for invalid token string")
	}
}

func TestBuildJWKSet(t *testing.T) {
	_, pubKey := generateTestKey(t)
	kid := ComputeKID(pubKey)
	jwkSet := BuildJWKSet(pubKey, kid)

	if len(jwkSet.Keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(jwkSet.Keys))
	}
	k := jwkSet.Keys[0]
	if k.KTY != "EC" {
		t.Fatalf("expected EC, got %s", k.KTY)
	}
	if k.CRV != "P-256" {
		t.Fatalf("expected P-256, got %s", k.CRV)
	}
	if k.ALG != "ES256" {
		t.Fatalf("expected ES256, got %s", k.ALG)
	}
	if k.Use != "sig" {
		t.Fatalf("expected sig, got %s", k.Use)
	}
	if k.KID != kid {
		t.Fatalf("expected kid %s, got %s", kid, k.KID)
	}
	if k.X == "" || k.Y == "" {
		t.Fatal("expected X and Y to be non-empty")
	}
}

func TestComputeKID_Deterministic(t *testing.T) {
	_, pubKey := generateTestKey(t)
	kid1 := ComputeKID(pubKey)
	kid2 := ComputeKID(pubKey)
	if kid1 != kid2 {
		t.Fatalf("KID should be deterministic, got %s and %s", kid1, kid2)
	}
	if kid1 == "" {
		t.Fatal("KID should not be empty")
	}
}

func TestComputeKID_DifferentKeys(t *testing.T) {
	_, pub1 := generateTestKey(t)
	_, pub2 := generateTestKey(t)
	kid1 := ComputeKID(pub1)
	kid2 := ComputeKID(pub2)
	if kid1 == kid2 {
		t.Fatal("different keys should produce different KIDs")
	}
}
