import { NextRequest, NextResponse } from "next/server";

// DAN Verify callback — dan.gerege.mn redirects here with citizen data as query params
export async function GET(req: NextRequest) {
  const params = req.nextUrl.searchParams;
  const regNo = params.get("reg_no") || "";
  const givenName = params.get("given_name") || "";
  const familyName = params.get("family_name") || "";

  if (!regNo) {
    return NextResponse.redirect(
      new URL("/auth/login?error=dan_failed", req.url)
    );
  }

  // Redirect to a page showing the result
  const resultURL = new URL("/auth/dan-result", req.url);
  params.forEach((value, key) => {
    resultURL.searchParams.set(key, value);
  });

  return NextResponse.redirect(resultURL);
}
