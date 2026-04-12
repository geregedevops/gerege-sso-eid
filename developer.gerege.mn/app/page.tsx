import Link from "next/link";

export default function LandingPage() {
  return (
    <main className="min-h-[calc(100vh-56px)]">
      {/* Hero */}
      <section className="max-w-4xl mx-auto px-6 pt-24 pb-16 text-center">
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary/10 border border-primary/20 text-primary text-xs font-semibold mb-6">
          OpenID Connect Provider
        </div>
        <h1 className="text-4xl sm:text-5xl font-extrabold text-white leading-tight mb-4">
          Gerege <span className="text-primary">Developer Portal</span>
        </h1>
        <p className="text-lg text-slate-400 max-w-xl mx-auto mb-10">
          Gerege platform-д app бүтээж эхлэ.
          POS plugin, social commerce, payment — нэг API-аар.
        </p>
        <div className="flex items-center justify-center gap-4">
          <Link
            href="/auth/login"
            className="px-6 py-3 bg-primary text-white font-semibold rounded-xl hover:bg-primary-light transition-colors"
          >
            GeregeID-р нэвтрэх
          </Link>
          <Link
            href="/docs/quickstart"
            className="px-6 py-3 border border-white/15 text-white font-medium rounded-xl hover:bg-white/5 transition-colors"
          >
            Quickstart
          </Link>
        </div>
      </section>

      {/* Features */}
      <section className="max-w-5xl mx-auto px-6 pb-24">
        <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-6">
          <FeatureCard
            title="POS Plugin API"
            desc="Борлуулалт, захиалга, бараа бүртгэл"
          />
          <FeatureCard
            title="Social Commerce API"
            desc="Лайв худалдаа, product feed, нийтлэл"
          />
          <FeatureCard
            title="Payment API"
            desc="QPay, SocialPay, eBarimt нэгтгэл"
          />
          <FeatureCard
            title="GeregeID нэвтрэлт"
            desc="SmartID + X.509, eIDAS High"
          />
        </div>
      </section>

      {/* Quick code */}
      <section className="max-w-3xl mx-auto px-6 pb-24">
        <h2 className="text-xl font-bold text-white mb-4 text-center">3 мөрөнд нэгтгэх</h2>
        <pre className="bg-surface border border-white/10 rounded-xl p-6 text-sm text-slate-300 overflow-x-auto">
{`// next-auth.config.ts
providers: [{
  id: "gerege-sso",
  name: "GeregeID",
  type: "oidc",
  issuer: "https://sso.gerege.mn",
  clientId: process.env.EID_CLIENT_ID,
  clientSecret: process.env.EID_CLIENT_SECRET,
}]`}
        </pre>
      </section>

      {/* Footer */}
      <footer className="border-t border-white/10 py-8 text-center text-xs text-slate-500">
        Gerege Developer Portal &middot; Powered by GeregeID &middot; sso.gerege.mn
      </footer>
    </main>
  );
}

function FeatureCard({ title, desc }: { title: string; desc: string }) {
  return (
    <div className="bg-surface border border-white/10 rounded-xl p-6 hover:border-primary/30 transition-colors">
      <h3 className="text-white font-semibold mb-2">{title}</h3>
      <p className="text-sm text-slate-400 leading-relaxed">{desc}</p>
    </div>
  );
}
