import { signIn } from "@/lib/auth";

export default async function LoginPage() {
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
          href={`https://dan.gerege.mn/verify?client_id=${process.env.DAN_CLIENT_ID || "dan_a088dd8fac47c7aeb9654b9563ac8d67"}&callback_url=${encodeURIComponent((process.env.NEXT_PUBLIC_APP_URL || "https://test.gerege.mn") + "/api/dan/callback")}`}
          className="block w-full py-3 bg-blue-600 text-white font-semibold rounded-xl hover:bg-blue-500 transition-colors"
        >
          DAN Verify (sso.gov.mn)
        </a>
      </div>
    </main>
  );
}
