"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

const GEREGE_SCOPES = [
  { id: "openid", label: "OpenID", description: "Үндсэн нэвтрэлт", locked: true },
  { id: "profile", label: "Profile", description: "Нэр, регистрийн мэдээлэл" },
  { id: "pos", label: "POS", description: "Борлуулалтын систем, захиалга, бараа" },
  { id: "social", label: "Social Commerce", description: "Нийгмийн сүлжээ, лайв худалдаа" },
  { id: "payment", label: "Payment", description: "QPay, SocialPay, eBarimt" },
];

export default function NewAppPage() {
  const router = useRouter();
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [redirectUris, setRedirectUris] = useState([""]);
  const [scopes, setScopes] = useState(["openid", "profile"]);
  const [loading, setLoading] = useState(false);
  const [credentials, setCredentials] = useState<{ clientId: string; clientSecret: string } | null>(null);
  const [error, setError] = useState("");

  const toggleScope = (scope: string) => {
    if (scope === "openid") return;
    setScopes((prev) => prev.includes(scope) ? prev.filter((s) => s !== scope) : [...prev, scope]);
  };

  const submit = async () => {
    if (!name.trim()) return;
    const uris = redirectUris.filter((u) => u.trim());
    if (uris.length === 0) { setError("Redirect URI оруулна уу"); return; }

    setLoading(true);
    setError("");
    try {
      const res = await fetch("/api/apps", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: name.trim(), description: description.trim(), redirectUris: uris, scopes }),
      });
      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || "Failed");
      }
      const data = await res.json();
      setCredentials({ clientId: data.clientId, clientSecret: data.clientSecret });
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  if (credentials) {
    return (
      <div className="max-w-lg mx-auto space-y-6">
        <div className="bg-yellow-500/10 border border-yellow-500/30 rounded-xl p-4 text-yellow-300 text-sm">
          Credentials-ийг нэг л удаа харуулна. Заавал хадгалаарай!
        </div>
        <div className="bg-surface border border-white/10 rounded-xl p-6 space-y-4">
          <CredRow label="client_id" value={credentials.clientId} />
          <CredRow label="client_secret" value={credentials.clientSecret} />
        </div>
        <button
          onClick={() => router.push("/dashboard/apps")}
          className="w-full py-3 bg-primary text-white font-semibold rounded-xl"
        >
          Ойлголоо, хаах
        </button>
      </div>
    );
  }

  return (
    <div className="max-w-lg mx-auto space-y-6">
      <h1 className="text-2xl font-bold text-white">Шинэ App бүртгүүлэх</h1>

      <div className="bg-surface border border-white/10 rounded-xl p-6 space-y-5">
        <Field label="App нэр *">
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full px-3 py-2 bg-bg border border-white/10 rounded-lg text-white text-sm focus:ring-1 focus:ring-primary outline-none"
            placeholder="My Application"
          />
        </Field>

        <Field label="Тайлбар">
          <input
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            className="w-full px-3 py-2 bg-bg border border-white/10 rounded-lg text-white text-sm focus:ring-1 focus:ring-primary outline-none"
            placeholder="Юу хийдэг app вэ?"
          />
        </Field>

        <Field label="Redirect URIs *">
          {redirectUris.map((uri, i) => (
            <div key={i} className="flex gap-2 mb-2">
              <input
                value={uri}
                onChange={(e) => {
                  const copy = [...redirectUris];
                  copy[i] = e.target.value;
                  setRedirectUris(copy);
                }}
                className="flex-1 px-3 py-2 bg-bg border border-white/10 rounded-lg text-white text-sm font-mono focus:ring-1 focus:ring-primary outline-none"
                placeholder="https://myapp.mn/callback"
              />
              {redirectUris.length > 1 && (
                <button
                  onClick={() => setRedirectUris(redirectUris.filter((_, j) => j !== i))}
                  className="px-2 text-red-400 hover:text-red-300"
                >
                  x
                </button>
              )}
            </div>
          ))}
          <button
            onClick={() => setRedirectUris([...redirectUris, ""])}
            className="text-xs text-primary hover:text-primary-light"
          >
            + URI нэмэх
          </button>
        </Field>

        <Field label="Scopes">
          <div className="space-y-2">
            {GEREGE_SCOPES.map((s) => (
              <label key={s.id} className="flex items-start gap-2 text-sm">
                <input
                  type="checkbox"
                  checked={scopes.includes(s.id)}
                  onChange={() => toggleScope(s.id)}
                  disabled={s.locked}
                  className="accent-primary mt-0.5"
                />
                <div>
                  <span className="text-white font-medium">{s.label}</span>
                  <p className="text-xs text-slate-500">{s.description}</p>
                </div>
              </label>
            ))}
          </div>
        </Field>

        {error && <p className="text-sm text-red-400">{error}</p>}

        <button
          onClick={submit}
          disabled={!name.trim() || loading}
          className="w-full py-3 bg-primary text-white font-semibold rounded-xl hover:bg-primary-light transition-colors disabled:opacity-40"
        >
          {loading ? "Үүсгэж байна..." : "App үүсгэх"}
        </button>
      </div>
    </div>
  );
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div>
      <label className="block text-xs text-slate-400 mb-1.5 font-medium">{label}</label>
      {children}
    </div>
  );
}

function CredRow({ label, value }: { label: string; value: string }) {
  const [copied, setCopied] = useState(false);
  const copy = () => {
    navigator.clipboard.writeText(value);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };
  return (
    <div>
      <p className="text-xs text-slate-500 mb-1">{label}</p>
      <div className="flex items-center gap-2">
        <code className="flex-1 text-sm font-mono text-white bg-bg px-3 py-2 rounded-lg border border-white/10 break-all">
          {value}
        </code>
        <button onClick={copy} className="px-3 py-2 text-xs text-primary hover:text-primary-light border border-white/10 rounded-lg">
          {copied ? "Copied" : "Хуулах"}
        </button>
      </div>
    </div>
  );
}
