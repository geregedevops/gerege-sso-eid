import Link from "next/link";

export default function DocsPage() {
  return (
    <main className="max-w-4xl mx-auto px-6 py-12 space-y-8">
      <div>
        <h1 className="text-3xl font-bold text-white mb-2">Documentation</h1>
        <p className="text-slate-400">Gerege SSO OIDC integration гайд.</p>
      </div>

      <div>
        <h2 className="text-lg font-bold text-white mb-3">3-р талын нэгтгэл</h2>
        <div className="grid sm:grid-cols-2 gap-4 mb-8">
          <DocCard href="/docs/guides/sso-integration" title="e-ID SSO нэгтгэх" desc="OIDC Authorization Code Flow бүрэн заавар" highlight />
          <DocCard href="/docs/guides/dan-integration" title="DAN Verify нэгтгэх" desc="Иргэний мэдээлэл + зураг авах" highlight />
        </div>
      </div>

      <div>
        <h2 className="text-lg font-bold text-white mb-3">Хэл / Framework</h2>
        <div className="grid sm:grid-cols-2 gap-4 mb-8">
          <DocCard href="/docs/quickstart" title="Quickstart" desc="5 минутад нэгтгэх" />
          <DocCard href="/docs/api-reference" title="API Reference" desc="OIDC endpoint-ууд" />
          <DocCard href="/docs/guides/nextjs" title="Next.js Guide" desc="NextAuth.js + OIDC" />
          <DocCard href="/docs/guides/go" title="Go Guide" desc="golang.org/x/oauth2" />
        </div>
      </div>

      <div>
        <h2 className="text-lg font-bold text-white mb-3">Gerege API</h2>
        <div className="grid sm:grid-cols-2 gap-4">
          <DocCard href="/docs/guides/pos-plugin" title="POS Plugin" desc="POS API нэгтгэх" />
          <DocCard href="/docs/guides/social" title="Social Commerce" desc="Лайв худалдаа, product feed" />
          <DocCard href="/docs/guides/payment" title="Payment" desc="QPay, SocialPay, eBarimt" />
        </div>
      </div>
    </main>
  );
}

function DocCard({ href, title, desc, highlight }: { href: string; title: string; desc: string; highlight?: boolean }) {
  return (
    <Link href={href} className={`rounded-xl p-5 hover:border-primary/30 transition-colors block ${highlight ? "bg-primary/5 border border-primary/20" : "bg-surface border border-white/10"}`}>
      <h3 className="font-semibold text-white mb-1">{title}</h3>
      <p className="text-sm text-slate-400">{desc}</p>
    </Link>
  );
}
