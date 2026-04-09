import Link from "next/link";
import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";

export default async function DocumentsPage() {
  const session = await auth();
  const sub = (session?.user as any)?.sub;
  const user = sub ? await prisma.user.findUnique({ where: { sub } }) : null;
  if (!user) return null;

  const documents = await prisma.document.findMany({
    where: { uploadedById: user.id },
    include: { organization: true, signatures: true },
    orderBy: { createdAt: "desc" },
    take: 50,
  });

  const statusColors: Record<string, string> = {
    signed: "bg-green-500/10 text-green-400",
    signing: "bg-yellow-500/10 text-yellow-400",
    uploaded: "bg-blue-500/10 text-blue-400",
    failed: "bg-red-500/10 text-red-400",
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-white">Баримтууд</h1>
        <Link href="/dashboard/sign" className="px-4 py-2 bg-primary text-white text-sm font-semibold rounded-lg hover:bg-primary-light transition-colors">
          + Гарын үсэг зурах
        </Link>
      </div>

      {documents.length === 0 ? (
        <div className="bg-surface border border-white/10 rounded-xl p-12 text-center">
          <p className="text-slate-400">Баримт бичиг байхгүй байна.</p>
        </div>
      ) : (
        <div className="space-y-2">
          {documents.map((d) => (
            <Link key={d.id} href={`/dashboard/documents/${d.id}`}
              className="bg-surface border border-white/10 rounded-lg p-4 flex items-center justify-between hover:border-primary/30 transition-colors block">
              <div>
                <p className="text-white text-sm font-medium">{d.name}</p>
                <p className="text-xs text-slate-400">{d.organization.name} &middot; {(d.fileSize / 1024).toFixed(0)} KB &middot; {d.createdAt.toISOString().split("T")[0]}</p>
              </div>
              <span className={`text-xs px-2 py-1 rounded-full ${statusColors[d.status] || "bg-slate-500/10 text-slate-400"}`}>
                {d.status}
              </span>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
