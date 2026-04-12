import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";
import { notFound } from "next/navigation";
import Link from "next/link";

export default async function AppDetailPage({ params }: { params: { appId: string } }) {
  const session = await auth();
  const sub = (session?.user as any)?.sub;

  const developer = sub ? await prisma.developer.findUnique({ where: { sub } }) : null;
  if (!developer) notFound();

  const app = await prisma.app.findFirst({
    where: { id: params.appId, developerId: developer.id },
  });
  if (!app) notFound();

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <div className="flex items-center gap-2">
            <span className={`w-2.5 h-2.5 rounded-full ${app.isActive ? 'bg-green-500' : 'bg-red-500'}`} />
            <h1 className="text-2xl font-bold text-white">{app.name}</h1>
          </div>
          {app.description && <p className="text-slate-400 text-sm mt-1">{app.description}</p>}
        </div>
        <Link href="/dashboard/apps" className="text-sm text-slate-400 hover:text-white">
          Буцах
        </Link>
      </div>

      {/* Credentials */}
      <Section title="Credentials">
        <InfoRow label="client_id" value={app.clientId} mono />
        <InfoRow label="client_secret" value="****" />
        <p className="text-xs text-slate-500 mt-2">Secret зөвхөн үүсгэх үед нэг удаа харуулна.</p>
      </Section>

      {/* Redirect URIs */}
      <Section title="Redirect URIs">
        {app.redirectUris.map((uri, i) => (
          <div key={i} className="text-sm font-mono text-slate-300 bg-bg px-3 py-2 rounded-lg border border-white/10 mb-2">
            {uri}
          </div>
        ))}
      </Section>

      {/* Scopes */}
      <Section title="Scopes">
        <div className="flex gap-2 flex-wrap">
          {app.scopes.map((s) => (
            <span key={s} className="px-3 py-1 bg-primary/10 text-primary text-xs font-medium rounded-full border border-primary/20">
              {s}
            </span>
          ))}
        </div>
      </Section>

      {/* Integration */}
      <Section title="Integration">
        <pre className="bg-bg border border-white/10 rounded-lg p-4 text-xs text-slate-300 overflow-x-auto">
{`// NextAuth.js
providers: [{
  id: "gerege-sso",
  name: "GeregeID",
  type: "oidc",
  issuer: "https://sso.gerege.mn",
  clientId: "${app.clientId}",
  clientSecret: process.env.EID_CLIENT_SECRET,
}]`}
        </pre>
      </Section>

      {/* Metadata */}
      <Section title="Мэдээлэл">
        <InfoRow label="App ID" value={app.id} mono />
        <InfoRow label="Үүсгэсэн" value={app.createdAt.toLocaleString("mn-MN")} />
        <InfoRow label="Шинэчилсэн" value={app.updatedAt.toLocaleString("mn-MN")} />
      </Section>
    </div>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="bg-surface border border-white/10 rounded-xl p-5">
      <h2 className="text-sm font-semibold text-slate-400 uppercase tracking-wider mb-3">{title}</h2>
      {children}
    </div>
  );
}

function InfoRow({ label, value, mono }: { label: string; value: string; mono?: boolean }) {
  return (
    <div className="flex items-center justify-between py-1.5">
      <span className="text-xs text-slate-500">{label}</span>
      <span className={`text-sm text-white ${mono ? "font-mono" : ""}`}>{value}</span>
    </div>
  );
}
