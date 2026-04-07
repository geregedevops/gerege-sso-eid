"use client";

import { useState, useEffect } from "react";

interface WebhookEvent {
  id: string;
  time: string;
  type: string;
  payload: any;
}

export default function WebhookInspectorPage() {
  const [webhookUrl, setWebhookUrl] = useState("");
  const [events, setEvents] = useState<WebhookEvent[]>([]);

  useEffect(() => {
    const uuid = crypto.randomUUID();
    const url = `${window.location.origin}/api/webhook/${uuid}`;
    setWebhookUrl(url);

    // Simulate some test events
    const testEvents: WebhookEvent[] = [
      { id: "1", time: new Date().toLocaleTimeString(), type: "payment.completed", payload: { amount: 50000, payment_id: "pay_sandbox_001" } },
      { id: "2", time: new Date(Date.now() - 3000).toLocaleTimeString(), type: "order.created", payload: { id: "ord_sandbox_123", items: 2 } },
      { id: "3", time: new Date(Date.now() - 60000).toLocaleTimeString(), type: "pos.sale", payload: { total: 25000, items: 1 } },
    ];
    setEvents(testEvents);
  }, []);

  const [copied, setCopied] = useState(false);
  const copyUrl = () => {
    navigator.clipboard.writeText(webhookUrl);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-white">Webhook Inspector</h1>

      <div className="bg-surface border border-white/10 rounded-xl p-5">
        <label className="text-xs text-slate-400 mb-2 block">Таны webhook endpoint:</label>
        <div className="flex items-center gap-2">
          <code className="flex-1 text-sm font-mono text-white bg-bg px-3 py-2 rounded-lg border border-white/10 break-all">
            {webhookUrl || "Loading..."}
          </code>
          <button onClick={copyUrl} className="px-3 py-2 text-xs text-primary hover:text-primary-light border border-white/10 rounded-lg">
            {copied ? "Copied" : "Хуулах"}
          </button>
        </div>
        <p className="text-xs text-slate-500 mt-2">Энэ URL-г app-ийн webhook endpoint-д тохируулна.</p>
      </div>

      <div className="bg-surface border border-white/10 rounded-xl p-5">
        <h2 className="text-sm font-semibold text-slate-400 uppercase tracking-wider mb-4">Incoming Events</h2>
        {events.length === 0 ? (
          <p className="text-slate-500 text-sm text-center py-8">Webhook event хүлээж байна...</p>
        ) : (
          <div className="space-y-2">
            {events.map((ev) => (
              <div key={ev.id} className="bg-bg border border-white/10 rounded-lg p-3 flex items-start gap-3">
                <span className="text-xs text-slate-500 font-mono w-20 flex-shrink-0">{ev.time}</span>
                <span className="px-2 py-0.5 bg-primary/10 text-primary text-xs font-medium rounded">{ev.type}</span>
                <pre className="text-xs text-slate-300 flex-1 overflow-x-auto">{JSON.stringify(ev.payload)}</pre>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
