import Link from "next/link";
import { auth, signOut } from "@/lib/auth";

export async function Navbar() {
  const session = await auth();

  return (
    <nav className="border-b border-white/10 bg-bg/80 backdrop-blur-md sticky top-0 z-50">
      <div className="max-w-6xl mx-auto px-6 h-14 flex items-center justify-between">
        <div className="flex items-center gap-8">
          <Link href="/" className="flex items-center gap-2.5 font-bold text-white">
            <div className="w-7 h-7 bg-primary rounded-md flex items-center justify-center text-white text-sm font-black">B</div>
            <span>d-business<span className="text-primary">.mn</span></span>
          </Link>
          <div className="hidden sm:flex items-center gap-6 text-sm text-slate-400">
            <Link href="/dashboard" className="hover:text-white transition-colors">Dashboard</Link>
            <Link href="/verify" className="hover:text-white transition-colors">Шалгах</Link>
          </div>
        </div>
        {session?.user ? (
          <div className="flex items-center gap-4">
            <span className="text-sm text-slate-400">{session.user.name}</span>
            <form action={async () => { "use server"; await signOut(); }}>
              <button className="text-sm px-3 py-1 rounded-lg text-red-400 hover:bg-red-500/10 transition-colors">Гарах</button>
            </form>
          </div>
        ) : (
          <Link href="/auth/login" className="text-sm px-4 py-1.5 rounded-lg bg-primary/10 text-primary border border-primary/20 hover:bg-primary/20 transition-colors">
            Нэвтрэх
          </Link>
        )}
      </div>
    </nav>
  );
}
