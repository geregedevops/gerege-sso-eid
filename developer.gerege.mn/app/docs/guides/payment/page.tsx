export default function PaymentGuidePage() {
  return (
    <main className="max-w-3xl mx-auto px-6 py-12 space-y-8">
      <h1 className="text-3xl font-bold text-white">QPay / SocialPay нэгтгэх</h1>
      <p className="text-slate-400">Gerege Payment API-г ашиглан төлбөр хүлээн авах.</p>

      <Code title="QPay QR код үүсгэх">{`POST api.gerege.mn/payment/v1/qpay
Authorization: Bearer {access_token}
{
  "amount": 50000,
  "description": "Захиалга #123",
  "callback_url": "https://myapp.mn/payment/callback"
}

// Response:
{
  "qr_code": "data:image/png;base64,...",
  "payment_id": "pay_abc123",
  "expires_in": 300
}`}</Code>

      <Code title="SocialPay">{`POST api.gerege.mn/payment/v1/socialpay
Authorization: Bearer {access_token}
{
  "amount": 25000,
  "description": "Бараа худалдаа",
  "callback_url": "https://myapp.mn/payment/callback"
}`}</Code>

      <Code title="eBarimt НӨАТ баримт">{`POST api.gerege.mn/payment/v1/ebarimt
Authorization: Bearer {access_token}
{
  "amount": 50000,
  "vat": 5000,
  "items": [{"name": "Бараа", "qty": 1, "price": 50000}]
}

// Response:
{
  "lottery": "abc123",
  "qr_data": "...",
  "bill_id": "..."
}`}</Code>

      <Code title="Webhook — төлбөр амжилттай">{`POST https://myapp.mn/payment/callback
{
  "event": "payment.completed",
  "payment_id": "pay_abc123",
  "amount": 50000,
  "status": "success"
}`}</Code>
    </main>
  );
}

function Code({ title, children }: { title: string; children: string }) {
  return (
    <div>
      <h3 className="text-sm font-semibold text-white mb-2">{title}</h3>
      <pre className="bg-bg border border-white/10 rounded-xl p-4 text-xs text-slate-300 overflow-x-auto">{children}</pre>
    </div>
  );
}
