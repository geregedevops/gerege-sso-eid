"use client";

import { useState, useEffect, useRef } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Suspense } from "react";

export default function SignPage() {
  return (
    <Suspense fallback={<div className="text-slate-400">Уншиж байна...</div>}>
      <SignContent />
    </Suspense>
  );
}

function SignContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const preselectedOrg = searchParams.get("org") || "";

  const [orgs, setOrgs] = useState<any[]>([]);
  const [orgId, setOrgId] = useState(preselectedOrg);
  const [file, setFile] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const fileRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    fetch("/api/org").then((r) => r.json()).then((data) => {
      const signable = data.filter((o: any) => ["owner", "admin", "signer"].includes(o.role));
      setOrgs(signable);
      if (!orgId && signable.length === 1) setOrgId(signable[0].id);
    });
  }, []);

  async function handleSign() {
    if (!orgId || !file) { setError("Байгууллага сонгож, PDF файл оруулна уу"); return; }
    if (file.size > 10 * 1024 * 1024) { setError("Файл хэт том (max 10MB)"); return; }

    setLoading(true);
    setError("");

    try {
      const base64 = await fileToBase64(file);
      const res = await fetch("/api/sign", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ organizationId: orgId, documentName: file.name, document: base64 }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || "Алдаа");
      router.push(`/dashboard/sign/${data.signatureId}`);
    } catch (err: any) {
      setError(err.message);
      setLoading(false);
    }
  }

  return (
    <div className="max-w-lg mx-auto space-y-6">
      <h1 className="text-2xl font-bold text-white">Гарын үсэг зурах</h1>

      <div>
        <label className="block text-sm text-slate-400 mb-1">Байгууллага *</label>
        <select value={orgId} onChange={(e) => setOrgId(e.target.value)}
          className="w-full px-4 py-3 bg-bg border border-white/10 rounded-xl text-white text-sm outline-none focus:border-primary">
          <option value="">Сонгох...</option>
          {orgs.map((o) => <option key={o.id} value={o.id}>{o.name} ({o.registrationNumber})</option>)}
        </select>
      </div>

      <div>
        <label className="block text-sm text-slate-400 mb-1">PDF файл * (max 10MB)</label>
        <div
          onClick={() => fileRef.current?.click()}
          className="w-full p-8 bg-bg border-2 border-dashed border-white/10 rounded-xl text-center cursor-pointer hover:border-primary/30 transition-colors"
        >
          {file ? (
            <div>
              <p className="text-white font-medium">{file.name}</p>
              <p className="text-xs text-slate-400 mt-1">{(file.size / 1024).toFixed(0)} KB</p>
            </div>
          ) : (
            <p className="text-slate-500 text-sm">PDF файл сонгохын тулд энд дарна уу</p>
          )}
        </div>
        <input ref={fileRef} type="file" accept=".pdf" className="hidden" onChange={(e) => setFile(e.target.files?.[0] || null)} />
      </div>

      {error && <p className="text-red-400 text-sm">{error}</p>}

      <button onClick={handleSign} disabled={loading || !orgId || !file}
        className="w-full py-3 bg-primary text-white font-semibold rounded-xl hover:bg-primary-light transition-colors disabled:opacity-50">
        {loading ? "SmartID хүсэлт илгээж байна..." : "SmartID PIN2-оор гарын үсэг зурах"}
      </button>

      <p className="text-xs text-slate-500 text-center">
        SmartID апп дээр PIN2 оруулах хүсэлт илгээгдэнэ.
      </p>
    </div>
  );
}

function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => {
      const result = reader.result as string;
      resolve(result.split(",")[1]); // Remove data:... prefix
    };
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
}
