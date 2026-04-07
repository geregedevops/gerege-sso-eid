export default function GoGuidePage() {
  return (
    <main className="max-w-3xl mx-auto px-6 py-12 space-y-8">
      <h1 className="text-3xl font-bold text-white">Go + OIDC</h1>
      <p className="text-slate-400">Go application дээр sso.gerege.mn OIDC нэгтгэх.</p>

      <Code title="1. Dependencies">{`go get github.com/coreos/go-oidc/v3/oidc
go get golang.org/x/oauth2`}</Code>

      <Code title="2. OIDC Provider Setup">{`package main

import (
    "context"
    "os"

    "github.com/coreos/go-oidc/v3/oidc"
    "golang.org/x/oauth2"
)

func main() {
    ctx := context.Background()

    provider, err := oidc.NewProvider(ctx, "https://sso.gerege.mn")
    if err != nil {
        panic(err)
    }

    config := oauth2.Config{
        ClientID:     os.Getenv("EID_CLIENT_ID"),
        ClientSecret: os.Getenv("EID_CLIENT_SECRET"),
        Endpoint:     provider.Endpoint(),
        RedirectURL:  "https://myapp.mn/callback",
        Scopes:       []string{oidc.ScopeOpenID, "profile", "pos"},
    }

    state := "random-csrf-state"
    url := config.AuthCodeURL(state)
    // Redirect user to: url
}`}</Code>

      <Code title="3. Callback Handler">{`func callbackHandler(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")

    token, err := config.Exchange(ctx, code)
    if err != nil {
        http.Error(w, "token exchange failed", 500)
        return
    }

    rawIDToken, ok := token.Extra("id_token").(string)
    if !ok {
        http.Error(w, "no id_token", 500)
        return
    }

    verifier := provider.Verifier(&oidc.Config{ClientID: config.ClientID})
    idToken, err := verifier.Verify(ctx, rawIDToken)
    if err != nil {
        http.Error(w, "invalid id_token", 500)
        return
    }

    var claims struct {
        Sub        string   \`json:"sub"\`
        Name       string   \`json:"name"\`
        TenantID   string   \`json:"tenant_id"\`
        TenantRole string   \`json:"tenant_role"\`
        Plan       string   \`json:"plan"\`
        AMR        []string \`json:"amr"\`
    }
    idToken.Claims(&claims)
}`}</Code>
    </main>
  );
}

function Code({ title, children }: { title: string; children: string }) {
  return (
    <div>
      <h3 className="text-sm font-semibold text-white mb-2">{title}</h3>
      <pre className="bg-bg border border-white/10 rounded-xl p-4 text-xs text-slate-300 overflow-x-auto">{children}</pre>
    </div>
  );
}
