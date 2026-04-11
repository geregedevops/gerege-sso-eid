import { listDANClients, deactivateDANClient } from "@/lib/api";

export default async function AdminPage() {
  const clients = await listDANClients();

  return (
    <div className="max-w-4xl mx-auto px-6 py-10">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-white">DAN Clients</h1>
          <p className="text-sm text-slate-400 mt-1">
            Бүртгэлтэй DAN verify client-ийн жагсаалт
          </p>
        </div>
        <a
          href="/admin/clients/new"
          className="px-5 py-2.5 bg-primary hover:bg-primary-light text-white font-bold text-sm rounded-xl transition-all"
        >
          + Шинэ Client
        </a>
      </div>

      <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-white/[0.06]">
              <th className="text-left px-5 py-3 text-slate-400 font-semibold">ID</th>
              <th className="text-left px-5 py-3 text-slate-400 font-semibold">Нэр</th>
              <th className="text-left px-5 py-3 text-slate-400 font-semibold">Callback URLs</th>
              <th className="text-left px-5 py-3 text-slate-400 font-semibold">Төлөв</th>
              <th className="text-right px-5 py-3 text-slate-400 font-semibold">Үйлдэл</th>
            </tr>
          </thead>
          <tbody>
            {clients.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-5 py-10 text-center text-slate-500">
                  Client бүртгэгдээгүй байна
                </td>
              </tr>
            ) : (
              clients.map((c: any) => (
                <tr key={c.id} className="border-b border-white/[0.03] hover:bg-white/[0.02]">
                  <td className="px-5 py-3 font-mono text-xs text-primary">{c.id}</td>
                  <td className="px-5 py-3 text-white">{c.name}</td>
                  <td className="px-5 py-3 text-slate-400 text-xs">
                    {c.callback_urls?.join(", ")}
                  </td>
                  <td className="px-5 py-3">
                    {c.active ? (
                      <span className="px-2 py-0.5 bg-green-500/10 text-green-400 text-xs rounded-md font-medium">
                        Идэвхтэй
                      </span>
                    ) : (
                      <span className="px-2 py-0.5 bg-red-500/10 text-red-400 text-xs rounded-md font-medium">
                        Идэвхгүй
                      </span>
                    )}
                  </td>
                  <td className="px-5 py-3 text-right">
                    {c.active && (
                      <form
                        action={async () => {
                          "use server";
                          await deactivateDANClient(c.id);
                          const { redirect } = await import("next/navigation");
                          redirect("/admin");
                        }}
                      >
                        <button
                          type="submit"
                          className="text-xs text-red-400 hover:text-red-300 font-medium"
                        >
                          Идэвхгүйжүүлэх
                        </button>
                      </form>
                    )}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
