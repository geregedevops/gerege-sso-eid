import { NextRequest, NextResponse } from "next/server";
import { danImageCache } from "@/lib/dan-image-cache";

const APP_URL = process.env.NEXT_PUBLIC_APP_URL || "https://test.gerege.mn";

// In-memory store for full citizen data received via POST
const citizenDataCache = new Map<string, { data: Record<string, string>; expires: number }>();

// Cleanup
if (typeof setInterval !== "undefined") {
  setInterval(() => {
    const now = Date.now();
    citizenDataCache.forEach((val, key) => {
      if (now > val.expires) citizenDataCache.delete(key);
    });
  }, 60_000);
}

// POST — dan.gerege.mn sends full citizen data (including base64 image) as JSON
export async function POST(req: NextRequest) {
  try {
    const data = await req.json() as Record<string, string>;
    const regNo = data.reg_no || "";
    if (!regNo) {
      return NextResponse.json({ error: "reg_no required" }, { status: 400 });
    }

    // Store image separately
    if (data.image) {
      const imgKey = regNo.replace(/[^a-zA-Z0-9]/g, "") + "_" + Date.now();
      danImageCache.set(imgKey, { data: data.image, expires: Date.now() + 5 * 60 * 1000 });
      data._img_key = imgKey;
    }

    // Store full citizen data (without image to save memory)
    const storeData = { ...data };
    delete storeData.image;
    const cacheKey = regNo.replace(/[^a-zA-Z0-9]/g, "") + "_data";
    citizenDataCache.set(cacheKey, { data: storeData, expires: Date.now() + 5 * 60 * 1000 });

    return NextResponse.json({ status: "ok" });
  } catch {
    return NextResponse.json({ error: "parse error" }, { status: 400 });
  }
}

// GET — browser redirect from dan.gerege.mn with ?status=ok&reg_no=...
export async function GET(req: NextRequest) {
  const params = req.nextUrl.searchParams;
  const regNo = params.get("reg_no") || "";
  const status = params.get("status") || "";

  if (!regNo || status !== "ok") {
    return NextResponse.redirect(new URL("/auth/login?error=dan_failed", APP_URL));
  }

  // Look up full data from POST cache
  const cacheKey = regNo.replace(/[^a-zA-Z0-9]/g, "") + "_data";
  const cached = citizenDataCache.get(cacheKey);

  const resultURL = new URL("/auth/dan-result", APP_URL);

  if (cached && Date.now() <= cached.expires) {
    // Use full data from POST
    citizenDataCache.delete(cacheKey);
    for (const [key, value] of Object.entries(cached.data)) {
      if (key !== "_img_key" || key === "_img_key") {
        resultURL.searchParams.set(key === "_img_key" ? "img_key" : key, value);
      }
    }
  } else {
    // Fallback: use query params from GET redirect
    params.forEach((value, key) => {
      if (key !== "status") {
        resultURL.searchParams.set(key, value);
      }
    });
  }

  return NextResponse.redirect(resultURL);
}
