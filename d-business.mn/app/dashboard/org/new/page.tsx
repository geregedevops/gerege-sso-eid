"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

const ORG_TYPES = ["ХХК", "ХК", "ТББ", "ТӨҮГ", "Бусад"];

export default function NewOrgPage() {
  const router = useRouter();
  const [name, setName] = useState("");
  const [regNo, setRegNo] = useState("");
  const [type, setType] = useState("ХХК");
  const [address, setAddress] = useState("");
  const [phone, setPhone] = useState("");
  const [email, setEmail] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  function validate(): string | null {
    if (name.trim().length < 2) return "Байгууллагын нэр хэт богино";
    if (!/^\d{5,10}$/.test(regNo.trim())) return "Регистрийн дугаар 5-10 оронтой тоо байх ёстой";
    if (email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) return "Имэйл хаяг буруу формат";
    if (phone && !/^[\d\-+() ]{4,20}$/.test(phone)) return "Утасны дугаар буруу формат";
    return null;
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    const validationError = validate();
    if (validationError) { setError(validationError); return; }

    setLoading(true);

    try {
      const res = await fetch("/api/org", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name, registrationNumber: regNo, type, address, phone, email }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || "Алдаа гарлаа");
      router.push(`/dashboard/org/${data.id}`);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="max-w-lg mx-auto space-y-6">
      <h1 className="text-2xl font-bold text-white">Байгууллага бүртгүүлэх</h1>

      <form onSubmit={handleSubmit} className="space-y-4">
        <Field label="Байгууллагын нэр *">
          <input value={name} onChange={(e) => setName(e.target.value)} required
            className="w-full px-4 py-3 bg-bg border border-white/10 rounded-xl text-white text-sm outline-none focus:border-primary" placeholder="Жишээ ХХК" />
        </Field>

        <Field label="Регистрийн дугаар *">
          <input value={regNo} onChange={(e) => setRegNo(e.target.value)} required
            className="w-full px-4 py-3 bg-bg border border-white/10 rounded-xl text-white text-sm outline-none focus:border-primary" placeholder="1234567" />
        </Field>

        <Field label="Төрөл *">
          <select value={type} onChange={(e) => setType(e.target.value)}
            className="w-full px-4 py-3 bg-bg border border-white/10 rounded-xl text-white text-sm outline-none focus:border-primary">
            {ORG_TYPES.map((t) => <option key={t} value={t}>{t}</option>)}
          </select>
        </Field>

        <Field label="Хаяг">
          <input value={address} onChange={(e) => setAddress(e.target.value)}
            className="w-full px-4 py-3 bg-bg border border-white/10 rounded-xl text-white text-sm outline-none focus:border-primary" placeholder="Улаанбаатар, Баянгол дүүрэг" />
        </Field>

        <Field label="Утас">
          <input value={phone} onChange={(e) => setPhone(e.target.value)}
            className="w-full px-4 py-3 bg-bg border border-white/10 rounded-xl text-white text-sm outline-none focus:border-primary" placeholder="7700-1234" />
        </Field>

        <Field label="Имэйл">
          <input value={email} onChange={(e) => setEmail(e.target.value)} type="email"
            className="w-full px-4 py-3 bg-bg border border-white/10 rounded-xl text-white text-sm outline-none focus:border-primary" placeholder="info@example.mn" />
        </Field>

        {error && <p className="text-red-400 text-sm">{error}</p>}

        <button type="submit" disabled={loading}
          className="w-full py-3 bg-primary text-white font-semibold rounded-xl hover:bg-primary-light transition-colors disabled:opacity-50">
          {loading ? "Бүртгэж байна..." : "Бүртгүүлэх"}
        </button>
      </form>
    </div>
  );
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div>
      <label className="block text-sm text-slate-400 mb-1">{label}</label>
      {children}
    </div>
  );
}
