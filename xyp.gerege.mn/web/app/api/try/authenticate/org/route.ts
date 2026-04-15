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
    const input = ceo_reg_no.toLowerCase().trim();

    // Check CEO match
    const actualCeo = (r.generalR?.regnum || "").toLowerCase().trim();
    const ceoMatch = actualCeo === input;

    // Find largest shareholder
    const activeFounders = (r.founder || []).filter((f: any) => f.status === "Тийм");
    let topOwner: any = null;
    let topPct = 0;
    for (const f of activeFounders) {
      const pct = parseFloat(f.sharePercent || "0");
      if (!topOwner || pct > topPct) {
        topOwner = f;
        topPct = pct;
      }
    }

    const ownerMatch = topOwner && (topOwner.stakeHolderRegnum || "").toLowerCase().trim() === input;

    // Either one must match
    if (!ceoMatch && !ownerMatch) {
      return NextResponse.json({
        authenticated: false,
        reason: "ceo_reg_no does not match director or largest shareholder",
      });
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

    const result: any = {
      authenticated: true,
      organization: {
        reg_no: r.changeName?.[0]?.companyRegnum || reg_no,
        name,
        type,
        ceo,
        ceo_reg_no: r.generalR?.regnum || "",
        ceo_position: r.generalR?.positionName || "",
      },
    };

    if (topOwner) {
      result.owner = {
        name: `${topOwner.lastName || ""} ${topOwner.firstName || ""}`.trim(),
        reg_no: topOwner.stakeHolderRegnum || "",
        type: topOwner.stakeHolderTypeName || "",
        share_percent: topOwner.sharePercent || "",
      };
    }

    return NextResponse.json(result);
  } catch (e: any) {
    return NextResponse.json({ error: e.message || "upstream error" }, { status: 502 });
  }
}
