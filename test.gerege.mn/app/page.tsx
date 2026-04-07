import Link from "next/link";

export default function LandingPage() {
  return (
    <main className="min-h-[calc(100vh-56px)]">
      <section className="max-w-4xl mx-auto px-6 pt-24 pb-16 text-center">
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary/10 border border-primary/20 text-primary text-xs font-semibold mb-6">
          API Sandbox
        </div>
        <h1 className="text-4xl sm:text-5xl font-extrabold text-white leading-tight mb-4">
          Gerege <span className="text-primary">API Sandbox</span>
        </h1>
        <p className="text-lg text-slate-400 max-w-xl mx-auto mb-10">
          Gerege platform-ийн API-г бодит transaction хийлгүй туршина.
        </p>
        <Link
          href="/auth/login"
          className="px-6 py-3 bg-primary text-white font-semibold rounded-xl hover:bg-primary-light transition-colors"
        >
          Sandbox нэвтрэх
        </Link>
      </section>

      <section className="max-w-5xl mx-auto px-6 pb-24">
        <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-6">
          <FeatureCard title="POS Sandbox" desc="Бараа, захиалга, гүйлгээ туршах" href="/sandbox/pos" />
          <FeatureCard title="Social Sandbox" desc="Product feed, лайв худалдаа" href="/sandbox/social" />
          <FeatureCard title="Payment Simulator" desc="QPay, SocialPay, eBarimt" href="/sandbox/payment" />
          <FeatureCard title="Webhook Inspector" desc="Event log, endpoint test" href="/sandbox/webhook" />
        </div>
      </section>

      <footer className="border-t border-white/10 py-8 text-center text-xs text-slate-500">
        Gerege API Sandbox &middot; Powered by e-ID Mongolia &middot; sso.gerege.mn
      </footer>
    </main>
  );
}

function FeatureCard({ title, desc, href }: { title: string; desc: string; href: string }) {
  return (
    <Link href={href} className="bg-surface border border-white/10 rounded-xl p-6 hover:border-primary/30 transition-colors block">
      <h3 className="text-white font-semibold mb-2">{title}</h3>
      <p className="text-sm text-slate-400 leading-relaxed">{desc}</p>
    </Link>
  );
}
