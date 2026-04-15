export default function LoginPage({
  searchParams,
}: {
  searchParams: { error?: string };
}) {
  return (
    <div className="flex min-h-[80vh] items-center justify-center">
      <div className="w-full max-w-sm bg-surface border border-white/[0.06] rounded-2xl p-8">
        <div className="text-center">
          <div className="w-12 h-12 bg-primary rounded-xl flex items-center justify-center text-white font-bold text-sm mx-auto mb-6">
            VFY
          </div>
          <h1 className="text-xl font-bold text-white mb-2">Verify Admin</h1>
          <p className="text-sm text-slate-400 mb-8">
            Баталгаажуулах API удирдлагын самбарт нэвтрэх
          </p>
        </div>
        <form action="/api/auth/login" method="post" className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1.5">
              Admin password
            </label>
            <input
              type="password"
              name="password"
              required
              autoComplete="current-password"
              className="w-full px-4 py-2.5 bg-white/[0.05] border border-white/10 rounded-xl text-white text-sm focus:outline-none focus:border-primary"
            />
          </div>
          {searchParams.error && (
            <p className="text-xs text-red-400">Буруу нууц үг</p>
          )}
          <button
            type="submit"
            className="w-full py-3 bg-primary hover:bg-primary-light text-white font-bold rounded-xl transition-all"
          >
            Нэвтрэх
          </button>
        </form>
      </div>
    </div>
  );
}
