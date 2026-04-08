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
      </div>
    </main>
  );
}
