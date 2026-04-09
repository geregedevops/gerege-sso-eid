import Link from "next/link";
import { auth } from "@/lib/auth";
import { query, queryOne } from "@/lib/db";

export default async function OrgListPage() {
  const session = await auth();
  const sub = (session?.user as any)?.sub;
  const user = sub ? await queryOne<any>(`SELECT id FROM dbiz_users WHERE sub=$1`, [sub]) : null;

  const orgs = user ? await query<any>(
    `SELECT o.*, m.role FROM dbiz_organizations o JOIN dbiz_org_members m ON m."organizationId"=o.id WHERE m."userId"=$1`,
    [user.id]
  ) : [];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-white">Байгууллагууд</h1>
        <Link href="/dashboard/org/new" className="px-4 py-2 bg-primary text-white text-sm font-semibold rounded-lg">+ Шинэ</Link>
      </div>
      {orgs.length === 0 ? (
        <div className="bg-surface border border-white/10 rounded-xl p-12 text-center">
          <p className="text-slate-400 mb-4">Байгууллага бүртгэгдээгүй.</p>
          <Link href="/dashboard/org/new" className="text-primary text-sm">Бүртгүүлэх</Link>
        </div>
      ) : (
        <div className="grid sm:grid-cols-2 gap-4">
          {orgs.map((o: any) => (
            <Link key={o.id} href={`/dashboard/org/${o.id}`} className="bg-surface border border-white/10 rounded-xl p-5 hover:border-primary/30 transition-colors block">
              <h3 className="font-semibold text-white">{o.name}</h3>
              <p className="text-xs text-slate-400">РД: {o.registrationNumber} &middot; {o.type} &middot; {o.role}</p>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
