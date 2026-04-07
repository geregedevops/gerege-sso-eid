"use client";

import { useState } from "react";

const POS_SCENARIOS = [
  { name: "Бараа жагсаалт авах", method: "GET", path: "/pos/v1/products" },
  { name: "Захиалга үүсгэх", method: "POST", path: "/pos/v1/orders", body: { items: [{ product_id: "prod_001", qty: 2 }], payment_method: "qpay" } },
  { name: "Гүйлгээ харах", method: "GET", path: "/pos/v1/transactions" },
  { name: "Өдрийн тайлан", method: "GET", path: "/pos/v1/reports/daily" },
];

export default function POSSandboxPage() {
  const [results, setResults] = useState<Record<string, any>>({});
  const [loading, setLoading] = useState<Record<string, boolean>>({});

  const run = async (scenario: typeof POS_SCENARIOS[0]) => {
    setLoading((l) => ({ ...l, [scenario.name]: true }));
    try {
      const res = await fetch("/api/sandbox", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ method: scenario.method, path: scenario.path, body: scenario.body }),
      });
      const data = await res.json();
      setResults((r) => ({ ...r, [scenario.name]: data }));
    } catch (e: any) {
      setResults((r) => ({ ...r, [scenario.name]: { error: e.message } }));
    } finally {
      setLoading((l) => ({ ...l, [scenario.name]: false }));
    }
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-white">POS API Sandbox</h1>
      <p className="text-slate-400 text-sm">Pre-built test scenarios — бодит transaction хийгдэхгүй.</p>

      <div className="space-y-4">
        {POS_SCENARIOS.map((s) => (
          <div key={s.name} className="bg-surface border border-white/10 rounded-xl p-5">
            <div className="flex items-center justify-between mb-3">
              <div>
                <span className="font-medium text-white">{s.name}</span>
                <p className="text-xs text-slate-500 font-mono mt-1">{s.method} {s.path}</p>
              </div>
              <button
                onClick={() => run(s)}
                disabled={loading[s.name]}
                className="px-4 py-2 bg-primary text-white text-sm font-semibold rounded-lg hover:bg-primary-light disabled:opacity-40"
              >
                {loading[s.name] ? "..." : "Ажиллуулах"}
              </button>
            </div>
            {results[s.name] && (
              <pre className="bg-bg border border-white/10 rounded-lg p-3 text-xs text-slate-300 overflow-x-auto mt-3">
                {JSON.stringify(results[s.name], null, 2)}
              </pre>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
