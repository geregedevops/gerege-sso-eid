"use client";

import { useState } from "react";

export default function PaymentSimulatorPage() {
  const [amount, setAmount] = useState("50000");
  const [paymentResult, setPaymentResult] = useState<any>(null);
  const [loading, setLoading] = useState(false);

  const simulateQPay = async (success: boolean) => {
    setLoading(true);
    try {
      const res = await fetch("/api/sandbox", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          method: "POST",
          path: "/payment/v1/qpay",
          body: { amount: parseInt(amount), simulate_result: success ? "success" : "failed", description: "Sandbox test" },
        }),
      });
      const data = await res.json();
      setPaymentResult(data);
    } catch (e: any) {
      setPaymentResult({ error: e.message });
    } finally {
      setLoading(false);
    }
  };

  const simulateEbarimt = async () => {
    setLoading(true);
    try {
      const res = await fetch("/api/sandbox", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          method: "POST",
          path: "/payment/v1/ebarimt",
          body: { amount: parseInt(amount), vat: Math.floor(parseInt(amount) * 0.1), items: [{ name: "Sandbox бараа", qty: 1, price: parseInt(amount) }] },
        }),
      });
      const data = await res.json();
      setPaymentResult(data);
    } catch (e: any) {
      setPaymentResult({ error: e.message });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-white">Payment Simulator</h1>
      <p className="text-slate-400 text-sm">QPay, SocialPay, eBarimt simulate — бодит төлбөр хийгдэхгүй.</p>

      <div className="bg-surface border border-white/10 rounded-xl p-5 space-y-4">
        <div>
          <label className="text-xs text-slate-400 mb-1 block">Дүн (MNT)</label>
          <input
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            className="w-full px-3 py-2 bg-bg border border-white/10 rounded-lg text-white text-sm font-mono focus:ring-1 focus:ring-primary outline-none"
          />
        </div>

        <div className="grid sm:grid-cols-3 gap-3">
          <button
            onClick={() => simulateQPay(true)}
            disabled={loading}
            className="py-3 bg-green-600 text-white font-semibold rounded-xl text-sm hover:bg-green-500 disabled:opacity-40"
          >
            QPay Амжилттай
          </button>
          <button
            onClick={() => simulateQPay(false)}
            disabled={loading}
            className="py-3 bg-red-600 text-white font-semibold rounded-xl text-sm hover:bg-red-500 disabled:opacity-40"
          >
            QPay Амжилтгүй
          </button>
          <button
            onClick={simulateEbarimt}
            disabled={loading}
            className="py-3 bg-blue-600 text-white font-semibold rounded-xl text-sm hover:bg-blue-500 disabled:opacity-40"
          >
            eBarimt Preview
          </button>
        </div>
      </div>

      {paymentResult && (
        <div className="bg-surface border border-white/10 rounded-xl p-5">
          <h2 className="text-sm font-semibold text-slate-400 uppercase mb-3">Response</h2>
          <pre className="bg-bg border border-white/10 rounded-lg p-3 text-xs text-slate-300 overflow-x-auto">
            {JSON.stringify(paymentResult, null, 2)}
          </pre>
        </div>
      )}
    </div>
  );
}
