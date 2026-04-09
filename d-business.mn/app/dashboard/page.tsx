import Link from "next/link";
import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";

export default async function DashboardPage() {
  const session = await auth();
  const sub = (session?.user as any)?.sub;
  if (!sub) return null;

  const user = await prisma.user.findUnique({ where: { sub }, include: { memberships: { include: { organization: true } } } });
  if (!user) return null;

  const orgCount = user.memberships.length;
  const docCount = await prisma.document.count({ where: { uploadedById: user.id } });
  const sigCount = await prisma.signature.count({ where: { signedById: user.id, status: "complete" } });
  const certCount = await prisma.certificate.count({
    where: { organizationId: { in: user.memberships.map((m) => m.organizationId) }, status: "active" },
  });

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-white">Сайн байна уу, {user.givenName || user.name}</h1>
        <p className="text-sm text-slate-400 mt-1">Байгууллагын цахим тамга платформ</p>
      </div>

      <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Stat label="Байгууллагууд" value={orgCount} />
        <Stat label="Баримтууд" value={docCount} />
        <Stat label="Гарын үсэг" value={sigCount} />
        <Stat label="Сертификат" value={certCount} />
      </div>

      <div className="flex gap-3">
        <Link href="/dashboard/org/new" className="px-4 py-2 bg-primary text-white text-sm font-semibold rounded-lg hover:bg-primary-light transition-colors">
          + Байгууллага бүртгүүлэх
        </Link>
        <Link href="/dashboard/sign" className="px-4 py-2 border border-white/15 text-white text-sm font-medium rounded-lg hover:bg-white/5 transition-colors">
          Гарын үсэг зурах
        </Link>
      </div>

      {user.memberships.length > 0 && (
        <div>
          <h2 className="text-lg font-bold text-white mb-3">Миний байгууллагууд</h2>
          <div className="grid sm:grid-cols-2 gap-4">
            {user.memberships.map((m) => (
              <Link key={m.organizationId} href={`/dashboard/org/${m.organizationId}`}
                className="bg-surface border border-white/10 rounded-xl p-5 hover:border-primary/30 transition-colors block">
                <div className="flex items-center gap-3 mb-2">
                  <div className={`w-2 h-2 rounded-full ${m.organization.isActive ? "bg-green-400" : "bg-slate-500"}`} />
                  <h3 className="font-semibold text-white">{m.organization.name}</h3>
                </div>
                <p className="text-xs text-slate-400">РД: {m.organization.registrationNumber} &middot; {m.organization.type} &middot; {m.role}</p>
              </Link>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

function Stat({ label, value }: { label: string; value: number }) {
  return (
    <div className="bg-surface border border-white/10 rounded-xl p-5">
      <p className="text-2xl font-bold text-white">{value}</p>
      <p className="text-xs text-slate-400 mt-1">{label}</p>
    </div>
  );
}
