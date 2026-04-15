"use client";

import { useState } from "react";

type CitizenResult = {
  found: boolean;
  citizen: {
    reg_no: string;
    last_name: string;
    first_name: string;
    surname: string;
    gender: string;
    birth_date: string;
    nationality: string;
  } | null;
};

type OrgFounder = {
  name: string;
  reg_no: string;
  type: string;
  share_percent: string;
};

type OrgStakeHolder = {
  name: string;
  reg_no: string;
  position: string;
};

type OrgResult = {
  found: boolean;
  organization: {
    reg_no: string;
    name: string;
    type: string;
    capital: string;
    ceo: string;
    ceo_reg_no: string;
    ceo_position: string;
    phone: string;
    address: string;
    industry: string[];
    founders: OrgFounder[];
    stake_holders: OrgStakeHolder[];
  } | null;
};

export default function HomePage() {
  const [tab, setTab] = useState<"citizen" | "org">("citizen");
  const [regNo, setRegNo] = useState("");
  const [loading, setLoading] = useState(false);
  const [citizenResult, setCitizenResult] = useState<CitizenResult | null>(null);
  const [orgResult, setOrgResult] = useState<OrgResult | null>(null);
  const [error, setError] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!regNo.trim()) return;

    setLoading(true);
    setError("");
    setCitizenResult(null);
    setOrgResult(null);

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
      } else if (tab === "citizen") {
        setCitizenResult(data);
      } else {
        setOrgResult(data);
      }
    } catch {
      setError("Сервертэй холбогдож чадсангүй");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div>
      {/* Hero */}
      <section className="max-w-4xl mx-auto px-6 pt-16 pb-12 text-center">
        <div className="inline-flex items-center gap-2 px-4 py-1.5 bg-primary/10 border border-primary/20 rounded-full text-primary text-sm font-medium mb-6">
          Gerege Verify API
        </div>
        <h1 className="text-4xl md:text-5xl font-bold text-white mb-4 leading-tight">
          Иргэн & Байгууллагын<br />мэдээлэл баталгаажуулах
        </h1>
        <p className="text-lg text-slate-400 max-w-2xl mx-auto mb-10">
          Регистрийн дугаараар иргэн болон байгууллагын мэдээллийг шууд тулгаж шалгана.
          REST API-аар 3-р талын системүүд холбогдож ашиглах боломжтой.
        </p>

        {/* Feature cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-16 text-left">
          <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl p-5">
            <div className="text-2xl mb-2">🔍</div>
            <h3 className="text-white font-semibold mb-1">Иргэн шалгах</h3>
            <p className="text-sm text-slate-400">Регистрийн дугаараар иргэний нэр, овог, хүйс, төрсөн он зэрэг мэдээллийг авна</p>
          </div>
          <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl p-5">
            <div className="text-2xl mb-2">🏢</div>
            <h3 className="text-white font-semibold mb-1">Байгууллага шалгах</h3>
            <p className="text-sm text-slate-400">Регистрийн дугаараар байгууллагын нэр, захирал, утас, хаяг, үйл ажиллагааны чиглэл</p>
          </div>
          <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl p-5">
            <div className="text-2xl mb-2">🔐</div>
            <h3 className="text-white font-semibold mb-1">API хандалт</h3>
            <p className="text-sm text-slate-400">Client ID + Secret-ээр холбогдож, Basic Auth-аар баталгаажуулж ашиглана</p>
          </div>
        </div>
      </section>

      {/* Try it */}
      <section id="try" className="max-w-2xl mx-auto px-6 pb-20">
        <h2 className="text-2xl font-bold text-white mb-6 text-center">Туршиж үзэх</h2>

        {/* Tabs */}
        <div className="flex gap-2 mb-6 justify-center">
          <button
            onClick={() => { setTab("citizen"); setCitizenResult(null); setOrgResult(null); setError(""); }}
            className={`px-5 py-2 rounded-xl text-sm font-medium transition-all ${
              tab === "citizen"
                ? "bg-primary text-white"
                : "bg-white/[0.05] text-slate-400 hover:text-white"
            }`}
          >
            Иргэн
          </button>
          <button
            onClick={() => { setTab("org"); setCitizenResult(null); setOrgResult(null); setError(""); }}
            className={`px-5 py-2 rounded-xl text-sm font-medium transition-all ${
              tab === "org"
                ? "bg-primary text-white"
                : "bg-white/[0.05] text-slate-400 hover:text-white"
            }`}
          >
            Байгууллага
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="flex gap-3">
          <input
            type="text"
            value={regNo}
            onChange={(e) => setRegNo(e.target.value)}
            placeholder={tab === "citizen" ? "Регистрийн дугаар (жнь: МА74101813)" : "Регистрийн дугаар (жнь: 6235972)"}
            className="flex-1 px-4 py-3 bg-white/[0.05] border border-white/10 rounded-xl text-white text-sm focus:outline-none focus:border-primary"
          />
          <button
            type="submit"
            disabled={loading}
            className="px-6 py-3 bg-primary hover:bg-primary-light text-white font-bold text-sm rounded-xl transition-all disabled:opacity-50"
          >
            {loading ? "Хайж байна..." : "Шалгах"}
          </button>
        </form>

        {/* Error */}
        {error && (
          <div className="mt-4 p-4 bg-red-500/10 border border-red-500/20 rounded-xl">
            <p className="text-sm text-red-400">{error}</p>
          </div>
        )}

        {/* Citizen result */}
        {citizenResult && (
          <div className="mt-6 bg-white/[0.03] border border-white/[0.06] rounded-2xl overflow-hidden">
            {!citizenResult.found ? (
              <div className="p-6 text-center text-slate-500">Иргэн олдсонгүй</div>
            ) : citizenResult.citizen && (
              <div className="divide-y divide-white/[0.06]">
                <div className="px-5 py-3 flex justify-between">
                  <span className="text-slate-400 text-sm">Регистр</span>
                  <span className="text-white font-mono text-sm">{citizenResult.citizen.reg_no}</span>
                </div>
                <div className="px-5 py-3 flex justify-between">
                  <span className="text-slate-400 text-sm">Овог</span>
                  <span className="text-white">{citizenResult.citizen.last_name}</span>
                </div>
                <div className="px-5 py-3 flex justify-between">
                  <span className="text-slate-400 text-sm">Нэр</span>
                  <span className="text-white font-semibold">{citizenResult.citizen.first_name}</span>
                </div>
                {citizenResult.citizen.surname && (
                  <div className="px-5 py-3 flex justify-between">
                    <span className="text-slate-400 text-sm">Ургийн овог</span>
                    <span className="text-white">{citizenResult.citizen.surname}</span>
                  </div>
                )}
                <div className="px-5 py-3 flex justify-between">
                  <span className="text-slate-400 text-sm">Хүйс</span>
                  <span className="text-white">{citizenResult.citizen.gender}</span>
                </div>
                <div className="px-5 py-3 flex justify-between">
                  <span className="text-slate-400 text-sm">Төрсөн огноо</span>
                  <span className="text-white">{citizenResult.citizen.birth_date}</span>
                </div>
                <div className="px-5 py-3 flex justify-between">
                  <span className="text-slate-400 text-sm">Үндэс угсаа</span>
                  <span className="text-white">{citizenResult.citizen.nationality}</span>
                </div>
              </div>
            )}
          </div>
        )}

        {/* Org result */}
        {orgResult && (
          <div className="mt-6 bg-white/[0.03] border border-white/[0.06] rounded-2xl overflow-hidden">
            {!orgResult.found ? (
              <div className="p-6 text-center text-slate-500">Байгууллага олдсонгүй</div>
            ) : orgResult.organization && (
              <div className="divide-y divide-white/[0.06]">
                <div className="px-5 py-3 flex justify-between">
                  <span className="text-slate-400 text-sm">Регистр</span>
                  <span className="text-white font-mono text-sm">{orgResult.organization.reg_no}</span>
                </div>
                <div className="px-5 py-3 flex justify-between">
                  <span className="text-slate-400 text-sm">Нэр</span>
                  <span className="text-white font-semibold">{orgResult.organization.name}</span>
                </div>
                <div className="px-5 py-3 flex justify-between">
                  <span className="text-slate-400 text-sm">Төрөл</span>
                  <span className="text-white text-sm">{orgResult.organization.type}</span>
                </div>
                {orgResult.organization.capital && (
                  <div className="px-5 py-3 flex justify-between">
                    <span className="text-slate-400 text-sm">Дүрмийн сан</span>
                    <span className="text-white">{Number(orgResult.organization.capital).toLocaleString()}₮</span>
                  </div>
                )}
                {orgResult.organization.ceo && (
                  <div className="px-5 py-3 flex justify-between">
                    <span className="text-slate-400 text-sm">
                      {orgResult.organization.ceo_position || "Захирал"}
                    </span>
                    <span className="text-white">{orgResult.organization.ceo}</span>
                  </div>
                )}
                {orgResult.organization.ceo_reg_no && (
                  <div className="px-5 py-3 flex justify-between">
                    <span className="text-slate-400 text-sm">Захирлын РД</span>
                    <span className="text-white font-mono text-sm">{orgResult.organization.ceo_reg_no}</span>
                  </div>
                )}
                {orgResult.organization.phone && (
                  <div className="px-5 py-3 flex justify-between">
                    <span className="text-slate-400 text-sm">Утас</span>
                    <span className="text-white">{orgResult.organization.phone}</span>
                  </div>
                )}
                {orgResult.organization.address && (
                  <div className="px-5 py-3 flex justify-between">
                    <span className="text-slate-400 text-sm">Хаяг</span>
                    <span className="text-white text-sm text-right max-w-[60%]">{orgResult.organization.address}</span>
                  </div>
                )}
                {orgResult.organization.industry?.length > 0 && (
                  <div className="px-5 py-3">
                    <span className="text-slate-400 text-sm block mb-2">Үйл ажиллагааны чиглэл</span>
                    <div className="flex flex-wrap gap-2">
                      {orgResult.organization.industry.map((ind, i) => (
                        <span key={i} className="px-3 py-1 bg-primary/10 text-primary text-xs rounded-lg">
                          {ind}
                        </span>
                      ))}
                    </div>
                  </div>
                )}
                {orgResult.organization.founders?.length > 0 && (
                  <div className="px-5 py-3">
                    <span className="text-slate-400 text-sm block mb-2">Үүсгэн байгуулагчид</span>
                    <div className="space-y-2">
                      {orgResult.organization.founders.map((f, i) => (
                        <div key={i} className="flex items-center justify-between bg-white/[0.02] rounded-lg px-3 py-2">
                          <div>
                            <span className="text-white text-sm">{f.name}</span>
                            <span className="text-slate-500 text-xs ml-2">({f.type})</span>
                            <span className="text-slate-500 font-mono text-xs ml-2">{f.reg_no}</span>
                          </div>
                          <span className="text-primary font-semibold text-sm">{f.share_percent}%</span>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
                {orgResult.organization.stake_holders?.length > 0 && (
                  <div className="px-5 py-3">
                    <span className="text-slate-400 text-sm block mb-2">ТУЗ гишүүд</span>
                    <div className="space-y-2">
                      {orgResult.organization.stake_holders.map((s, i) => (
                        <div key={i} className="flex items-center justify-between bg-white/[0.02] rounded-lg px-3 py-2">
                          <div>
                            <span className="text-white text-sm">{s.name}</span>
                            <span className="text-slate-500 font-mono text-xs ml-2">{s.reg_no}</span>
                          </div>
                          <span className="text-slate-400 text-xs">{s.position}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        )}
      </section>

      {/* API section */}
      <section className="max-w-4xl mx-auto px-6 pb-20">
        <h2 className="text-2xl font-bold text-white mb-6 text-center">API ашиглах</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl p-5">
            <h3 className="text-white font-semibold mb-3">1. Client авах</h3>
            <p className="text-sm text-slate-400 mb-3">
              Admin-аас client_id болон client_secret авна. Secret зөвхөн нэг удаа харагдана.
            </p>
            <code className="block bg-black/30 rounded-lg p-3 text-xs text-slate-300 font-mono overflow-x-auto">
              client_id: vfy_abc123...<br />
              client_secret: YLash...
            </code>
          </div>
          <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl p-5">
            <h3 className="text-white font-semibold mb-3">2. API дуудах</h3>
            <p className="text-sm text-slate-400 mb-3">
              HTTP Basic Auth ашиглан POST хүсэлт илгээнэ.
            </p>
            <code className="block bg-black/30 rounded-lg p-3 text-xs text-green-400 font-mono overflow-x-auto whitespace-pre">{`curl -u client_id:secret \\
  -X POST /v1/org/lookup \\
  -H "Content-Type: application/json" \\
  -d '{"reg_no":"6235972"}'`}</code>
          </div>
        </div>

        <div className="mt-6 bg-white/[0.03] border border-white/[0.06] rounded-2xl p-5">
          <h3 className="text-white font-semibold mb-3">Endpoint-ууд</h3>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-white/[0.06]">
                  <th className="text-left py-2 pr-4 text-slate-400 font-semibold">Method</th>
                  <th className="text-left py-2 pr-4 text-slate-400 font-semibold">Endpoint</th>
                  <th className="text-left py-2 text-slate-400 font-semibold">Тайлбар</th>
                </tr>
              </thead>
              <tbody className="text-slate-300">
                <tr className="border-b border-white/[0.03]">
                  <td className="py-2 pr-4"><span className="px-2 py-0.5 bg-blue-500/10 text-blue-400 text-xs rounded font-mono">POST</span></td>
                  <td className="py-2 pr-4 font-mono text-xs text-primary">/v1/citizen/lookup</td>
                  <td className="py-2 text-sm">Иргэний мэдээлэл хайх</td>
                </tr>
                <tr className="border-b border-white/[0.03]">
                  <td className="py-2 pr-4"><span className="px-2 py-0.5 bg-blue-500/10 text-blue-400 text-xs rounded font-mono">POST</span></td>
                  <td className="py-2 pr-4 font-mono text-xs text-primary">/v1/citizen/verify</td>
                  <td className="py-2 text-sm">Иргэний нэр тулгах</td>
                </tr>
                <tr className="border-b border-white/[0.03]">
                  <td className="py-2 pr-4"><span className="px-2 py-0.5 bg-blue-500/10 text-blue-400 text-xs rounded font-mono">POST</span></td>
                  <td className="py-2 pr-4 font-mono text-xs text-primary">/v1/org/lookup</td>
                  <td className="py-2 text-sm">Байгууллагын мэдээлэл хайх</td>
                </tr>
                <tr>
                  <td className="py-2 pr-4"><span className="px-2 py-0.5 bg-blue-500/10 text-blue-400 text-xs rounded font-mono">POST</span></td>
                  <td className="py-2 pr-4 font-mono text-xs text-primary">/v1/org/verify</td>
                  <td className="py-2 text-sm">Байгууллагын нэр тулгах</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-white/5 py-8 text-center text-sm text-slate-500">
        Gerege Systems &copy; {new Date().getFullYear()} &middot; xyp.gerege.mn
      </footer>
    </div>
  );
}
