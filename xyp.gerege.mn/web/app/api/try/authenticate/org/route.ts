import { NextResponse } from "next/server";

const UPSTREAM = process.env.UPSTREAM_API_URL || "http://10.0.0.187:8000";

export async function POST(req: Request) {
  const body = await req.json();
  const { reg_no, ceo_reg_no } = body;

  if (!reg_no) return NextResponse.json({ error: "reg_no шаардлагатай" }, { status: 400 });
  if (!ceo_reg_no) return NextResponse.json({ error: "ceo_reg_no шаардлагатай" }, { status: 400 });

  try {
    const res = await fetch(`${UPSTREAM}/legalentity/info`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ reg_no }),
    });

    const data = await res.json();
    if (data.status !== "success" || data.result?.resultCode !== 200) {
      return NextResponse.json({ authenticated: false, reason: "organization not found" });
    }

    const r = data.result;
    const actualCeoRegNo = (r.generalR?.regnum || "").toLowerCase().trim();
    const inputCeoRegNo = ceo_reg_no.toLowerCase().trim();

    if (actualCeoRegNo !== inputCeoRegNo) {
      return NextResponse.json({ authenticated: false, reason: "ceo_reg_no does not match" });
    }

    let name = "";
    let type = "";
    if (r.changeName?.length > 0) {
      name = r.changeName[0].requestedName || "";
      type = r.changeName[0].companyType || "";
    }

    let ceo = "";
    if (r.generalR?.firstName) {
      ceo = `${r.generalR.lastName || ""} ${r.generalR.firstName}`.trim();
    }

    return NextResponse.json({
      authenticated: true,
      organization: {
        reg_no: r.changeName?.[0]?.companyRegnum || reg_no,
        name,
        type,
        ceo,
        ceo_reg_no: r.generalR?.regnum || "",
        ceo_position: r.generalR?.positionName || "",
      },
    });
  } catch (e: any) {
    return NextResponse.json({ error: e.message || "upstream error" }, { status: 502 });
  }
}
