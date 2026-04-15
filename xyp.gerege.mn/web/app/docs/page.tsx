export default function DocsPage() {
  return (
    <div className="max-w-3xl mx-auto px-6 py-10">
      <h1 className="text-3xl font-bold text-white mb-2">API Заавар</h1>
      <p className="text-slate-400 mb-10">xyp.gerege.mn REST API ашиглах заавар</p>

      {/* Auth */}
      <section className="mb-10">
        <h2 className="text-xl font-bold text-white mb-4">Баталгаажуулалт</h2>
        <p className="text-sm text-slate-400 mb-3">
          Бүх verification endpoint-ууд HTTP Basic Auth шаарддаг.
          Admin-аас авсан <code className="text-primary">client_id</code> болон <code className="text-primary">client_secret</code>-ээ ашиглана.
        </p>
        <div className="bg-black/30 rounded-xl p-4 font-mono text-sm text-slate-300 overflow-x-auto">
          Authorization: Basic base64(client_id:client_secret)
        </div>
      </section>

      {/* Citizen Lookup */}
      <section className="mb-10">
        <h2 className="text-xl font-bold text-white mb-4">Иргэн хайх</h2>
        <div className="flex items-center gap-2 mb-3">
          <span className="px-2 py-0.5 bg-blue-500/10 text-blue-400 text-xs rounded font-mono font-bold">POST</span>
          <code className="text-primary text-sm">/v1/citizen/lookup</code>
        </div>
        <p className="text-sm text-slate-400 mb-3">Регистрийн дугаараар иргэний мэдээлэл хайна.</p>
        <h4 className="text-sm font-semibold text-slate-300 mb-2">Request</h4>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-green-400 font-mono overflow-x-auto mb-3">{`curl -u $CLIENT_ID:$SECRET \\
  -X POST https://xyp.gerege.mn/v1/citizen/lookup \\
  -H "Content-Type: application/json" \\
  -d '{"reg_no": "МА74101813"}'`}</pre>
        <h4 className="text-sm font-semibold text-slate-300 mb-2">Response</h4>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-slate-300 font-mono overflow-x-auto">{`{
  "found": true,
  "citizen": {
    "reg_no": "ма74101813",
    "last_name": "Цэнддорж",
    "first_name": "ЭРДЭНЭБАТ",
    "surname": "Харчин",
    "gender": "Эрэгтэй",
    "birth_date": "1974-10-18 00:00",
    "nationality": "Халх"
  }
}`}</pre>
      </section>

      {/* Citizen Verify */}
      <section className="mb-10">
        <h2 className="text-xl font-bold text-white mb-4">Иргэний нэр тулгах</h2>
        <div className="flex items-center gap-2 mb-3">
          <span className="px-2 py-0.5 bg-blue-500/10 text-blue-400 text-xs rounded font-mono font-bold">POST</span>
          <code className="text-primary text-sm">/v1/citizen/verify</code>
        </div>
        <p className="text-sm text-slate-400 mb-3">Регистрийн дугаар + нэрийг тулгаж шалгана. Зөвхөн match true/false буцаана.</p>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-green-400 font-mono overflow-x-auto mb-3">{`curl -u $CLIENT_ID:$SECRET \\
  -X POST https://xyp.gerege.mn/v1/citizen/verify \\
  -H "Content-Type: application/json" \\
  -d '{"reg_no":"МА74101813","first_name":"ЭРДЭНЭБАТ","last_name":"Цэнддорж"}'`}</pre>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-slate-300 font-mono overflow-x-auto">{`{ "match": true, "reg_no": "МА74101813" }`}</pre>
      </section>

      {/* Org Lookup */}
      <section className="mb-10">
        <h2 className="text-xl font-bold text-white mb-4">Байгууллага хайх</h2>
        <div className="flex items-center gap-2 mb-3">
          <span className="px-2 py-0.5 bg-blue-500/10 text-blue-400 text-xs rounded font-mono font-bold">POST</span>
          <code className="text-primary text-sm">/v1/org/lookup</code>
        </div>
        <p className="text-sm text-slate-400 mb-3">Регистрийн дугаараар байгууллагын мэдээлэл хайна.</p>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-green-400 font-mono overflow-x-auto mb-3">{`curl -u $CLIENT_ID:$SECRET \\
  -X POST https://xyp.gerege.mn/v1/org/lookup \\
  -H "Content-Type: application/json" \\
  -d '{"reg_no": "6235972"}'`}</pre>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-slate-300 font-mono overflow-x-auto">{`{
  "found": true,
  "organization": {
    "reg_no": "6235972",
    "name": "Гэрэгэ системс",
    "type": "Хязгаарлагдмал хариуцлагатай компани",
    "ceo": "Нацагдорж Энхжаргал",
    "phone": "99102856",
    "address": "Улаанбаатар, Сүхбаатар, ...",
    "industry": ["Программ хангамжийн үйлчилгээ", ...]
  }
}`}</pre>
      </section>

      {/* Org Verify */}
      <section className="mb-10">
        <h2 className="text-xl font-bold text-white mb-4">Байгууллагын нэр тулгах</h2>
        <div className="flex items-center gap-2 mb-3">
          <span className="px-2 py-0.5 bg-blue-500/10 text-blue-400 text-xs rounded font-mono font-bold">POST</span>
          <code className="text-primary text-sm">/v1/org/verify</code>
        </div>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-green-400 font-mono overflow-x-auto mb-3">{`curl -u $CLIENT_ID:$SECRET \\
  -X POST https://xyp.gerege.mn/v1/org/verify \\
  -H "Content-Type: application/json" \\
  -d '{"reg_no":"6235972","name":"Гэрэгэ системс"}'`}</pre>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-slate-300 font-mono overflow-x-auto">{`{ "match": true, "reg_no": "6235972" }`}</pre>
      </section>

      {/* Error codes */}
      <section className="mb-10">
        <h2 className="text-xl font-bold text-white mb-4">Алдааны кодууд</h2>
        <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-white/[0.06]">
                <th className="text-left px-5 py-3 text-slate-400 font-semibold">Код</th>
                <th className="text-left px-5 py-3 text-slate-400 font-semibold">Тайлбар</th>
              </tr>
            </thead>
            <tbody className="text-slate-300">
              <tr className="border-b border-white/[0.03]">
                <td className="px-5 py-2 font-mono text-yellow-400">401</td>
                <td className="px-5 py-2">Буруу client_id эсвэл client_secret</td>
              </tr>
              <tr className="border-b border-white/[0.03]">
                <td className="px-5 py-2 font-mono text-yellow-400">403</td>
                <td className="px-5 py-2">Client идэвхгүйжсэн эсвэл scope хүрэлцэхгүй</td>
              </tr>
              <tr className="border-b border-white/[0.03]">
                <td className="px-5 py-2 font-mono text-yellow-400">429</td>
                <td className="px-5 py-2">Rate limit хэтэрсэн (default: 100 req/min)</td>
              </tr>
              <tr>
                <td className="px-5 py-2 font-mono text-yellow-400">502</td>
                <td className="px-5 py-2">Гадаад системийн алдаа</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <div className="text-center text-sm text-slate-500 pb-10">
        <a href="/" className="hover:text-white">← Нүүр хуудас руу буцах</a>
      </div>
    </div>
  );
}
