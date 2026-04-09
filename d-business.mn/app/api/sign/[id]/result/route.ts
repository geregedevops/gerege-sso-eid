import { NextRequest, NextResponse } from "next/server";
import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";
import { getSignResult } from "@/lib/api-client";

export async function GET(_req: NextRequest, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const session = await auth();
  if (!session) return NextResponse.json({ error: "Unauthorized" }, { status: 401 });

  const accessToken = (session.user as any)?.accessToken;
  const sub = (session.user as any)?.sub;
  const user = sub ? await prisma.user.findUnique({ where: { sub } }) : null;
  if (!user) return NextResponse.json({ error: "User not found" }, { status: 404 });

  const signature = await prisma.signature.findUnique({ where: { id }, include: { document: true } });
  if (!signature || signature.signedById !== user.id || signature.status !== "complete") {
    return NextResponse.json({ error: "Not found or not complete" }, { status: 404 });
  }

  try {
    const { buffer, filename } = await getSignResult(accessToken, signature.sessionId!);
    return new NextResponse(new Uint8Array(buffer), {
      headers: {
        "Content-Type": "application/pdf",
        "Content-Disposition": `attachment; filename="${filename}"`,
      },
    });
  } catch (err: any) {
    return NextResponse.json({ error: err.message }, { status: 502 });
  }
}
