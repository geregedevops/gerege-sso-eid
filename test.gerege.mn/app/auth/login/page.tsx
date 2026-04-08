import { signIn } from "@/lib/auth";

export default async function LoginPage({ searchParams }: { searchParams: Promise<{ method?: string }> }) {
  const { method } = await searchParams;
  const ssoURL = process.env.NEXT_PUBLIC_SSO_URL || "https://sso.gerege.mn";
  const clientId = process.env.EID_CLIENT_ID || "";
  const appURL = process.env.NEXT_PUBLIC_APP_URL || "https://test.gerege.mn";

  const danAuthorizeURL = `${ssoURL}/oauth/authorize?client_id=${clientId}&redirect_uri=${encodeURIComponent(appURL + "/api/auth/callback/gerege-sso")}&response_type=code&scope=${encodeURIComponent("openid profile pos social payment")}&auth_method=dan`;

  if (method === "dan") {
    const { redirect } = await import("next/navigation");
    redirect(danAuthorizeURL);
  }

  return (
    <main className="min-h-[calc(100vh-56px)] flex items-center justify-center p-6">
      <div className="max-w-sm w-full bg-surface border border-white/10 rounded-2xl p-8 text-center space-y-6">
        <div className="w-14 h-14 bg-primary/10 rounded-xl flex items-center justify-center mx-auto">
          <span className="text-primary text-2xl font-black">G</span>
        </div>
        <div>
          <h1 className="text-xl font-bold text-white mb-2">API Sandbox</h1>
          <p className="text-sm text-slate-400">e-ID Mongolia-р нэвтрэн Gerege API-г sandbox орчинд туршина.</p>
        </div>
        <form
          action={async () => {
            "use server";
            await signIn("gerege-sso", { redirectTo: "/sandbox" });
          }}
        >
          <button
            type="submit"
            className="w-full py-3 bg-primary text-white font-semibold rounded-xl hover:bg-primary-light transition-colors"
          >
            e-ID Mongolia-р нэвтрэх
          </button>
        </form>
        <div className="relative">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-white/10"></div>
          </div>
          <div className="relative flex justify-center text-xs">
            <span className="bg-surface px-2 text-slate-500">эсвэл</span>
          </div>
        </div>
        <a
          href={danAuthorizeURL}
          className="block w-full py-3 bg-blue-600 text-white font-semibold rounded-xl hover:bg-blue-500 transition-colors"
        >
          DAN нэвтрэх (sso.gov.mn)
        </a>
      </div>
    </main>
  );
}
