import Link from "next/link";
import { redirect } from "next/navigation";
import { auth } from "@/lib/auth";

export default async function DashboardLayout({ children }: { children: React.ReactNode }) {
  const session = await auth();
  if (!session) redirect("/auth/login");

  const links = [
    { href: "/dashboard", label: "Тойм" },
    { href: "/dashboard/org", label: "Байгууллагууд" },
    { href: "/dashboard/sign", label: "Гарын үсэг зурах" },
    { href: "/dashboard/documents", label: "Баримтууд" },
    { href: "/dashboard/settings", label: "Тохиргоо" },
  ];

  return (
    <div className="flex min-h-[calc(100vh-56px)]">
      <aside className="hidden md:block w-56 border-r border-white/6 p-4 space-y-1">
        {links.map((l) => (
          <Link key={l.href} href={l.href} className="block px-3 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-white/5 transition-colors">
            {l.label}
          </Link>
        ))}
        <div className="border-t border-white/6 my-3" />
        <Link href="/verify" className="block px-3 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-white/5 transition-colors">
          Баримт шалгах
        </Link>
      </aside>
      <main className="flex-1 p-6 max-w-5xl">{children}</main>
    </div>
  );
}
