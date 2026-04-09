export default function DANIntegrationGuidePage() {
  return (
    <main className="max-w-3xl mx-auto px-6 py-12 space-y-10">
      <div>
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-500/10 border border-blue-500/20 text-blue-400 text-xs font-semibold mb-4">
          DAN Gateway
        </div>
        <h1 className="text-3xl font-bold text-white mb-2">DAN Verify нэгтгэх заавар</h1>
        <p className="text-slate-400">
          dan.gerege.mn gateway-р дамжуулан sso.gov.mn-ийн ДАН системээс
          иргэний бүртгэлийн мэдээлэл (регистрийн дугаар, нэр, хаяг, зураг) авна.
        </p>
      </div>

      <Section title="SSO нэвтрэлт vs DAN Verify">
        <div className="grid sm:grid-cols-2 gap-4">
          <div className="p-4 bg-primary/5 border border-primary/15 rounded-xl">
            <h4 className="font-semibold text-white text-sm mb-2">SSO нэвтрэлт (sso.gerege.mn)</h4>
            <ul className="text-xs text-slate-400 space-y-1 list-disc list-inside">
              <li>e-ID Mongolia смарт картаар нэвтрэх</li>
              <li>OIDC стандарт (access_token, id_token)</li>
              <li>Session удирдлага</li>
              <li>Зориулалт: app нэвтрэлт</li>
            </ul>
          </div>
          <div className="p-4 bg-blue-500/5 border border-blue-500/15 rounded-xl">
            <h4 className="font-semibold text-white text-sm mb-2">DAN Verify (dan.gerege.mn)</h4>
            <ul className="text-xs text-slate-400 space-y-1 list-disc list-inside">
              <li>sso.gov.mn ДАН системээр баталгаажуулах</li>
              <li>Иргэний бүтэн мэдээлэл (РД, хаяг, зураг)</li>
              <li>Нэг удаагийн шалгалт</li>
              <li>Зориулалт: KYC, бүртгэл, баталгаажуулалт</li>
            </ul>
          </div>
        </div>
      </Section>

      <Section title="Хоёр горим">
        <div className="space-y-4">
          <div className="p-4 bg-surface border border-white/10 rounded-xl">
            <h4 className="font-semibold text-white text-sm mb-1">Горим 1: DAN Verify (зургүй)</h4>
            <p className="text-xs text-slate-400">
              <code className="text-primary">/verify</code> — Иргэний мэдээлэл callback URL-д query param-р дамжина.
              Зураг дамжуулахгүй (URL-д багтахгүй).
            </p>
          </div>
          <div className="p-4 bg-surface border border-white/10 rounded-xl">
            <h4 className="font-semibold text-white text-sm mb-1">Горим 2: DAN Verify Full (зураг бүхий)</h4>
            <p className="text-xs text-slate-400">
              <code className="text-primary">/verify-full</code> — Callback URL-д зөвхөн <code className="text-blue-400">token</code> дамжина.
              Бүтэн data + зургийг <code className="text-primary">/api/citizen?token=xxx</code> endpoint-р JSON-р авна.
            </p>
          </div>
        </div>
      </Section>

      <Section title="Алхам 1: Client бүртгэл">
        <p>
          DAN Admin-аас client бүртгүүлнэ:{" "}
          <a href="https://dan.gerege.mn/admin" className="text-primary hover:underline font-semibold">dan.gerege.mn/admin</a>
        </p>
        <div className="mt-3 space-y-2">
          <Field label="Name" desc="Таны системийн нэр" />
          <Field label="Callback URLs" desc="DAN-аас мэдээлэл буцаах URL. Жишээ: https://myapp.mn/api/dan/callback" />
        </div>
        <p className="mt-3 text-xs text-slate-400">
          Бүртгүүлсний дараа <code className="text-primary">client_id</code> болон <code className="text-primary">client_secret</code> авна.
          HMAC signature шалгахад secret ашиглана.
        </p>
      </Section>

      <Section title="Алхам 2: Хэрэглэгчийг DAN руу чиглүүлэх">
        <p className="font-semibold text-white text-sm mb-2">Горим 1 — Зургүй:</p>
        <Code>{`GET https://dan.gerege.mn/verify
  ?client_id=YOUR_CLIENT_ID
  &callback_url=https://myapp.mn/api/dan/callback`}</Code>

        <p className="font-semibold text-white text-sm mb-2 mt-4">Горим 2 — Зураг бүхий:</p>
        <Code>{`GET https://dan.gerege.mn/verify-full
  ?client_id=YOUR_CLIENT_ID
  &callback_url=https://myapp.mn/api/dan/callback-full`}</Code>
      </Section>

      <Section title="Алхам 3: Callback хүлээн авах">
        <p className="font-semibold text-white text-sm mb-2">Горим 1 — Зургүй callback:</p>
        <p className="text-xs text-slate-400 mb-2">Иргэний мэдээлэл query param-р дамжина:</p>
        <Code>{`GET https://myapp.mn/api/dan/callback
  ?reg_no=РД98012345
  &given_name=ГАНБААТАР
  &family_name=БАТБОЛД
  &civil_id=12345678
  &gender=male
  &birth_date=1998-01-23
  &aimag_name=Улаанбаатар
  &sum_name=Баянгол
  &timestamp=1744200000
  &client_id=YOUR_CLIENT_ID
  &signature=HMAC_SHA256_SIGNATURE`}</Code>

        <p className="font-semibold text-white text-sm mb-2 mt-6">Горим 2 — Зураг бүхий callback:</p>
        <p className="text-xs text-slate-400 mb-2">Зөвхөн token дамжина:</p>
        <Code>{`GET https://myapp.mn/api/dan/callback-full
  ?token=ONE_TIME_TOKEN
  &timestamp=1744200000
  &client_id=YOUR_CLIENT_ID
  &signature=HMAC_SHA256_SIGNATURE`}</Code>
        <p className="text-xs text-slate-400 mt-2">Дараа нь token ашиглан бүтэн data авна:</p>
        <Code>{`GET https://dan.gerege.mn/api/citizen?token=ONE_TIME_TOKEN

// Хариу:
{
  "success": true,
  "citizen": {
    "reg_no": "РД98012345",
    "given_name": "ГАНБААТАР",
    "family_name": "БАТБОЛД",
    "civil_id": "12345678",
    "gender": "male",
    "birth_date": "1998-01-23",
    "aimag_name": "Улаанбаатар",
    "sum_name": "Баянгол",
    "bag_name": "5-р хороо",
    "address_detail": "Баянгол, 5-р хороо, 25-р байр",
    "image": "base64_encoded_jpeg..."
  }
}`}</Code>
        <div className="mt-3 p-3 bg-amber-500/5 border border-amber-500/15 rounded-xl text-xs">
          <p className="text-amber-400 font-semibold">Token зөвхөн нэг удаа, 5 минутад хүчинтэй.</p>
        </div>
      </Section>

      <Section title="Алхам 4: HMAC Signature шалгах">
        <p>Callback дээр ирсэн мэдээллийн бүрэн бүтэн байдлыг HMAC-SHA256 signature-р шалгана:</p>
        <Code>{`// Go
func verifySignature(params url.Values, secret string) bool {
    expected := params.Get("signature")
    keys := []string{}
    for k := range params {
        if k != "signature" { keys = append(keys, k) }
    }
    sort.Strings(keys)

    var buf strings.Builder
    for i, k := range keys {
        if i > 0 { buf.WriteByte('&') }
        buf.WriteString(url.QueryEscape(k) + "=" + url.QueryEscape(params.Get(k)))
    }

    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(buf.String()))
    return hex.EncodeToString(mac.Sum(nil)) == expected
}`}</Code>

        <Code>{`# Python
import hmac, hashlib, urllib.parse

def verify(params: dict, secret: str) -> bool:
    sig = params.pop("signature", "")
    canonical = "&".join(
        f"{urllib.parse.quote(k)}={urllib.parse.quote(params[k])}"
        for k in sorted(params)
    )
    expected = hmac.new(
        secret.encode(), canonical.encode(), hashlib.sha256
    ).hexdigest()
    return sig == expected`}</Code>

        <Code>{`// Node.js
const crypto = require("crypto");
function verify(params, secret) {
  const sig = params.signature;
  const keys = Object.keys(params).filter(k => k !== "signature").sort();
  const canonical = keys
    .map(k => encodeURIComponent(k) + "=" + encodeURIComponent(params[k]))
    .join("&");
  const expected = crypto
    .createHmac("sha256", secret)
    .update(canonical)
    .digest("hex");
  return sig === expected;
}`}</Code>
      </Section>

      <Section title="Callback параметрүүд (бүрэн жагсаалт)">
        <table className="w-full text-xs">
          <thead>
            <tr className="border-b border-white/10">
              <th className="text-left py-2 text-slate-400">Параметр</th>
              <th className="text-left py-2 text-slate-400">Тайлбар</th>
            </tr>
          </thead>
          <tbody className="text-slate-300">
            <Row label="reg_no" value="Регистрийн дугаар" />
            <Row label="given_name" value="Нэр" />
            <Row label="family_name" value="Овог" />
            <Row label="surname" value="Ургийн овог" />
            <Row label="civil_id" value="Иргэний ID" />
            <Row label="gender" value="Хүйс (male/female)" />
            <Row label="birth_date" value="Төрсөн огноо" />
            <Row label="birth_place" value="Төрсөн газар" />
            <Row label="nationality" value="Үндэс/угсаа" />
            <Row label="aimag_name" value="Аймаг/Хот" />
            <Row label="sum_name" value="Сум/Дүүрэг" />
            <Row label="bag_name" value="Баг/Хороо" />
            <Row label="address_detail" value="Дэлгэрэнгүй хаяг" />
            <Row label="image" value="Зураг (base64 JPEG) — зөвхөн /verify-full горимд" />
            <Row label="timestamp" value="Unix timestamp" />
            <Row label="client_id" value="Таны client ID" />
            <Row label="signature" value="HMAC-SHA256 signature" />
          </tbody>
        </table>
      </Section>

      <Section title="Аюулгүй байдал">
        <ul className="list-disc list-inside text-slate-400 space-y-2 text-xs">
          <li><strong className="text-white">HMAC шалгалт заавал</strong> — Signature тулгахгүй бол мэдээллийг хүлээж авахгүй.</li>
          <li><strong className="text-white">Timestamp шалгалт</strong> — 5 минутаас хэтэрсэн timestamp бүхий хүсэлтийг татгалзана (replay attack).</li>
          <li><strong className="text-white">HTTPS заавал</strong> — Callback URL заавал HTTPS байна.</li>
          <li><strong className="text-white">Token нэг удаа</strong> — /verify-full token нэг удаа fetch хийгдэнэ.</li>
        </ul>
      </Section>

      <div className="p-4 bg-surface border border-white/10 rounded-xl text-sm text-slate-400">
        <p className="font-semibold text-white mb-2">Тусламж</p>
        <p>DAN client бүртгэлийг <a href="https://dan.gerege.mn/admin" className="text-primary hover:underline">dan.gerege.mn/admin</a> хуудаснаас хийнэ.
        Техникийн асуудлаар <a href="mailto:dev@gerege.mn" className="text-primary hover:underline">dev@gerege.mn</a> хаягаар холбогдоно.</p>
      </div>
    </main>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <section className="space-y-3">
      <h2 className="text-xl font-bold text-white">{title}</h2>
      <div className="text-sm text-slate-400 leading-relaxed">{children}</div>
    </section>
  );
}

function Code({ children }: { children: string }) {
  return (
    <pre className="bg-bg border border-white/10 rounded-xl p-4 text-xs text-slate-300 overflow-x-auto whitespace-pre-wrap">{children}</pre>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <tr className="border-b border-white/5">
      <td className="py-2 pr-4 text-slate-400 font-mono whitespace-nowrap">{label}</td>
      <td className="py-2 text-slate-300">{value}</td>
    </tr>
  );
}

function Field({ label, desc }: { label: string; desc: string }) {
  return (
    <div className="p-3 bg-bg border border-white/10 rounded-lg">
      <p className="text-sm font-semibold text-white">{label}</p>
      <p className="text-xs text-slate-400 mt-1">{desc}</p>
    </div>
  );
}
