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

  const { name, registrationNumber, type, address, phone, email } = await req.json();
  if (!name || !registrationNumber || !type) {
    return NextResponse.json({ error: "name, registrationNumber, type шаардлагатай" }, { status: 400 });
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
