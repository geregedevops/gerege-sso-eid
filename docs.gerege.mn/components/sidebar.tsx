import Link from "next/link";
import { sections } from "@/content/docs";

export function Sidebar({ currentSlug }: { currentSlug: string }) {
  return (
    <nav className="w-64 min-w-[16rem] border-r border-white/6 p-6 space-y-6 overflow-y-auto h-[calc(100vh-56px)] sticky top-14">
      {sections.map((section) => (
        <div key={section.title}>
          <h3 className="text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">
            {section.title}
          </h3>
          <ul className="space-y-1">
            {section.pages.map((page) => {
              const isActive = page.slug === currentSlug;
              return (
                <li key={page.slug}>
                  <Link
                    href={`/${page.slug}`}
                    className={`block px-3 py-1.5 rounded-lg text-sm transition-colors ${
                      isActive
                        ? "bg-primary/10 text-primary font-medium"
                        : "text-slate-400 hover:text-white hover:bg-white/5"
                    }`}
                  >
                    {page.title}
                  </Link>
                </li>
              );
            })}
          </ul>
        </div>
      ))}
    </nav>
  );
}
