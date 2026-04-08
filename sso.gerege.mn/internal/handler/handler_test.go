package handler

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"sso.gerege.mn/internal/model"
	"sso.gerege.mn/internal/token"
	"golang.org/x/crypto/bcrypt"
)

// --- Mock DB ---

type mockDB struct {
	clients       map[string]*model.Client
	tenantMembers map[string]string // "tenantID:sub" -> role
	tenantPlans   map[string]string // tenantID -> plan
	recordedTokens []recordedToken
}

type recordedToken struct {
	ClientID  string
	Sub       string
	Scope     string
	ExpiresAt time.Time
}

func newMockDB() *mockDB {
	return &mockDB{
		clients:       make(map[string]*model.Client),
		tenantMembers: make(map[string]string),
		tenantPlans:   make(map[string]string),
	}
}

func (m *mockDB) GetClient(_ context.Context, clientID string) (*model.Client, error) {
	c, ok := m.clients[clientID]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (m *mockDB) GetTenantMember(_ context.Context, tenantID, sub string) (string, error) {
	return m.tenantMembers[tenantID+":"+sub], nil
}

func (m *mockDB) GetTenantPlan(_ context.Context, tenantID string) (string, error) {
	return m.tenantPlans[tenantID], nil
}

func (m *mockDB) RecordIssuedToken(_ context.Context, clientID, sub, scope string, expiresAt time.Time) error {
	m.recordedTokens = append(m.recordedTokens, recordedToken{clientID, sub, scope, expiresAt})
	return nil
}

// --- Mock Cache ---

type mockCache struct {
	data map[string][]byte
}

func newMockCache() *mockCache {
	return &mockCache{data: make(map[string][]byte)}
}

func (m *mockCache) Set(_ context.Context, key string, value any, _ time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	m.data[key] = b
	return nil
}

func (m *mockCache) Get(_ context.Context, key string, dest any) error {
	b, ok := m.data[key]
	if !ok {
		return fmt.Errorf("key not found: %s", key)
	}
	return json.Unmarshal(b, dest)
}

func (m *mockCache) Del(_ context.Context, key string) error {
	delete(m.data, key)
	return nil
}

func (m *mockCache) GetAndDel(_ context.Context, key string, dest any) error {
	b, ok := m.data[key]
	if !ok {
		return fmt.Errorf("key not found: %s", key)
	}
	delete(m.data, key)
	return json.Unmarshal(b, dest)
}

func (m *mockCache) Incr(_ context.Context, key string, _ time.Duration) (int64, error) {
	var count int64
	if b, ok := m.data[key]; ok {
		json.Unmarshal(b, &count)
	}
	count++
	b, _ := json.Marshal(count)
	m.data[key] = b
	return count, nil
}

func (m *mockCache) SetString(_ context.Context, key, value string, _ time.Duration) error {
	m.data[key] = []byte(value)
	return nil
}

func (m *mockCache) GetString(_ context.Context, key string) (string, error) {
	b, ok := m.data[key]
	if !ok {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return string(b), nil
}

// --- Test helpers ---

func testHandler(t *testing.T) (*Handler, *mockDB, *mockCache, *ecdsa.PrivateKey) {
	t.Helper()
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	kid := token.ComputeKID(&privKey.PublicKey)
	db := newMockDB()
	cache := newMockCache()
	issuer := token.NewIssuer(privKey, kid, "https://sso.gerege.mn")

	h := New(Config{
		Issuer:           "https://sso.gerege.mn",
		EIDBaseURL:       "https://e-id.mn",
		PrivKey:          privKey,
		PubKey:           &privKey.PublicKey,
		KID:              kid,
		DB:               db,
		Cache:            cache,
		TokenIssuer:      issuer,
		DANClientID:     "test-dan-client",
		DANClientSecret: "test-secret",
		DANScope:        "test-scope",
		DANCallbackURI:  "http://dan.gerege.mn/authorized",
		DANTokenURL:     "https://sso.gov.mn/oauth2/token",
		DANServiceURL:   "https://sso.gov.mn/oauth2/api/v1/service",
	})
	return h, db, cache, privKey
}

func addTestClient(db *mockDB, clientID, secret string, redirectURIs []string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.MinCost)
	db.clients[clientID] = &model.Client{
		ID:           clientID,
		SecretHash:   string(hash),
		Name:         "Test Client",
		RedirectURIs: redirectURIs,
		Scopes:       []string{"openid", "profile"},
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func parseJSON(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var result map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("parse JSON: %v, body: %s", err, rec.Body.String())
	}
	return result
}

// =============================================================
// Tests
// =============================================================

// --- Health ---

func TestHealth(t *testing.T) {
	h, _, _, _ := testHandler(t)
	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	h.Health(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := parseJSON(t, rec)
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %v", body["status"])
	}
	if body["issuer"] != "https://sso.gerege.mn" {
		t.Fatalf("expected issuer https://sso.gerege.mn, got %v", body["issuer"])
	}
}

// --- Discovery ---

func TestDiscovery(t *testing.T) {
	h, _, _, _ := testHandler(t)
	req := httptest.NewRequest("GET", "/.well-known/openid-configuration", nil)
	rec := httptest.NewRecorder()
	h.Discovery(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := parseJSON(t, rec)
	if body["issuer"] != "https://sso.gerege.mn" {
		t.Fatalf("wrong issuer: %v", body["issuer"])
	}
	if body["authorization_endpoint"] != "https://sso.gerege.mn/oauth/authorize" {
		t.Fatalf("wrong authorization_endpoint: %v", body["authorization_endpoint"])
	}
	if body["token_endpoint"] != "https://sso.gerege.mn/oauth/token" {
		t.Fatalf("wrong token_endpoint: %v", body["token_endpoint"])
	}
	if body["jwks_uri"] != "https://sso.gerege.mn/.well-known/jwks.json" {
		t.Fatalf("wrong jwks_uri: %v", body["jwks_uri"])
	}

	// Check scopes
	scopes, ok := body["scopes_supported"].([]any)
	if !ok || len(scopes) != 5 {
		t.Fatalf("expected 5 scopes, got %v", body["scopes_supported"])
	}

	// Check auth methods
	methods, ok := body["auth_methods_supported"].([]any)
	if !ok || len(methods) != 2 {
		t.Fatalf("expected 2 auth methods, got %v", body["auth_methods_supported"])
	}
}

// --- JWKS ---

func TestJWKS(t *testing.T) {
	h, _, _, _ := testHandler(t)
	req := httptest.NewRequest("GET", "/.well-known/jwks.json", nil)
	rec := httptest.NewRecorder()
	h.JWKS(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := parseJSON(t, rec)
	keys, ok := body["keys"].([]any)
	if !ok || len(keys) != 1 {
		t.Fatalf("expected 1 key, got %v", body["keys"])
	}
	key := keys[0].(map[string]any)
	if key["kty"] != "EC" {
		t.Fatalf("expected EC key type, got %v", key["kty"])
	}
	if key["crv"] != "P-256" {
		t.Fatalf("expected P-256, got %v", key["crv"])
	}
	if key["alg"] != "ES256" {
		t.Fatalf("expected ES256, got %v", key["alg"])
	}
	if key["use"] != "sig" {
		t.Fatalf("expected sig, got %v", key["use"])
	}
}

// --- Authorize ---

func TestAuthorize_UnknownClient(t *testing.T) {
	h, _, _, _ := testHandler(t)
	req := httptest.NewRequest("GET", "/oauth/authorize?client_id=unknown&redirect_uri=http://x&response_type=code&scope=openid", nil)
	rec := httptest.NewRecorder()
	h.Authorize(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	body := parseJSON(t, rec)
	if body["error"] != "invalid_request" {
		t.Fatalf("expected invalid_request, got %v", body["error"])
	}
}

func TestAuthorize_BadRedirectURI(t *testing.T) {
	h, db, _, _ := testHandler(t)
	addTestClient(db, "c1", "secret", []string{"http://app.test/callback"})

	req := httptest.NewRequest("GET", "/oauth/authorize?client_id=c1&redirect_uri=http://evil.com&response_type=code&scope=openid", nil)
	rec := httptest.NewRecorder()
	h.Authorize(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestAuthorize_UnsupportedResponseType(t *testing.T) {
	h, db, _, _ := testHandler(t)
	addTestClient(db, "c1", "secret", []string{"http://app.test/callback"})

	req := httptest.NewRequest("GET", "/oauth/authorize?client_id=c1&redirect_uri=http://app.test/callback&response_type=token&scope=openid", nil)
	rec := httptest.NewRecorder()
	h.Authorize(rec, req)

	// Should redirect with error
	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if !strings.Contains(loc, "unsupported_response_type") {
		t.Fatalf("expected unsupported_response_type in redirect, got %s", loc)
	}
}

func TestAuthorize_MissingOpenIDScope(t *testing.T) {
	h, db, _, _ := testHandler(t)
	addTestClient(db, "c1", "secret", []string{"http://app.test/callback"})

	req := httptest.NewRequest("GET", "/oauth/authorize?client_id=c1&redirect_uri=http://app.test/callback&response_type=code&scope=profile", nil)
	rec := httptest.NewRecorder()
	h.Authorize(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if !strings.Contains(loc, "invalid_scope") {
		t.Fatalf("expected invalid_scope, got %s", loc)
	}
}

func TestAuthorize_EID_RedirectsToEID(t *testing.T) {
	h, db, _, _ := testHandler(t)
	addTestClient(db, "c1", "secret", []string{"http://app.test/callback"})

	req := httptest.NewRequest("GET", "/oauth/authorize?client_id=c1&redirect_uri=http://app.test/callback&response_type=code&scope=openid%20profile", nil)
	rec := httptest.NewRecorder()
	h.Authorize(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if !strings.HasPrefix(loc, "https://e-id.mn/auth") {
		t.Fatalf("expected redirect to e-id.mn, got %s", loc)
	}
}

func TestAuthorize_DAN_RedirectsToSSOGovMN(t *testing.T) {
	h, db, _, _ := testHandler(t)
	addTestClient(db, "c1", "secret", []string{"http://app.test/callback"})

	req := httptest.NewRequest("GET", "/oauth/authorize?client_id=c1&redirect_uri=http://app.test/callback&response_type=code&scope=openid%20profile&auth_method=dan", nil)
	rec := httptest.NewRecorder()
	h.Authorize(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if !strings.HasPrefix(loc, "https://sso.gov.mn/login") {
		t.Fatalf("expected redirect to sso.gov.mn, got %s", loc)
	}
	if !strings.Contains(loc, "client_id=test-dan-client") {
		t.Fatalf("expected DAN client_id in redirect, got %s", loc)
	}
}

// --- EID Callback ---

func TestEIDCallback_MissingSession(t *testing.T) {
	h, _, _, _ := testHandler(t)
	req := httptest.NewRequest("GET", "/callback/eid", nil)
	rec := httptest.NewRecorder()
	h.EIDCallback(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestEIDCallback_ExpiredSession(t *testing.T) {
	h, _, _, _ := testHandler(t)
	req := httptest.NewRequest("GET", "/callback/eid?session=nonexistent&sub=test", nil)
	rec := httptest.NewRecorder()
	h.EIDCallback(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestEIDCallback_Success(t *testing.T) {
	h, _, cache, _ := testHandler(t)

	// Store a session
	session := model.AuthSession{
		ClientID:    "c1",
		RedirectURI: "http://app.test/callback",
		Scope:       "openid profile",
		State:       "mystate",
		Nonce:       "mynonce",
	}
	cache.Set(context.Background(), "sso:sess123", session, 10*time.Minute)

	req := httptest.NewRequest("GET", "/callback/eid?session=sess123&sub=user1&name=Test&given_name=Test&family_name=User&cert_serial=ABC123", nil)
	rec := httptest.NewRecorder()
	h.EIDCallback(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if !strings.HasPrefix(loc, "http://app.test/callback?code=") {
		t.Fatalf("expected redirect to app callback with code, got %s", loc)
	}
	if !strings.Contains(loc, "state=mystate") {
		t.Fatalf("expected state in redirect, got %s", loc)
	}

	// Session should be deleted
	var s model.AuthSession
	if err := cache.Get(context.Background(), "sso:sess123", &s); err == nil {
		t.Fatal("session should have been deleted")
	}
}

func TestEIDCallback_UserCancelled(t *testing.T) {
	h, _, cache, _ := testHandler(t)

	session := model.AuthSession{
		ClientID:    "c1",
		RedirectURI: "http://app.test/callback",
		State:       "mystate",
	}
	cache.Set(context.Background(), "sso:sess123", session, 10*time.Minute)

	// No sub param = user cancelled
	req := httptest.NewRequest("GET", "/callback/eid?session=sess123", nil)
	rec := httptest.NewRecorder()
	h.EIDCallback(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if !strings.Contains(loc, "error=access_denied") {
		t.Fatalf("expected access_denied error, got %s", loc)
	}
}

// --- Token ---

func TestToken_UnsupportedGrantType(t *testing.T) {
	h, _, _, _ := testHandler(t)
	form := url.Values{"grant_type": {"client_credentials"}, "client_id": {"c1"}, "client_secret": {"s"}}
	req := httptest.NewRequest("POST", "/oauth/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Token(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	body := parseJSON(t, rec)
	if body["error"] != "unsupported_grant_type" {
		t.Fatalf("expected unsupported_grant_type, got %v", body["error"])
	}
}

func TestToken_MissingClientCredentials(t *testing.T) {
	h, _, _, _ := testHandler(t)
	form := url.Values{"grant_type": {"authorization_code"}, "code": {"abc"}}
	req := httptest.NewRequest("POST", "/oauth/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Token(rec, req)

	if rec.Code != 401 {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestToken_InvalidClient(t *testing.T) {
	h, _, _, _ := testHandler(t)
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {"abc"},
		"client_id":     {"unknown"},
		"client_secret": {"secret"},
	}
	req := httptest.NewRequest("POST", "/oauth/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Token(rec, req)

	if rec.Code != 401 {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestToken_WrongSecret(t *testing.T) {
	h, db, _, _ := testHandler(t)
	addTestClient(db, "c1", "correct-secret", []string{"http://app.test/callback"})

	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {"abc"},
		"client_id":     {"c1"},
		"client_secret": {"wrong-secret"},
	}
	req := httptest.NewRequest("POST", "/oauth/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Token(rec, req)

	if rec.Code != 401 {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestToken_InvalidCode(t *testing.T) {
	h, db, _, _ := testHandler(t)
	addTestClient(db, "c1", "secret", []string{"http://app.test/callback"})

	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {"nonexistent"},
		"client_id":     {"c1"},
		"client_secret": {"secret"},
	}
	req := httptest.NewRequest("POST", "/oauth/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Token(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	body := parseJSON(t, rec)
	if body["error"] != "invalid_grant" {
		t.Fatalf("expected invalid_grant, got %v", body["error"])
	}
}

func TestToken_CodeBelongsToDifferentClient(t *testing.T) {
	h, db, cache, _ := testHandler(t)
	addTestClient(db, "c1", "secret", []string{"http://app.test/callback"})

	// Store code for different client
	codeData := model.AuthCode{
		Sub:         "user1",
		ClientID:    "other-client",
		RedirectURI: "http://app.test/callback",
		Scope:       "openid",
	}
	cache.Set(context.Background(), "code:testcode", codeData, 5*time.Minute)

	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {"testcode"},
		"client_id":     {"c1"},
		"client_secret": {"secret"},
	}
	req := httptest.NewRequest("POST", "/oauth/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Token(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	body := parseJSON(t, rec)
	if body["error"] != "invalid_grant" {
		t.Fatalf("expected invalid_grant, got %v", body["error"])
	}
}

func TestToken_RedirectURIMismatch(t *testing.T) {
	h, db, cache, _ := testHandler(t)
	addTestClient(db, "c1", "secret", []string{"http://app.test/callback"})

	codeData := model.AuthCode{
		Sub:         "user1",
		ClientID:    "c1",
		RedirectURI: "http://app.test/callback",
		Scope:       "openid",
	}
	cache.Set(context.Background(), "code:testcode", codeData, 5*time.Minute)

	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {"testcode"},
		"redirect_uri":  {"http://evil.com/callback"},
		"client_id":     {"c1"},
		"client_secret": {"secret"},
	}
	req := httptest.NewRequest("POST", "/oauth/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Token(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestToken_Success_EID(t *testing.T) {
	h, db, cache, privKey := testHandler(t)
	addTestClient(db, "c1", "secret", []string{"http://app.test/callback"})

	codeData := model.AuthCode{
		Sub:         "user1",
		Name:        "Test User",
		GivenName:   "Test",
		FamilyName:  "User",
		CertSerial:  "CERT123",
		ClientID:    "c1",
		RedirectURI: "http://app.test/callback",
		Scope:       "openid profile",
		Nonce:       "nonce123",
	}
	cache.Set(context.Background(), "code:validcode", codeData, 5*time.Minute)

	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {"validcode"},
		"redirect_uri":  {"http://app.test/callback"},
		"client_id":     {"c1"},
		"client_secret": {"secret"},
	}
	req := httptest.NewRequest("POST", "/oauth/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Token(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	body := parseJSON(t, rec)

	if body["token_type"] != "Bearer" {
		t.Fatalf("expected Bearer, got %v", body["token_type"])
	}
	if body["access_token"] == nil || body["access_token"] == "" {
		t.Fatal("expected access_token")
	}
	if body["id_token"] == nil || body["id_token"] == "" {
		t.Fatal("expected id_token")
	}
	if body["scope"] != "openid profile" {
		t.Fatalf("expected scope openid profile, got %v", body["scope"])
	}

	// Verify ID token
	idTokenStr := body["id_token"].(string)
	issuer := token.NewIssuer(privKey, "", "https://sso.gerege.mn")
	claims, err := issuer.VerifyIDToken(idTokenStr, &privKey.PublicKey)
	if err != nil {
		t.Fatalf("verify id_token: %v", err)
	}
	if claims.Subject != "user1" {
		t.Fatalf("expected sub user1, got %s", claims.Subject)
	}
	if claims.Name != "Test User" {
		t.Fatalf("expected name Test User, got %s", claims.Name)
	}
	if claims.Nonce != "nonce123" {
		t.Fatalf("expected nonce nonce123, got %s", claims.Nonce)
	}

	// Code should be consumed (single use)
	var c model.AuthCode
	if err := cache.Get(context.Background(), "code:validcode", &c); err == nil {
		t.Fatal("code should have been deleted after use")
	}

	// Token audit should be recorded
	if len(db.recordedTokens) != 1 {
		t.Fatalf("expected 1 recorded token, got %d", len(db.recordedTokens))
	}
}

func TestToken_Success_DAN(t *testing.T) {
	h, db, cache, _ := testHandler(t)
	addTestClient(db, "c1", "secret", []string{"http://app.test/callback"})

	codeData := model.AuthCode{
		Sub:         "AA12345678",
		Name:        "Бат",
		GivenName:   "Бат",
		FamilyName:  "",
		RegNo:       "AA12345678",
		ClientID:    "c1",
		RedirectURI: "http://app.test/callback",
		Scope:       "openid profile",
		Nonce:       "nonce456",
		Surname:     "Дорж",
		Gender:      "male",
		BirthDate:   "1990-01-01",
	}
	cache.Set(context.Background(), "code:dancode", codeData, 5*time.Minute)

	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {"dancode"},
		"client_id":     {"c1"},
		"client_secret": {"secret"},
	}
	req := httptest.NewRequest("POST", "/oauth/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Token(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	body := parseJSON(t, rec)
	if body["access_token"] == nil {
		t.Fatal("expected access_token")
	}
}

func TestToken_BasicAuth(t *testing.T) {
	h, db, cache, _ := testHandler(t)
	addTestClient(db, "c1", "secret", []string{"http://app.test/callback"})

	codeData := model.AuthCode{
		Sub:         "user1",
		ClientID:    "c1",
		RedirectURI: "http://app.test/callback",
		Scope:       "openid",
	}
	cache.Set(context.Background(), "code:basiccode", codeData, 5*time.Minute)

	form := url.Values{
		"grant_type": {"authorization_code"},
		"code":       {"basiccode"},
	}
	req := httptest.NewRequest("POST", "/oauth/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("c1", "secret")
	rec := httptest.NewRecorder()
	h.Token(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestToken_RateLimit(t *testing.T) {
	h, db, cache, _ := testHandler(t)
	addTestClient(db, "c1", "secret", []string{"http://app.test/callback"})

	// Simulate 10 previous requests
	for i := 0; i < 10; i++ {
		cache.Incr(context.Background(), "rl:token:c1", time.Minute)
	}

	codeData := model.AuthCode{
		Sub:      "user1",
		ClientID: "c1",
		Scope:    "openid",
	}
	cache.Set(context.Background(), "code:rlcode", codeData, 5*time.Minute)

	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {"rlcode"},
		"client_id":     {"c1"},
		"client_secret": {"secret"},
	}
	req := httptest.NewRequest("POST", "/oauth/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Token(rec, req)

	if rec.Code != 429 {
		t.Fatalf("expected 429 rate limit, got %d", rec.Code)
	}
}

// --- UserInfo ---

func TestUserInfo_MissingToken(t *testing.T) {
	h, _, _, _ := testHandler(t)
	req := httptest.NewRequest("GET", "/oauth/userinfo", nil)
	rec := httptest.NewRecorder()
	h.UserInfo(rec, req)

	if rec.Code != 401 {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	if rec.Header().Get("WWW-Authenticate") != "Bearer" {
		t.Fatal("expected WWW-Authenticate: Bearer")
	}
}

func TestUserInfo_InvalidToken(t *testing.T) {
	h, _, _, _ := testHandler(t)
	req := httptest.NewRequest("GET", "/oauth/userinfo", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	rec := httptest.NewRecorder()
	h.UserInfo(rec, req)

	if rec.Code != 401 {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestUserInfo_Success_ProfileScope(t *testing.T) {
	h, _, cache, _ := testHandler(t)

	atData := model.AccessTokenData{
		Sub:        "user1",
		ClientID:   "c1",
		Scope:      "openid profile",
		Name:       "Test User",
		GivenName:  "Test",
		FamilyName: "User",
		IssuedAt:   time.Now().Unix(),
		ExpiresAt:  time.Now().Add(time.Hour).Unix(),
	}
	cache.Set(context.Background(), "at:valid-at", atData, time.Hour)

	req := httptest.NewRequest("GET", "/oauth/userinfo", nil)
	req.Header.Set("Authorization", "Bearer valid-at")
	rec := httptest.NewRecorder()
	h.UserInfo(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := parseJSON(t, rec)
	if body["sub"] != "user1" {
		t.Fatalf("expected sub user1, got %v", body["sub"])
	}
	if body["name"] != "Test User" {
		t.Fatalf("expected name Test User, got %v", body["name"])
	}
	if body["locale"] != "mn-MN" {
		t.Fatalf("expected locale mn-MN, got %v", body["locale"])
	}
}

func TestUserInfo_Success_POSScope(t *testing.T) {
	h, _, cache, _ := testHandler(t)

	atData := model.AccessTokenData{
		Sub:        "user1",
		ClientID:   "c1",
		Scope:      "openid pos",
		TenantID:   "tenant1",
		TenantRole: "admin",
		Plan:       "pro",
		IssuedAt:   time.Now().Unix(),
		ExpiresAt:  time.Now().Add(time.Hour).Unix(),
	}
	cache.Set(context.Background(), "at:pos-at", atData, time.Hour)

	req := httptest.NewRequest("GET", "/oauth/userinfo", nil)
	req.Header.Set("Authorization", "Bearer pos-at")
	rec := httptest.NewRecorder()
	h.UserInfo(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := parseJSON(t, rec)
	if body["tenant_id"] != "tenant1" {
		t.Fatalf("expected tenant_id tenant1, got %v", body["tenant_id"])
	}
	if body["tenant_role"] != "admin" {
		t.Fatalf("expected tenant_role admin, got %v", body["tenant_role"])
	}
	if body["plan"] != "pro" {
		t.Fatalf("expected plan pro, got %v", body["plan"])
	}
}

// --- Introspect ---

func TestIntrospect_MissingCredentials(t *testing.T) {
	h, _, _, _ := testHandler(t)
	form := url.Values{"token": {"abc"}}
	req := httptest.NewRequest("POST", "/oauth/introspect", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Introspect(rec, req)

	if rec.Code != 401 {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestIntrospect_EmptyToken(t *testing.T) {
	h, db, _, _ := testHandler(t)
	addTestClient(db, "c1", "secret", []string{"http://app.test/callback"})

	form := url.Values{"client_id": {"c1"}, "client_secret": {"secret"}}
	req := httptest.NewRequest("POST", "/oauth/introspect", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Introspect(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := parseJSON(t, rec)
	if body["active"] != false {
		t.Fatalf("expected active false, got %v", body["active"])
	}
}

func TestIntrospect_InvalidToken(t *testing.T) {
	h, db, _, _ := testHandler(t)
	addTestClient(db, "c1", "secret", []string{"http://app.test/callback"})

	form := url.Values{"client_id": {"c1"}, "client_secret": {"secret"}, "token": {"badtoken"}}
	req := httptest.NewRequest("POST", "/oauth/introspect", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Introspect(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := parseJSON(t, rec)
	if body["active"] != false {
		t.Fatalf("expected active false for invalid token")
	}
}

func TestIntrospect_ValidToken(t *testing.T) {
	h, db, cache, _ := testHandler(t)
	addTestClient(db, "c1", "secret", []string{"http://app.test/callback"})

	now := time.Now()
	atData := model.AccessTokenData{
		Sub:        "user1",
		ClientID:   "c1",
		Scope:      "openid profile",
		Name:       "Test",
		TenantID:   "t1",
		TenantRole: "owner",
		Plan:       "enterprise",
		IssuedAt:   now.Unix(),
		ExpiresAt:  now.Add(time.Hour).Unix(),
	}
	cache.Set(context.Background(), "at:goodtoken", atData, time.Hour)

	form := url.Values{"client_id": {"c1"}, "client_secret": {"secret"}, "token": {"goodtoken"}}
	req := httptest.NewRequest("POST", "/oauth/introspect", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Introspect(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := parseJSON(t, rec)
	if body["active"] != true {
		t.Fatalf("expected active true, got %v", body["active"])
	}
	if body["sub"] != "user1" {
		t.Fatalf("expected sub user1, got %v", body["sub"])
	}
	if body["tenant_id"] != "t1" {
		t.Fatalf("expected tenant_id t1, got %v", body["tenant_id"])
	}
}

// --- Revoke ---

func TestRevoke_MissingCredentials(t *testing.T) {
	h, _, _, _ := testHandler(t)
	form := url.Values{"token": {"abc"}}
	req := httptest.NewRequest("POST", "/oauth/revoke", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Revoke(rec, req)

	if rec.Code != 401 {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestRevoke_UnknownClient_Returns200(t *testing.T) {
	h, _, _, _ := testHandler(t)
	form := url.Values{"client_id": {"unknown"}, "client_secret": {"s"}, "token": {"abc"}}
	req := httptest.NewRequest("POST", "/oauth/revoke", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Revoke(rec, req)

	// RFC 7009: always 200
	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestRevoke_Success(t *testing.T) {
	h, db, cache, _ := testHandler(t)
	addTestClient(db, "c1", "secret", []string{"http://app.test/callback"})

	// Store a token
	atData := model.AccessTokenData{Sub: "user1", ClientID: "c1", Scope: "openid"}
	cache.Set(context.Background(), "at:revokeme", atData, time.Hour)

	form := url.Values{"client_id": {"c1"}, "client_secret": {"secret"}, "token": {"revokeme"}}
	req := httptest.NewRequest("POST", "/oauth/revoke", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Revoke(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	// Token should be deleted
	var at model.AccessTokenData
	if err := cache.Get(context.Background(), "at:revokeme", &at); err == nil {
		t.Fatal("token should have been deleted")
	}
}

// --- DAN Callback ---

func TestDANCallback_MissingSession(t *testing.T) {
	h, _, _, _ := testHandler(t)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /callback/dan/{session}", h.DANCallback)

	req := httptest.NewRequest("GET", "/callback/dan/nonexistent?reg_no=AA12345678", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestDANCallback_Success(t *testing.T) {
	h, _, cache, _ := testHandler(t)

	session := model.AuthSession{
		ClientID:    "c1",
		RedirectURI: "http://app.test/callback",
		Scope:       "openid profile",
		State:       "danstate",
		Nonce:       "dannonce",
		AuthMethod:  "dan",
	}
	cache.Set(context.Background(), "sso:dansess", session, 10*time.Minute)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /callback/dan/{session}", h.DANCallback)

	req := httptest.NewRequest("GET", "/callback/dan/dansess?reg_no=AA12345678&surname=Дорж&given_name=Бат&gender=male&birth_date=1990-01-01", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d: %s", rec.Code, rec.Body.String())
	}
	loc := rec.Header().Get("Location")
	if !strings.HasPrefix(loc, "http://app.test/callback?code=") {
		t.Fatalf("expected redirect with code, got %s", loc)
	}
	if !strings.Contains(loc, "state=danstate") {
		t.Fatalf("expected state, got %s", loc)
	}
}

func TestDANCallback_NoRegNo(t *testing.T) {
	h, _, cache, _ := testHandler(t)

	session := model.AuthSession{
		ClientID:    "c1",
		RedirectURI: "http://app.test/callback",
		State:       "st",
	}
	cache.Set(context.Background(), "sso:dansess2", session, 10*time.Minute)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /callback/dan/{session}", h.DANCallback)

	req := httptest.NewRequest("GET", "/callback/dan/dansess2", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if !strings.Contains(loc, "error=access_denied") {
		t.Fatalf("expected access_denied, got %s", loc)
	}
}

// --- DAN Gateway ---

func TestDANGatewayAuthorized_MissingCode(t *testing.T) {
	h, _, _, _ := testHandler(t)
	req := httptest.NewRequest("GET", "/authorized?state=abc", nil)
	rec := httptest.NewRecorder()
	h.DANGatewayAuthorized(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestDANGatewayAuthorized_MissingState(t *testing.T) {
	h, _, _, _ := testHandler(t)
	req := httptest.NewRequest("GET", "/authorized?code=abc", nil)
	rec := httptest.NewRecorder()
	h.DANGatewayAuthorized(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestDANGatewayAuthorized_InvalidState(t *testing.T) {
	h, _, _, _ := testHandler(t)
	req := httptest.NewRequest("GET", "/authorized?code=abc&state=not-base64!!!", nil)
	rec := httptest.NewRecorder()
	h.DANGatewayAuthorized(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

// --- Index ---

func TestIndex(t *testing.T) {
	h, _, _, _ := testHandler(t)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	h.Index(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Fatalf("expected text/html content type, got %s", ct)
	}
	if !strings.Contains(rec.Body.String(), "sso.gerege.mn") {
		t.Fatal("expected sso.gerege.mn in HTML body")
	}
}

func TestIndex_NotFound(t *testing.T) {
	h, _, _, _ := testHandler(t)
	req := httptest.NewRequest("GET", "/nonexistent", nil)
	rec := httptest.NewRecorder()
	h.Index(rec, req)

	if rec.Code != 404 {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

// --- Favicon ---

func TestFavicon(t *testing.T) {
	h, _, _, _ := testHandler(t)
	req := httptest.NewRequest("GET", "/favicon.ico", nil)
	rec := httptest.NewRecorder()
	h.Favicon(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "image/svg+xml" {
		t.Fatalf("expected image/svg+xml, got %s", ct)
	}
}

// --- Token with tenant context ---

func TestToken_WithTenantContext(t *testing.T) {
	h, db, cache, _ := testHandler(t)
	addTestClient(db, "c1", "secret", []string{"http://app.test/callback"})
	db.clients["c1"].TenantID = "tenant1"
	db.clients["c1"].Scopes = []string{"openid", "profile", "pos"}
	db.tenantMembers["tenant1:user1"] = "admin"
	db.tenantPlans["tenant1"] = "enterprise"

	codeData := model.AuthCode{
		Sub:         "user1",
		Name:        "User",
		ClientID:    "c1",
		RedirectURI: "http://app.test/callback",
		Scope:       "openid profile pos",
		Nonce:       "n1",
	}
	cache.Set(context.Background(), "code:tenantcode", codeData, 5*time.Minute)

	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {"tenantcode"},
		"client_id":     {"c1"},
		"client_secret": {"secret"},
	}
	req := httptest.NewRequest("POST", "/oauth/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Token(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// Verify the access token has tenant data
	body := parseJSON(t, rec)
	accessToken := body["access_token"].(string)

	var atData model.AccessTokenData
	if err := cache.Get(context.Background(), "at:"+accessToken, &atData); err != nil {
		t.Fatalf("get access token data: %v", err)
	}
	if atData.TenantID != "tenant1" {
		t.Fatalf("expected tenant_id tenant1, got %s", atData.TenantID)
	}
	if atData.TenantRole != "admin" {
		t.Fatalf("expected tenant_role admin, got %s", atData.TenantRole)
	}
	if atData.Plan != "enterprise" {
		t.Fatalf("expected plan enterprise, got %s", atData.Plan)
	}
}
