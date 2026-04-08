import { NextRequest, NextResponse } from "next/server";
import { imageStore } from "@/lib/image-store";

const APP_URL = process.env.NEXT_PUBLIC_APP_URL || "https://test.gerege.mn";
const DAN_URL = process.env.DAN_URL || "https://dan.gerege.mn";

// DAN Verify Full callback — receives a token, fetches full citizen data (with photo) from dan.gerege.mn
export async function GET(req: NextRequest) {
  const token = req.nextUrl.searchParams.get("token") || "";

  if (!token) {
    return NextResponse.redirect(new URL("/auth/login?error=dan_no_token", APP_URL));
  }

  try {
    // Fetch full citizen data from dan.gerege.mn using one-time token
    const res = await fetch(`${DAN_URL}/api/citizen?token=${encodeURIComponent(token)}`, {
      cache: "no-store",
    });

    if (!res.ok) {
      return NextResponse.redirect(new URL("/auth/login?error=dan_token_expired", APP_URL));
    }

    const data = await res.json();
    const citizen = data.citizen || {};

    if (!citizen.reg_no) {
      return NextResponse.redirect(new URL("/auth/login?error=dan_no_data", APP_URL));
    }

    // Store full citizen data (with image) in a temporary server-side endpoint
    // We pass non-image fields as query params and image via a fetch endpoint
    const resultURL = new URL("/auth/dan-result-full", APP_URL);
    let imageBase64 = "";
    for (const [key, value] of Object.entries(citizen)) {
      if (key === "image" && typeof value === "string") {
        imageBase64 = value;
      } else if (typeof value === "string" && value) {
        resultURL.searchParams.set(key, value);
      }
    }

    // If there's an image, store it temporarily and pass a fetch key
    if (imageBase64) {
      const imgKey = Math.random().toString(36).slice(2, 18);
      imageStore.set(imgKey, { data: imageBase64, expires: Date.now() + 5 * 60 * 1000 });
      resultURL.searchParams.set("img_key", imgKey);
    }

    return NextResponse.redirect(resultURL);
  } catch (err) {
    console.error("DAN full callback error:", err);
    return NextResponse.redirect(new URL("/auth/login?error=dan_fetch_failed", APP_URL));
  }
}

