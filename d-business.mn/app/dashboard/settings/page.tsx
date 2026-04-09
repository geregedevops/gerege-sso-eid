import { auth } from "@/lib/auth";
import { query, queryOne } from "@/lib/db";

export default async function SettingsPage() {
  const session = await auth();
  const sub = (session?.user as any)?.sub;
  const user = sub ? await queryOne<any>(
    `SELECT id, sub, name, "givenName", "familyName", "certSerial", "createdAt" FROM dbiz_users WHERE sub=$1`, [sub]
  ) : null;
  if (!user) return null;

  const memberships = await query<any>(
    `SELECT o.name, o."registrationNumber", m.role FROM dbiz_org_members m JOIN dbiz_organizations o ON o.id=m."organizationId" WHERE m."userId"=$1`, [user.id]
  );

  return (
    <div className="max-w-2xl space-y-6">
      <h1 className="text-2xl font-bold text-white">Тохиргоо</h1>
      <div className="bg-surface border border-white/10 rounded-xl p-6 space-y-3 text-sm">
        <Row label="Нэр" value={user.name} />
        <Row label="Овог" value={user.familyName} />
        <Row label="Нэр" value={user.givenName} />
        <Row label="Sub" value={user.sub} mono />
        <Row label="Cert" value={user.certSerial || "—"} mono />
      </div>
      {memberships.length > 0 && (
        <div className="bg-surface border border-white/10 rounded-xl p-6 space-y-3">
          <h2 className="font-semibold text-white">Байгууллагууд</h2>
          {memberships.map((m: any, i: number) => (
            <div key={i} className="flex items-center justify-between py-2 border-b border-white/5 last:border-0">
              <div><p className="text-white text-sm">{m.name}</p><p className="text-xs text-slate-400">РД: {m.registrationNumber}</p></div>
              <span className="text-xs px-2 py-0.5 bg-primary/10 text-primary rounded-full">{m.role}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

function Row({ label, value, mono }: { label: string; value: string; mono?: boolean }) {
  return (
    <div className="flex gap-4"><span className="text-slate-400 w-20">{label}</span><span className={`text-white ${mono ? "font-mono text-xs" : ""}`}>{value}</span></div>
  );
}
