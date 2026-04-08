import { NextRequest, NextResponse } from "next/server";

const APP_URL = process.env.NEXT_PUBLIC_APP_URL || "https://test.gerege.mn";

// DAN Verify callback — dan.gerege.mn redirects here with citizen data as query params
export async function GET(req: NextRequest) {
  const params = req.nextUrl.searchParams;
  const regNo = params.get("reg_no") || "";

  if (!regNo) {
    return NextResponse.redirect(new URL("/auth/login?error=dan_failed", APP_URL));
  }

  const resultURL = new URL("/auth/dan-result", APP_URL);
  params.forEach((value, key) => {
    resultURL.searchParams.set(key, value);
  });

  return NextResponse.redirect(resultURL);
}
