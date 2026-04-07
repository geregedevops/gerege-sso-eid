import Link from "next/link";

export default function DocsPage() {
  return (
    <main className="max-w-4xl mx-auto px-6 py-12 space-y-8">
      <div>
        <h1 className="text-3xl font-bold text-white mb-2">Documentation</h1>
        <p className="text-slate-400">Gerege SSO OIDC integration гайд.</p>
      </div>

      <div className="grid sm:grid-cols-2 gap-4">
        <DocCard href="/docs/quickstart" title="Quickstart" desc="5 минутад нэгтгэх" />
        <DocCard href="/docs/api-reference" title="API Reference" desc="OIDC endpoint-ууд" />
        <DocCard href="/docs/guides/nextjs" title="Next.js Guide" desc="NextAuth.js + OIDC" />
        <DocCard href="/docs/guides/go" title="Go Guide" desc="golang.org/x/oauth2" />
        <DocCard href="/docs/guides/pos-plugin" title="POS Plugin" desc="POS API нэгтгэх" />
        <DocCard href="/docs/guides/social" title="Social Commerce" desc="Лайв худалдаа, product feed" />
        <DocCard href="/docs/guides/payment" title="Payment" desc="QPay, SocialPay, eBarimt" />
      </div>
    </main>
  );
}

function DocCard({ href, title, desc }: { href: string; title: string; desc: string }) {
  return (
    <Link href={href} className="bg-surface border border-white/10 rounded-xl p-5 hover:border-primary/30 transition-colors block">
      <h3 className="font-semibold text-white mb-1">{title}</h3>
      <p className="text-sm text-slate-400">{desc}</p>
    </Link>
  );
}
