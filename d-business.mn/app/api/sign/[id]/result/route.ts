import { NextRequest, NextResponse } from "next/server";
import { auth } from "@/lib/auth";
import { queryOne } from "@/lib/db";
import { getSignResult } from "@/lib/api-client";

export async function GET(_req: NextRequest, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const session = await auth();
  if (!session) return NextResponse.json({ error: "Unauthorized" }, { status: 401 });

  const accessToken = (session.user as any)?.accessToken;
  const sub = (session.user as any)?.sub;
  const user = await queryOne<{ id: string }>(`SELECT id FROM dbiz_users WHERE sub=$1`, [sub]);
  if (!user) return NextResponse.json({ error: "User not found" }, { status: 404 });

  const sig = await queryOne<{ signedById: string; sessionId: string; status: string }>(
    `SELECT "signedById", "sessionId", status FROM dbiz_signatures WHERE id=$1`, [id]
  );
  if (!sig || sig.signedById !== user.id || sig.status !== "complete") {
    return NextResponse.json({ error: "Not found" }, { status: 404 });
  }

  try {
    const { buffer, filename } = await getSignResult(accessToken, sig.sessionId);
    return new NextResponse(new Uint8Array(buffer), {
      headers: { "Content-Type": "application/pdf", "Content-Disposition": `attachment; filename="${filename}"` },
    });
  } catch (err: any) {
    return NextResponse.json({ error: err.message }, { status: 502 });
  }
}
