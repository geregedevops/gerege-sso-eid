"use client";

import { useState } from "react";

export default function APIExplorerPage() {
  const [method, setMethod] = useState("GET");
  const [path, setPath] = useState("/pos/v1/products");
  const [body, setBody] = useState("");
  const [response, setResponse] = useState<any>(null);
  const [loading, setLoading] = useState(false);

  const send = async () => {
    setLoading(true);
    setResponse(null);
    try {
      const start = Date.now();
      const res = await fetch("/api/sandbox", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ method, path, body: body ? JSON.parse(body) : undefined }),
      });
      const data = await res.json();
      data.duration = Date.now() - start;
      setResponse(data);
    } catch (e: any) {
      setResponse({ error: e.message });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-white">API Explorer</h1>

      <div className="bg-surface border border-white/10 rounded-xl p-5 space-y-4">
        <div className="flex gap-3">
          <select
            value={method}
            onChange={(e) => setMethod(e.target.value)}
            className="px-3 py-2 bg-bg border border-white/10 rounded-lg text-white text-sm"
          >
            <option>GET</option>
            <option>POST</option>
            <option>PUT</option>
            <option>DELETE</option>
          </select>
          <input
            value={path}
            onChange={(e) => setPath(e.target.value)}
            className="flex-1 px-3 py-2 bg-bg border border-white/10 rounded-lg text-white text-sm font-mono focus:ring-1 focus:ring-primary outline-none"
            placeholder="/pos/v1/products"
          />
          <button
            onClick={send}
            disabled={loading}
            className="px-5 py-2 bg-primary text-white font-semibold rounded-lg text-sm hover:bg-primary-light disabled:opacity-40"
          >
            {loading ? "..." : "Илгээх"}
          </button>
        </div>

        <p className="text-xs text-slate-500">
          Base URL: sandbox.gerege.mn &middot; Authorization: Bearer token (автоматаар)
        </p>

        {method !== "GET" && (
          <div>
            <label className="text-xs text-slate-400 mb-1 block">Body (JSON)</label>
            <textarea
              value={body}
              onChange={(e) => setBody(e.target.value)}
              rows={4}
              className="w-full px-3 py-2 bg-bg border border-white/10 rounded-lg text-white text-sm font-mono focus:ring-1 focus:ring-primary outline-none"
              placeholder='{"key": "value"}'
            />
          </div>
        )}
      </div>

      {response && (
        <div className="bg-surface border border-white/10 rounded-xl p-5">
          <div className="flex items-center gap-3 mb-3">
            <span className={`px-2 py-0.5 rounded text-xs font-bold ${
              response.status < 300 ? "bg-green-500/15 text-green-400" :
              response.status < 500 ? "bg-yellow-500/15 text-yellow-400" :
              "bg-red-500/15 text-red-400"
            }`}>
              {response.status || "ERR"}
            </span>
            {response.duration && (
              <span className="text-xs text-slate-500">{response.duration}ms</span>
            )}
          </div>
          <pre className="text-xs text-slate-300 overflow-x-auto whitespace-pre-wrap">
            {JSON.stringify(response.body || response, null, 2)}
          </pre>
        </div>
      )}
    </div>
  );
}
