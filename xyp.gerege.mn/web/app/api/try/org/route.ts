import { NextResponse } from "next/server";

const UPSTREAM = process.env.UPSTREAM_API_URL || "http://10.0.0.187:8000";

export async function POST(req: Request) {
  const body = await req.json();
  const regNo = body.reg_no;

  if (!regNo) {
    return NextResponse.json({ error: "reg_no шаардлагатай" }, { status: 400 });
  }

  try {
    const res = await fetch(`${UPSTREAM}/legalentity/info`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ reg_no: regNo }),
    });

    const data = await res.json();

    if (data.status !== "success" || data.result?.resultCode !== 200) {
      return NextResponse.json({ found: false, organization: null });
    }

    const r = data.result;

    // Company name from changeName
    let name = "";
    let type = "";
    let companyRegNo = regNo;
    if (r.changeName?.length > 0) {
      name = r.changeName[0].requestedName || "";
      type = r.changeName[0].companyType || "";
      companyRegNo = r.changeName[0].companyRegnum || regNo;
    }

    // CEO
    let ceo = "";
    if (r.generalR?.firstName) {
      ceo = `${r.generalR.lastName || ""} ${r.generalR.firstName}`.trim();
    }

    // Phone and address from active address
    let phone = "";
    let address = "";
    if (r.address) {
      const active = r.address.find((a: any) => a.addressStatus === "Тийм");
      if (active) {
        phone = active.phoneNumber || "";
        const parts = [
          active.stateCity?.name,
          active.soumDistrict?.name,
          active.bagKhoroo?.name,
          active.region?.name,
          active.door,
        ].filter(Boolean);
        address = parts.join(", ");
      }
    }

    // Active industries
    const industry = (r.induty || [])
      .filter((i: any) => i.industryStatus === "Тийм")
      .map((i: any) => i.industryName);

    const ceoRegNo = r.generalR?.regnum || "";
    const ceoPosition = r.generalR?.positionName || "";

    // Capital (active)
    const capitalEntry = (r.capital || []).find((c: any) => c.rowStatusName === "Тийм");
    const capital = capitalEntry?.totalAmount || "";

    // Active founders
    const founders = (r.founder || [])
      .filter((f: any) => f.status === "Тийм")
      .map((f: any) => ({
        name: `${f.lastName || ""} ${f.firstName || ""}`.trim(),
        reg_no: f.stakeHolderRegnum || "",
        type: f.stakeHolderTypeName || "",
        share_percent: f.sharePercent || "",
      }));

    // Active board members
    const stake_holders = (r.stakeHolders || [])
      .filter((s: any) => s.status === "Тийм")
      .map((s: any) => ({
        name: `${s.lastname || ""} ${s.firstname || ""}`.trim(),
        reg_no: s.stateRegnum || "",
        position: s.positionName || "",
      }));

    return NextResponse.json({
      found: true,
      organization: {
        reg_no: companyRegNo,
        name,
        type,
        capital,
        ceo,
        ceo_reg_no: ceoRegNo,
        ceo_position: ceoPosition,
        phone,
        address,
        industry,
        founders,
        stake_holders,
      },
    });
  } catch (e: any) {
    return NextResponse.json({ error: e.message || "upstream error" }, { status: 502 });
  }
}
