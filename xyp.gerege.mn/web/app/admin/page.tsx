import { listClients, getUsage } from "@/lib/api";

export default async function DashboardPage() {
  const [clients, usage] = await Promise.all([listClients(), getUsage()]);

  const activeClients = clients.filter((c: any) => c.active).length;
  const totalCalls = usage.reduce((sum: number, u: any) => sum + u.total_calls, 0);

  return (
    <div className="max-w-4xl mx-auto px-6 py-10">
      <h1 className="text-2xl font-bold text-white mb-2">Dashboard</h1>
      <p className="text-sm text-slate-400 mb-8">
        Verify API ашиглалтын тойм
      </p>

      <div className="mb-6 flex gap-3 flex-wrap">
        <a
          href="/admin/verify"
          className="inline-flex items-center gap-2 px-5 py-2.5 bg-primary hover:bg-primary-light text-white font-bold text-sm rounded-xl transition-all"
        >
          Баталгаажуулалт →
        </a>
        <a
          href="/admin/lookup"
          className="inline-flex items-center gap-2 px-5 py-2.5 bg-white/[0.05] hover:bg-white/[0.1] text-white font-bold text-sm rounded-xl transition-all"
        >
          Дэлгэрэнгүй хайлт →
        </a>
        <a
          href="/admin/clients"
          className="inline-flex items-center gap-2 px-5 py-2.5 bg-white/[0.05] hover:bg-white/[0.1] text-white font-bold text-sm rounded-xl transition-all"
        >
          Client удирдлага →
        </a>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-10">
        <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl p-6">
          <p className="text-sm text-slate-400 mb-1">Нийт Client</p>
          <p className="text-3xl font-bold text-white">{clients.length}</p>
        </div>
        <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl p-6">
          <p className="text-sm text-slate-400 mb-1">Идэвхтэй Client</p>
          <p className="text-3xl font-bold text-primary">{activeClients}</p>
        </div>
        <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl p-6">
          <p className="text-sm text-slate-400 mb-1">Нийт API дуудлага</p>
          <p className="text-3xl font-bold text-white">{totalCalls}</p>
        </div>
      </div>

      {usage.length > 0 && (
        <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl overflow-hidden">
          <div className="px-5 py-4 border-b border-white/[0.06]">
            <h2 className="font-semibold text-white">Хэрэглээ (endpoint-аар)</h2>
          </div>
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-white/[0.06]">
                <th className="text-left px-5 py-3 text-slate-400 font-semibold">Client</th>
                <th className="text-left px-5 py-3 text-slate-400 font-semibold">Endpoint</th>
                <th className="text-right px-5 py-3 text-slate-400 font-semibold">Дуудлага</th>
              </tr>
            </thead>
            <tbody>
              {usage.map((u: any, i: number) => (
                <tr key={i} className="border-b border-white/[0.03] hover:bg-white/[0.02]">
                  <td className="px-5 py-3 text-white">{u.client_name}</td>
                  <td className="px-5 py-3 font-mono text-xs text-primary">{u.endpoint}</td>
                  <td className="px-5 py-3 text-right text-white font-medium">{u.total_calls}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
