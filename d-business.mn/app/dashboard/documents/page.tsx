import Link from "next/link";
import { auth } from "@/lib/auth";
import { query, queryOne } from "@/lib/db";

export default async function DocumentsPage() {
  const session = await auth();
  const sub = (session?.user as any)?.sub;
  const user = sub ? await queryOne<any>(`SELECT id FROM dbiz_users WHERE sub=$1`, [sub]) : null;
  if (!user) return null;

  const docs = await query<any>(
    `SELECT d.*, o.name as org_name FROM dbiz_documents d
     JOIN dbiz_organizations o ON o.id=d."organizationId"
     WHERE d."uploadedById"=$1 ORDER BY d."createdAt" DESC LIMIT 50`,
    [user.id]
  );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-white">Баримтууд</h1>
        <Link href="/dashboard/sign" className="px-4 py-2 bg-primary text-white text-sm font-semibold rounded-lg">+ Гарын үсэг</Link>
      </div>
      {docs.length === 0 ? (
        <div className="bg-surface border border-white/10 rounded-xl p-12 text-center"><p className="text-slate-400">Баримт байхгүй.</p></div>
      ) : (
        <div className="space-y-2">
          {docs.map((d: any) => (
            <Link key={d.id} href={`/dashboard/documents/${d.id}`}
              className="bg-surface border border-white/10 rounded-lg p-4 flex items-center justify-between hover:border-primary/30 block">
              <div>
                <p className="text-white text-sm">{d.name}</p>
                <p className="text-xs text-slate-400">{d.org_name} &middot; {Math.round(d.fileSize / 1024)} KB</p>
              </div>
              <span className={`text-xs px-2 py-1 rounded-full ${d.status === "signed" ? "bg-green-500/10 text-green-400" : "bg-slate-500/10 text-slate-400"}`}>{d.status}</span>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
