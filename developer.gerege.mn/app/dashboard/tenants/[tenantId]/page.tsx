import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";
import { notFound } from "next/navigation";
import Link from "next/link";

export default async function TenantDetailPage({ params }: { params: { tenantId: string } }) {
  const session = await auth();
  const sub = (session?.user as any)?.sub;

  const developer = sub ? await prisma.developer.findUnique({ where: { sub } }) : null;
  if (!developer) notFound();

  const membership = await prisma.tenantMember.findUnique({
    where: { tenantId_developerId: { tenantId: params.tenantId, developerId: developer.id } },
    include: {
      tenant: { include: { members: { include: { developer: true } }, apps: true } },
    },
  });
  if (!membership) notFound();

  const tenant = membership.tenant;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white">{tenant.name}</h1>
          <p className="text-slate-500 text-sm font-mono mt-1">{tenant.slug}</p>
        </div>
        <Link href="/dashboard/tenants" className="text-sm text-slate-400 hover:text-white">
          Буцах
        </Link>
      </div>

      <div className="grid sm:grid-cols-3 gap-4">
        <div className="bg-surface border border-white/10 rounded-xl p-5">
          <p className="text-xs text-slate-500 uppercase">Plan</p>
          <p className="text-xl font-bold text-white mt-1 capitalize">{tenant.plan}</p>
        </div>
        <div className="bg-surface border border-white/10 rounded-xl p-5">
          <p className="text-xs text-slate-500 uppercase">Members</p>
          <p className="text-xl font-bold text-white mt-1">{tenant.members.length}</p>
        </div>
        <div className="bg-surface border border-white/10 rounded-xl p-5">
          <p className="text-xs text-slate-500 uppercase">Apps</p>
          <p className="text-xl font-bold text-white mt-1">{tenant.apps.length}</p>
        </div>
      </div>

      <Section title="Members">
        {tenant.members.map((m) => (
          <div key={m.developerId} className="flex items-center justify-between py-2">
            <span className="text-sm text-white">{m.developer.name || m.developer.sub}</span>
            <span className="px-2 py-1 bg-primary/10 text-primary text-xs font-medium rounded-full">{m.role}</span>
          </div>
        ))}
      </Section>

      {tenant.apps.length > 0 && (
        <Section title="Apps">
          {tenant.apps.map((app) => (
            <Link
              key={app.id}
              href={`/dashboard/apps/${app.id}`}
              className="block py-2 text-sm text-white hover:text-primary"
            >
              {app.name} <span className="text-slate-500 font-mono text-xs">({app.clientId})</span>
            </Link>
          ))}
        </Section>
      )}
    </div>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="bg-surface border border-white/10 rounded-xl p-5">
      <h2 className="text-sm font-semibold text-slate-400 uppercase tracking-wider mb-3">{title}</h2>
      {children}
    </div>
  );
}
