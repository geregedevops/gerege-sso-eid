export default function NextJSGuidePage() {
  return (
    <main className="max-w-3xl mx-auto px-6 py-12 space-y-8">
      <h1 className="text-3xl font-bold text-white">Next.js + NextAuth.js</h1>
      <p className="text-slate-400">Next.js App Router дээр sso.gerege.mn OIDC нэгтгэх.</p>

      <Code title="1. Install">{`npm install next-auth@beta`}</Code>

      <Code title="2. lib/auth.ts">{`import NextAuth from "next-auth"

export const { handlers, signIn, signOut, auth } = NextAuth({
  providers: [{
    id: "gerege-sso",
    name: "e-ID Mongolia",
    type: "oidc",
    issuer: "https://sso.gerege.mn",
    clientId: process.env.EID_CLIENT_ID!,
    clientSecret: process.env.EID_CLIENT_SECRET!,
    authorization: { params: { scope: "openid profile pos" } },
  }],
  callbacks: {
    async jwt({ token, profile }) {
      if (profile) {
        token.sub = profile.sub
        token.name = profile.name
        token.certSerial = profile.cert_serial
        token.tenantId = profile.tenant_id
        token.tenantRole = profile.tenant_role
      }
      return token
    },
  },
})`}</Code>

      <Code title="3. app/api/auth/[...nextauth]/route.ts">{`import { handlers } from "@/lib/auth"
export const { GET, POST } = handlers`}</Code>

      <Code title="4. .env.local">{`NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-random-secret
EID_CLIENT_ID=your-client-id
EID_CLIENT_SECRET=your-client-secret`}</Code>

      <Code title="5. Sign In Button">{`import { signIn } from "@/lib/auth"

export default function LoginPage() {
  return (
    <form action={async () => {
      "use server"
      await signIn("gerege-sso")
    }}>
      <button type="submit">e-ID Mongolia-р нэвтрэх</button>
    </form>
  )
}`}</Code>

      <Code title="6. Protected Page">{`import { auth } from "@/lib/auth"
import { redirect } from "next/navigation"

export default async function DashboardPage() {
  const session = await auth()
  if (!session) redirect("/api/auth/signin")

  return <h1>Hello {session.user?.name}</h1>
}`}</Code>
    </main>
  );
}

function Code({ title, children }: { title: string; children: string }) {
  return (
    <div>
      <h3 className="text-sm font-semibold text-white mb-2">{title}</h3>
      <pre className="bg-bg border border-white/10 rounded-xl p-4 text-xs text-slate-300 overflow-x-auto">{children}</pre>
    </div>
  );
}
