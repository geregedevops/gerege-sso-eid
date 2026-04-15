"use client";

import { useState } from "react";

export default function AdminVerifyPage() {
  const [tab, setTab] = useState<"citizen" | "org">("citizen");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<any>(null);
  const [error, setError] = useState("");

  // Citizen fields
  const [regNo, setRegNo] = useState("");
  const [phone, setPhone] = useState("");

  // Org fields
  const [orgRegNo, setOrgRegNo] = useState("");
  const [ceoRegNo, setCeoRegNo] = useState("");

  function reset() {
    setResult(null);
    setError("");
  }

  async function handleCitizen(e: React.FormEvent) {
    e.preventDefault();
    if (!regNo.trim() || !phone.trim()) return;
    setLoading(true);
    reset();
    try {
      const res = await fetch("/api/try/authenticate/citizen", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ reg_no: regNo.trim(), phone: phone.trim() }),
      });
      const data = await res.json();
      if (!res.ok) setError(data.error || "Алдаа");
      else setResult(data);
    } catch {
      setError("Сервертэй холбогдож чадсангүй");
    } finally {
      setLoading(false);
    }
  }

  async function handleOrg(e: React.FormEvent) {
    e.preventDefault();
    if (!orgRegNo.trim() || !ceoRegNo.trim()) return;
    setLoading(true);
    reset();
    try {
      const res = await fetch("/api/try/authenticate/org", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ reg_no: orgRegNo.trim(), ceo_reg_no: ceoRegNo.trim() }),
      });
      const data = await res.json();
      if (!res.ok) setError(data.error || "Алдаа");
      else setResult(data);
    } catch {
      setError("Сервертэй холбогдож чадсангүй");
    } finally {
      setLoading(false);
    }
  }

  const citizen = result?.citizen;
  const org = result?.organization;

  return (
    <div className="max-w-2xl mx-auto px-6 py-10">
      <h1 className="text-2xl font-bold text-white mb-2">Баталгаажуулалт</h1>
      <p className="text-sm text-slate-400 mb-8">
        Иргэн болон байгууллагын мэдээллийг тулгаж баталгаажуулна
      </p>

      {/* Tabs */}
      <div className="flex gap-2 mb-6">
        <button
          onClick={() => { setTab("citizen"); reset(); }}
          className={`px-5 py-2 rounded-xl text-sm font-medium transition-all ${
            tab === "citizen" ? "bg-primary text-white" : "bg-white/[0.05] text-slate-400 hover:text-white"
          }`}
        >
          Иргэн
        </button>
        <button
          onClick={() => { setTab("org"); reset(); }}
          className={`px-5 py-2 rounded-xl text-sm font-medium transition-all ${
            tab === "org" ? "bg-primary text-white" : "bg-white/[0.05] text-slate-400 hover:text-white"
          }`}
        >
          Байгууллага
        </button>
      </div>

      {/* Citizen form */}
      {tab === "citizen" && (
        <form onSubmit={handleCitizen} className="space-y-4 mb-8">
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1.5">Регистрийн дугаар</label>
            <input
              type="text"
              value={regNo}
              onChange={(e) => setRegNo(e.target.value)}
              placeholder="МА74101813"
              required
              className="w-full px-4 py-3 bg-white/[0.05] border border-white/10 rounded-xl text-white text-sm focus:outline-none focus:border-primary"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1.5">Утасны дугаар</label>
            <input
              type="text"
              value={phone}
              onChange={(e) => setPhone(e.target.value)}
              placeholder="99112233"
              required
              className="w-full px-4 py-3 bg-white/[0.05] border border-white/10 rounded-xl text-white text-sm focus:outline-none focus:border-primary"
            />
          </div>
          <button
            type="submit"
            disabled={loading}
            className="w-full py-3 bg-primary hover:bg-primary-light text-white font-bold rounded-xl transition-all disabled:opacity-50"
          >
            {loading ? "Шалгаж байна..." : "Баталгаажуулах"}
          </button>
        </form>
      )}

      {/* Org form */}
      {tab === "org" && (
        <form onSubmit={handleOrg} className="space-y-4 mb-8">
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1.5">Байгууллагын регистр</label>
            <input
              type="text"
              value={orgRegNo}
              onChange={(e) => setOrgRegNo(e.target.value)}
              placeholder="6235972"
              required
              className="w-full px-4 py-3 bg-white/[0.05] border border-white/10 rounded-xl text-white text-sm focus:outline-none focus:border-primary"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1.5">Захирлын регистрийн дугаар</label>
            <input
              type="text"
              value={ceoRegNo}
              onChange={(e) => setCeoRegNo(e.target.value)}
              placeholder="уш72060800"
              required
              className="w-full px-4 py-3 bg-white/[0.05] border border-white/10 rounded-xl text-white text-sm focus:outline-none focus:border-primary"
            />
          </div>
          <button
            type="submit"
            disabled={loading}
            className="w-full py-3 bg-primary hover:bg-primary-light text-white font-bold rounded-xl transition-all disabled:opacity-50"
          >
            {loading ? "Шалгаж байна..." : "Баталгаажуулах"}
          </button>
        </form>
      )}

      {/* Error */}
      {error && (
        <div className="mb-6 p-4 bg-red-500/10 border border-red-500/20 rounded-xl">
          <p className="text-sm text-red-400">{error}</p>
        </div>
      )}

      {/* Auth failed */}
      {result && !result.authenticated && (
        <div className="p-6 bg-red-500/10 border border-red-500/20 rounded-2xl text-center">
          <div className="text-3xl mb-2">✗</div>
          <p className="text-red-400 font-semibold mb-1">Баталгаажуулалт амжилтгүй</p>
          <p className="text-sm text-slate-400">{result.reason}</p>
        </div>
      )}

      {/* Citizen result */}
      {result?.authenticated && citizen && (
        <div className="bg-white/[0.03] border border-primary/30 rounded-2xl overflow-hidden">
          <div className="px-5 py-3 bg-primary/10 border-b border-primary/20 flex items-center gap-2">
            <span className="text-primary text-lg">✓</span>
            <span className="text-primary font-semibold text-sm">Баталгаажсан</span>
          </div>
          <div className="flex">
            {citizen.image && (
              <div className="p-5 border-r border-white/[0.06] flex-shrink-0">
                <img
                  src={`data:image/jpeg;base64,${citizen.image}`}
                  alt=""
                  className="w-28 h-36 object-cover rounded-xl"
                />
              </div>
            )}
            <div className="flex-1 divide-y divide-white/[0.06]">
              <Row label="Регистр" value={citizen.reg_no} mono />
              <Row label="ИҮ дугаар" value={citizen.civil_id} mono />
              <Row label="Овог" value={citizen.last_name} />
              <Row label="Нэр" value={citizen.first_name} bold />
              <Row label="Хүйс" value={citizen.gender} />
              <Row label="Төрсөн огноо" value={citizen.birth_date} />
            </div>
          </div>
        </div>
      )}

      {/* Org result */}
      {result?.authenticated && org && (
        <div className="bg-white/[0.03] border border-primary/30 rounded-2xl overflow-hidden">
          <div className="px-5 py-3 bg-primary/10 border-b border-primary/20 flex items-center gap-2">
            <span className="text-primary text-lg">✓</span>
            <span className="text-primary font-semibold text-sm">Баталгаажсан — Захирлын РД таарч байна</span>
          </div>
          <div className="divide-y divide-white/[0.06]">
            <Row label="Регистр" value={org.reg_no} mono />
            <Row label="Нэр" value={org.name} bold />
            <Row label="Төрөл" value={org.type} />
            <Row label={org.ceo_position || "Захирал"} value={org.ceo} />
            <Row label="Захирлын РД" value={org.ceo_reg_no} mono />
          </div>
        </div>
      )}
    </div>
  );
}

function Row({ label, value, mono, bold }: { label: string; value?: string; mono?: boolean; bold?: boolean }) {
  if (!value) return null;
  return (
    <div className="px-5 py-2.5 flex justify-between">
      <span className="text-slate-400 text-sm">{label}</span>
      <span className={`text-white text-sm ${mono ? "font-mono" : ""} ${bold ? "font-semibold" : ""}`}>
        {value}
      </span>
    </div>
  );
}
