"use client";

import { useState, useRef } from "react";

export default function VerifyPage() {
  const [file, setFile] = useState<File | null>(null);
  const [result, setResult] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const fileRef = useRef<HTMLInputElement>(null);

  async function handleVerify() {
    if (!file) return;
    setLoading(true);
    setError("");
    setResult(null);

    try {
      const formData = new FormData();
      formData.append("file", file);
      const res = await fetch("/api/verify", { method: "POST", body: formData });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || "Шалгалт амжилтгүй");
      setResult(data);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="min-h-[calc(100vh-56px)] flex items-center justify-center p-6">
      <div className="max-w-md w-full space-y-6">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-white mb-2">Баримт шалгах</h1>
          <p className="text-sm text-slate-400">Гарын үсэг зурсан PDF файлыг upload хийж бүрэн бүтэн байдлыг шалгана.</p>
        </div>

        <div
          onClick={() => fileRef.current?.click()}
          className="w-full p-10 bg-surface border-2 border-dashed border-white/10 rounded-xl text-center cursor-pointer hover:border-primary/30 transition-colors"
        >
          {file ? (
            <div>
              <p className="text-white font-medium">{file.name}</p>
              <p className="text-xs text-slate-400 mt-1">{(file.size / 1024).toFixed(0)} KB</p>
            </div>
          ) : (
            <p className="text-slate-500 text-sm">PDF файл сонгох</p>
          )}
        </div>
        <input ref={fileRef} type="file" accept=".pdf" className="hidden" onChange={(e) => { setFile(e.target.files?.[0] || null); setResult(null); }} />

        <button onClick={handleVerify} disabled={loading || !file}
          className="w-full py-3 bg-primary text-white font-semibold rounded-xl hover:bg-primary-light transition-colors disabled:opacity-50">
          {loading ? "Шалгаж байна..." : "Шалгах"}
        </button>

        {error && <p className="text-red-400 text-sm text-center">{error}</p>}

        {result && (
          <div className={`p-5 rounded-xl border ${result.valid ? "bg-green-500/5 border-green-500/20" : "bg-red-500/5 border-red-500/20"}`}>
            <div className="flex items-center gap-3 mb-3">
              <span className={`text-2xl ${result.valid ? "text-green-400" : "text-red-400"}`}>
                {result.valid ? "\u2713" : "\u2717"}
              </span>
              <h3 className={`font-bold ${result.valid ? "text-green-400" : "text-red-400"}`}>
                {result.valid ? "Баталгаажсан" : "Баталгаажаагүй"}
              </h3>
            </div>
            {result.signatures?.map((s: any, i: number) => (
              <div key={i} className="text-sm text-slate-300 space-y-1 mt-2">
                {s.signerName && <p>Гарын үсэг: <strong className="text-white">{s.signerName}</strong></p>}
                {s.organizationName && <p>Байгууллага: {s.organizationName}</p>}
                {s.signedAt && <p>Огноо: {s.signedAt}</p>}
                {s.certSerial && <p className="font-mono text-xs text-slate-400">Cert: {s.certSerial}</p>}
              </div>
            ))}
            {result.message && <p className="text-sm text-slate-400 mt-2">{result.message}</p>}
          </div>
        )}
      </div>
    </main>
  );
}
