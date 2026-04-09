import Link from "next/link";

export default function LandingPage() {
  return (
    <main className="min-h-[calc(100vh-56px)]">
      <section className="max-w-4xl mx-auto px-6 pt-24 pb-16 text-center">
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary/10 border border-primary/20 text-primary text-xs font-semibold mb-6">
          Байгууллагын цахим тамга
        </div>
        <h1 className="text-4xl sm:text-5xl font-extrabold text-white leading-tight mb-4">
          d-business<span className="text-primary">.mn</span>
        </h1>
        <p className="text-lg text-slate-400 max-w-xl mx-auto mb-4">
          Монголын байгууллагуудад зориулсан цахим тамга (e-Seal) платформ.
          Баримт бичигт e-ID-р гарын үсэг зурж, баталгаажуулна.
        </p>
        <p className="text-sm text-slate-500 max-w-lg mx-auto mb-10">
          Estonia-ийн e-Stamp загварыг Монголын e-ID дэд бүтцэд нийцүүлсэн.
          SmartID PIN2-оор баталгаажуулсан цахим гарын үсэг.
        </p>
        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
          <Link href="/auth/login" className="px-8 py-3 bg-primary text-white font-semibold rounded-xl hover:bg-primary-light transition-colors">
            e-ID Mongolia-р нэвтрэх
          </Link>
          <Link href="/verify" className="px-8 py-3 border border-white/15 text-white font-medium rounded-xl hover:bg-white/5 transition-colors">
            Баримт шалгах
          </Link>
        </div>
      </section>

      <section className="max-w-4xl mx-auto px-6 pb-16">
        <div className="flex justify-center gap-3 flex-wrap mb-16">
          {[
            { num: "1", title: "Бүртгэл", desc: "Байгууллага бүртгүүлэх" },
            { num: "2", title: "Сертификат", desc: "e-Seal сертификат авах" },
            { num: "3", title: "Гарын үсэг", desc: "PDF-д SmartID-р зурах" },
            { num: "4", title: "Шалгалт", desc: "Баримтыг баталгаажуулах" },
          ].map((s) => (
            <div key={s.num} className="bg-surface border border-white/10 rounded-xl p-5 w-44 text-center">
              <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center font-bold text-white text-sm mx-auto mb-3">{s.num}</div>
              <h3 className="text-white font-semibold text-sm mb-1">{s.title}</h3>
              <p className="text-xs text-slate-400">{s.desc}</p>
            </div>
          ))}
        </div>
      </section>

      <section className="max-w-5xl mx-auto px-6 pb-24">
        <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-6">
          <FeatureCard title="e-ID баталгаажуулалт" desc="SmartID + X.509 сертификатаар иргэний хувийн мэдээллийг баталгаажуулна." />
          <FeatureCard title="PDF гарын үсэг" desc="Байгууллагын нэрийн тамгаар баримт бичгийг цахимаар баталгаажуулна." />
          <FeatureCard title="Сертификат удирдлага" desc="Байгууллагын e-Seal сертификат олгох, хяналт, хүчингүй болгох." />
          <FeatureCard title="Баримт шалгалт" desc="Гарын үсэг зурсан баримтын бүрэн бүтэн байдлыг шалгана." />
        </div>
      </section>

      <section className="max-w-3xl mx-auto px-6 pb-16">
        <div className="bg-surface border border-white/10 rounded-2xl p-8 text-center">
          <h2 className="text-xl font-bold text-white mb-4">Яаж ажилладаг вэ?</h2>
          <div className="text-sm text-slate-400 space-y-3 text-left max-w-lg mx-auto">
            <p><strong className="text-white">1. Бүртгэл:</strong> e-ID Mongolia-р нэвтэрч байгууллагаа бүртгүүлнэ (нэр, регистрийн дугаар, төрөл).</p>
            <p><strong className="text-white">2. Сертификат:</strong> Байгууллагын e-Seal сертификат хүсэлт илгээнэ. Баталгаажсны дараа цахим тамга ашиглах боломжтой.</p>
            <p><strong className="text-white">3. Гарын үсэг:</strong> PDF баримт upload хийж, SmartID PIN2 оруулан цахим гарын үсэг зурна. api.gerege.mn-р дамжуулан баталгаажуулалт хийгдэнэ.</p>
            <p><strong className="text-white">4. Шалгалт:</strong> Хэн ч гарын үсэг зурсан баримтыг upload хийж бүрэн бүтэн байдлыг шалгаж болно.</p>
          </div>
        </div>
      </section>

      <footer className="border-t border-white/10 py-8 text-center text-xs text-slate-500">
        d-business.mn &middot; Powered by <a href="https://sso.gerege.mn" className="text-primary hover:underline">sso.gerege.mn</a> &middot; <a href="https://gerege.mn" className="text-primary hover:underline">Gerege Systems</a>
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
