import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "XYP — xyp.gerege.mn",
  description: "Иргэн & байгууллага баталгаажуулах API",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="mn" className="dark">
      <body className="font-sans">
        <nav className="flex items-center justify-between px-8 py-4 border-b border-white/5">
          <a href="/" className="flex items-center gap-3">
            <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center text-white font-bold text-xs">
              VFY
            </div>
            <span className="font-bold text-white">Verify API</span>
          </a>
          <div className="flex gap-6 text-sm text-slate-400 items-center">
            <a href="/" className="hover:text-white">Нүүр</a>
            <a href="/docs" className="hover:text-white">API Заавар</a>
            <a href="/admin" className="hover:text-white">Admin</a>
          </div>
        </nav>
        {children}
      </body>
    </html>
  );
}
