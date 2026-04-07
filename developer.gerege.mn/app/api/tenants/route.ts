import { NextResponse } from "next/server";
import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";

export async function POST(req: Request) {
  const session = await auth();
  const sub = (session?.user as any)?.sub;
  if (!sub) {
    return NextResponse.json({ error: "unauthorized" }, { status: 401 });
  }

  const developer = await prisma.developer.findUnique({ where: { sub } });
  if (!developer) {
    return NextResponse.json({ error: "developer not found" }, { status: 404 });
  }

  const body = await req.json();
  const { name, slug, plan } = body;

  if (!name || !slug) {
    return NextResponse.json({ error: "name and slug required" }, { status: 400 });
  }

  // Check slug uniqueness
  const existing = await prisma.tenant.findUnique({ where: { slug } });
  if (existing) {
    return NextResponse.json({ error: "slug already taken" }, { status: 409 });
  }

  // Create tenant + owner membership
  const tenant = await prisma.tenant.create({
    data: {
      name,
      slug,
      plan: plan || "starter",
      members: {
        create: {
          developerId: developer.id,
          role: "owner",
        },
      },
    },
  });

  // Sync to gerege_tenants (shared DB with SSO)
  try {
    await prisma.$executeRawUnsafe(
      `INSERT INTO gerege_tenants (id, name, plan, is_active)
       VALUES ($1, $2, $3, true)
       ON CONFLICT (id) DO UPDATE SET name=$2, plan=$3`,
      tenant.id,
      name,
      plan || "starter"
    );
    await prisma.$executeRawUnsafe(
      `INSERT INTO tenant_members (tenant_id, sub, role)
       VALUES ($1, $2, 'owner')
       ON CONFLICT (tenant_id, sub) DO NOTHING`,
      tenant.id,
      sub
    );
  } catch (e) {
    console.error("Failed to sync gerege_tenants:", e);
  }

  return NextResponse.json({ id: tenant.id, slug: tenant.slug }, { status: 201 });
}
