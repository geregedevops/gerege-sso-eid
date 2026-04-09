import { NextRequest, NextResponse } from "next/server";
import { auth } from "@/lib/auth";
import { prisma } from "@/lib/db";
import { initiateSign } from "@/lib/api-client";
import { canSign, getOrgMembership } from "@/lib/permissions";
import { createHash } from "crypto";

export async function POST(req: NextRequest) {
  const session = await auth();
  if (!session) return NextResponse.json({ error: "Unauthorized" }, { status: 401 });

  const sub = (session.user as any)?.sub;
  const accessToken = (session.user as any)?.accessToken;
  const certSerial = (session.user as any)?.certSerial;

  if (!accessToken) {
    return NextResponse.json({ error: "Access token байхгүй. Дахин нэвтэрнэ үү." }, { status: 401 });
  }

  const user = await prisma.user.findUnique({ where: { sub } });
  if (!user) return NextResponse.json({ error: "User not found" }, { status: 404 });

  const body = await req.json();
  const { organizationId, documentName, document: docBase64 } = body;

  if (!organizationId || !documentName || !docBase64) {
    return NextResponse.json({ error: "organizationId, documentName, document шаардлагатай" }, { status: 400 });
  }

  // Check membership & signing permission
  const membership = await getOrgMembership(user.id, organizationId);
  if (!membership || !canSign(membership.role)) {
    return NextResponse.json({ error: "Гарын үсэг зурах эрхгүй" }, { status: 403 });
  }

  // Decode and hash document
  const docBuffer = Buffer.from(docBase64, "base64");
  const fileHash = createHash("sha256").update(docBuffer).digest("hex");

  // Create document record
  const doc = await prisma.document.create({
    data: {
      organizationId,
      uploadedById: user.id,
      name: documentName,
      fileName: documentName,
      fileSize: docBuffer.length,
      fileHash,
      status: "signing",
    },
  });

  try {
    // Call api.gerege.mn to initiate signing
    const result = await initiateSign(accessToken, certSerial || sub, documentName, docBase64);

    // Create signature record
    const signature = await prisma.signature.create({
      data: {
        documentId: doc.id,
        organizationId,
        signedById: user.id,
        sessionId: result.session_id,
        verificationCode: result.verification_code,
        status: "pending",
      },
    });

    return NextResponse.json({
      signatureId: signature.id,
      sessionId: result.session_id,
      verificationCode: result.verification_code,
    });
  } catch (err: any) {
    await prisma.document.update({ where: { id: doc.id }, data: { status: "failed" } });
    return NextResponse.json({ error: err.message || "Signing request failed" }, { status: 502 });
  }
}
