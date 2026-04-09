import { NextRequest, NextResponse } from "next/server";
import { danImageCache } from "@/lib/dan-image-cache";

// GET /api/dan/photo?key=xxx — serve cached citizen photo
export async function GET(req: NextRequest) {
  const key = req.nextUrl.searchParams.get("key") || "";
  if (!key) return NextResponse.json({ error: "key required" }, { status: 400 });

  const entry = danImageCache.get(key);
  if (!entry || Date.now() > entry.expires) {
    danImageCache.delete(key);
    return NextResponse.json({ error: "not found" }, { status: 404 });
  }

  // One-time use
  danImageCache.delete(key);

  const buffer = Buffer.from(entry.data, "base64");
  return new NextResponse(new Uint8Array(buffer), {
    headers: { "Content-Type": "image/jpeg", "Cache-Control": "no-store" },
  });
}
