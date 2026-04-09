import { NextRequest, NextResponse } from "next/server";
import { auth } from "@/lib/auth";
import { query, queryOne } from "@/lib/db";
import { getSignStatus } from "@/lib/api-client";

export async function GET(_req: NextRequest, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const session = await auth();
  if (!session) return NextResponse.json({ error: "Unauthorized" }, { status: 401 });

  const accessToken = (session.user as any)?.accessToken;
  const sub = (session.user as any)?.sub;
  const user = await queryOne<{ id: string }>(`SELECT id FROM dbiz_users WHERE sub=$1`, [sub]);
  if (!user) return NextResponse.json({ error: "User not found" }, { status: 404 });

  const sig = await queryOne<{ id: string; signedById: string; sessionId: string; status: string; signerName: string; certSerial: string; signedAt: string; documentId: string }>(
    `SELECT * FROM dbiz_signatures WHERE id=$1`, [id]
  );
  if (!sig || sig.signedById !== user.id) return NextResponse.json({ error: "Not found" }, { status: 404 });

  if (sig.status === "complete" || sig.status === "failed") {
    return NextResponse.json({ status: sig.status, signerName: sig.signerName, certSerial: sig.certSerial, signedAt: sig.signedAt });
  }

  try {
    const result = await getSignStatus(accessToken, sig.sessionId);
    if (result.status === "COMPLETE") {
      await query(`UPDATE dbiz_signatures SET status='complete', "signerName"=$1, "certSerial"=$2, "signedAt"=now() WHERE id=$3`, [result.signer_name, result.cert_serial, id]);
      await query(`UPDATE dbiz_documents SET status='signed' WHERE id=$1`, [sig.documentId]);
    } else if (["ERROR", "EXPIRED", "CANCELLED"].includes(result.status)) {
      await query(`UPDATE dbiz_signatures SET status='failed' WHERE id=$1`, [id]);
      await query(`UPDATE dbiz_documents SET status='failed' WHERE id=$1`, [sig.documentId]);
    }
    return NextResponse.json({
      status: result.status === "COMPLETE" ? "complete" : ["ERROR", "EXPIRED", "CANCELLED"].includes(result.status) ? "failed" : "pending",
      signerName: result.signer_name, certSerial: result.cert_serial,
    });
  } catch (err: any) {
    return NextResponse.json({ status: "pending", error: err.message });
  }
}
