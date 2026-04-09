import Link from "next/link";
import { notFound } from "next/navigation";
import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";

export default async function OrgDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const session = await auth();
  const sub = (session?.user as any)?.sub;
  const user = sub ? await prisma.user.findUnique({ where: { sub } }) : null;
  if (!user) notFound();

  const membership = await prisma.orgMember.findUnique({
    where: { organizationId_userId: { organizationId: id, userId: user.id } },
    include: { organization: { include: { members: { include: { user: true } }, certificates: true, documents: { take: 5, orderBy: { createdAt: "desc" }, include: { signatures: true } } } } },
  });
  if (!membership) notFound();

  const org = membership.organization;

  return (
    <div className="space-y-8">
      <div>
        <div className="flex items-center gap-3 mb-2">
          <h1 className="text-2xl font-bold text-white">{org.name}</h1>
          {org.isVerified && <span className="text-xs px-2 py-0.5 bg-green-500/10 text-green-400 rounded-full">Баталгаажсан</span>}
          <span className="text-xs px-2 py-0.5 bg-primary/10 text-primary rounded-full">{membership.role}</span>
        </div>
        <p className="text-sm text-slate-400">РД: {org.registrationNumber} &middot; {org.type} {org.address ? `&middot; ${org.address}` : ""}</p>
      </div>

      <div className="grid sm:grid-cols-3 gap-4">
        <Stat label="Гишүүд" value={org.members.length} href={`/dashboard/org/${id}/members`} />
        <Stat label="Сертификат" value={org.certificates.length} href={`/dashboard/org/${id}/certificates`} />
        <Stat label="Баримтууд" value={org.documents.length} />
      </div>

      <div className="flex gap-3">
        <Link href={`/dashboard/sign?org=${id}`} className="px-4 py-2 bg-primary text-white text-sm font-semibold rounded-lg hover:bg-primary-light transition-colors">
          Гарын үсэг зурах
        </Link>
        <Link href={`/dashboard/org/${id}/members`} className="px-4 py-2 border border-white/15 text-white text-sm rounded-lg hover:bg-white/5 transition-colors">
          Гишүүд удирдах
        </Link>
        <Link href={`/dashboard/org/${id}/certificates`} className="px-4 py-2 border border-white/15 text-white text-sm rounded-lg hover:bg-white/5 transition-colors">
          Сертификат
        </Link>
      </div>

      {org.members.length > 0 && (
        <div>
          <h2 className="text-lg font-bold text-white mb-3">Гишүүд</h2>
          <div className="bg-surface border border-white/10 rounded-xl overflow-hidden">
            <table className="w-full text-sm">
              <thead><tr className="border-b border-white/5"><th className="text-left px-4 py-3 text-slate-400 font-medium">Нэр</th><th className="text-left px-4 py-3 text-slate-400 font-medium">Үүрэг</th></tr></thead>
              <tbody>
                {org.members.map((m) => (
                  <tr key={m.userId} className="border-b border-white/5">
                    <td className="px-4 py-3 text-white">{m.user.name}</td>
                    <td className="px-4 py-3 text-slate-400">{m.role}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {org.documents.length > 0 && (
        <div>
          <h2 className="text-lg font-bold text-white mb-3">Сүүлийн баримтууд</h2>
          <div className="space-y-2">
            {org.documents.map((d) => (
              <Link key={d.id} href={`/dashboard/documents/${d.id}`}
                className="bg-surface border border-white/10 rounded-lg p-4 flex items-center justify-between hover:border-primary/30 transition-colors block">
                <div>
                  <p className="text-white text-sm font-medium">{d.name}</p>
                  <p className="text-xs text-slate-400">{d.fileName} &middot; {(d.fileSize / 1024).toFixed(0)} KB</p>
                </div>
                <span className={`text-xs px-2 py-1 rounded-full ${d.status === "signed" ? "bg-green-500/10 text-green-400" : d.status === "signing" ? "bg-yellow-500/10 text-yellow-400" : "bg-slate-500/10 text-slate-400"}`}>
                  {d.status}
                </span>
              </Link>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

function Stat({ label, value, href }: { label: string; value: number; href?: string }) {
  const content = (
    <div className="bg-surface border border-white/10 rounded-xl p-5 hover:border-primary/30 transition-colors">
      <p className="text-2xl font-bold text-white">{value}</p>
      <p className="text-xs text-slate-400 mt-1">{label}</p>
    </div>
  );
  return href ? <Link href={href} className="block">{content}</Link> : content;
}
