import Link from "next/link";
import { notFound } from "next/navigation";
import { auth } from "@/lib/auth";
import { query, queryOne } from "@/lib/db";

export default async function OrgDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const session = await auth();
  const sub = (session?.user as any)?.sub;
  const user = sub ? await queryOne<any>(`SELECT id FROM dbiz_users WHERE sub=$1`, [sub]) : null;
  if (!user) notFound();

  const membership = await queryOne<any>(
    `SELECT role FROM dbiz_org_members WHERE "organizationId"=$1 AND "userId"=$2`, [id, user.id]
  );
  if (!membership) notFound();

  const org = await queryOne<any>(`SELECT * FROM dbiz_organizations WHERE id=$1`, [id]);
  if (!org) notFound();

  const members = await query<any>(
    `SELECT u.name, m.role FROM dbiz_org_members m JOIN dbiz_users u ON u.id=m."userId" WHERE m."organizationId"=$1`, [id]
  );
  const docs = await query<any>(
    `SELECT id, name, "fileName", "fileSize", status, "createdAt" FROM dbiz_documents WHERE "organizationId"=$1 ORDER BY "createdAt" DESC LIMIT 5`, [id]
  );

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-white">{org.name}</h1>
        <p className="text-sm text-slate-400">РД: {org.registrationNumber} &middot; {org.type}</p>
      </div>

      <div className="flex gap-3">
        <Link href={`/dashboard/sign?org=${id}`} className="px-4 py-2 bg-primary text-white text-sm font-semibold rounded-lg">Гарын үсэг зурах</Link>
      </div>

      <div>
        <h2 className="text-lg font-bold text-white mb-3">Гишүүд ({members.length})</h2>
        <div className="bg-surface border border-white/10 rounded-xl overflow-hidden">
          <table className="w-full text-sm">
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

      {docs.length > 0 && (
        <div>
          <h2 className="text-lg font-bold text-white mb-3">Баримтууд</h2>
          <div className="space-y-2">
            {docs.map((d: any) => (
              <Link key={d.id} href={`/dashboard/documents/${d.id}`}
                className="bg-surface border border-white/10 rounded-lg p-4 flex items-center justify-between hover:border-primary/30 block">
                <div>
                  <p className="text-white text-sm">{d.name}</p>
                  <p className="text-xs text-slate-400">{Math.round(d.fileSize / 1024)} KB</p>
                </div>
                <span className={`text-xs px-2 py-1 rounded-full ${d.status === "signed" ? "bg-green-500/10 text-green-400" : "bg-slate-500/10 text-slate-400"}`}>{d.status}</span>
              </Link>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
