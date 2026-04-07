"use client";

import { useState } from "react";

const SOCIAL_SCENARIOS = [
  { name: "Product feed авах", method: "GET", path: "/social/v1/products" },
  { name: "Лайв худалдаа эхлүүлэх", method: "POST", path: "/social/v1/live", body: { title: "Хаврын хямдрал", product_ids: ["prod_001"], platform: "facebook" } },
  { name: "Нийтлэл хуваалцах", method: "POST", path: "/social/v1/posts", body: { content: "Шинэ бараа ирлээ!", product_ids: ["prod_001"], platforms: ["facebook", "instagram"] } },
];

export default function SocialSandboxPage() {
  const [results, setResults] = useState<Record<string, any>>({});
  const [loading, setLoading] = useState<Record<string, boolean>>({});

  const run = async (scenario: typeof SOCIAL_SCENARIOS[0]) => {
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
      <h1 className="text-2xl font-bold text-white">Social Commerce Sandbox</h1>
      <p className="text-slate-400 text-sm">Social API test scenarios.</p>

      <div className="space-y-4">
        {SOCIAL_SCENARIOS.map((s) => (
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
