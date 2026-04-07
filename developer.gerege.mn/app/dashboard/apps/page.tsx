import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";
import Link from "next/link";

export default async function AppsPage() {
  const session = await auth();
  const sub = (session?.user as any)?.sub;

  const developer = sub ? await prisma.developer.findUnique({
    where: { sub },
    include: { apps: { orderBy: { createdAt: "desc" } } },
  }) : null;

  const apps = developer?.apps || [];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-white">Apps</h1>
        <Link
          href="/dashboard/apps/new"
          className="px-4 py-2 bg-primary text-white font-semibold rounded-lg text-sm hover:bg-primary-light transition-colors"
        >
          + Шинэ App
        </Link>
      </div>

      {apps.length === 0 ? (
        <div className="bg-surface border border-white/10 rounded-xl p-12 text-center">
          <p className="text-slate-400 mb-4">Одоогоор app бүртгүүлээгүй байна.</p>
          <Link
            href="/dashboard/apps/new"
            className="px-5 py-2.5 bg-primary text-white font-semibold rounded-xl text-sm"
          >
            Эхний App үүсгэх
          </Link>
        </div>
      ) : (
        <div className="space-y-3">
          {apps.map((app) => (
            <Link
              key={app.id}
              href={`/dashboard/apps/${app.id}`}
              className="block bg-surface border border-white/10 rounded-xl p-5 hover:border-primary/30 transition-colors"
            >
              <div className="flex items-center justify-between">
                <div>
                  <div className="flex items-center gap-2 mb-1">
                    <span className={`w-2 h-2 rounded-full ${app.isActive ? 'bg-green-500' : 'bg-red-500'}`} />
                    <span className="font-semibold text-white">{app.name}</span>
                  </div>
                  <p className="text-xs text-slate-500 font-mono">client_id: {app.clientId}</p>
                  {app.description && <p className="text-sm text-slate-400 mt-1">{app.description}</p>}
                </div>
                <div className="text-right">
                  <p className="text-xs text-slate-500">Redirect URIs: {app.redirectUris.length}</p>
                  <p className="text-xs text-slate-500 mt-1">{app.scopes.join(", ")}</p>
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
