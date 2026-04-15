import { NextResponse } from "next/server";

const UPSTREAM = process.env.UPSTREAM_API_URL || "http://10.0.0.187:8000";

export async function POST(req: Request) {
  const body = await req.json();
  const regNo = body.reg_no;

  if (!regNo) {
    return NextResponse.json({ error: "reg_no шаардлагатай" }, { status: 400 });
  }

  try {
    const res = await fetch(`${UPSTREAM}/user/validate`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ reg_no: regNo }),
    });

    const data = await res.json();

    if (data.status !== "success" || data.result?.resultCode !== 200) {
      return NextResponse.json({ found: false, citizen: null });
    }

    const r = data.result;
    return NextResponse.json({
      found: true,
      citizen: {
        reg_no: r.regnum || "",
        last_name: r.lastname || "",
        first_name: r.firstname || "",
        surname: r.surname || "",
        gender: r.gender || "",
        birth_date: r.birthDateAsText || "",
        nationality: r.nationality || "",
      },
    });
  } catch (e: any) {
    return NextResponse.json({ error: e.message || "upstream error" }, { status: 502 });
  }
}
