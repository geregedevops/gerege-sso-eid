import { NextRequest, NextResponse } from "next/server";
import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";
import { getSignStatus } from "@/lib/api-client";

export async function GET(_req: NextRequest, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const session = await auth();
  if (!session) return NextResponse.json({ error: "Unauthorized" }, { status: 401 });

  const accessToken = (session.user as any)?.accessToken;
  const sub = (session.user as any)?.sub;
  const user = sub ? await prisma.user.findUnique({ where: { sub } }) : null;
  if (!user) return NextResponse.json({ error: "User not found" }, { status: 404 });

  const signature = await prisma.signature.findUnique({ where: { id } });
  if (!signature || signature.signedById !== user.id) {
    return NextResponse.json({ error: "Not found" }, { status: 404 });
  }

  // If already terminal, return cached status
  if (signature.status === "complete" || signature.status === "failed") {
    return NextResponse.json({
      status: signature.status,
      signerName: signature.signerName,
      certSerial: signature.certSerial,
      signedAt: signature.signedAt,
    });
  }

  // Poll api.gerege.mn
  try {
    const result = await getSignStatus(accessToken, signature.sessionId!);

    if (result.status === "COMPLETE") {
      await prisma.signature.update({
        where: { id },
        data: {
          status: "complete",
          signerName: result.signer_name,
          certSerial: result.cert_serial,
          signedAt: new Date(),
        },
      });
      await prisma.document.update({
        where: { id: signature.documentId },
        data: { status: "signed" },
      });
    } else if (result.status === "ERROR" || result.status === "EXPIRED" || result.status === "CANCELLED") {
      await prisma.signature.update({ where: { id }, data: { status: "failed" } });
      await prisma.document.update({ where: { id: signature.documentId }, data: { status: "failed" } });
    }

    return NextResponse.json({
      status: result.status === "COMPLETE" ? "complete" : result.status === "ERROR" || result.status === "EXPIRED" || result.status === "CANCELLED" ? "failed" : "pending",
      signerName: result.signer_name,
      certSerial: result.cert_serial,
    });
  } catch (err: any) {
    return NextResponse.json({ status: "pending", error: err.message });
  }
}
