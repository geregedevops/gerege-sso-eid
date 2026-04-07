import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";
import Link from "next/link";

export default async function DashboardPage() {
  const session = await auth();
  const sub = (session?.user as any)?.sub;

  const developer = sub ? await prisma.developer.findUnique({
    where: { sub },
    include: { apps: true, tenants: { include: { tenant: true } } },
  }) : null;

  const appCount = developer?.apps.length || 0;
  const tenantCount = developer?.tenants.length || 0;

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-white">Dashboard</h1>
        <p className="text-slate-400 text-sm mt-1">
          Сайн байна уу, {session?.user?.name || "Developer"}
        </p>
      </div>

      <div className="grid sm:grid-cols-4 gap-4">
        <StatCard label="Миний App-ууд" value={appCount} />
        <StatCard label="Tenants" value={tenantCount} />
        <StatCard label="Нийт Scopes" value={developer?.apps.reduce((a, b) => a + b.scopes.length, 0) || 0} />
        <StatCard label="Active" value={developer?.apps.filter(a => a.isActive).length || 0} />
      </div>

      <div className="flex gap-4">
        <Link
          href="/dashboard/apps/new"
          className="px-5 py-2.5 bg-primary text-white font-semibold rounded-xl text-sm hover:bg-primary-light transition-colors"
        >
          + Шинэ App нэмэх
        </Link>
        <Link
          href="/dashboard/tenants/new"
          className="px-5 py-2.5 border border-primary/30 text-primary font-medium rounded-xl text-sm hover:bg-primary/10 transition-colors"
        >
          + Tenant үүсгэх
        </Link>
        <Link
          href="/docs/quickstart"
          className="px-5 py-2.5 border border-white/15 text-white font-medium rounded-xl text-sm hover:bg-white/5 transition-colors"
        >
          Quickstart
        </Link>
      </div>

      {appCount > 0 && (
        <div>
          <h2 className="text-lg font-semibold text-white mb-3">Apps</h2>
          <div className="space-y-3">
            {developer?.apps.map((app) => (
              <Link
                key={app.id}
                href={`/dashboard/apps/${app.id}`}
                className="block bg-surface border border-white/10 rounded-xl p-4 hover:border-primary/30 transition-colors"
              >
                <div className="flex items-center justify-between">
                  <div>
                    <div className="flex items-center gap-2">
                      <span className={`w-2 h-2 rounded-full ${app.isActive ? 'bg-green-500' : 'bg-red-500'}`} />
                      <span className="font-medium text-white">{app.name}</span>
                    </div>
                    <p className="text-xs text-slate-500 mt-1 font-mono">{app.clientId}</p>
                  </div>
                  <div className="text-xs text-slate-500">
                    {app.scopes.join(", ")}
                  </div>
                </div>
              </Link>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

function StatCard({ label, value }: { label: string; value: number }) {
  return (
    <div className="bg-surface border border-white/10 rounded-xl p-5">
      <p className="text-xs text-slate-500 uppercase tracking-wider">{label}</p>
      <p className="text-3xl font-bold text-white mt-1">{value}</p>
    </div>
  );
}
