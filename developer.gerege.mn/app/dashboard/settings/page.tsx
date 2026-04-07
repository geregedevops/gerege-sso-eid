import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";

export default async function SettingsPage() {
  const session = await auth();
  const sub = (session?.user as any)?.sub;
  const developer = sub ? await prisma.developer.findUnique({ where: { sub } }) : null;

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-white">Тохиргоо</h1>

      <div className="bg-surface border border-white/10 rounded-xl p-5 space-y-4">
        <h2 className="text-sm font-semibold text-slate-400 uppercase tracking-wider">Профайл</h2>
        <InfoRow label="Нэр" value={developer?.name || "-"} />
        <InfoRow label="Овог" value={developer?.familyName || "-"} />
        <InfoRow label="Нэр" value={developer?.givenName || "-"} />
        <InfoRow label="Sub" value={developer?.sub || "-"} mono />
        <InfoRow label="Certificate Serial" value={developer?.certSerial || "-"} mono />
        <InfoRow label="Бүртгүүлсэн" value={developer?.createdAt.toLocaleString("mn-MN") || "-"} />
      </div>
    </div>
  );
}

function InfoRow({ label, value, mono }: { label: string; value: string; mono?: boolean }) {
  return (
    <div className="flex items-center justify-between py-1.5">
      <span className="text-xs text-slate-500">{label}</span>
      <span className={`text-sm text-white ${mono ? "font-mono text-xs" : ""}`}>{value}</span>
    </div>
  );
}
