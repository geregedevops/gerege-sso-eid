import { auth } from "@/lib/auth";
import { redirect } from "next/navigation";
import Link from "next/link";

export default async function SandboxLayout({ children }: { children: React.ReactNode }) {
  const session = await auth();
  if (!session?.user) redirect("/auth/login");

  return (
    <div className="max-w-6xl mx-auto px-6 py-8 flex gap-8">
      <aside className="hidden md:block w-48 flex-shrink-0">
        <nav className="space-y-1 sticky top-20">
          <SideLink href="/sandbox">API Explorer</SideLink>
          <SideLink href="/sandbox/pos">POS Sandbox</SideLink>
          <SideLink href="/sandbox/payment">Payment</SideLink>
          <SideLink href="/sandbox/social">Social</SideLink>
          <SideLink href="/sandbox/webhook">Webhook</SideLink>
        </nav>
      </aside>
      <main className="flex-1 min-w-0">{children}</main>
    </div>
  );
}

function SideLink({ href, children }: { href: string; children: React.ReactNode }) {
  return (
    <Link
      href={href}
      className="block px-3 py-2 text-sm text-slate-400 rounded-lg hover:bg-white/5 hover:text-white transition-colors"
    >
      {children}
    </Link>
  );
}
