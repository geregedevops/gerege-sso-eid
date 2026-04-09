import { notFound } from "next/navigation";
import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";

export default async function DocumentDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const session = await auth();
  const sub = (session?.user as any)?.sub;
  const user = sub ? await prisma.user.findUnique({ where: { sub } }) : null;
  if (!user) notFound();

  const doc = await prisma.document.findUnique({
    where: { id },
    include: { organization: true, uploadedBy: true, signatures: { include: { signedBy: true } } },
  });
  if (!doc || doc.uploadedById !== user.id) notFound();

  return (
    <div className="max-w-2xl space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-white">{doc.name}</h1>
        <p className="text-sm text-slate-400 mt-1">{doc.organization.name} &middot; {doc.fileName}</p>
      </div>

      <div className="bg-surface border border-white/10 rounded-xl p-5 space-y-3">
        <Row label="Файлын нэр" value={doc.fileName} />
        <Row label="Хэмжээ" value={`${(doc.fileSize / 1024).toFixed(0)} KB`} />
        <Row label="SHA256" value={doc.fileHash} mono />
        <Row label="Статус" value={doc.status} />
        <Row label="Огноо" value={doc.createdAt.toISOString().split("T")[0]} />
        <Row label="Upload хийсэн" value={doc.uploadedBy.name} />
      </div>

      {doc.signatures.length > 0 && (
        <div>
          <h2 className="text-lg font-bold text-white mb-3">Гарын үсгүүд</h2>
          <div className="space-y-3">
            {doc.signatures.map((s) => (
              <div key={s.id} className="bg-surface border border-white/10 rounded-xl p-5">
                <div className="flex items-center justify-between mb-2">
                  <p className="text-white font-medium">{s.signerName || s.signedBy.name}</p>
                  <span className={`text-xs px-2 py-0.5 rounded-full ${s.status === "complete" ? "bg-green-500/10 text-green-400" : s.status === "failed" ? "bg-red-500/10 text-red-400" : "bg-yellow-500/10 text-yellow-400"}`}>
                    {s.status}
                  </span>
                </div>
                <p className="text-xs text-slate-400">
                  {s.certSerial ? `Cert: ${s.certSerial}` : ""}
                  {s.signedAt ? ` &middot; ${s.signedAt.toISOString().split("T")[0]}` : ""}
                </p>
                {s.status === "complete" && (
                  <a href={`/api/sign/${s.id}/result`} className="inline-block mt-3 px-4 py-2 bg-primary text-white text-xs font-semibold rounded-lg">
                    Татаж авах
                  </a>
                )}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

function Row({ label, value, mono }: { label: string; value: string; mono?: boolean }) {
  return (
    <div className="flex items-start gap-4">
      <span className="text-xs text-slate-400 w-32 shrink-0">{label}</span>
      <span className={`text-sm text-white ${mono ? "font-mono text-xs break-all" : ""}`}>{value}</span>
    </div>
  );
}
