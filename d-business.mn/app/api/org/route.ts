import { NextRequest, NextResponse } from "next/server";
import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";

export async function GET() {
  const session = await auth();
  if (!session) return NextResponse.json({ error: "Unauthorized" }, { status: 401 });

  const sub = (session.user as any)?.sub;
  const user = await prisma.user.findUnique({ where: { sub } });
  if (!user) return NextResponse.json({ error: "User not found" }, { status: 404 });

  const memberships = await prisma.orgMember.findMany({
    where: { userId: user.id },
    include: { organization: true },
  });

  return NextResponse.json(memberships.map((m) => ({ ...m.organization, role: m.role })));
}

export async function POST(req: NextRequest) {
  const session = await auth();
  if (!session) return NextResponse.json({ error: "Unauthorized" }, { status: 401 });

  const sub = (session.user as any)?.sub;
  const user = await prisma.user.findUnique({ where: { sub } });
  if (!user) return NextResponse.json({ error: "User not found" }, { status: 404 });

  const body = await req.json();
  const { name, registrationNumber, type, address, phone, email } = body;

  if (!name || !registrationNumber || !type) {
    return NextResponse.json({ error: "name, registrationNumber, type шаардлагатай" }, { status: 400 });
  }

  const existing = await prisma.organization.findUnique({ where: { registrationNumber } });
  if (existing) {
    return NextResponse.json({ error: "Энэ регистрийн дугаар бүртгэгдсэн байна" }, { status: 409 });
  }

  const org = await prisma.organization.create({
    data: { name, registrationNumber, type, address, phone, email },
  });

  await prisma.orgMember.create({
    data: { organizationId: org.id, userId: user.id, role: "owner" },
  });

  return NextResponse.json(org, { status: 201 });
}
