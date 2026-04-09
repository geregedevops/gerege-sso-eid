import { NextRequest, NextResponse } from "next/server";
import { prisma } from "@/lib/db";
import { createHash } from "crypto";

export async function POST(req: NextRequest) {
  try {
    const formData = await req.formData();
    const file = formData.get("file") as File;
    if (!file) return NextResponse.json({ error: "PDF файл шаардлагатай" }, { status: 400 });

    const buffer = Buffer.from(await file.arrayBuffer());
    const fileHash = createHash("sha256").update(buffer).digest("hex");

    // Look up document by hash
    const doc = await prisma.document.findFirst({
      where: { fileHash, status: "signed" },
      include: {
        organization: true,
        signatures: { where: { status: "complete" }, include: { signedBy: true } },
      },
    });

    if (!doc || doc.signatures.length === 0) {
      return NextResponse.json({
        valid: false,
        message: "Энэ файлын гарын үсгийн мэдээлэл олдсонгүй. Файл өөрчлөгдсөн эсвэл бүртгэгдээгүй байж магадгүй.",
      });
    }

    return NextResponse.json({
      valid: true,
      documentName: doc.name,
      signatures: doc.signatures.map((s) => ({
        signerName: s.signerName || s.signedBy.name,
        organizationName: doc.organization.name,
        certSerial: s.certSerial,
        signedAt: s.signedAt?.toISOString().split("T")[0],
      })),
    });
  } catch (err: any) {
    return NextResponse.json({ error: err.message }, { status: 500 });
  }
}
