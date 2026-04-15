import { NextResponse } from "next/server";

const UPSTREAM = process.env.UPSTREAM_API_URL || "http://10.0.0.187:8000";

export async function POST(req: Request) {
  const body = await req.json();
  const { reg_no, phone } = body;

  if (!reg_no) return NextResponse.json({ error: "reg_no шаардлагатай" }, { status: 400 });
  if (!phone) return NextResponse.json({ error: "phone шаардлагатай" }, { status: 400 });

  try {
    const res = await fetch(`${UPSTREAM}/user/validate`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ reg_no }),
    });

    const data = await res.json();
    if (data.status !== "success" || data.result?.resultCode !== 200) {
      return NextResponse.json({ authenticated: false, reason: "citizen not found" });
    }

    const r = data.result;
    return NextResponse.json({
      authenticated: true,
      citizen: {
        reg_no: r.regnum || "",
        civil_id: r.civilId || "",
        last_name: r.lastname || "",
        first_name: r.firstname || "",
        gender: r.gender || "",
        birth_date: r.birthDateAsText || "",
        image: r.image || "",
      },
    });
  } catch (e: any) {
    return NextResponse.json({ error: e.message || "upstream error" }, { status: 502 });
  }
}
