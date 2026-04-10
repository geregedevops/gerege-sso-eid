import { signIn } from "@/lib/auth";

export default function LoginPage() {
  return (
    <div className="flex min-h-[80vh] items-center justify-center">
      <div className="w-full max-w-sm bg-surface border border-white/[0.06] rounded-2xl p-8 text-center">
        <div className="w-12 h-12 bg-primary rounded-xl flex items-center justify-center text-white font-bold text-sm mx-auto mb-6">
          DAN
        </div>
        <h1 className="text-xl font-bold text-white mb-2">DAN Admin</h1>
        <p className="text-sm text-slate-400 mb-8">
          e-ID Mongolia-р нэвтэрч DAN client удирдлагын самбарт хандана
        </p>
        <form
          action={async () => {
            "use server";
            await signIn("gerege-sso", { redirectTo: "/dashboard" });
          }}
        >
          <button
            type="submit"
            className="w-full py-3 bg-primary hover:bg-primary-light text-white font-bold rounded-xl transition-all"
          >
            e-ID-р нэвтрэх
          </button>
        </form>
      </div>
    </div>
  );
}
