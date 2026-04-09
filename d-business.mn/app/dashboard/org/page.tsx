import Link from "next/link";
import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";

export default async function OrgListPage() {
  const session = await auth();
  const sub = (session?.user as any)?.sub;
  const user = sub ? await prisma.user.findUnique({ where: { sub }, include: { memberships: { include: { organization: true } } } }) : null;
  const orgs = user?.memberships || [];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-white">Байгууллагууд</h1>
        <Link href="/dashboard/org/new" className="px-4 py-2 bg-primary text-white text-sm font-semibold rounded-lg hover:bg-primary-light transition-colors">
          + Шинэ
        </Link>
      </div>

      {orgs.length === 0 ? (
        <div className="bg-surface border border-white/10 rounded-xl p-12 text-center">
          <p className="text-slate-400 mb-4">Байгууллага бүртгэгдээгүй байна.</p>
          <Link href="/dashboard/org/new" className="text-primary hover:underline text-sm">Байгууллага бүртгүүлэх</Link>
        </div>
      ) : (
        <div className="grid sm:grid-cols-2 gap-4">
          {orgs.map((m) => (
            <Link key={m.organizationId} href={`/dashboard/org/${m.organizationId}`}
              className="bg-surface border border-white/10 rounded-xl p-5 hover:border-primary/30 transition-colors block">
              <div className="flex items-center gap-3 mb-2">
                <div className={`w-2 h-2 rounded-full ${m.organization.isActive ? "bg-green-400" : "bg-slate-500"}`} />
                <h3 className="font-semibold text-white">{m.organization.name}</h3>
                {m.organization.isVerified && <span className="text-xs px-2 py-0.5 bg-green-500/10 text-green-400 rounded-full">Баталгаажсан</span>}
              </div>
              <p className="text-xs text-slate-400">РД: {m.organization.registrationNumber} &middot; {m.organization.type} &middot; {m.role}</p>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
