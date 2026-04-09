import type { Metadata } from "next";
import Link from "next/link";
import "./globals.css";

export const metadata: Metadata = {
  title: "Gerege Docs",
  description: "Gerege platform documentation wiki",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="mn">
      <body className="font-sans antialiased">
        <nav className="border-b border-white/6 bg-bg/80 backdrop-blur-md sticky top-0 z-50">
          <div className="max-w-7xl mx-auto px-6 h-14 flex items-center justify-between">
            <Link href="/" className="flex items-center gap-2.5 font-bold text-white">
              <div className="w-7 h-7 bg-primary rounded-md flex items-center justify-center text-white text-sm font-black">G</div>
              <span>Gerege <span className="text-primary">Docs</span></span>
            </Link>
            <div className="flex items-center gap-6 text-sm text-slate-400">
              <a href="https://developer.gerege.mn" className="hover:text-white transition-colors">Developer Portal</a>
              <a href="https://sso.gerege.mn" className="hover:text-white transition-colors">SSO</a>
              <a href="https://dan.gerege.mn" className="hover:text-white transition-colors">DAN</a>
              <a href="https://gsign.gerege.mn" className="hover:text-white transition-colors">G-Sign</a>
            </div>
          </div>
        </nav>
        {children}
      </body>
    </html>
  );
}
