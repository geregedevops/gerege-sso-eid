import { NextRequest, NextResponse } from "next/server";
import { imageStore } from "../callback-full/route";

// Serve citizen photo from temporary in-memory store
// GET /api/dan/photo?key=xxx
export async function GET(req: NextRequest) {
  const key = req.nextUrl.searchParams.get("key") || "";

  if (!key) {
    return NextResponse.json({ error: "key required" }, { status: 400 });
  }

  const entry = imageStore.get(key);
  if (!entry || Date.now() > entry.expires) {
    imageStore.delete(key);
    return NextResponse.json({ error: "expired or not found" }, { status: 404 });
  }

  // One-time use: delete after serving
  imageStore.delete(key);

  // Return as JPEG image
  const buffer = Buffer.from(entry.data, "base64");
  return new NextResponse(buffer, {
    headers: {
      "Content-Type": "image/jpeg",
      "Cache-Control": "no-store",
    },
  });
}
