export default function SocialGuidePage() {
  return (
    <main className="max-w-3xl mx-auto px-6 py-12 space-y-8">
      <h1 className="text-3xl font-bold text-white">Social Commerce нэгтгэх</h1>
      <p className="text-slate-400">Gerege Social API-г ашиглан лайв худалдаа, product feed нэгтгэх.</p>

      <Code title="Product feed авах">{`GET api.gerege.mn/social/v1/products
Authorization: Bearer {access_token}

// Response:
{
  "products": [
    {
      "id": "prod_123",
      "name": "Кашемир цамц",
      "price": 180000,
      "image_url": "https://...",
      "stock": 25
    }
  ]
}`}</Code>

      <Code title="Лайв худалдаа эхлүүлэх">{`POST api.gerege.mn/social/v1/live
Authorization: Bearer {access_token}
{
  "title": "Хаврын хямдрал",
  "product_ids": ["prod_123", "prod_456"],
  "platform": "facebook"
}`}</Code>

      <Code title="Нийтлэл хуваалцах">{`POST api.gerege.mn/social/v1/posts
Authorization: Bearer {access_token}
{
  "content": "Шинэ бараа ирлээ!",
  "product_ids": ["prod_123"],
  "platforms": ["facebook", "instagram"]
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
