export default function APIReferencePage() {
  return (
    <main className="max-w-4xl mx-auto px-6 py-12 space-y-8">
      <h1 className="text-3xl font-bold text-white">API Reference</h1>
      <p className="text-slate-400">sso.gerege.mn OIDC endpoint-ууд</p>

      <Endpoint method="GET" path="/.well-known/openid-configuration" desc="OIDC Discovery document" />
      <Endpoint method="GET" path="/.well-known/jwks.json" desc="EC P-256 JWK Set (ES256)" />
      <Endpoint method="GET" path="/oauth/authorize" desc="Authorization endpoint — browser redirect">
        <Params items={[
          ["client_id", "string", "Бүртгэлтэй client ID"],
          ["redirect_uri", "string", "Exact match бүртгэлтэй URI"],
          ["response_type", "string", "code"],
          ["scope", "string", "openid profile [pos] [social] [payment]"],
          ["state", "string", "CSRF protection"],
          ["nonce", "string", "ID token-д хадгалагдана"],
        ]} />
      </Endpoint>

      <Endpoint method="POST" path="/oauth/token" desc="Token endpoint — code exchange">
        <Params items={[
          ["grant_type", "string", "authorization_code"],
          ["code", "string", "Auth code (нэг удаа)"],
          ["client_id", "string", ""],
          ["client_secret", "string", "Basic auth эсвэл form"],
          ["redirect_uri", "string", "Exact match"],
        ]} />
        <h4 className="text-sm font-semibold text-white mt-4 mb-2">Response</h4>
        <pre className="bg-bg border border-white/10 rounded-lg p-3 text-xs text-slate-300">{`{
  "access_token": "opaque_token",
  "token_type": "Bearer",
  "expires_in": 3600,
  "id_token": "eyJhbGciOiJFUzI1NiJ9...",
  "scope": "openid profile pos"
}`}</pre>
      </Endpoint>

      <Endpoint method="GET" path="/oauth/userinfo" desc="User info — Bearer token шаардана">
        <h4 className="text-sm font-semibold text-white mt-3 mb-2">Response</h4>
        <pre className="bg-bg border border-white/10 rounded-lg p-3 text-xs text-slate-300">{`{
  "sub": "sha256_of_national_id",
  "name": "Батаа Дорж",
  "given_name": "Дорж",
  "family_name": "Батаа",
  "locale": "mn-MN",
  "tenant_id": "restaurant-govi",
  "tenant_role": "owner",
  "plan": "pro"
}`}</pre>
      </Endpoint>

      <Endpoint method="POST" path="/oauth/revoke" desc="Token revoke — RFC 7009" />
      <Endpoint method="POST" path="/oauth/introspect" desc="Token introspection — active/inactive" />

      <div className="border-t border-white/10 pt-8">
        <h2 className="text-xl font-bold text-white mb-3">ID Token Claims</h2>
        <div className="bg-surface border border-white/10 rounded-xl overflow-hidden">
          <table className="w-full text-sm">
            <thead><tr className="border-b border-white/10 text-left text-xs text-slate-500">
              <th className="px-4 py-2">Claim</th><th className="px-4 py-2">Тайлбар</th>
            </tr></thead>
            <tbody className="text-slate-300">
              {[
                ["sub", "Регистрийн дугаарын SHA-256 hex"],
                ["name", "Бүтэн нэр"],
                ["given_name", "Нэр"],
                ["family_name", "Овог"],
                ["cert_serial", "X.509 certificate serial"],
                ["identity_assurance_level", "high"],
                ["amr", '["smartid", "pin1", "x509"]'],
                ["tenant_id", "Tenant slug (pos/social scope)"],
                ["tenant_role", "owner | admin | member"],
                ["plan", "starter | pro | enterprise"],
                ["locale", "mn-MN"],
              ].map(([k, v]) => (
                <tr key={k} className="border-b border-white/5">
                  <td className="px-4 py-2 font-mono text-primary text-xs">{k}</td>
                  <td className="px-4 py-2 text-xs">{v}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </main>
  );
}

function Endpoint({ method, path, desc, children }: { method: string; path: string; desc: string; children?: React.ReactNode }) {
  const color = method === "GET" ? "bg-green-500/15 text-green-400" : "bg-yellow-500/15 text-yellow-400";
  return (
    <div className="bg-surface border border-white/10 rounded-xl p-5">
      <div className="flex items-center gap-3 mb-2">
        <span className={`px-2 py-0.5 rounded text-xs font-bold ${color}`}>{method}</span>
        <code className="text-sm font-mono text-white">{path}</code>
      </div>
      <p className="text-sm text-slate-400">{desc}</p>
      {children}
    </div>
  );
}

function Params({ items }: { items: string[][] }) {
  return (
    <div className="mt-3">
      <h4 className="text-sm font-semibold text-white mb-2">Parameters</h4>
      <div className="space-y-1">
        {items.map(([name, type, desc]) => (
          <div key={name} className="flex items-center gap-3 text-xs">
            <code className="text-primary font-mono w-28">{name}</code>
            <span className="text-slate-500 w-12">{type}</span>
            <span className="text-slate-400">{desc}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
