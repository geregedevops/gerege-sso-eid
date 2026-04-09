import Link from "next/link";
import { auth } from "@/lib/auth";
import { query, queryOne } from "@/lib/db";

export default async function DashboardPage() {
  const session = await auth();
  const sub = (session?.user as any)?.sub;
  if (!sub) return null;

  const user = await queryOne<any>(`SELECT id, name, "givenName" FROM dbiz_users WHERE sub=$1`, [sub]);
  if (!user) return null;

  const orgs = await query<any>(
    `SELECT o.*, m.role FROM dbiz_organizations o JOIN dbiz_org_members m ON m."organizationId"=o.id WHERE m."userId"=$1`,
    [user.id]
  );
  const docCount = await queryOne<any>(`SELECT count(*)::int as c FROM dbiz_documents WHERE "uploadedById"=$1`, [user.id]);
  const sigCount = await queryOne<any>(`SELECT count(*)::int as c FROM dbiz_signatures WHERE "signedById"=$1 AND status='complete'`, [user.id]);

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-white">Сайн байна уу, {user.givenName || user.name}</h1>
        <p className="text-sm text-slate-400 mt-1">Байгууллагын цахим тамга платформ</p>
      </div>

      <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Stat label="Байгууллагууд" value={orgs.length} />
        <Stat label="Баримтууд" value={docCount?.c || 0} />
        <Stat label="Гарын үсэг" value={sigCount?.c || 0} />
      </div>

      <div className="flex gap-3">
        <Link href="/dashboard/org/new" className="px-4 py-2 bg-primary text-white text-sm font-semibold rounded-lg">+ Байгууллага</Link>
        <Link href="/dashboard/sign" className="px-4 py-2 border border-white/15 text-white text-sm rounded-lg">Гарын үсэг зурах</Link>
      </div>

      {orgs.length > 0 && (
        <div>
          <h2 className="text-lg font-bold text-white mb-3">Миний байгууллагууд</h2>
          <div className="grid sm:grid-cols-2 gap-4">
            {orgs.map((o: any) => (
              <Link key={o.id} href={`/dashboard/org/${o.id}`}
                className="bg-surface border border-white/10 rounded-xl p-5 hover:border-primary/30 transition-colors block">
                <h3 className="font-semibold text-white">{o.name}</h3>
                <p className="text-xs text-slate-400">РД: {o.registrationNumber} &middot; {o.type} &middot; {o.role}</p>
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
