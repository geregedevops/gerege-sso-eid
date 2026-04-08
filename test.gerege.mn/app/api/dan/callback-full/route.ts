import { NextRequest, NextResponse } from "next/server";

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

    // Pass citizen data to result page via query params
    // For the image, we store it in a cookie since it's too large for URL
    const resultURL = new URL("/auth/dan-result-full", APP_URL);
    for (const [key, value] of Object.entries(citizen)) {
      if (key !== "image" && typeof value === "string" && value) {
        resultURL.searchParams.set(key, value);
      }
    }

    const response = NextResponse.redirect(resultURL);

    // Store image in a short-lived cookie (base64, ~50KB typical)
    if (citizen.image) {
      response.cookies.set("dan_photo", citizen.image, {
        maxAge: 300, // 5 minutes
        httpOnly: false,
        path: "/auth/dan-result-full",
        sameSite: "lax",
      });
    }

    return response;
  } catch (err) {
    console.error("DAN full callback error:", err);
    return NextResponse.redirect(new URL("/auth/login?error=dan_fetch_failed", APP_URL));
  }
}
