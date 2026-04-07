export default function POSPluginGuidePage() {
  return (
    <main className="max-w-3xl mx-auto px-6 py-12 space-y-8">
      <h1 className="text-3xl font-bold text-white">Gerege POS Plugin хийх</h1>
      <p className="text-slate-400">POS системд plugin нэгтгэх гайд.</p>

      <Step n={1} title="App бүртгүүл">
        <p className="text-slate-400 text-sm">
          developer.gerege.mn дээр app үүсгэж, <code className="text-primary">pos</code> scope сонгоно.
        </p>
      </Step>

      <Step n={2} title="sso.gerege.mn-р нэвтрэх">
        <Code>{`// scope: openid profile pos
authorization: {
  params: { scope: "openid profile pos" }
}`}</Code>
      </Step>

      <Step n={3} title="tenant_id claim авах">
        <Code>{`// ID Token-д tenant_id, tenant_role, plan орно
{
  "sub": "...",
  "tenant_id": "restaurant-govi",
  "tenant_role": "owner",
  "plan": "pro"
}`}</Code>
      </Step>

      <Step n={4} title="POS API дуудах">
        <Code>{`// Бараа жагсаалт
GET api.gerege.mn/pos/v1/products
Authorization: Bearer {access_token}

// Захиалга үүсгэх
POST api.gerege.mn/pos/v1/orders
Authorization: Bearer {access_token}
{
  "items": [{"product_id": "...", "qty": 2}],
  "payment_method": "qpay"
}

// Гүйлгээ
GET api.gerege.mn/pos/v1/transactions

// Тайлан
GET api.gerege.mn/pos/v1/reports/daily`}</Code>
      </Step>
    </main>
  );
}

function Step({ n, title, children }: { n: number; title: string; children: React.ReactNode }) {
  return (
    <div>
      <h2 className="text-lg font-semibold text-white flex items-center gap-3 mb-3">
        <span className="w-7 h-7 bg-primary/20 text-primary text-sm font-bold rounded-full flex items-center justify-center">{n}</span>
        {title}
      </h2>
      {children}
    </div>
  );
}

function Code({ children }: { children: string }) {
  return (
    <pre className="bg-bg border border-white/10 rounded-xl p-4 text-xs text-slate-300 overflow-x-auto">
      {children}
    </pre>
  );
}
