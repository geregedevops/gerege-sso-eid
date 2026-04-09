import { NextRequest, NextResponse } from "next/server";
import { auth } from "@/lib/auth";
import { query, queryOne, genId } from "@/lib/db";
import { initiateSign } from "@/lib/api-client";
import { canSign } from "@/lib/permissions";
import { createHash } from "crypto";

export async function POST(req: NextRequest) {
  const session = await auth();
  if (!session) return NextResponse.json({ error: "Unauthorized" }, { status: 401 });

  const sub = (session.user as any)?.sub;
  const accessToken = (session.user as any)?.accessToken;
  const certSerial = (session.user as any)?.certSerial;

  if (!accessToken) return NextResponse.json({ error: "Access token байхгүй" }, { status: 401 });

  const user = await queryOne<{ id: string }>(`SELECT id FROM dbiz_users WHERE sub=$1`, [sub]);
  if (!user) return NextResponse.json({ error: "User not found" }, { status: 404 });

  const { organizationId, documentName, document: docBase64 } = await req.json();
  if (!organizationId || !documentName || !docBase64) {
    return NextResponse.json({ error: "organizationId, documentName, document шаардлагатай" }, { status: 400 });
  }

  const membership = await queryOne<{ role: string }>(
    `SELECT role FROM dbiz_org_members WHERE "organizationId"=$1 AND "userId"=$2`,
    [organizationId, user.id]
  );
  if (!membership || !canSign(membership.role)) {
    return NextResponse.json({ error: "Гарын үсэг зурах эрхгүй" }, { status: 403 });
  }

  const docBuffer = Buffer.from(docBase64, "base64");
  const fileHash = createHash("sha256").update(docBuffer).digest("hex");

  const docId = genId();
  await query(
    `INSERT INTO dbiz_documents (id, "organizationId", "uploadedById", name, "fileName", "fileSize", "fileHash", status)
     VALUES ($1,$2,$3,$4,$5,$6,$7,'signing')`,
    [docId, organizationId, user.id, documentName, documentName, docBuffer.length, fileHash]
  );

  try {
    const result = await initiateSign(accessToken, certSerial || sub, documentName, docBase64);

    const sigId = genId();
    await query(
      `INSERT INTO dbiz_signatures (id, "documentId", "organizationId", "signedById", "sessionId", "verificationCode", status)
       VALUES ($1,$2,$3,$4,$5,$6,'pending')`,
      [sigId, docId, organizationId, user.id, result.session_id, result.verification_code]
    );

    return NextResponse.json({ signatureId: sigId, sessionId: result.session_id, verificationCode: result.verification_code });
  } catch (err: any) {
    await query(`UPDATE dbiz_documents SET status='failed' WHERE id=$1`, [docId]);
    return NextResponse.json({ error: err.message || "Signing failed" }, { status: 502 });
  }
}
