import { notFound } from "next/navigation";
import { auth } from "@/lib/auth";
import { query, queryOne } from "@/lib/db";

export default async function MembersPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const session = await auth();
  const sub = (session?.user as any)?.sub;
  const user = sub ? await queryOne<any>(`SELECT id FROM dbiz_users WHERE sub=$1`, [sub]) : null;
  if (!user) notFound();

  const org = await queryOne<any>(`SELECT name FROM dbiz_organizations WHERE id=$1`, [id]);
  if (!org) notFound();

  const members = await query<any>(
    `SELECT u.name, m.role, m."createdAt" FROM dbiz_org_members m JOIN dbiz_users u ON u.id=m."userId" WHERE m."organizationId"=$1`, [id]
  );

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-white">{org.name} — Гишүүд</h1>
      <div className="bg-surface border border-white/10 rounded-xl overflow-hidden">
        <table className="w-full text-sm">
          <thead><tr className="border-b border-white/10"><th className="text-left px-4 py-3 text-slate-400">Нэр</th><th className="text-left px-4 py-3 text-slate-400">Үүрэг</th></tr></thead>
          <tbody>
            {members.map((m: any, i: number) => (
              <tr key={i} className="border-b border-white/5">
                <td className="px-4 py-3 text-white">{m.name}</td>
                <td className="px-4 py-3 text-slate-400">{m.role}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
