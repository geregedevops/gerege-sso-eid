"use client";

import { useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

const fields = [
  { key: "reg_no", label: "Регистрийн дугаар" },
  { key: "family_name", label: "Овог" },
  { key: "given_name", label: "Нэр" },
  { key: "surname", label: "Ургийн овог" },
  { key: "civil_id", label: "Иргэний ID" },
  { key: "gender", label: "Хүйс" },
  { key: "birth_date", label: "Төрсөн огноо" },
  { key: "birth_place", label: "Төрсөн газар" },
  { key: "nationality", label: "Үндэс" },
  { key: "aimag_name", label: "Аймаг/Хот" },
  { key: "sum_name", label: "Сум/Дүүрэг" },
  { key: "bag_name", label: "Баг/Хороо" },
  { key: "address_detail", label: "Хаяг" },
  { key: "passport_address", label: "Паспортын хаяг" },
  { key: "apartment_name", label: "Байр" },
  { key: "street_name", label: "Гудамж" },
  { key: "passport_issue_date", label: "Паспорт олгосон" },
  { key: "passport_expire_date", label: "Паспорт дуусах" },
];

export default function DANResultFullPage() {
  const searchParams = useSearchParams();
  const [photo, setPhoto] = useState<string | null>(null);

  const regNo = searchParams.get("reg_no") || "";
  const givenName = searchParams.get("given_name") || "";
  const familyName = searchParams.get("family_name") || "";

  useEffect(() => {
    // Read photo from cookie
    const match = document.cookie.match(/(?:^|; )dan_photo=([^;]*)/);
    if (match) {
      setPhoto(decodeURIComponent(match[1]));
      // Clear cookie after reading
      document.cookie = "dan_photo=; path=/auth/dan-result-full; max-age=0";
    }
  }, []);

  if (!regNo) {
    return (
      <main className="min-h-[calc(100vh-56px)] flex items-center justify-center p-6">
        <div className="max-w-sm w-full bg-surface border border-white/10 rounded-2xl p-8 text-center space-y-4">
          <h1 className="text-xl font-bold text-red-400">DAN Verify алдаа</h1>
          <p className="text-sm text-slate-400">Иргэний мэдээлэл ирсэнгүй.</p>
          <a href="/auth/login" className="block py-3 bg-primary text-white font-semibold rounded-xl">Буцах</a>
        </div>
      </main>
    );
  }

  return (
    <main className="min-h-[calc(100vh-56px)] flex items-center justify-center p-6">
      <div className="max-w-lg w-full bg-surface border border-white/10 rounded-2xl p-8 space-y-6">
        <div className="text-center space-y-2">
          <div className="w-14 h-14 bg-green-500/10 rounded-xl flex items-center justify-center mx-auto">
            <span className="text-green-400 text-2xl">&#10003;</span>
          </div>
          <h1 className="text-xl font-bold text-white">DAN Verified (Full)</h1>
          <p className="text-sm text-slate-400">{givenName} {familyName} ({regNo})</p>
          <span className="inline-block px-3 py-1 bg-blue-500/10 text-blue-400 text-xs font-semibold rounded-full">
            Зураг бүхий
          </span>
        </div>

        {photo && (
          <div className="flex justify-center">
            <img
              src={`data:image/jpeg;base64,${photo}`}
              alt="Иргэний зураг"
              className="w-40 h-52 object-cover rounded-xl border-2 border-white/10"
            />
          </div>
        )}

        <div className="bg-white/5 rounded-xl overflow-hidden">
          <table className="w-full text-sm">
            <tbody>
              {fields.map(({ key, label }) => {
                const value = searchParams.get(key);
                if (!value) return null;
                return (
                  <tr key={key} className="border-b border-white/5">
                    <td className="px-4 py-3 text-slate-400 font-medium whitespace-nowrap">{label}</td>
                    <td className="px-4 py-3 text-white">{value}</td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>

        <div className="flex gap-3">
          <a href="/auth/login" className="flex-1 block py-3 text-center bg-white/10 text-white font-semibold rounded-xl hover:bg-white/20 transition-colors">Буцах</a>
          <a href="/sandbox" className="flex-1 block py-3 text-center bg-primary text-white font-semibold rounded-xl hover:bg-primary-light transition-colors">Sandbox</a>
        </div>
      </div>
    </main>
  );
}
