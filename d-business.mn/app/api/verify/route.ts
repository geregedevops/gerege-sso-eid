import { NextRequest, NextResponse } from "next/server";
import { query } from "@/lib/db";
import { createHash } from "crypto";

export async function POST(req: NextRequest) {
  try {
    const formData = await req.formData();
    const file = formData.get("file") as File;
    if (!file) return NextResponse.json({ error: "PDF файл шаардлагатай" }, { status: 400 });

    const buffer = Buffer.from(await file.arrayBuffer());
    const fileHash = createHash("sha256").update(buffer).digest("hex");

    const docs = await query(
      `SELECT d.name, o.name as org_name, s."signerName", s."certSerial", s."signedAt"
       FROM dbiz_documents d
       JOIN dbiz_organizations o ON o.id = d."organizationId"
       JOIN dbiz_signatures s ON s."documentId" = d.id AND s.status = 'complete'
       WHERE d."fileHash" = $1 AND d.status = 'signed'`,
      [fileHash]
    );

    if (docs.length === 0) {
      return NextResponse.json({ valid: false, message: "Файлын гарын үсгийн мэдээлэл олдсонгүй." });
    }

    return NextResponse.json({
      valid: true,
      documentName: docs[0].name,
      signatures: docs.map((d: any) => ({
        signerName: d.signerName,
        organizationName: d.org_name,
        certSerial: d.certSerial,
        signedAt: d.signedAt?.toISOString?.()?.split("T")[0] || d.signedAt,
      })),
    });
  } catch (err: any) {
    return NextResponse.json({ error: err.message }, { status: 500 });
  }
}
