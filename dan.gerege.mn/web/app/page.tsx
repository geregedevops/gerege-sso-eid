export default function HomePage() {
  return (
    <div className="max-w-3xl mx-auto px-6 py-16 text-center">
      <div className="inline-flex items-center gap-2 px-4 py-1.5 bg-yellow-500/10 border border-yellow-500/25 rounded-full text-xs text-yellow-400 font-semibold mb-8">
        Gerege Systems
      </div>
      <h1 className="text-4xl font-extrabold text-white mb-4">
        DAN <span className="text-primary">Verify</span> Gateway
      </h1>
      <p className="text-slate-400 text-base mb-3 max-w-xl mx-auto leading-relaxed">
        sso.gov.mn-ийн ДАН системээр иргэний бүртгэлийн мэдээлэл баталгаажуулах OAuth2 gateway.
        Регистрийн дугаар, нэр, хаяг, зураг зэргийг авна.
      </p>
      <p className="text-slate-500 text-sm mb-10">
        sso.gov.mn OAuth2 &middot; POST callback &middot; HMAC-SHA256
      </p>
      <div className="flex flex-wrap justify-center gap-3">
        <a href="/try" className="px-8 py-3 bg-green-600 hover:bg-green-500 text-white font-bold rounded-xl transition-all shadow-lg shadow-green-600/30">
          DAN Verify
        </a>
        <a href="/docs" className="px-8 py-3 bg-primary hover:bg-primary-light text-white font-bold rounded-xl transition-all">
          Холболтын заавар
        </a>
        <a href="/admin" className="px-8 py-3 border border-white/15 text-white font-bold rounded-xl hover:bg-white/5 transition-all">
          Admin
        </a>
      </div>
      <p className="text-slate-600 text-xs mt-3">
        DAN Verify товч дарвал шууд sso.gov.mn-р нэвтэрч иргэний мэдээллийг харна
      </p>

      <div className="mt-16 grid grid-cols-1 md:grid-cols-3 gap-4 text-left">
        <div className="bg-white/[0.03] border border-white/[0.06] rounded-xl p-5">
          <h3 className="text-sm font-bold text-white mb-1">ДАН баталгаажуулалт</h3>
          <p className="text-xs text-slate-400 leading-relaxed">sso.gov.mn OAuth2 flow-р иргэний бүрэн мэдээлэл авна</p>
        </div>
        <div className="bg-white/[0.03] border border-white/[0.06] rounded-xl p-5">
          <h3 className="text-sm font-bold text-white mb-1">HMAC-SHA256</h3>
          <p className="text-xs text-slate-400 leading-relaxed">Callback дата бүрэн бүтэн, өөрчлөгдөөгүйг signature-р баталгаажуулна</p>
        </div>
        <div className="bg-white/[0.03] border border-white/[0.06] rounded-xl p-5">
          <h3 className="text-sm font-bold text-white mb-1">Зураг + Мэдээлэл</h3>
          <p className="text-xs text-slate-400 leading-relaxed">Иргэний цээж зураг, хаяг, бүртгэл зэрэг бүрэн мэдээлэл POST-р ирнэ</p>
        </div>
      </div>
    </div>
  );
}
