import { NextResponse } from "next/server";
import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";
import { randomBytes } from "crypto";
import bcrypt from "bcryptjs";

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

  // Rate limit: 10 apps per developer
  const appCount = await prisma.app.count({ where: { developerId: developer.id } });
  if (appCount >= 10) {
    return NextResponse.json({ error: "app limit reached (max 10)" }, { status: 429 });
  }

  const body = await req.json();
  const { name, description, redirectUris, scopes, tenantId } = body;

  if (!name || !redirectUris || redirectUris.length === 0) {
    return NextResponse.json({ error: "name and redirect_uris required" }, { status: 400 });
  }

  // No wildcards in redirect URIs
  for (const uri of redirectUris) {
    if (uri.includes("*")) {
      return NextResponse.json({ error: "wildcard redirect_uri not allowed" }, { status: 400 });
    }
  }

  // Generate credentials
  const clientSecret = randomBytes(32).toString("base64url");
  const secretHash = await bcrypt.hash(clientSecret, 12);

  // Create app in dev_apps
  const app = await prisma.app.create({
    data: {
      name,
      description: description || null,
      secretHash,
      redirectUris,
      scopes: scopes || ["openid", "profile"],
      developerId: developer.id,
      tenantId: tenantId || null,
    },
  });

  // Sync to sso_clients (shared DB)
  try {
    await prisma.$executeRawUnsafe(
      `INSERT INTO sso_clients (id, secret_hash, name, redirect_uris, scopes, tenant_id, is_active)
       VALUES ($1, $2, $3, $4::text[], $5::text[], $6, true)
       ON CONFLICT (id) DO UPDATE SET secret_hash=$2, name=$3, redirect_uris=$4::text[], scopes=$5::text[], tenant_id=$6`,
      app.clientId,
      secretHash,
      name,
      redirectUris,
      scopes || ["openid", "profile"],
      tenantId || null
    );
  } catch (e) {
    console.error("Failed to sync sso_clients:", e);
  }

  return NextResponse.json({
    id: app.id,
    clientId: app.clientId,
    clientSecret, // plaintext — shown once only
  }, { status: 201 });
}
