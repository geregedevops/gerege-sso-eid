import { notFound } from "next/navigation";
import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";

export default async function CertificatesPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const session = await auth();
  const sub = (session?.user as any)?.sub;
  const user = sub ? await prisma.user.findUnique({ where: { sub } }) : null;
  if (!user) notFound();

  const org = await prisma.organization.findUnique({ where: { id }, include: { certificates: { orderBy: { createdAt: "desc" } } } });
  if (!org) notFound();

  const statusColors: Record<string, string> = {
    active: "bg-green-500/10 text-green-400",
    pending: "bg-yellow-500/10 text-yellow-400",
    revoked: "bg-red-500/10 text-red-400",
    expired: "bg-slate-500/10 text-slate-400",
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-white">{org.name} — Сертификат</h1>
      </div>

      {org.certificates.length === 0 ? (
        <div className="bg-surface border border-white/10 rounded-xl p-12 text-center">
          <p className="text-slate-400">Сертификат байхгүй байна.</p>
        </div>
      ) : (
        <div className="space-y-3">
          {org.certificates.map((c) => (
            <div key={c.id} className="bg-surface border border-white/10 rounded-xl p-5">
              <div className="flex items-center justify-between mb-2">
                <h3 className="font-semibold text-white">{c.commonName}</h3>
                <span className={`text-xs px-2 py-0.5 rounded-full ${statusColors[c.status] || statusColors.pending}`}>{c.status}</span>
              </div>
              <p className="text-xs text-slate-400">
                {c.purpose} &middot; {c.serialNumber || "Serial pending"}
                {c.issuedAt ? ` &middot; Олгосон: ${c.issuedAt.toISOString().split("T")[0]}` : ""}
                {c.expiresAt ? ` &middot; Дуусах: ${c.expiresAt.toISOString().split("T")[0]}` : ""}
              </p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
