import { notFound } from "next/navigation";
import { auth } from "@/lib/auth";
import { query, queryOne } from "@/lib/db";

export default async function DocumentDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const session = await auth();
  const sub = (session?.user as any)?.sub;
  const user = sub ? await queryOne<any>(`SELECT id FROM dbiz_users WHERE sub=$1`, [sub]) : null;
  if (!user) notFound();

  const doc = await queryOne<any>(
    `SELECT d.*, o.name as org_name FROM dbiz_documents d JOIN dbiz_organizations o ON o.id=d."organizationId" WHERE d.id=$1`, [id]
  );
  if (!doc || doc.uploadedById !== user.id) notFound();

  const sigs = await query<any>(
    `SELECT s.*, u.name as signer_user_name FROM dbiz_signatures s JOIN dbiz_users u ON u.id=s."signedById" WHERE s."documentId"=$1`, [id]
  );

  return (
    <div className="max-w-2xl space-y-6">
      <h1 className="text-2xl font-bold text-white">{doc.name}</h1>
      <p className="text-sm text-slate-400">{doc.org_name} &middot; {doc.fileName}</p>

      <div className="bg-surface border border-white/10 rounded-xl p-5 space-y-2 text-sm">
        <Row label="Хэмжээ" value={`${Math.round(doc.fileSize / 1024)} KB`} />
        <Row label="SHA256" value={doc.fileHash} mono />
        <Row label="Статус" value={doc.status} />
      </div>

      {sigs.length > 0 && (
        <div>
          <h2 className="text-lg font-bold text-white mb-3">Гарын үсгүүд</h2>
          {sigs.map((s: any) => (
            <div key={s.id} className="bg-surface border border-white/10 rounded-xl p-5 mb-3">
              <p className="text-white font-medium">{s.signerName || s.signer_user_name}</p>
              <p className="text-xs text-slate-400">{s.status} {s.certSerial ? `&middot; ${s.certSerial}` : ""}</p>
              {s.status === "complete" && (
                <a href={`/api/sign/${s.id}/result`} className="inline-block mt-2 px-4 py-2 bg-primary text-white text-xs rounded-lg">Татах</a>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

function Row({ label, value, mono }: { label: string; value: string; mono?: boolean }) {
  return (
    <div className="flex gap-4">
      <span className="text-slate-400 w-24">{label}</span>
      <span className={`text-white ${mono ? "font-mono text-xs break-all" : ""}`}>{value}</span>
    </div>
  );
}
