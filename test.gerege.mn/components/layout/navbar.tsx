import Link from "next/link";

export function Navbar() {
  return (
    <nav className="border-b border-white/10 bg-bg/80 backdrop-blur-md sticky top-0 z-50">
      <div className="max-w-6xl mx-auto px-6 h-14 flex items-center justify-between">
        <div className="flex items-center gap-8">
          <Link href="/" className="flex items-center gap-2.5 font-bold text-white">
            <div className="w-7 h-7 bg-primary rounded-md flex items-center justify-center text-white text-sm font-black">G</div>
            <span>Gerege <span className="text-primary">Sandbox</span></span>
          </Link>
          <div className="hidden sm:flex items-center gap-6 text-sm text-slate-400">
            <Link href="/sandbox" className="hover:text-white transition-colors">API Explorer</Link>
            <Link href="/sandbox/pos" className="hover:text-white transition-colors">POS</Link>
            <Link href="/sandbox/payment" className="hover:text-white transition-colors">Payment</Link>
            <Link href="/sandbox/webhook" className="hover:text-white transition-colors">Webhook</Link>
          </div>
        </div>
        <Link
          href="/auth/login"
          className="text-sm px-4 py-1.5 rounded-lg bg-primary/10 text-primary border border-primary/20 hover:bg-primary/20 transition-colors"
        >
          Нэвтрэх
        </Link>
      </div>
    </nav>
  );
}
