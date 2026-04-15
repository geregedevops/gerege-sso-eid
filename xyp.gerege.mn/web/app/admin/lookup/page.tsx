"use client";

import { useState } from "react";

export default function AdminLookupPage() {
  const [tab, setTab] = useState<"citizen" | "org">("citizen");
  const [regNo, setRegNo] = useState("");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<any>(null);
  const [error, setError] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!regNo.trim()) return;
    setLoading(true);
    setError("");
    setResult(null);

    try {
      const endpoint = tab === "citizen" ? "/api/try/citizen" : "/api/try/org";
      const res = await fetch(endpoint, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ reg_no: regNo.trim() }),
      });
      const data = await res.json();
      if (!res.ok) {
        setError(data.error || "Алдаа гарлаа");
      } else {
        setResult(data);
      }
    } catch {
      setError("Сервертэй холбогдож чадсангүй");
    } finally {
      setLoading(false);
    }
  }

  const citizen = result?.citizen;
  const org = result?.organization;

  return (
    <div className="max-w-4xl mx-auto px-6 py-10">
      <h1 className="text-2xl font-bold text-white mb-2">Дэлгэрэнгүй хайлт</h1>
      <p className="text-sm text-slate-400 mb-8">Бүрэн мэдээлэл — зөвхөн admin хандалт</p>

      {/* Tabs */}
      <div className="flex gap-2 mb-6">
        <button
          onClick={() => { setTab("citizen"); setResult(null); setError(""); }}
          className={`px-5 py-2 rounded-xl text-sm font-medium transition-all ${
            tab === "citizen" ? "bg-primary text-white" : "bg-white/[0.05] text-slate-400 hover:text-white"
          }`}
        >
          Иргэн
        </button>
        <button
          onClick={() => { setTab("org"); setResult(null); setError(""); }}
          className={`px-5 py-2 rounded-xl text-sm font-medium transition-all ${
            tab === "org" ? "bg-primary text-white" : "bg-white/[0.05] text-slate-400 hover:text-white"
          }`}
        >
          Байгууллага
        </button>
      </div>

      {/* Search */}
      <form onSubmit={handleSubmit} className="flex gap-3 mb-8">
        <input
          type="text"
          value={regNo}
          onChange={(e) => setRegNo(e.target.value)}
          placeholder={tab === "citizen" ? "Иргэний РД (жнь: МА74101813)" : "Байгууллагын регистр (жнь: 6235972)"}
          className="flex-1 px-4 py-3 bg-white/[0.05] border border-white/10 rounded-xl text-white text-sm focus:outline-none focus:border-primary"
        />
        <button
          type="submit"
          disabled={loading}
          className="px-6 py-3 bg-primary hover:bg-primary-light text-white font-bold text-sm rounded-xl transition-all disabled:opacity-50"
        >
          {loading ? "Хайж байна..." : "Хайх"}
        </button>
      </form>

      {error && (
        <div className="mb-6 p-4 bg-red-500/10 border border-red-500/20 rounded-xl">
          <p className="text-sm text-red-400">{error}</p>
        </div>
      )}

      {/* Citizen full result */}
      {tab === "citizen" && citizen && (
        <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl overflow-hidden">
          <div className="px-5 py-4 border-b border-white/[0.06]">
            <h2 className="font-semibold text-white">Иргэний бүрэн мэдээлэл</h2>
          </div>
          <div className="divide-y divide-white/[0.06]">
            <Row label="Регистр" value={citizen.reg_no} mono />
            <Row label="Овог" value={citizen.last_name} />
            <Row label="Нэр" value={citizen.first_name} bold />
            <Row label="Ургийн овог" value={citizen.surname} />
            <Row label="Хүйс" value={citizen.gender} />
            <Row label="Төрсөн огноо" value={citizen.birth_date} />
            <Row label="Үндэс угсаа" value={citizen.nationality} />
          </div>
        </div>
      )}

      {tab === "citizen" && result && !citizen && (
        <div className="p-6 text-center text-slate-500 bg-white/[0.03] border border-white/[0.06] rounded-2xl">
          Иргэн олдсонгүй
        </div>
      )}

      {/* Org full result */}
      {tab === "org" && org && (
        <div className="space-y-6">
          {/* Basic info */}
          <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl overflow-hidden">
            <div className="px-5 py-4 border-b border-white/[0.06]">
              <h2 className="font-semibold text-white">Байгууллагын мэдээлэл</h2>
            </div>
            <div className="divide-y divide-white/[0.06]">
              <Row label="Регистр" value={org.reg_no} mono />
              <Row label="Нэр" value={org.name} bold />
              <Row label="Төрөл" value={org.type} />
              <Row label="Дүрмийн сан" value={org.capital ? `${Number(org.capital).toLocaleString()}₮` : ""} />
              <Row label={org.ceo_position || "Захирал"} value={org.ceo} />
              <Row label="Захирлын РД" value={org.ceo_reg_no} mono />
              <Row label="Утас" value={org.phone} />
              <Row label="Хаяг" value={org.address} />
            </div>
          </div>

          {/* Industry */}
          {org.industry?.length > 0 && (
            <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl overflow-hidden">
              <div className="px-5 py-4 border-b border-white/[0.06]">
                <h2 className="font-semibold text-white">Үйл ажиллагааны чиглэл</h2>
              </div>
              <div className="px-5 py-4 flex flex-wrap gap-2">
                {org.industry.map((ind: string, i: number) => (
                  <span key={i} className="px-3 py-1.5 bg-primary/10 text-primary text-xs rounded-lg">
                    {ind}
                  </span>
                ))}
              </div>
            </div>
          )}

          {/* Founders */}
          {org.founders?.length > 0 && (
            <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl overflow-hidden">
              <div className="px-5 py-4 border-b border-white/[0.06]">
                <h2 className="font-semibold text-white">Үүсгэн байгуулагчид</h2>
              </div>
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-white/[0.06]">
                    <th className="text-left px-5 py-2.5 text-slate-400 font-semibold">Нэр</th>
                    <th className="text-left px-5 py-2.5 text-slate-400 font-semibold">РД</th>
                    <th className="text-left px-5 py-2.5 text-slate-400 font-semibold">Төрөл</th>
                    <th className="text-right px-5 py-2.5 text-slate-400 font-semibold">Хувь</th>
                  </tr>
                </thead>
                <tbody>
                  {org.founders.map((f: any, i: number) => (
                    <tr key={i} className="border-b border-white/[0.03] hover:bg-white/[0.02]">
                      <td className="px-5 py-2.5 text-white">{f.name}</td>
                      <td className="px-5 py-2.5 font-mono text-xs text-slate-300">{f.reg_no}</td>
                      <td className="px-5 py-2.5 text-slate-400 text-xs">{f.type}</td>
                      <td className="px-5 py-2.5 text-right text-primary font-semibold">{f.share_percent}%</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}

          {/* Board members */}
          {org.stake_holders?.length > 0 && (
            <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl overflow-hidden">
              <div className="px-5 py-4 border-b border-white/[0.06]">
                <h2 className="font-semibold text-white">ТУЗ гишүүд</h2>
              </div>
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-white/[0.06]">
                    <th className="text-left px-5 py-2.5 text-slate-400 font-semibold">Нэр</th>
                    <th className="text-left px-5 py-2.5 text-slate-400 font-semibold">РД</th>
                    <th className="text-left px-5 py-2.5 text-slate-400 font-semibold">Албан тушаал</th>
                  </tr>
                </thead>
                <tbody>
                  {org.stake_holders.map((s: any, i: number) => (
                    <tr key={i} className="border-b border-white/[0.03] hover:bg-white/[0.02]">
                      <td className="px-5 py-2.5 text-white">{s.name}</td>
                      <td className="px-5 py-2.5 font-mono text-xs text-slate-300">{s.reg_no}</td>
                      <td className="px-5 py-2.5 text-slate-400">{s.position}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}

      {tab === "org" && result && !org && (
        <div className="p-6 text-center text-slate-500 bg-white/[0.03] border border-white/[0.06] rounded-2xl">
          Байгууллага олдсонгүй
        </div>
      )}
    </div>
  );
}

function Row({ label, value, mono, bold }: { label: string; value?: string; mono?: boolean; bold?: boolean }) {
  if (!value) return null;
  return (
    <div className="px-5 py-3 flex justify-between">
      <span className="text-slate-400 text-sm">{label}</span>
      <span className={`text-white text-sm text-right max-w-[65%] ${mono ? "font-mono" : ""} ${bold ? "font-semibold" : ""}`}>
        {value}
      </span>
    </div>
  );
}
