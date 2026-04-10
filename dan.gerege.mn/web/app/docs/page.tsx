export default function DocsPage() {
  return (
    <div className="max-w-3xl mx-auto px-6 py-10">
      <h1 className="text-2xl font-bold text-white mb-6">DAN Gateway API</h1>

      <div className="space-y-8">
        <section>
          <h2 className="text-lg font-bold text-white mb-3">Ерөнхий ойлголт</h2>
          <p className="text-sm text-slate-400 leading-relaxed mb-4">
            DAN Gateway нь sso.gov.mn-ийн ДАН системтэй OAuth2 authorization code flow-р холбогдож,
            иргэний бүртгэлийн мэдээлэл (регистрийн дугаар, нэр, хаяг, зураг) авч таны callback URL руу
            POST-р дамжуулна.
          </p>
        </section>

        <section>
          <h2 className="text-lg font-bold text-white mb-3">Flow</h2>
          <div className="bg-white/[0.03] border border-white/[0.06] rounded-xl p-5 font-mono text-xs text-slate-300 leading-relaxed space-y-1">
            <p>1. Хэрэглэгчийн browser → GET /verify?client_id=X&callback_url=Y</p>
            <p>2. DAN Gateway → redirect → sso.gov.mn/login</p>
            <p>3. Хэрэглэгч ДАН-р нэвтрэх</p>
            <p>4. sso.gov.mn → redirect → /authorized?code=Z&state=S</p>
            <p>5. DAN Gateway → sso.gov.mn token exchange</p>
            <p>6. DAN Gateway → sso.gov.mn citizen data API</p>
            <p>7. DAN Gateway → POST callback_url (JSON: citizen data + image + signature)</p>
            <p>8. DAN Gateway → redirect browser → callback_url?status=ok&reg_no=...</p>
          </div>
        </section>

        <section>
          <h2 className="text-lg font-bold text-white mb-3">Endpoints</h2>
          <div className="bg-white/[0.03] border border-white/[0.06] rounded-xl overflow-hidden">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-white/[0.06]">
                  <th className="text-left px-4 py-2.5 text-slate-400">Method</th>
                  <th className="text-left px-4 py-2.5 text-slate-400">Path</th>
                  <th className="text-left px-4 py-2.5 text-slate-400">Тайлбар</th>
                </tr>
              </thead>
              <tbody className="text-xs">
                <tr className="border-b border-white/[0.03]">
                  <td className="px-4 py-2.5"><span className="px-2 py-0.5 bg-blue-500/10 text-blue-400 rounded font-mono font-bold">GET</span></td>
                  <td className="px-4 py-2.5 font-mono text-white">/verify</td>
                  <td className="px-4 py-2.5 text-slate-400">DAN flow эхлүүлэх</td>
                </tr>
                <tr className="border-b border-white/[0.03]">
                  <td className="px-4 py-2.5"><span className="px-2 py-0.5 bg-blue-500/10 text-blue-400 rounded font-mono font-bold">GET</span></td>
                  <td className="px-4 py-2.5 font-mono text-white">/authorized</td>
                  <td className="px-4 py-2.5 text-slate-400">sso.gov.mn callback</td>
                </tr>
                <tr>
                  <td className="px-4 py-2.5"><span className="px-2 py-0.5 bg-blue-500/10 text-blue-400 rounded font-mono font-bold">GET</span></td>
                  <td className="px-4 py-2.5 font-mono text-white">/health</td>
                  <td className="px-4 py-2.5 text-slate-400">Health check</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        <section>
          <h2 className="text-lg font-bold text-white mb-3">Callback POST формат</h2>
          <div className="bg-black/30 rounded-xl p-5 font-mono text-xs text-slate-300 leading-relaxed overflow-x-auto">
            <pre>{`{
  "reg_no": "УБ12345678",
  "given_name": "Ганбаатар",
  "family_name": "БАТБОЛД",
  "surname": "Батболд",
  "gender": "male",
  "birth_date": "1990-01-15",
  "aimag_name": "Улаанбаатар",
  "image": "base64...",
  "timestamp": "1712345678",
  "client_id": "dan_abc123",
  "signature": "hmac-sha256-hex"
}`}</pre>
          </div>
        </section>

        <section>
          <h2 className="text-lg font-bold text-white mb-3">HMAC Signature шалгах</h2>
          <p className="text-sm text-slate-400 leading-relaxed mb-3">
            Callback-р ирсэн датаг шалгахдаа <code className="px-1.5 py-0.5 bg-white/5 rounded text-primary text-xs">signature</code> болон{" "}
            <code className="px-1.5 py-0.5 bg-white/5 rounded text-primary text-xs">image</code> field-ийг хасаж,
            бусад field-ийг key-р нь sort хийж <code className="px-1.5 py-0.5 bg-white/5 rounded text-primary text-xs">key=value&key=value</code> форматаар
            HMAC-SHA256 тооцоолно.
          </p>
          <div className="bg-black/30 rounded-xl p-5 font-mono text-xs text-slate-300 leading-relaxed">
            <pre>{`HMAC-SHA256(
  key: your_hmac_key,
  data: "birth_date=1990-01-15&client_id=dan_abc123&..."
) == signature`}</pre>
          </div>
        </section>
      </div>
    </div>
  );
}
