import { getUsage } from "@/lib/api";

export default async function UsagePage({
  searchParams,
}: {
  searchParams: { client_id?: string; from?: string; to?: string };
}) {
  const usage = await getUsage(searchParams);

  const totalCalls = usage.reduce((sum: number, u: any) => sum + u.total_calls, 0);

  return (
    <div className="max-w-4xl mx-auto px-6 py-10">
      <h1 className="text-2xl font-bold text-white mb-2">API хэрэглээ</h1>
      <p className="text-sm text-slate-400 mb-8">
        Баталгаажуулах API-ийн дуудлагын дэлгэрэнгүй статистик
      </p>

      <form className="flex gap-3 mb-8 flex-wrap">
        <input
          name="client_id"
          placeholder="Client ID (бүгд)"
          defaultValue={searchParams.client_id || ""}
          className="px-4 py-2 bg-white/[0.05] border border-white/10 rounded-xl text-white text-sm focus:outline-none focus:border-primary"
        />
        <input
          name="from"
          type="date"
          defaultValue={searchParams.from || ""}
          className="px-4 py-2 bg-white/[0.05] border border-white/10 rounded-xl text-white text-sm focus:outline-none focus:border-primary"
        />
        <input
          name="to"
          type="date"
          defaultValue={searchParams.to || ""}
          className="px-4 py-2 bg-white/[0.05] border border-white/10 rounded-xl text-white text-sm focus:outline-none focus:border-primary"
        />
        <button
          type="submit"
          className="px-5 py-2 bg-primary hover:bg-primary-light text-white font-bold text-sm rounded-xl transition-all"
        >
          Шүүх
        </button>
      </form>

      <div className="mb-6">
        <p className="text-sm text-slate-400">
          Нийт дуудлага: <span className="text-white font-bold">{totalCalls}</span>
        </p>
      </div>

      <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-white/[0.06]">
              <th className="text-left px-5 py-3 text-slate-400 font-semibold">Client</th>
              <th className="text-left px-5 py-3 text-slate-400 font-semibold">Endpoint</th>
              <th className="text-right px-5 py-3 text-slate-400 font-semibold">Дуудлага</th>
            </tr>
          </thead>
          <tbody>
            {usage.length === 0 ? (
              <tr>
                <td colSpan={3} className="px-5 py-10 text-center text-slate-500">
                  Мэдээлэл олдсонгүй
                </td>
              </tr>
            ) : (
              usage.map((u: any, i: number) => (
                <tr key={i} className="border-b border-white/[0.03] hover:bg-white/[0.02]">
                  <td className="px-5 py-3 text-white">{u.client_name}</td>
                  <td className="px-5 py-3 font-mono text-xs text-primary">{u.endpoint}</td>
                  <td className="px-5 py-3 text-right text-white font-medium">{u.total_calls}</td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
