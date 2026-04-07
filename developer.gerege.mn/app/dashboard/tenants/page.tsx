import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";
import Link from "next/link";

export default async function TenantsPage() {
  const session = await auth();
  const sub = (session?.user as any)?.sub;

  const developer = sub ? await prisma.developer.findUnique({
    where: { sub },
    include: { tenants: { include: { tenant: true } } },
  }) : null;

  const memberships = developer?.tenants || [];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-white">Tenants</h1>
        <Link
          href="/dashboard/tenants/new"
          className="px-4 py-2 bg-primary text-white font-semibold rounded-lg text-sm hover:bg-primary-light transition-colors"
        >
          + Шинэ Tenant
        </Link>
      </div>

      {memberships.length === 0 ? (
        <div className="bg-surface border border-white/10 rounded-xl p-12 text-center">
          <p className="text-slate-400 mb-4">Одоогоор tenant үүсгээгүй байна.</p>
          <Link
            href="/dashboard/tenants/new"
            className="px-5 py-2.5 bg-primary text-white font-semibold rounded-xl text-sm"
          >
            Tenant үүсгэх
          </Link>
        </div>
      ) : (
        <div className="space-y-3">
          {memberships.map((m) => (
            <Link
              key={m.tenantId}
              href={`/dashboard/tenants/${m.tenantId}`}
              className="block bg-surface border border-white/10 rounded-xl p-5 hover:border-primary/30 transition-colors"
            >
              <div className="flex items-center justify-between">
                <div>
                  <span className="font-semibold text-white">{m.tenant.name}</span>
                  <p className="text-xs text-slate-500 font-mono mt-1">{m.tenant.slug}</p>
                </div>
                <div className="text-right">
                  <span className="px-2 py-1 bg-primary/10 text-primary text-xs font-medium rounded-full">{m.role}</span>
                  <p className="text-xs text-slate-500 mt-1">Plan: {m.tenant.plan}</p>
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
