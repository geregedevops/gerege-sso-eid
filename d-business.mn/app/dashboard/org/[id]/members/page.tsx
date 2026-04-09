import { notFound } from "next/navigation";
import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";

export default async function MembersPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const session = await auth();
  const sub = (session?.user as any)?.sub;
  const user = sub ? await prisma.user.findUnique({ where: { sub } }) : null;
  if (!user) notFound();

  const members = await prisma.orgMember.findMany({
    where: { organizationId: id },
    include: { user: true },
    orderBy: { createdAt: "asc" },
  });

  const org = await prisma.organization.findUnique({ where: { id } });
  if (!org) notFound();

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-white">{org.name} — Гишүүд</h1>
      <div className="bg-surface border border-white/10 rounded-xl overflow-hidden">
        <table className="w-full text-sm">
          <thead><tr className="border-b border-white/10"><th className="text-left px-4 py-3 text-slate-400">Нэр</th><th className="text-left px-4 py-3 text-slate-400">Үүрэг</th><th className="text-left px-4 py-3 text-slate-400">Нэгдсэн</th></tr></thead>
          <tbody>
            {members.map((m) => (
              <tr key={m.userId} className="border-b border-white/5">
                <td className="px-4 py-3 text-white">{m.user.name}</td>
                <td className="px-4 py-3"><span className={`text-xs px-2 py-0.5 rounded-full ${m.role === "owner" ? "bg-primary/10 text-primary" : "bg-slate-500/10 text-slate-400"}`}>{m.role}</span></td>
                <td className="px-4 py-3 text-slate-400">{m.createdAt.toISOString().split("T")[0]}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
