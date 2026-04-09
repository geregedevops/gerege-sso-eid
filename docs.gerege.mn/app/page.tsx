import Link from "next/link";
import { sections } from "@/content/docs";

export default function HomePage() {
  return (
    <main className="max-w-4xl mx-auto px-6 py-16 space-y-12">
      <div className="text-center space-y-4">
        <h1 className="text-4xl font-extrabold text-white">Gerege Docs</h1>
        <p className="text-lg text-slate-400 max-w-xl mx-auto">
          Gerege platform-ийн бүх service-ийн documentation, заавар, тохиргоо.
        </p>
      </div>

      <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {sections.map((section) => (
          <div key={section.title} className="bg-surface border border-white/8 rounded-xl p-5 space-y-3">
            <h2 className="font-bold text-white text-sm">{section.title}</h2>
            <ul className="space-y-1">
              {section.pages.map((page) => (
                <li key={page.slug}>
                  <Link
                    href={`/${page.slug}`}
                    className="text-sm text-slate-400 hover:text-primary transition-colors block py-0.5"
                  >
                    {page.title}
                    {page.description && (
                      <span className="text-slate-600 ml-1">— {page.description}</span>
                    )}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </div>

      <div className="text-center text-xs text-slate-600">
        Gerege Systems LLC &middot; <a href="https://gerege.mn" className="text-primary hover:underline">gerege.mn</a>
      </div>
    </main>
  );
}
