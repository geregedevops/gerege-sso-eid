import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";

export default async function SettingsPage() {
  const session = await auth();
  const sub = (session?.user as any)?.sub;
  const user = sub ? await prisma.user.findUnique({ where: { sub }, include: { memberships: { include: { organization: true } } } }) : null;
  if (!user) return null;

  return (
    <div className="max-w-2xl space-y-6">
      <h1 className="text-2xl font-bold text-white">Тохиргоо</h1>

      <div className="bg-surface border border-white/10 rounded-xl p-6 space-y-4">
        <h2 className="font-semibold text-white">Хэрэглэгчийн мэдээлэл</h2>
        <Row label="Нэр" value={user.name} />
        <Row label="Овог" value={user.familyName} />
        <Row label="Нэр" value={user.givenName} />
        <Row label="Sub" value={user.sub} mono />
        <Row label="Cert Serial" value={user.certSerial || "—"} mono />
        <Row label="Бүртгүүлсэн" value={user.createdAt.toISOString().split("T")[0]} />
      </div>

      {user.memberships.length > 0 && (
        <div className="bg-surface border border-white/10 rounded-xl p-6 space-y-3">
          <h2 className="font-semibold text-white">Миний байгууллагууд</h2>
          {user.memberships.map((m) => (
            <div key={m.organizationId} className="flex items-center justify-between py-2 border-b border-white/5 last:border-0">
              <div>
                <p className="text-white text-sm">{m.organization.name}</p>
                <p className="text-xs text-slate-400">РД: {m.organization.registrationNumber}</p>
              </div>
              <span className="text-xs px-2 py-0.5 bg-primary/10 text-primary rounded-full">{m.role}</span>
            </div>
          ))}
        </div>
      )}

      <p className="text-xs text-slate-500">Мэдээлэл e-ID Mongolia-р баталгаажсан тул өөрчлөх боломжгүй.</p>
    </div>
  );
}

function Row({ label, value, mono }: { label: string; value: string; mono?: boolean }) {
  return (
    <div className="flex items-start gap-4 py-1">
      <span className="text-xs text-slate-400 w-28 shrink-0">{label}</span>
      <span className={`text-sm text-white ${mono ? "font-mono text-xs" : ""}`}>{value}</span>
    </div>
  );
}
