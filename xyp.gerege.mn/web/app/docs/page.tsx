export default function DocsPage() {
  return (
    <div className="max-w-3xl mx-auto px-6 py-10">
      <h1 className="text-3xl font-bold text-white mb-2">API Заавар</h1>
      <p className="text-slate-400 mb-10">xyp.gerege.mn REST API ашиглах заавар</p>

      {/* Auth */}
      <section className="mb-10">
        <h2 className="text-xl font-bold text-white mb-4">Баталгаажуулалт</h2>
        <p className="text-sm text-slate-400 mb-3">
          Бүх API endpoint-ууд HTTP Basic Auth шаарддаг.
          Admin-аас авсан <code className="text-primary">client_id</code> болон <code className="text-primary">client_secret</code>-ээ ашиглана.
        </p>
        <pre className="bg-black/30 rounded-xl p-4 font-mono text-sm text-slate-300 overflow-x-auto">{`Authorization: Basic base64(client_id:client_secret)`}</pre>
      </section>

      {/* Endpoints table */}
      <section className="mb-10">
        <h2 className="text-xl font-bold text-white mb-4">Endpoint-ууд</h2>
        <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-white/[0.06]">
                <th className="text-left px-5 py-3 text-slate-400 font-semibold">Method</th>
                <th className="text-left px-5 py-3 text-slate-400 font-semibold">Endpoint</th>
                <th className="text-left px-5 py-3 text-slate-400 font-semibold">Тайлбар</th>
              </tr>
            </thead>
            <tbody className="text-slate-300">
              <tr className="border-b border-white/[0.03]">
                <td className="px-5 py-2"><Post /></td>
                <td className="px-5 py-2 font-mono text-xs text-primary">/v1/citizen/authenticate</td>
                <td className="px-5 py-2 text-sm">Иргэн баталгаажуулах (РД + утас)</td>
              </tr>
              <tr className="border-b border-white/[0.03]">
                <td className="px-5 py-2"><Post /></td>
                <td className="px-5 py-2 font-mono text-xs text-primary">/v1/org/authenticate</td>
                <td className="px-5 py-2 text-sm">Байгууллага баталгаажуулах (регистр + захирлын РД)</td>
              </tr>
              <tr className="border-b border-white/[0.03] bg-white/[0.01]">
                <td className="px-5 py-2"><Post /></td>
                <td className="px-5 py-2 font-mono text-xs text-slate-400">/v1/citizen/lookup</td>
                <td className="px-5 py-2 text-sm text-slate-500">Иргэний мэдээлэл хайх</td>
              </tr>
              <tr className="border-b border-white/[0.03] bg-white/[0.01]">
                <td className="px-5 py-2"><Post /></td>
                <td className="px-5 py-2 font-mono text-xs text-slate-400">/v1/citizen/verify</td>
                <td className="px-5 py-2 text-sm text-slate-500">Иргэний нэр тулгах</td>
              </tr>
              <tr className="border-b border-white/[0.03] bg-white/[0.01]">
                <td className="px-5 py-2"><Post /></td>
                <td className="px-5 py-2 font-mono text-xs text-slate-400">/v1/org/lookup</td>
                <td className="px-5 py-2 text-sm text-slate-500">Байгууллагын мэдээлэл хайх</td>
              </tr>
              <tr className="bg-white/[0.01]">
                <td className="px-5 py-2"><Post /></td>
                <td className="px-5 py-2 font-mono text-xs text-slate-400">/v1/org/verify</td>
                <td className="px-5 py-2 text-sm text-slate-500">Байгууллагын нэр тулгах</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <hr className="border-white/[0.06] mb-10" />

      {/* Citizen Authenticate */}
      <section className="mb-10">
        <h2 className="text-xl font-bold text-white mb-4">Иргэн баталгаажуулах</h2>
        <div className="flex items-center gap-2 mb-3">
          <Post />
          <code className="text-primary text-sm">/v1/citizen/authenticate</code>
        </div>
        <p className="text-sm text-slate-400 mb-3">
          Иргэний регистрийн дугаар + утасны дугаараар баталгаажуулна.
          Амжилттай бол иргэний зураг болон үндсэн мэдээлэл буцаана.
        </p>
        <h4 className="text-sm font-semibold text-slate-300 mb-2">Request</h4>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-green-400 font-mono overflow-x-auto mb-3">{`curl -u $CLIENT_ID:$SECRET \\
  -X POST https://xyp.gerege.mn/v1/citizen/authenticate \\
  -H "Content-Type: application/json" \\
  -d '{"reg_no": "МА74101813", "phone": "99102856"}'`}</pre>
        <h4 className="text-sm font-semibold text-slate-300 mb-2">Response (амжилттай)</h4>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-slate-300 font-mono overflow-x-auto mb-3">{`{
  "authenticated": true,
  "citizen": {
    "reg_no": "ма74101813",
    "civil_id": "111949212017",
    "last_name": "Цэнддорж",
    "first_name": "ЭРДЭНЭБАТ",
    "gender": "Эрэгтэй",
    "birth_date": "1974-10-18 00:00",
    "image": "/9j/4AAQSkZJRg..."
  }
}`}</pre>
        <h4 className="text-sm font-semibold text-slate-300 mb-2">Response (олдоогүй)</h4>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-slate-300 font-mono overflow-x-auto">{`{
  "authenticated": false,
  "reason": "citizen not found"
}`}</pre>
      </section>

      {/* Org Authenticate */}
      <section className="mb-10">
        <h2 className="text-xl font-bold text-white mb-4">Байгууллага баталгаажуулах</h2>
        <div className="flex items-center gap-2 mb-3">
          <Post />
          <code className="text-primary text-sm">/v1/org/authenticate</code>
        </div>
        <p className="text-sm text-slate-400 mb-3">
          Байгууллагын регистр + захирлын регистрийн дугаараар баталгаажуулна.
          Захирлын РД таарвал байгууллагын мэдээллийг буцаана.
        </p>
        <h4 className="text-sm font-semibold text-slate-300 mb-2">Request</h4>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-green-400 font-mono overflow-x-auto mb-3">{`curl -u $CLIENT_ID:$SECRET \\
  -X POST https://xyp.gerege.mn/v1/org/authenticate \\
  -H "Content-Type: application/json" \\
  -d '{"reg_no": "6235972", "ceo_reg_no": "уш72060800"}'`}</pre>
        <h4 className="text-sm font-semibold text-slate-300 mb-2">Response (амжилттай)</h4>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-slate-300 font-mono overflow-x-auto mb-3">{`{
  "authenticated": true,
  "organization": {
    "reg_no": "6235972",
    "name": "Гэрэгэ системс",
    "type": "Хязгаарлагдмал хариуцлагатай компани",
    "ceo": "Нацагдорж Энхжаргал",
    "ceo_reg_no": "уш72060800",
    "ceo_position": "Гүйцэтгэх  Захирал"
  }
}`}</pre>
        <h4 className="text-sm font-semibold text-slate-300 mb-2">Response (таараагүй)</h4>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-slate-300 font-mono overflow-x-auto">{`{
  "authenticated": false,
  "reason": "ceo_reg_no does not match"
}`}</pre>
      </section>

      <hr className="border-white/[0.06] mb-10" />

      {/* Citizen Lookup */}
      <section className="mb-10">
        <h2 className="text-xl font-bold text-white mb-4">Иргэн хайх</h2>
        <div className="flex items-center gap-2 mb-3">
          <Post />
          <code className="text-primary text-sm">/v1/citizen/lookup</code>
        </div>
        <p className="text-sm text-slate-400 mb-3">Регистрийн дугаараар иргэний бүрэн мэдээлэл хайна.</p>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-green-400 font-mono overflow-x-auto mb-3">{`curl -u $CLIENT_ID:$SECRET \\
  -X POST https://xyp.gerege.mn/v1/citizen/lookup \\
  -H "Content-Type: application/json" \\
  -d '{"reg_no": "МА74101813"}'`}</pre>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-slate-300 font-mono overflow-x-auto">{`{
  "found": true,
  "citizen": {
    "reg_no": "ма74101813",
    "last_name": "Цэнддорж",
    "first_name": "ЭРДЭНЭБАТ",
    "surname": "Харчин",
    "gender": "Эрэгтэй",
    "birth_date": "1974-10-18 00:00",
    "birth_place": "Улаанбаатар,Сүхбаатар",
    "nationality": "Халх",
    "civil_id": "111949212017",
    "passport_num": "PE 0305079",
    "passport_address": "УБ, Хан-Уул, ...",
    "image": "/9j/4AAQSkZJRg..."
  }
}`}</pre>
      </section>

      {/* Citizen Verify */}
      <section className="mb-10">
        <h2 className="text-xl font-bold text-white mb-4">Иргэний нэр тулгах</h2>
        <div className="flex items-center gap-2 mb-3">
          <Post />
          <code className="text-primary text-sm">/v1/citizen/verify</code>
        </div>
        <p className="text-sm text-slate-400 mb-3">РД + нэрийг тулгаж шалгана. Зөвхөн match true/false буцаана.</p>
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
          <Post />
          <code className="text-primary text-sm">/v1/org/lookup</code>
        </div>
        <p className="text-sm text-slate-400 mb-3">Регистрийн дугаараар байгууллагын бүрэн мэдээлэл хайна.</p>
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
    "capital": "2630666",
    "ceo": "Нацагдорж Энхжаргал",
    "ceo_reg_no": "уш72060800",
    "ceo_position": "Гүйцэтгэх  Захирал",
    "phone": "99102856",
    "address": "Улаанбаатар, Сүхбаатар, ...",
    "industry": ["Программ хангамжийн үйлчилгээ", ...],
    "founders": [
      {"name": "цэнддорж эрдэнэбат", "reg_no": "ма74101813",
       "type": "Иргэн", "share_percent": "42"}, ...
    ],
    "stake_holders": [
      {"name": "Цэнддорж Эрдэнэбат", "reg_no": "ма74101813",
       "position": "ТУЗ-ийн дарга"}, ...
    ]
  }
}`}</pre>
      </section>

      {/* Org Verify */}
      <section className="mb-10">
        <h2 className="text-xl font-bold text-white mb-4">Байгууллагын нэр тулгах</h2>
        <div className="flex items-center gap-2 mb-3">
          <Post />
          <code className="text-primary text-sm">/v1/org/verify</code>
        </div>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-green-400 font-mono overflow-x-auto mb-3">{`curl -u $CLIENT_ID:$SECRET \\
  -X POST https://xyp.gerege.mn/v1/org/verify \\
  -H "Content-Type: application/json" \\
  -d '{"reg_no":"6235972","name":"Гэрэгэ системс"}'`}</pre>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-slate-300 font-mono overflow-x-auto">{`{ "match": true, "reg_no": "6235972" }`}</pre>
      </section>

      <hr className="border-white/[0.06] mb-10" />

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
                <td className="px-5 py-2 font-mono text-yellow-400">400</td>
                <td className="px-5 py-2">Буруу хүсэлт (шаардлагатай талбар дутуу)</td>
              </tr>
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
                <td className="px-5 py-2">Гадаад системийн алдаа (XYP хариу өгөхгүй)</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      {/* Rate limit */}
      <section className="mb-10">
        <h2 className="text-xl font-bold text-white mb-4">Rate Limiting</h2>
        <p className="text-sm text-slate-400 mb-3">
          Client тус бүр минутад 100 хүсэлт илгээх боломжтой (тохируулж болно).
          Хэтэрсэн тохиолдолд <code className="text-yellow-400">429</code> статус кодтой хариу буцаана.
        </p>
        <pre className="bg-black/30 rounded-xl p-4 text-sm text-slate-300 font-mono overflow-x-auto">{`HTTP/1.1 429 Too Many Requests
Retry-After: 60
{"error": "rate limit exceeded"}`}</pre>
      </section>

      <div className="text-center text-sm text-slate-500 pb-10">
        <a href="/" className="hover:text-white">← Нүүр хуудас руу буцах</a>
      </div>
    </div>
  );
}

function Post() {
  return (
    <span className="px-2 py-0.5 bg-blue-500/10 text-blue-400 text-xs rounded font-mono font-bold">POST</span>
  );
}
