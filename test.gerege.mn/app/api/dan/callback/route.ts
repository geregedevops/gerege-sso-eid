import { NextRequest, NextResponse } from "next/server";
import { danImageCache } from "@/lib/dan-image-cache";

const APP_URL = process.env.NEXT_PUBLIC_APP_URL || "https://test.gerege.mn";

// POST — server-to-server from dan.gerege.mn with full citizen data including image (base64)
export async function POST(req: NextRequest) {
  try {
    const data = await req.json();
    if (data.image && data.reg_no) {
      const imgKey = (data.reg_no as string).replace(/[^a-zA-Z0-9]/g, "") + "_" + Date.now();
      danImageCache.set(imgKey, { data: data.image, expires: Date.now() + 5 * 60 * 1000 });
    }
    return NextResponse.json({ status: "ok" });
  } catch {
    return NextResponse.json({ error: "parse error" }, { status: 400 });
  }
}

// GET — browser redirect from dan.gerege.mn with citizen data as query params
export async function GET(req: NextRequest) {
  const params = req.nextUrl.searchParams;
  const regNo = params.get("reg_no") || "";

  if (!regNo) {
    return NextResponse.redirect(new URL("/auth/login?error=dan_failed", APP_URL));
  }

  // Check if we have a cached image from the prior POST
  let imgKey = "";
  const prefix = regNo.replace(/[^a-zA-Z0-9]/g, "") + "_";
  danImageCache.forEach((_val, key) => {
    if (key.startsWith(prefix)) imgKey = key;
  });

  const resultURL = new URL("/auth/dan-result", APP_URL);
  params.forEach((value, key) => {
    resultURL.searchParams.set(key, value);
  });
  if (imgKey) {
    resultURL.searchParams.set("img_key", imgKey);
  }

  return NextResponse.redirect(resultURL);
}
