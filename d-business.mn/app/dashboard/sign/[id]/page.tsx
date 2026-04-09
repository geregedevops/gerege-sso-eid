"use client";

import { Suspense } from "react";
import { useParams } from "next/navigation";
import { useState, useEffect, useRef } from "react";

export default function SignSessionPage() {
  return (
    <Suspense fallback={<div className="text-slate-400">Уншиж байна...</div>}>
      <SignSessionContent />
    </Suspense>
  );
}

function SignSessionContent() {
  const { id } = useParams();
  const [status, setStatus] = useState("pending");
  const [signerName, setSigner] = useState("");
  const [error, setError] = useState("");
  const [verificationCode, setVerificationCode] = useState("");
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    // Load initial signature data
    pollStatus();
    intervalRef.current = setInterval(pollStatus, 3000);
    return () => { if (intervalRef.current) clearInterval(intervalRef.current); };
  }, []);

  async function pollStatus() {
    try {
      const res = await fetch(`/api/sign/${id}/status`);
      const data = await res.json();
      setStatus(data.status);
      if (data.signerName) setSigner(data.signerName);
      if (data.status === "complete" || data.status === "failed") {
        if (intervalRef.current) clearInterval(intervalRef.current);
      }
      if (data.error && data.status === "failed") setError(data.error);
    } catch {}
  }

  return (
    <div className="max-w-md mx-auto text-center space-y-6 py-12">
      {status === "pending" && (
        <>
          <div className="w-16 h-16 bg-primary/10 rounded-2xl flex items-center justify-center mx-auto">
            <div className="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin" />
          </div>
          <h1 className="text-2xl font-bold text-white">SmartID хүлээж байна</h1>
          <p className="text-slate-400 text-sm">SmartID апп дээр PIN2 оруулна уу.</p>
          {verificationCode && (
            <div className="bg-surface border border-white/10 rounded-xl p-6">
              <p className="text-xs text-slate-400 mb-2">Verification Code</p>
              <p className="text-4xl font-bold text-primary tracking-wider">{verificationCode}</p>
            </div>
          )}
        </>
      )}

      {status === "complete" && (
        <>
          <div className="w-16 h-16 bg-green-500/10 rounded-2xl flex items-center justify-center mx-auto">
            <span className="text-green-400 text-3xl">&#10003;</span>
          </div>
          <h1 className="text-2xl font-bold text-white">Гарын үсэг амжилттай!</h1>
          {signerName && <p className="text-slate-400 text-sm">Гарын үсэг зурсан: {signerName}</p>}
          <a href={`/api/sign/${id}/result`}
            className="inline-block px-6 py-3 bg-primary text-white font-semibold rounded-xl hover:bg-primary-light transition-colors">
            Татаж авах (PDF)
          </a>
        </>
      )}

      {status === "failed" && (
        <>
          <div className="w-16 h-16 bg-red-500/10 rounded-2xl flex items-center justify-center mx-auto">
            <span className="text-red-400 text-3xl">&#10007;</span>
          </div>
          <h1 className="text-2xl font-bold text-red-400">Алдаа гарлаа</h1>
          <p className="text-slate-400 text-sm">{error || "Гарын үсэг зурах үйлдэл цуцлагдсан эсвэл хугацаа дууссан."}</p>
          <a href="/dashboard/sign" className="inline-block px-6 py-3 border border-white/15 text-white font-medium rounded-xl hover:bg-white/5 transition-colors">
            Дахин оролдох
          </a>
        </>
      )}
    </div>
  );
}
