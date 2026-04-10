import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({ subsets: ["latin", "cyrillic"] });

export const metadata: Metadata = {
  title: "DAN Gateway — dan.gerege.mn",
  description: "ДАН баталгаажуулалтын gateway удирдлагын самбар",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="mn" className="dark">
      <body className={inter.className}>
        <nav className="flex items-center justify-between px-8 py-4 border-b border-white/5">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center text-white font-bold text-xs">
              DAN
            </div>
            <span className="font-bold text-white">DAN Gateway</span>
          </div>
          <div className="flex gap-6 text-sm text-slate-400">
            <a href="/" className="hover:text-white">Нүүр</a>
            <a href="/dashboard" className="hover:text-white">Dashboard</a>
            <a href="/docs" className="hover:text-white">Заавар</a>
          </div>
        </nav>
        {children}
      </body>
    </html>
  );
}
