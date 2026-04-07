export default function QuickstartPage() {
  return (
    <main className="max-w-3xl mx-auto px-6 py-12">
      <h1 className="text-3xl font-bold text-white mb-8">5 минутад нэгтгэх</h1>
      <div className="prose prose-invert max-w-none space-y-8">
        <Step n={1} title="App бүртгүүлэх">
          <p className="text-slate-400">
            <a href="/dashboard/apps/new" className="text-primary hover:underline">developer.gerege.mn/dashboard/apps/new</a> руу орж app үүсгэнэ.
            <code className="text-sm">client_id</code> болон <code className="text-sm">client_secret</code> авна.
          </p>
        </Step>

        <Step n={2} title="OIDC Discovery">
          <Code>{`GET https://sso.gerege.mn/.well-known/openid-configuration`}</Code>
        </Step>

        <Step n={3} title="Next.js + NextAuth">
          <Code>{`// lib/auth.ts
import NextAuth from "next-auth"

export const { handlers, signIn, signOut, auth } = NextAuth({
  providers: [{
    id: "gerege-sso",
    name: "e-ID Mongolia",
    type: "oidc",
    issuer: "https://sso.gerege.mn",
    clientId: process.env.EID_CLIENT_ID,
    clientSecret: process.env.EID_CLIENT_SECRET,
  }],
})`}</Code>
        </Step>

        <Step n={4} title="Go + golang.org/x/oauth2">
          <Code>{`provider, _ := oidc.NewProvider(ctx, "https://sso.gerege.mn")
config := oauth2.Config{
    ClientID:     os.Getenv("EID_CLIENT_ID"),
    ClientSecret: os.Getenv("EID_CLIENT_SECRET"),
    Endpoint:     provider.Endpoint(),
    RedirectURL:  "https://myapp.mn/callback",
    Scopes:       []string{oidc.ScopeOpenID, "profile", "pos"},
}`}</Code>
        </Step>

        <Step n={5} title="Тест хийх">
          <p className="text-slate-400">
            App-аа ажиллуулж &quot;e-ID Mongolia-р нэвтрэх&quot; товч дарна. SmartID апп-д push ирнэ.
            PIN1 оруулж баталгаажуулсны дараа ID Token + access_token авна.
          </p>
        </Step>
      </div>
    </main>
  );
}

function Step({ n, title, children }: { n: number; title: string; children: React.ReactNode }) {
  return (
    <div>
      <h2 className="text-lg font-semibold text-white flex items-center gap-3 mb-3">
        <span className="w-7 h-7 bg-primary/20 text-primary text-sm font-bold rounded-full flex items-center justify-center">{n}</span>
        {title}
      </h2>
      {children}
    </div>
  );
}

function Code({ children }: { children: string }) {
  return (
    <pre className="bg-bg border border-white/10 rounded-xl p-4 text-sm text-slate-300 overflow-x-auto">
      {children}
    </pre>
  );
}
