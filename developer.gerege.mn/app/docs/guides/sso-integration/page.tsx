export default function SSOIntegrationGuidePage() {
  return (
    <main className="max-w-3xl mx-auto px-6 py-12 space-y-10">
      <div>
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-green-500/10 border border-green-500/20 text-green-400 text-xs font-semibold mb-4">
          Нээлттэй — Бүх platform-д зориулсан
        </div>
        <h1 className="text-3xl font-bold text-white mb-2">e-ID SSO нэгтгэх заавар</h1>
        <p className="text-slate-400">
          Таны platform дээр sso.gerege.mn-р дамжуулан e-ID Mongolia нэвтрэлт нэмнэ.
          OpenID Connect (OIDC) стандартаар ажиллана. Аливаа 3-р талын систем чөлөөтэй холбогдох боломжтой.
        </p>
      </div>

      <Section title="Тойм">
        <div className="p-4 bg-green-500/5 border border-green-500/15 rounded-xl text-sm mb-4">
          <p className="font-semibold text-green-400">Бүх 3-р талын platform-д нээлттэй</p>
          <p className="text-slate-400 mt-1">sso.gerege.mn нь аливаа систем, platform, апп-д нээлттэй OIDC provider.
          developer.gerege.mn дээр app бүртгүүлж client_id авахад хангалттай.
          Хэлний болон framework-ийн хязгаарлалтгүй — стандарт OIDC дэмжсэн дурын технологи ашиглаж болно.</p>
        </div>
        <p>sso.gerege.mn нь OpenID Connect 1.0 provider бөгөөд e-ID Mongolia
        смарт картаар баталгаажуулсан иргэний мэдээллийг OAuth 2.0 Authorization Code Flow-р
        3-р талын системд дамжуулна.</p>
        <div className="mt-4 p-4 bg-primary/5 border border-primary/15 rounded-xl text-sm space-y-2">
          <p className="font-semibold text-white">Flow:</p>
          <ol className="list-decimal list-inside text-slate-400 space-y-1">
            <li>Хэрэглэгч таны сайт дээр &quot;e-ID-р нэвтрэх&quot; товч дарна</li>
            <li>sso.gerege.mn руу redirect хийгдэнэ</li>
            <li>e-ID Mongolia смарт картаар нэвтэрнэ</li>
            <li>Таны callback URL руу authorization code-тэй буцна</li>
            <li>Сервер талаас code → access_token + id_token солилцоно</li>
            <li>ID token-оос иргэний мэдээлэл уншина</li>
          </ol>
        </div>
      </Section>

      <Section title="Алхам 1: App бүртгүүлэх">
        <p>
          <a href="/dashboard/apps/new" className="text-primary hover:underline font-semibold">developer.gerege.mn/dashboard/apps/new</a> хуудаснаас шинэ app бүртгүүлнэ.
        </p>
        <div className="mt-3 space-y-3">
          <Field label="App нэр" desc="Таны системийн нэр (хэрэглэгчдэд харагдана)" />
          <Field label="Redirect URI" desc="Нэвтрэлт амжилттай болсны дараа буцах URL. Жишээ: https://myapp.mn/api/auth/callback/gerege-sso" />
          <Field label="Scopes" desc="openid (заавал) + profile, pos, social, payment" />
        </div>
        <div className="mt-4 p-4 bg-amber-500/5 border border-amber-500/15 rounded-xl text-sm">
          <p className="text-amber-400 font-semibold">Анхааруулга</p>
          <p className="text-slate-400 mt-1">
            <code className="text-amber-300">client_secret</code> зөвхөн нэг удаа харагдана.
            Хадгалаагүй бол шинэ app үүсгэх шаардлагатай.
          </p>
        </div>
      </Section>

      <Section title="Алхам 2: OIDC тохиргоо">
        <p>sso.gerege.mn нь стандарт OIDC Discovery endpoint дэмждэг:</p>
        <Code>{`https://sso.gerege.mn/.well-known/openid-configuration`}</Code>
        <div className="mt-4">
          <h4 className="text-sm font-semibold text-white mb-2">Endpoint-ууд:</h4>
          <table className="w-full text-sm">
            <tbody>
              <Row label="Authorization" value="https://sso.gerege.mn/oauth/authorize" />
              <Row label="Token" value="https://sso.gerege.mn/oauth/token" />
              <Row label="UserInfo" value="https://sso.gerege.mn/oauth/userinfo" />
              <Row label="JWKS" value="https://sso.gerege.mn/.well-known/jwks.json" />
              <Row label="Introspect" value="https://sso.gerege.mn/oauth/introspect" />
              <Row label="Revoke" value="https://sso.gerege.mn/oauth/revoke" />
            </tbody>
          </table>
        </div>
      </Section>

      <Section title="Алхам 3: Authorization Request">
        <p>Хэрэглэгчийг authorize endpoint руу redirect хийнэ:</p>
        <Code>{`GET https://sso.gerege.mn/oauth/authorize
  ?response_type=code
  &client_id=YOUR_CLIENT_ID
  &redirect_uri=https://myapp.mn/callback
  &scope=openid profile
  &state=random-csrf-token
  &nonce=random-nonce`}</Code>
        <div className="mt-3">
          <h4 className="text-sm font-semibold text-white mb-2">Параметрүүд:</h4>
          <table className="w-full text-sm">
            <tbody>
              <Row label="response_type" value="code (заавал)" />
              <Row label="client_id" value="Dashboard-аас авсан client_id" />
              <Row label="redirect_uri" value="Бүртгүүлсэн redirect URI" />
              <Row label="scope" value="openid profile (нэмэлт: pos, social, payment)" />
              <Row label="state" value="CSRF хамгаалалт — random string" />
              <Row label="nonce" value="Replay attack-аас хамгаалах" />
            </tbody>
          </table>
        </div>
      </Section>

      <Section title="Алхам 4: Token Exchange">
        <p>Callback URL дээр <code className="text-primary">code</code> авсны дараа token endpoint руу POST хийнэ:</p>
        <Code>{`POST https://sso.gerege.mn/oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code
&code=AUTHORIZATION_CODE
&redirect_uri=https://myapp.mn/callback
&client_id=YOUR_CLIENT_ID
&client_secret=YOUR_CLIENT_SECRET`}</Code>
        <p className="mt-3 text-sm text-slate-400">Хариу:</p>
        <Code>{`{
  "access_token": "eyJhbGciOiJFUzI1NiI...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "id_token": "eyJhbGciOiJFUzI1NiI...",
  "refresh_token": "dGhpcyBpcyBh..."
}`}</Code>
      </Section>

      <Section title="Алхам 5: ID Token Claims">
        <p>ID Token (JWT) decode хийхэд дараах мэдээлэл агуулагдана:</p>
        <Code>{`{
  "sub": "eid-12345678",
  "name": "БАТБОЛД Ганбаатар",
  "given_name": "Ганбаатар",
  "family_name": "БАТБОЛД",
  "cert_serial": "ABC123DEF456",
  "identity_assurance_level": "high",
  "amr": ["eid"],
  "locale": "mn-MN",

  // pos, social scope нэмсэн үед:
  "tenant_id": "t_abc123",
  "tenant_role": "owner",
  "plan": "pro",

  "iss": "https://sso.gerege.mn",
  "aud": "YOUR_CLIENT_ID",
  "exp": 1744200000,
  "iat": 1744196400,
  "nonce": "your-nonce"
}`}</Code>
      </Section>

      <Section title="Алхам 6: UserInfo Endpoint">
        <p>Access token ашиглан нэмэлт мэдээлэл авах:</p>
        <Code>{`GET https://sso.gerege.mn/oauth/userinfo
Authorization: Bearer ACCESS_TOKEN

// Хариу:
{
  "sub": "eid-12345678",
  "name": "БАТБОЛД Ганбаатар",
  "given_name": "Ганбаатар",
  "family_name": "БАТБОЛД",
  "cert_serial": "ABC123DEF456"
}`}</Code>
      </Section>

      <Section title="Жишээ: Next.js">
        <Code>{`// lib/auth.ts
import NextAuth from "next-auth"

export const { handlers, signIn, signOut, auth } = NextAuth({
  providers: [{
    id: "gerege-sso",
    name: "e-ID Mongolia",
    type: "oidc",
    issuer: "https://sso.gerege.mn",
    clientId: process.env.EID_CLIENT_ID!,
    clientSecret: process.env.EID_CLIENT_SECRET!,
    authorization: { params: { scope: "openid profile" } },
  }],
})

// .env.local
// EID_CLIENT_ID=dashboard-аас авсан client_id
// EID_CLIENT_SECRET=dashboard-аас авсан secret`}</Code>
      </Section>

      <Section title="Жишээ: Go">
        <Code>{`package main

import (
    "golang.org/x/oauth2"
    "github.com/coreos/go-oidc/v3/oidc"
)

func main() {
    ctx := context.Background()
    provider, _ := oidc.NewProvider(ctx, "https://sso.gerege.mn")

    oauth2Config := &oauth2.Config{
        ClientID:     os.Getenv("EID_CLIENT_ID"),
        ClientSecret: os.Getenv("EID_CLIENT_SECRET"),
        RedirectURL:  "https://myapp.mn/callback",
        Scopes:       []string{oidc.ScopeOpenID, "profile"},
        Endpoint:     provider.Endpoint(),
    }

    // 1. Redirect
    url := oauth2Config.AuthCodeURL("state", oidc.Nonce("nonce"))
    http.Redirect(w, r, url, http.StatusFound)

    // 2. Callback
    token, _ := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
    rawIDToken, _ := token.Extra("id_token").(string)

    verifier := provider.Verifier(&oidc.Config{ClientID: os.Getenv("EID_CLIENT_ID")})
    idToken, _ := verifier.Verify(ctx, rawIDToken)

    var claims struct {
        Sub       string \`json:"sub"\`
        Name      string \`json:"name"\`
        GivenName string \`json:"given_name"\`
    }
    idToken.Claims(&claims)
}`}</Code>
      </Section>

      <Section title="Жишээ: Python (Flask)">
        <Code>{`from authlib.integrations.flask_client import OAuth

oauth = OAuth(app)
oauth.register(
    name="gerege",
    server_metadata_url="https://sso.gerege.mn/.well-known/openid-configuration",
    client_id="YOUR_CLIENT_ID",
    client_secret="YOUR_CLIENT_SECRET",
    client_kwargs={"scope": "openid profile"},
)

@app.route("/login")
def login():
    redirect_uri = url_for("callback", _external=True)
    return oauth.gerege.authorize_redirect(redirect_uri)

@app.route("/callback")
def callback():
    token = oauth.gerege.authorize_access_token()
    userinfo = token["userinfo"]
    # userinfo["sub"], userinfo["name"], userinfo["given_name"]`}</Code>
      </Section>

      <Section title="Жишээ: PHP (Laravel)">
        <Code>{`// config/services.php
'gerege' => [
    'client_id' => env('EID_CLIENT_ID'),
    'client_secret' => env('EID_CLIENT_SECRET'),
    'redirect' => env('APP_URL') . '/auth/callback',
],

// Laravel Socialite custom provider
use Laravel\\Socialite\\Facades\\Socialite;

Route::get('/login', fn() =>
    Socialite::driver('gerege')
        ->scopes(['openid', 'profile'])
        ->redirect()
);

Route::get('/auth/callback', function () {
    $user = Socialite::driver('gerege')->user();
    // $user->getId()    — sub
    // $user->getName()  — full name
    // $user->getRaw()   — бүх claims
});`}</Code>
      </Section>

      <Section title="Scopes тайлбар">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-white/10">
              <th className="text-left py-2 text-slate-400 font-medium">Scope</th>
              <th className="text-left py-2 text-slate-400 font-medium">Тайлбар</th>
              <th className="text-left py-2 text-slate-400 font-medium">Нэмэлт claims</th>
            </tr>
          </thead>
          <tbody className="text-slate-300">
            <tr className="border-b border-white/5"><td className="py-2"><code className="text-primary">openid</code></td><td>Заавал</td><td>sub, iss, aud, exp, iat</td></tr>
            <tr className="border-b border-white/5"><td className="py-2"><code className="text-primary">profile</code></td><td>Хэрэглэгчийн мэдээлэл</td><td>name, given_name, family_name, cert_serial</td></tr>
            <tr className="border-b border-white/5"><td className="py-2"><code className="text-primary">pos</code></td><td>POS Plugin API</td><td>tenant_id, tenant_role, plan</td></tr>
            <tr className="border-b border-white/5"><td className="py-2"><code className="text-primary">social</code></td><td>Social Commerce API</td><td>tenant_id, tenant_role, plan</td></tr>
            <tr className="border-b border-white/5"><td className="py-2"><code className="text-primary">payment</code></td><td>Payment API</td><td>tenant_id, tenant_role, plan</td></tr>
          </tbody>
        </table>
      </Section>

      <Section title="Аюулгүй байдал">
        <ul className="list-disc list-inside text-slate-400 space-y-2 text-sm">
          <li><strong className="text-white">state параметр</strong> — CSRF халдлагаас хамгаална. Санамсаргүй string үүсгэж session-д хадгалаад callback дээр тулгана.</li>
          <li><strong className="text-white">nonce параметр</strong> — ID Token дотор буцаж ирнэ. Replay attack-аас хамгаалахад хэрэглэнэ.</li>
          <li><strong className="text-white">HTTPS заавал</strong> — Redirect URI заавал HTTPS байх ёстой (localhost-г эс тооцвол).</li>
          <li><strong className="text-white">client_secret хамгаалалт</strong> — Secret-г зөвхөн сервер талд хадгална. Frontend-д ДАМЖУУЛАХГҮЙ.</li>
          <li><strong className="text-white">Token шалгалт</strong> — ID Token-ийн iss, aud, exp, nonce зэргийг заавал шалгана.</li>
          <li><strong className="text-white">Redirect URI тулгалт</strong> — Dashboard-д бүртгүүлсэн URI-тай яг тохирох ёстой.</li>
        </ul>
      </Section>

      <Section title="Алдааны кодууд">
        <table className="w-full text-sm">
          <tbody className="text-slate-300">
            <Row label="invalid_client" value="Буруу client_id эсвэл client_secret" />
            <Row label="invalid_redirect_uri" value="Бүртгэлгүй redirect URI" />
            <Row label="invalid_scope" value="Зөвшөөрөгдөөгүй scope" />
            <Row label="invalid_grant" value="Хугацаа дууссан эсвэл ашиглагдсан code" />
            <Row label="access_denied" value="Хэрэглэгч нэвтрэлтийг цуцалсан" />
          </tbody>
        </table>
      </Section>

      <div className="p-4 bg-surface border border-white/10 rounded-xl text-sm text-slate-400">
        <p className="font-semibold text-white mb-2">Тусламж</p>
        <p>Асуудал гарвал <a href="mailto:dev@gerege.mn" className="text-primary hover:underline">dev@gerege.mn</a> хаягаар холбогдоно уу.
        App бүртгэлийн бүх тохиргоог <a href="/dashboard/apps" className="text-primary hover:underline">Dashboard → Apps</a> хуудаснаас хийнэ.</p>
      </div>
    </main>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <section className="space-y-3">
      <h2 className="text-xl font-bold text-white">{title}</h2>
      <div className="text-sm text-slate-400 leading-relaxed">{children}</div>
    </section>
  );
}

function Code({ children }: { children: string }) {
  return (
    <pre className="bg-bg border border-white/10 rounded-xl p-4 text-xs text-slate-300 overflow-x-auto whitespace-pre-wrap">{children}</pre>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <tr className="border-b border-white/5">
      <td className="py-2 pr-4 text-slate-400 font-mono text-xs whitespace-nowrap">{label}</td>
      <td className="py-2 text-slate-300 text-xs">{value}</td>
    </tr>
  );
}

function Field({ label, desc }: { label: string; desc: string }) {
  return (
    <div className="p-3 bg-bg border border-white/10 rounded-lg">
      <p className="text-sm font-semibold text-white">{label}</p>
      <p className="text-xs text-slate-400 mt-1">{desc}</p>
    </div>
  );
}
