import { NextRequest, NextResponse } from "next/server";
import { auth } from "@/lib/auth";
import { query, queryOne, genId } from "@/lib/db";

export async function GET() {
  const session = await auth();
  if (!session) return NextResponse.json({ error: "Unauthorized" }, { status: 401 });
  const sub = (session.user as any)?.sub;

  const user = await queryOne(`SELECT id FROM dbiz_users WHERE sub=$1`, [sub]);
  if (!user) return NextResponse.json({ error: "User not found" }, { status: 404 });

  const orgs = await query(
    `SELECT o.*, m.role FROM dbiz_organizations o
     JOIN dbiz_org_members m ON m."organizationId" = o.id
     WHERE m."userId" = $1 ORDER BY o."createdAt" DESC`,
    [user.id]
  );

  return NextResponse.json(orgs);
}

export async function POST(req: NextRequest) {
  const session = await auth();
  if (!session) return NextResponse.json({ error: "Unauthorized" }, { status: 401 });
  const sub = (session.user as any)?.sub;

  const user = await queryOne<{ id: string }>(`SELECT id FROM dbiz_users WHERE sub=$1`, [sub]);
  if (!user) return NextResponse.json({ error: "User not found" }, { status: 404 });

  const body = await req.json();
  const name = (body.name || "").trim();
  const registrationNumber = (body.registrationNumber || "").trim();
  const type = body.type;
  const address = body.address;
  const phone = body.phone;
  const email = body.email;

  if (!name || name.length < 2) {
    return NextResponse.json({ error: "Байгууллагын нэр хэт богино" }, { status: 400 });
  }
  if (!/^\d{5,10}$/.test(registrationNumber)) {
    return NextResponse.json({ error: "Регистрийн дугаар буруу формат (5-10 оронтой тоо)" }, { status: 400 });
  }
  const VALID_ORG_TYPES = ["ХХК", "ХК", "ТББ", "ТӨҮГ", "Бусад"];
  if (!type || !VALID_ORG_TYPES.includes(type)) {
    return NextResponse.json({ error: "Байгууллагын төрөл буруу" }, { status: 400 });
  }
  if (email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
    return NextResponse.json({ error: "Имэйл хаяг буруу формат" }, { status: 400 });
  }
  if (phone && !/^[\d\-+() ]{4,20}$/.test(phone)) {
    return NextResponse.json({ error: "Утасны дугаар буруу формат" }, { status: 400 });
  }

  // Check if user already owns too many organizations
  const ownedCount = await queryOne<{ count: string }>(
    `SELECT COUNT(*) as count FROM dbiz_org_members WHERE "userId"=$1 AND role='owner'`,
    [user.id]
  );
  if (ownedCount && parseInt(ownedCount.count) >= 10) {
    return NextResponse.json({ error: "Хэрэглэгч хамгийн ихдээ 10 байгууллагын эзэмшигч байж болно" }, { status: 400 });
  }

  const existing = await queryOne(`SELECT id FROM dbiz_organizations WHERE "registrationNumber"=$1`, [registrationNumber]);
  if (existing) return NextResponse.json({ error: "Энэ регистрийн дугаар бүртгэгдсэн" }, { status: 409 });

  const orgId = genId();
  await query(
    `INSERT INTO dbiz_organizations (id, name, "registrationNumber", type, address, phone, email) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
    [orgId, name, registrationNumber, type, address || null, phone || null, email || null]
  );
  await query(
    `INSERT INTO dbiz_org_members ("organizationId", "userId", role) VALUES ($1,$2,'owner')`,
    [orgId, user.id]
  );

  return NextResponse.json({ id: orgId, name, registrationNumber, type }, { status: 201 });
}
