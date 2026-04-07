"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

export default function NewTenantPage() {
  const router = useRouter();
  const [name, setName] = useState("");
  const [slug, setSlug] = useState("");
  const [plan, setPlan] = useState("starter");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const submit = async () => {
    if (!name.trim() || !slug.trim()) return;
    setLoading(true);
    setError("");
    try {
      const res = await fetch("/api/tenants", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: name.trim(), slug: slug.trim().toLowerCase(), plan }),
      });
      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || "Failed");
      }
      router.push("/dashboard/tenants");
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-lg mx-auto space-y-6">
      <h1 className="text-2xl font-bold text-white">Шинэ Tenant үүсгэх</h1>

      <div className="bg-surface border border-white/10 rounded-xl p-6 space-y-5">
        <Field label="Tenant нэр *">
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full px-3 py-2 bg-bg border border-white/10 rounded-lg text-white text-sm focus:ring-1 focus:ring-primary outline-none"
            placeholder="Миний ресторан"
          />
        </Field>

        <Field label="Slug *">
          <input
            value={slug}
            onChange={(e) => setSlug(e.target.value.replace(/[^a-z0-9-]/g, ""))}
            className="w-full px-3 py-2 bg-bg border border-white/10 rounded-lg text-white text-sm font-mono focus:ring-1 focus:ring-primary outline-none"
            placeholder="my-restaurant"
          />
          <p className="text-xs text-slate-500 mt-1">Жижиг үсэг, тоо, зураас (-) зөвхөн</p>
        </Field>

        <Field label="Plan">
          <select
            value={plan}
            onChange={(e) => setPlan(e.target.value)}
            className="w-full px-3 py-2 bg-bg border border-white/10 rounded-lg text-white text-sm focus:ring-1 focus:ring-primary outline-none"
          >
            <option value="starter">Starter</option>
            <option value="pro">Pro</option>
            <option value="enterprise">Enterprise</option>
          </select>
        </Field>

        {error && <p className="text-sm text-red-400">{error}</p>}

        <button
          onClick={submit}
          disabled={!name.trim() || !slug.trim() || loading}
          className="w-full py-3 bg-primary text-white font-semibold rounded-xl hover:bg-primary-light transition-colors disabled:opacity-40"
        >
          {loading ? "Үүсгэж байна..." : "Tenant үүсгэх"}
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
