import { NextRequest, NextResponse } from "next/server";
import {
  SESSION_COOKIE,
  SESSION_MAX_AGE,
  createSessionValue,
  constantTimeEqual,
} from "@/lib/session";

export async function POST(req: NextRequest) {
  const form = await req.formData();
  const password = (form.get("password") as string) || "";
  const expected = process.env.VERIFY_WEB_ADMIN_KEY || "";

  const ok = !!expected && constantTimeEqual(password, expected);
  const target = ok ? "/admin" : "/auth/login?error=1";

  // Build redirect URL from forwarded host (nginx) instead of internal 0.0.0.0
  const host = req.headers.get("x-forwarded-host") || req.headers.get("host") || "";
  const proto = req.headers.get("x-forwarded-proto") || "https";
  const baseUrl = `${proto}://${host}`;

  const res = NextResponse.redirect(new URL(target, baseUrl), { status: 303 });
  if (ok) {
    const value = await createSessionValue();
    res.cookies.set(SESSION_COOKIE, value, {
      httpOnly: true,
      secure: proto === "https",
      sameSite: "lax",
      path: "/",
      maxAge: SESSION_MAX_AGE,
    });
  }
  return res;
}
