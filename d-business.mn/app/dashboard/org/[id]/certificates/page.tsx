import { notFound } from "next/navigation";
import { auth } from "@/lib/auth";
import { query, queryOne } from "@/lib/db";

export default async function CertificatesPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const session = await auth();
  const sub = (session?.user as any)?.sub;
  const user = sub ? await queryOne<any>(`SELECT id FROM dbiz_users WHERE sub=$1`, [sub]) : null;
  if (!user) notFound();

  const org = await queryOne<any>(`SELECT name FROM dbiz_organizations WHERE id=$1`, [id]);
  if (!org) notFound();

  const certs = await query<any>(
    `SELECT * FROM dbiz_certificates WHERE "organizationId"=$1 ORDER BY "createdAt" DESC`, [id]
  );

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-white">{org.name} — Сертификат</h1>
      {certs.length === 0 ? (
        <div className="bg-surface border border-white/10 rounded-xl p-12 text-center">
          <p className="text-slate-400">Сертификат байхгүй.</p>
        </div>
      ) : (
        <div className="space-y-3">
          {certs.map((c: any) => (
            <div key={c.id} className="bg-surface border border-white/10 rounded-xl p-5">
              <h3 className="font-semibold text-white">{c.commonName}</h3>
              <p className="text-xs text-slate-400">{c.purpose} &middot; {c.status}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
